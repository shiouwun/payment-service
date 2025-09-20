# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=payment-service
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_IMAGE=payment-service
DOCKER_TAG=latest

# Database parameters
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=payment_service

.PHONY: all build clean test coverage deps run docker-build docker-run help

all: test build ## Run tests and build the application

build: ## Build the application
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server

clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

test: ## Run unit tests
	$(GOTEST) -v ./...

test-coverage: ## Run tests with coverage
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

integration-test: ## Run integration tests
	$(GOTEST) -tags=integration -v ./...

deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) verify

deps-update: ## Update dependencies
	$(GOMOD) tidy

run: ## Run the application
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/server
	./$(BINARY_NAME)

run-dev: ## Run the application in development mode
	$(GOCMD) run cmd/server/main.go

# Docker commands
docker-build: ## Build Docker image
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Run Docker container
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

docker-compose-up: ## Start services with docker-compose
	docker-compose up -d

docker-compose-down: ## Stop services with docker-compose
	docker-compose down

docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

# Database commands
db-up: ## Start database with docker-compose
	docker-compose up -d postgres

db-migrate: ## Run database migrations
	psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -f scripts/migrations/001_initial_schema.sql

db-reset: ## Reset database (WARNING: This will drop all data)
	docker-compose down postgres
	docker volume rm payment-service_postgres_data || true
	docker-compose up -d postgres
	sleep 10
	make db-migrate

# Linting and formatting
lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	$(GOCMD) fmt ./...

vet: ## Run go vet
	$(GOCMD) vet ./...

# Development helpers
dev-setup: ## Setup development environment
	cp .env.example .env
	make docker-compose-up
	sleep 15
	make db-migrate

install-tools: ## Install development tools
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Production build
build-linux: ## Build for Linux
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./cmd/server

# Help
help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'