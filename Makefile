.PHONY: all build run test clean generate swagger migrate migrate-dry docker-up docker-down fmt lint help

# Variables
APP_NAME := super-shiharai-api
GO := go
DOCKER_COMPOSE := docker compose

# Default target
all: generate build

# Build the application
build:
	$(GO) build -o bin/$(APP_NAME) ./cmd/api

# Run the application locally
run:
	$(GO) run ./cmd/api

# Run tests
test:
	$(GO) test -v ./...

# Run tests with coverage
test-coverage:
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out coverage.html

# Generate all (mocks + sqlc)
generate:
	$(GO) generate ./...
	cd db && sqlc generate

# Generate mocks only
generate-mocks:
	$(GO) generate ./...

# Generate sqlc only
generate-sqlc:
	cd db && sqlc generate

# Generate swagger docs
swagger:
	swag init -g cmd/api/main.go -o docs/swagger

# Run database migration (dry-run)
migrate-dry:
	psqldef -U postgres -h localhost -p 5432 super_shiharai --dry-run < db/schema.sql

# Run database migration
migrate:
	psqldef -U postgres -h localhost -p 5432 super_shiharai < db/schema.sql

# Start Docker containers
docker-up:
	$(DOCKER_COMPOSE) up -d

# Stop Docker containers
docker-down:
	$(DOCKER_COMPOSE) down

# Start Docker containers with rebuild
docker-up-build:
	$(DOCKER_COMPOSE) up -d --build

# View Docker logs
docker-logs:
	$(DOCKER_COMPOSE) logs -f

# Format code (golangci-lint v2)
fmt:
	golangci-lint fmt ./...

# Run linter (golangci-lint v2)
lint:
	golangci-lint run ./...

# Install development tools
tools:
	$(GO) install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	$(GO) install github.com/sqldef/sqldef/cmd/psqldef@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install go.uber.org/mock/mockgen@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Help
help:
	@echo "Available targets:"
	@echo "  all            - Generate and build"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application locally"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Clean build artifacts"
	@echo "  generate       - Generate mocks and sqlc code"
	@echo "  generate-mocks - Generate mocks only"
	@echo "  generate-sqlc  - Generate sqlc code only"
	@echo "  swagger        - Generate swagger documentation"
	@echo "  migrate-dry    - Preview database migration"
	@echo "  migrate        - Run database migration"
	@echo "  docker-up      - Start Docker containers"
	@echo "  docker-down    - Stop Docker containers"
	@echo "  docker-up-build- Rebuild and start Docker containers"
	@echo "  docker-logs    - View Docker logs"
	@echo "  fmt            - Format code (golangci-lint v2)"
	@echo "  lint           - Run linter (golangci-lint v2)"
	@echo "  tools          - Install development tools"
	@echo "  help           - Show this help message"
