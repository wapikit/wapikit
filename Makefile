GOPATH := $(HOME)/go
GOBIN  := $(GOPATH)/bin
ATLAS  ?= /usr/local/bin/atlas
STUFFBIN ?= $(GOBIN)/stuffbin
JET    ?= $(GOBIN)/jet
PNPM ?= $(shell command -v pnpm 2> /dev/null)
FRONTEND_DIR := ./frontend
FRONTEND_BUILD_DIR := $(FRONTEND_DIR)/.next
BIN := ./wapikit

$(ATLAS):
	curl -sSf https://atlasgo.sh | sh -s -- --yes 

$(JET):
	go install github.com/go-jet/jet/v2/cmd/jet@latest

$(STUFFBIN):
	go install github.com/knadh/stuffbin/...

$(PNPM):
	curl -fsSL https://get.pnpm.io/install.sh | sh -


.PHONY: frontend-codegen
frontend-codegen: $(PNPM)
	cd $(FRONTEND_DIR) && $(PNPM) install && $(PNPM) run codegen

.PHONY: build-frontend
build-frontend: frontend-codegen
	cd $(FRONTEND_DIR) && $(PNPM) install && $(PNPM) run build

# $(BIN): 
# 	go build -o $(BIN) cmd/*.go

# ! TODO: add build target
# .PHONY: build-backend
# build-backend: $(BIN)

.PHONY: build
build: build-frontend build-backend

# .PHONY: pack-bin
# pack-bin: $(STUFFBIN) 
# 	$(STUFFBIN) -in ./wapikit -out ./dist/wapikit -a

.PHONY: dist
# dist: $(STUFFBIN) build pack-bin

.PHONY: run_frontend
run_frontend: frontend-codegen 
	cd $(FRONTEND_DIR) && $(PNPM) install && $(PNPM) run dev

.PHONY: db-migrate
db-migrate: $(ATLAS)
	$(ATLAS) migrate diff --env gorm --var "database_url=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

.PHONY: db-apply
db-apply: $(ATLAS)
	$(ATLAS) migrate apply --env gorm --var "database_url=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

.PHONY: db-gen
db-gen: $(JET)
	$(JET) -dsn=postgres://sarthakjdev@localhost:5432/wapikit?sslmode=disable -path=./.db-generated


# ! TODO: add pre and post build targets
# ! TODO: add frontend builder, linter and prettier target
# ! TODO: add backend builder target
# ! TODO: add docker build target
# ! TODO: add binary build target
# ! TODO: development environment setup target