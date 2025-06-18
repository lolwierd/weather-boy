include .env

BINARY_NAME=weatherboyapi
MIGRATE_DSN=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: run
run: build
	@set -a && source .env && set +a && ./tmp/$(BINARY_NAME)

.PHONY: build
build:
	go build -o tmp/$(BINARY_NAME) ./cmd/weatherboyapi

.PHONY: test
test:
        go test -v ./...

.PHONY: migrate
migrate:
        migrate -path migrations -database $(MIGRATE_DSN) up

.PHONY: dev
dev:
        docker compose --env-file .env -f docker-compose.dev.yml up --build

.PHONY: clean
clean:
	go clean
	rm -f tmp/$(BINARY_NAME)
