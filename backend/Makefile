.PHONY: help build run test test-cover lint fmt vet clean docker-build docker-run deps swagger

# Variables
BINARY_NAME=doodle-clone
BUILD_DIR=bin
GO_FILES=$(shell find . -name '*.go' -type f)
GO_VERSION=1.24

# Docker variables
DOCKER_IMAGE=doodle-clone-backend
DOCKER_TAG=latest

help:
	@echo "Backend Makefile - Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make run           - Run the application (requires .env)"
	@echo "  make run-dev       - Run with hot reload using air"
	@echo ""
	@echo "Building:"
	@echo "  make build         - Build the application binary"
	@echo "  make build-linux   - Build for Linux"
	@echo "  make build-darwin   - Build for macOS"
	@echo ""
	@echo "Testing:"
	@echo "  make test          - Run tests"
	@echo "  make test-cover    - Run tests with coverage"
	@echo "  make test-race     - Run tests with race detection"
	@echo "  make bench         - Run benchmarks"
	@echo ""
	@echo "Documentation:"
	@echo "  make swagger       - Generate Swagger documentation"
	@echo ""
	@echo "Code Quality:"
	@echo "  make fmt           - Format Go code"
	@echo "  make vet           - Run go vet"
	@echo "  make lint          - Run linter"
	@echo "  make fix           - Fix linting issues automatically"
	@echo ""
	@echo "Database:"
	@echo "  make db-create     - Create database"
	@echo "  make db-drop       - Drop database"
	@echo "  make db-reset      - Reset database"
	@echo "  make db-console    - Open database console"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run Docker container"
	@echo ""
	@echo "Dependencies:"
	@echo "  make deps          - Download dependencies"
	@echo "  make deps-tidy     - Tidy dependencies"
	@echo "  make deps-update   - Update dependencies"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make mod           - Initialize Go module"

# =============================================================================
# Development
# =============================================================================

run:
	@if [ ! -f .env ]; then \
		echo "Error: .env file not found. Copy from .env.example first:"; \
		echo "  cp .env.example .env"; \
		exit 1; \
	fi
	@echo "Starting $(BINARY_NAME)..."
	@go run main.go

run-dev:
	@which air > /dev/null 2>&1 || (echo "Installing air..." && go install github.com/air-verse/air@latest)
	@air

# =============================================================================
# Building
# =============================================================================

build: $(BUILD_DIR)/$(BINARY_NAME)

$(BUILD_DIR)/$(BINARY_NAME): $(GO_FILES)
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux main.go

build-darwin:
	@echo "Building for macOS..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin main.go

build-arm:
	@echo "Building for ARM64..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-arm64 main.go

# =============================================================================
# Testing
# =============================================================================

test:
	@echo "Running tests..."
	@go test -v ./...

test-cover:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	@echo "Running tests with race detection..."
	@go test -race -v ./...

bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# =============================================================================
# Code Quality
# =============================================================================

fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@goimports -w . 2>/dev/null || true

vet:
	@echo "Running go vet..."
	@go vet ./...

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null 2>&1 || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	@golangci-lint run ./...

fix:
	@echo "Fixing linting issues..."
	@golangci-lint run --fix ./...

# =============================================================================
# Database helpers (requires psql)
# =============================================================================

DB_HOST?=localhost
DB_PORT?=5432
DB_NAME?=doodle_clone
DB_USER?=postgres
DB_PASSWORD?=postgres

db-create:
	@echo "Creating database..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -c "CREATE DATABASE $(DB_NAME);" || true

db-drop:
	@echo "Dropping database..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -c "DROP DATABASE IF EXISTS $(DB_NAME);"

db-reset: db-drop db-create
	@echo "Database reset complete"

db-console:
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME)

db-migrate:
	@echo "Running migrations..."
	@$(MAKE) run &
	@sleep 2
	@pkill -f "go run main.go" || true

# =============================================================================
# Docker
# =============================================================================

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# =============================================================================
# Dependencies
# =============================================================================

deps:
	@echo "Downloading dependencies..."
	@go mod download

deps-tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

mod:
	@echo "Initializing Go module..."
	@go mod init github.com/steph/doodle-clone 2>/dev/null || true
	@go mod tidy

# =============================================================================
# Maintenance
# =============================================================================

clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean

# =============================================================================
# Install tools
# =============================================================================

install-tools:
	@echo "Installing development tools..."
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Tools installed:"
	@echo "  air    - Hot reload for development"
	@echo "  golangci-lint - Linter"
	@echo "  goimports - Import management"
	@echo "  swag   - Swagger documentation generator"

# =============================================================================
# Swagger Documentation
# =============================================================================

swagger:
	@echo "Generating Swagger documentation..."
	@which swag > /dev/null 2>&1 || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	@swag init -g main.go -o docs
	@echo "Swagger documentation generated at docs/"

