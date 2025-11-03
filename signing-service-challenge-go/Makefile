.PHONY: help build test test-verbose test-coverage run clean install lint fmt check

# Variables
BINARY_NAME=signing-service
BUILD_DIR=./build
GO=go
GOTEST=$(GO) test
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOGET=$(GO) get
GOFMT=$(GO) fmt
GOVET=$(GO) vet

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	@echo "Installing dependencies..."
	$(GOGET) github.com/gin-gonic/gin
	$(GOGET) github.com/google/uuid
	$(GO) mod tidy
	@echo "Dependencies installed!"

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) -v
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

run: build ## Build and run the application
	@echo "Starting $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	$(GOTEST) -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-domain: ## Run domain tests only
	@echo "Running domain tests..."
	$(GOTEST) -v ./domain

test-api: ## Run API tests only
	@echo "Running API tests..."
	$(GOTEST) -v ./api

test-crypto: ## Run crypto tests only
	@echo "Running crypto tests..."
	$(GOTEST) -v ./crypto

test-persistence: ## Run persistence tests only
	@echo "Running persistence tests..."
	$(GOTEST) -v ./persistence

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) ./...
	@echo "Code formatted!"

lint: ## Run linter (requires golangci-lint)
	@echo "Running linter..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with: brew install golangci-lint"; \
		$(GOVET) ./...; \
	fi

vet: ## Run go vet
	@echo "Running go vet..."
	$(GOVET) ./...

check: fmt vet test ## Format, vet and test

clean: ## Clean build artifacts
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete!"

rebuild: clean build ## Clean and rebuild

all: clean install build test ## Install dependencies, build and test

dev: ## Run in development mode (with auto-reload if air is installed)
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not installed. Running normally..."; \
		echo "Install air with: go install github.com/cosmtrek/air@latest"; \
		make run; \
	fi

docker-build: ## Build docker image
	@echo "Building docker image..."
	docker build -t $(BINARY_NAME):latest .

docker-run: ## Run docker container
	@echo "Running docker container..."
	docker run -p 8080:8080 $(BINARY_NAME):latest

# Show current Go version and environment
info: ## Show Go environment info
	@echo "Go version:"
	@$(GO) version
	@echo "\nGo environment:"
	@$(GO) env

