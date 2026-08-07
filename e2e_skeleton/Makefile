.PHONY: build run test migrate-up migrate-down clean help

# Environment variables
include .env
export

APP_NAME=github.com/test/skeleton
BUILD_DIR=bin
MAIN_PATH=cmd/server/main.go
DB_DSN="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

build: ## Build the application
	@echo "Building ${APP_NAME}..."
	@go build -o ${BUILD_DIR}/${APP_NAME} ${MAIN_PATH}

run: ## Run the application
	@echo "Starting ${APP_NAME}..."
	@go run ${MAIN_PATH}

test: ## Run the test suite
	@go test -v -race ./...

migrate-up: ## Run database migrations up
	@echo "Running migrations up..."
	@migrate -path migrations -database ${DB_DSN} up

migrate-down: ## Run database migrations down
	@echo "Running migrations down..."
	@migrate -path migrations -database ${DB_DSN} down

clean: ## Clean build artifacts
	@rm -rf ${BUILD_DIR}

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)
