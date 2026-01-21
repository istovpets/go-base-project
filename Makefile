GOLANGCI_LINT_VERSION := 2.6.0
GOOSE=goose
MIGRATIONS_DIR=./migrations/postgres/main
DB_CONN_STR=postgres://postgres:postgres@localhost:5432/dadata_v2?sslmode=disable
DB_DRIVER=postgres

# lint

./bin/golangci-lint:
	mkdir -p ./bin
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "./bin" "v${GOLANGCI_LINT_VERSION}";

.PHONY: lint
lint:
	$(MAKE) ./bin/golangci-lint
	./bin/golangci-lint run ./...

# goose

.PHONY: migrate_new
migrate_new:
	$(GOOSE) -dir $(MIGRATIONS_DIR) create new sql

.PHONY: migrate_up
migrate_up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DB_DRIVER) "$(DB_CONN_STR)" up

# shugar

.PHONY: validateall
validate:
	go vet ./...
	$(MAKE) lint
	go test -race ./...
