# Load environment variables from .env file if it exists
# https://stackoverflow.com/a/70663753
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Makefile for managing database migrations and running the Go application

# Variables
GOOSE_CMD = goose
MIGRATIONS_DIR = migrations
GO_CMD = go

# Construct the database URL using environment variables
DB_URL = postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

# Targets
.PHONY: migrate-up migrate-down migrate-status run

migrate-up:
	@echo "Applying up migrations..."
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" up

migrate-down:
	@echo "Reverting the last migration..."
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" down

migrate-status:
	@echo "Checking migration status..."
	$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)" status

run:
	@echo "Running the Go application..."
	-$(GO_CMD) run ./cmd/fitbyte

create-migration:
	@if [ -z "$(name)" ]; then \
		echo "⚠️  WARNING: You didn't provide a name for table"; \
	else \
		echo "Creating migration '$(name)'..."; \
		$(GOOSE_CMD) -dir $(MIGRATIONS_DIR) create create_$(name)_table sql; \
	fi