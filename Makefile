GOPATH := $(HOME)/go
GOBIN  := $(GOPATH)/bin
ATLAS  ?= /usr/local/bin/atlas
STUFFBIN ?= $(GOBIN)/stuffbin
JET    ?= $(GOBIN)/jet

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

db-migrate: $(ATLAS)
	$(ATLAS) migrate diff --env gorm --var "database_url=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

db-apply: $(ATLAS)
	$(ATLAS) migrate apply --env gorm --var "database_url=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable"

db-gen: $(JET)
	$(JET) -dsn=postgres://sarthakjdev@localhost:5432/wapikit?sslmode=disable -path=./.generated
