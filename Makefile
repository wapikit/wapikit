GOPATH := $(HOME)/go
GOBIN  := $(GOPATH)/bin
ATLAS  ?= /usr/local/bin/atlas
STUFFBIN ?= $(GOBIN)/stuffbin
JET    ?= $(GOBIN)/jet
PNPM ?= $(shell command -v pnpm 2> /dev/null)
FRONTEND_DIR := ./frontend
FRONTEND_BUILD_DIR := $(FRONTEND_DIR)/.next

$(FRONTEND_BUILD_DIR):
	cd $(FRONTEND_DIR) && $(PNPM) run build

# Database Connection (Sensitive information, don't commit to version control)
DB_USER := your_user_name
DB_PASSWORD := your_password
DB_NAME := your_database_name
DB_HOST := localhost
DB_PORT := 5432


$(ATLAS):
	curl -sSf https://atlasgo.sh | sh -s -- --yes 

$(JET):
	go install github.com/go-jet/jet/v2/cmd/jet@latest

$(STUFFBIN):
	go install github.com/knadh/stuffbin/...

$(PNPM):
	curl -fsSL https://get.pnpm.io/install.sh | sh -


fontend_build: frontend_codegen
	cd $(FRONTEND_DIR) && $(PNPM) install && $(PNPM) run build

frontend_codegen: $(PNPM)
	cd $(FRONTEND_DIR) && $(PNPM) install && $(PNPM) run codegen

build: frontend_build db_gen 

db-migrate: $(ATLAS)
	$(ATLAS) migrate diff --env gorm --var "database_url=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

db-apply: $(ATLAS)
	$(ATLAS) migrate apply --env gorm --var "database_url=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

db-gen: $(JET)
	$(JET) -dsn=postgres://sarthakjdev@localhost:5432/wapikit?sslmode=disable -path=./.db-generated

dev: db-migrate db-apply db-gen
	go run ./cmd/main.go




# ! TODO: add pre and post build targets
# ! TODO: add frontend builder, linter and prettier target
# ! TODO: add backend builder target
# ! TODO: add docker build target
# ! TODO: add binary build target
# ! TODO: development environment setup target