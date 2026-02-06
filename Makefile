.PHONY: help dev dev-backend dev-frontend build install clean

help:
	@echo "Doodle Clone - Commands:"
	@echo ""
	@echo "  make dev           - Start backend + frontend (development)"
	@echo "  make dev-backend   - Start backend only"
	@echo "  make dev-frontend  - Start frontend only"
	@echo "  make build         - Build all"
	@echo "  make install        - Install dependencies"
	@echo "  make clean         - Clean build artifacts"

dev:
	@echo "Starting backend and frontend..."
	@make -C backend run &
	@make -C frontend dev

dev-backend:
	@echo "Starting backend..."
	@make -C backend run

dev-frontend:
	@echo "Starting frontend..."
	@make -C frontend dev

build:
	@echo "Building backend..."
	@make -C backend build
	@echo "Building frontend..."
	@make -C frontend build

install:
	@echo "Installing backend dependencies..."
	@make -C backend deps
	@echo "Installing frontend dependencies..."
	@make -C frontend install

clean:
	@echo "Cleaning..."
	@make -C backend clean
	@make -C frontend clean
