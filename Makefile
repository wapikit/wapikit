GOPATH := $(HOME)/go
GOBIN  := $(GOPATH)/bin
ATLAS  ?= /usr/local/bin/atlas
STUFFBIN ?= $(GOBIN)/stuffbin
JET    ?= $(GOBIN)/jet
OPI_CODEGEN ?= $(GOBIN)/oapi-codegen
PNPM ?= $(shell command -v pnpm 2> /dev/null)
FRONTEND_DIR := ./frontend

FRONTEND_BUILD_DIR := $(FRONTEND_DIR)/.next
BIN := wapikit

$(ATLAS):
	curl -sSf https://atlasgo.sh | sh -s -- --yes 

$(JET):
	go install github.com/go-jet/jet/v2/cmd/jet@latest

$(STUFFBIN):
	go install github.com/knadh/stuffbin/...

$(OPI_CODEGEN):
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

$(PNPM):
	curl -fsSL https://get.pnpm.io/install.sh | sh -

FRONTEND_MODULES = frontend/node_modules

FRONTEND_DEPS = \
	$(FRONTEND_MODULES) \
	frontend/package.json \
	frontend/tsconfig.json \
	frontend/.prettierrc.json \
	frontend/tailwind.config.ts \
	$(shell find frontend/src frontend/public frontend/src -type f)

$(FRONTEND_MODULES): frontend/package.json frontend/pnpm-lock.yaml
	cd frontend && $(PNPM) install
	touch -c $(FRONTEND_MODULES)

.PHONY: $(FRONTEND_DEPS) frontend-codegen
frontend-codegen: $(PNPM)
	cd $(FRONTEND_DIR) && $(PNPM) install && $(PNPM) run codegen

.PHONY: backend-codegen
backend-codegen: $(OPI_CODEGEN)
	$(OPI_CODEGEN) -package handlers -generate types -o handlers/types.gen.go spec.openapi.yaml

.PHONY: codegen
codegen: backend-codegen frontend-codegen

.PHONY: build-frontend
build-frontend: frontend-codegen
	cd $(FRONTEND_DIR) && $(PNPM) run build

STATIC := config.toml.sample \
	frontend/out:/ \

$(BIN): $(shell find . -type f -name "*.go") go.mod go.sum
	CGO_ENABLED=0 go build -o ${BIN} -ldflags="-s -w" cmd/*.go

# ! TODO: add build target
.PHONY: build-backend
build-backend: $(BIN)

.PHONY: pack-bin
pack-bin: build-frontend $(BIN) $(STUFFBIN)
	$(STUFFBIN) -a stuff -in $(BIN) -out ${BIN} ${STATIC}


# build everything frotnend and backend using stuffbin into ./wapikit
.PHONY: build
build: $(STUFFBIN) build-backend build-frontend pack-bin

.PHONY: run_frontend
run_frontend: frontend-codegen 
	cd $(FRONTEND_DIR) && $(PNPM) install && $(PNPM) run dev

.PHONY: db-migrate
db-migrate: $(ATLAS)
	$(ATLAS) migrate diff --env local

.PHONY: db-apply
db-apply: $(ATLAS)
	$(ATLAS) migrate apply --env local

.PHONY: db-gen
db-gen: $(JET)
	$(JET) -dsn=postgres://sarthakjdev@localhost:5432/wapikit?sslmode=disable -path=./.db-generated && mv ./.db-generated/wapikit/public/** ./.db-generated && rm -rf ./.db-generated/wapikit

.PHONY: db-init
db-init: db-apply

# dev mode targets for misc tasks
.PHONY: format
format: 
	 go fmt ./... && cd $(FRONTEND_DIR) && $(PNPM) run pretty 

.PHONY: api-doc
api-doc: $(PNPM)
	pnpm dlx @mintlify/scraping@latest openapi-file ./spec.openapi.yaml -o docs.wapikit.com/api-reference