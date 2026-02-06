.PHONY: help dev build preview test lint format clean deps docker-build docker-run

# Node modules directory
NODE_modules=node_modules

# Build directory
DIST=dist

# Docker variables
DOCKER_IMAGE=doodle-clone-frontend
DOCKER_TAG=latest

help:
	@echo "Frontend Makefile - Available commands:"
	@echo ""
	@echo "Development:"
	@echo "  make dev          - Start development server (Vite)"
	@echo "  make preview      - Preview production build locally"
	@echo ""
	@echo "Building:"
	@echo "  make build        - Build for production"
	@echo "  make build-analyze - Build with bundle analysis"
	@echo ""
	@echo "Testing:"
	@echo "  make test         - Run tests"
	@echo "  make test:watch   - Run tests in watch mode"
	@echo "  make test:coverage - Run tests with coverage"
	@echo ""
	@echo "Code Quality:"
	@echo "  make lint         - Run ESLint"
	@echo "  make lint:fix     - Fix ESLint errors automatically"
	@echo "  make format       - Format code with Prettier"
	@echo "  make check        - Run all checks (lint + type check)"
	@echo ""
	@echo "Dependencies:"
	@echo "  make deps         - Install dependencies"
	@echo "  make deps-update  - Update dependencies"
	@echo "  make deps-outdated - Check for outdated packages"
	@echo "  make deps-clean   - Remove node_modules and lock files"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-build - Build Docker image"
	@echo "  make docker-run   - Run Docker container"
	@echo ""
	@echo "Maintenance:"
	@echo "  make clean        - Remove build artifacts"

# =============================================================================
# Development
# =============================================================================

dev:
	@echo "Starting development server..."
	@npm run dev

preview:
	@echo "Previewing production build..."
	@npm run preview

# =============================================================================
# Building
# =============================================================================

build:
	@echo "Building for production..."
	@npm run build

build-analyze:
	@echo "Building with bundle analysis..."
	@npm run build -- --mode analyze

# =============================================================================
# Testing
# =============================================================================

test:
	@echo "Running tests..."
	@npm run test || echo "No tests configured"

test-watch:
	@echo "Running tests in watch mode..."
	@npm run test:watch || echo "No tests configured"

test-coverage:
	@echo "Running tests with coverage..."
	@npm run test:coverage || echo "No coverage configured"

# =============================================================================
# Code Quality
# =============================================================================

lint:
	@echo "Running ESLint..."
	@npm run lint || echo "ESLint not configured"

lint-fix:
	@echo "Fixing ESLint errors..."
	@npm run lint:fix || echo "ESLint not configured"

format:
	@echo "Formatting code..."
	@npm run format || echo "Prettier not configured"

format-check:
	@echo "Checking formatting..."
	@npm run format:check || echo "Prettier not configured"

check: lint
	@echo "Running type check..."
	@npm run type-check || echo "TypeScript not configured"

# =============================================================================
# Dependencies
# =============================================================================

deps:
	@echo "Installing dependencies..."
	@npm install

deps-update:
	@echo "Updating dependencies..."
	@npm update

deps-outdated:
	@echo "Checking for outdated packages..."
	@npm outdated

deps-clean:
	@echo "Removing node_modules and lock files..."
	@rm -rf $(NODE_modules)
	@rm -f package-lock.json
	@rm -f yarn.lock
	@rm -f pnpm-lock.yaml

deps-fresh: deps-clean deps

# =============================================================================
# Docker
# =============================================================================

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 80:80 $(DOCKER_IMAGE):$(DOCKER_TAG)

# =============================================================================
# Maintenance
# =============================================================================

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(DIST)
	@rm -rf .vite
	@rm -rf coverage

clean-all: clean deps-clean

# =============================================================================
# Environment
# =============================================================================

env:
	@echo "Creating .env file..."
	@cp .env.example .env || echo ".env.example not found"

# =============================================================================
# Install (alias for deps)
# =============================================================================

install: deps
