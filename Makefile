GOPATH ?= $(HOME)/go
GOBIN  ?= $(GOPATH)/bin
ATLAS  ?= /usr/local/bin/atlas
STUFFBIN ?= $(GOBIN)/stuffbin
JET    ?= $(GOBIN)/jet
GOLANGCI_LINT ?= $(GOBIN)/golangci-lint
OPI_CODEGEN ?= $(GOBIN)/oapi-codegen
AIR ?= $(GOBIN)/air
PNPM ?= $(shell command -v pnpm 2> /dev/null)
FRONTEND_DIR = frontend
FRONTEND_BUILD_DIR = frontend/out
BIN := wapikit
BIN_MANAGED := wapikit-cloud

STATIC := config.toml.sample \
	frontend/out:/ \
	internal/database/migrations:/migrations \

FRONTEND_MODULES = frontend/node_modules

$(ATLAS):
	curl -sSf https://atlasgo.sh | sh -s -- --yes 

$(JET):
	go install github.com/go-jet/jet/v2/cmd/jet@latest

$(GOLANGCI_LINT):
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

$(STUFFBIN):
	go install github.com/knadh/stuffbin/...

$(AIR):
	go install github.com/cosmtrek/air@latest

$(OPI_CODEGEN):
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

.PHONY: install-pnpm
install-pnpm:
	@if ! command -v pnpm > /dev/null; then \
		echo "PNPM is not installed. Installing..."; \
		curl -fsSL https://get.pnpm.io/install.sh | sh -; \
	fi

$(PNPM): install-pnpm

FRONTEND_DEPS = \
	$(FRONTEND_MODULES) \
	frontend/package.json \
	frontend/pnpm-lock.yaml \
	frontend/tsconfig.json \
	frontend/.generated.ts \
	frontend/.eslintrc.js \
	frontend/.eslintignore \
	frontend/postcss.config.mjs \
	frontend/tailwind.config.ts \
	$(shell find frontend/src frontend/public -type f)

$(FRONTEND_MODULES): frontend/package.json frontend/pnpm-lock.yaml
	cd frontend && $(PNPM) install
	touch -c $(FRONTEND_MODULES)

# Check for DB_URL environment variable
.PHONY: check-db-url
check-db-url:
	@if [ -z "$$DB_URL" ]; then \
		echo "Error: DB_URL environment variable is not set."; \
		exit 1; \
	fi

.PHONY: frontend-codegen
frontend-codegen: $(PNPM)
	cd $(FRONTEND_DIR) && $(PNPM) install && $(PNPM) run codegen

$(FRONTEND_BUILD_DIR): $(FRONTEND_DEPS)
	cd frontend && $(PNPM) run build
	touch -c $(FRONTEND_BUILD_DIR)

.PHONY: build-frontend
build-frontend: $(FRONTEND_BUILD_DIR)

.PHONY: dev-backend
dev-backend: $(AIR)
	air -c .air.toml

.PHONY: dev-backend-docker
dev-backend-docker:
		CGO_ENABLED=0 go run -ldflags="-s -w " cmd/*.go --config=dev/config.toml

.PHONY: dev-frontend
dev-frontend: $(PNPM) $(FRONTEND_MODULES) 
	cd $(FRONTEND_DIR) && $(PNPM) run dev

.PHONY: dev
dev: dev-backend dev-frontend

.PHONY: backend-codegen
backend-codegen: $(OPI_CODEGEN)
	$(OPI_CODEGEN) -package api_types -generate types -o api/api_types/types.go spec.openapi.yaml

.PHONY: codegen
codegen: backend-codegen frontend-codegen


$(BIN): $(shell find . -type f -name "*.go") go.mod go.sum
	CGO_ENABLED=0 GOFLAGS="-tags=community_edition" go build -o ${BIN} -ldflags="-s -w" ./cmd/

$(BIN_MANAGED): $(shell find . -type f -name "*.go") go.mod go.sum
	CGO_ENABLED=0 GOFLAGS="-tags=managed_cloud" go build -o ${BIN_MANAGED} -ldflags="-s -w" ./cmd/

.PHONY: build-backend
build-backend: $(BIN)

.PHONY: build-cloud-edition-backend
build-cloud-edition-backend: $(BIN_MANAGED)

.PHONY: build-cloud-edition-frontend
build-cloud-edition-frontend: $(FRONTEND_DEPS)
	cd frontend && $(PNPM) run build:cloud
	touch -c $(FRONTEND_BUILD_DIR)

.PHONY: build-cloud-edition	
build-cloud-edition: build-cloud-edition-backend build-cloud-edition-frontend

.PHONY: dist
dist: build-frontend $(BIN) $(STUFFBIN)
	$(STUFFBIN) -a stuff -in $(BIN) -out ${BIN} ${STATIC}

.PHONY: build
build: build-backend build-frontend

.PHONY: run_frontend
run_frontend: frontend-codegen
	cd $(FRONTEND_DIR) && $(PNPM) install && $(PNPM) run dev

.PHONY: db-migrate
db-migrate: check-db-url $(ATLAS)
	$(ATLAS) migrate diff --env global --var DB_URL=$$DB_URL

.PHONY: cloud-db-migrate
cloud-db-migrate: check-db-url $(ATLAS)
	$(ATLAS) migrate diff --env managed_cloud --var DB_URL=$$DB_URL 

.PHONY: db-apply
db-apply: check-db-url $(ATLAS)
	$(ATLAS) migrate apply --env global --var DB_URL=$$DB_URL

.PHONY: cloud-db-apply
cloud-db-apply: check-db-url $(ATLAS)
	$(ATLAS) migrate apply --env managed_cloud --var DB_URL=$$DB_URL

.PHONY: db-gen
db-gen: check-db-url $(JET)
	$(JET) -dsn=$$DB_URL -path=./.db-generated && rm -rf ./.db-generated/model ./.db-generated/table ./.db-generated/enum && mv ./.db-generated/wapikit/public/** ./.db-generated && rm -rf ./.db-generated/wapikit

.PHONY: cloud-db-gen
cloud-db-gen: check-db-url $(JET)
	$(JET) -dsn=$$DB_URL -path=./.enterprise/.db-generated && rm -rf ./.enterprise/.db-generated/model ./.enterprise/.db-generated/table ./.enterprise/.db-generated/enum && mv ./.enterprise/.db-generated/wapikit/public/** ./.enterprise/.db-generated && rm -rf ./.enterprise/.db-generated/wapikit

.PHONY: db-init
db-init: db-apply

# dev mode targets for misc tasks
.PHONY: format
format: 
	 go fmt ./... && cd $(FRONTEND_DIR) && $(PNPM) run pretty 

.PHONY: lint
lint: $(JET) $(GOLANGCI_LINT) $(PNPM)
	 $(GOLANGCI_LINT) run --build-tags managed_cloud && cd $(FRONTEND_DIR) && $(PNPM) run lint 

.PHONY: api-doc
api-doc: $(PNPM)
	pnpm dlx @mintlify/scraping@latest openapi-file ./spec.openapi.yaml -o docs.wapikit.com/api-reference

