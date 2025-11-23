.PHONY: help dev dev-backend dev-frontend build run clean

help:
	@echo "Available commands:"
	@echo "  make dev          - Run both backend and frontend in development mode"
	@echo "  make dev-backend  - Run backend server only"
	@echo "  make dev-frontend - Run frontend only"
	@echo "  make build        - Build both backend and frontend for production"
	@echo "  make run          - Run the production build"
	@echo "  make clean        - Clean build artifacts"

dev:
	@echo "Starting development servers..."
	@echo "Backend will run on http://localhost:8080"
	@echo "Frontend will run on http://localhost:3000"
	@$(MAKE) -j2 dev-backend dev-frontend

dev-backend:
	@echo "Starting Go backend server..."
	@cd . && go run cmd/server/main.go

dev-frontend:
	@echo "Starting React frontend..."
	@cd web && npm start

build:
	@echo "Building backend..."
	@go build -o bin/email-client cmd/server/main.go
	@echo "Building frontend..."
	@cd web && npm run build
	@echo "Build complete! Binary: bin/email-client, Frontend: web/build/"

run:
	@./bin/email-client

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -rf web/build/
	@rm -rf web/node_modules/
	@echo "Clean complete!"

install-deps:
	@echo "Installing Go dependencies..."
	@go mod download
	@echo "Installing frontend dependencies..."
	@cd web && npm install
	@echo "Dependencies installed!"
