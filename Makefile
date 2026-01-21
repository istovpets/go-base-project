GOLANGCI_LINT_VERSION := 2.6.0
GOOSE=goose
MIGRATIONS_DIR=./migrations/postgres/main
DB_CONN_STR=postgres://postgres:postgres@localhost:5432/dadata_v2?sslmode=disable
DB_DRIVER=postgres

ifneq (,$(wildcard .env))
	include .env
	export
endif

# lint

./bin/golangci-lint:
	mkdir -p ./bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "./bin" "v${GOLANGCI_LINT_VERSION}";

.PHONY: lint
lint:
	$(MAKE) ./bin/golangci-lint
	./bin/golangci-lint run ./...

.PHONY: test
test:
	go vet ./...
	go test -race ./...

.PHONY: generate
generate:
	go generate ./...

# goose

.PHONY: migrate_new
migrate_new:
	$(GOOSE) -dir $(MIGRATIONS_DIR) create new sql

.PHONY: migrate_up
migrate_up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_CONN_STR)" up

# shugar

.PHONY: validate
validate:
	$(MAKE) lint
	$(MAKE) test

.PHONY: all
all:
	$(MAKE) generate
	$(MAKE) validate
