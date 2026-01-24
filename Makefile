# Common variables
GOLANGCI_LINT_VERSION := 2.6.0
MIGRATIONS_DIR=./migrations/postgres/main


# DB variables
DB_CONTAINER_NAME := local-postgres
DB_IMAGE := postgres:16-alpine

DB_USER := postgres
DB_PASSWORD := postgres
DB_NAME := base_project
DB_PORT := 6432

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable

PGDATA_DIR := $(PWD)/.pgdata

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

.PHONY: migrate-new
migrate-new:
	goose -dir $(MIGRATIONS_DIR) create new sql

.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

# DB

.PHONY: db-up
db-up: db-down
	docker run -d \
		--name $(DB_CONTAINER_NAME) \
		-e POSTGRES_USER=$(DB_USER) \
		-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
		-e POSTGRES_DB=$(DB_NAME) \
		-p $(DB_PORT):5432 \
		-v $(PGDATA_DIR):/var/lib/postgresql/data \
		$(DB_IMAGE)

.PHONY: db-wait
db-wait:
	@echo "Waiting for postgres to be ready..."
	@until docker exec $(DB_CONTAINER_NAME) pg_isready -U $(DB_USER) >/dev/null 2>&1; do \
		sleep 1; \
	done

.PHONY: db-start
db-start: db-up db-wait migrate-up

.PHONY: db-down
db-down:
	docker stop $(DB_CONTAINER_NAME) || true
	docker rm $(DB_CONTAINER_NAME) || true

.PHONY: db-clean
db-clean: db-down
	rm -rf $(PGDATA_DIR)

# shugar

.PHONY: dep-up
dep-up: db-start

.PHONY: validate
validate: lint test

.PHONY: all
all: generate validate
