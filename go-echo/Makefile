# Makefile for Go Echo API
# Common commands for development and deployment

.PHONY: help build run test clean docker-build docker-run docker-compose docker-clean

# Default target
help:
	@echo "Available commands:"
	@echo "  build         - Build the application"
	@echo "  run           - Run the application locally"
	@echo "  test          - Run all tests"
	@echo "  test-cover    - Run tests with coverage"
	@echo "  clean         - Clean build artifacts"
	@echo "  docker-build  - Build Docker image"
	@echo "  docker-run    - Run with Docker"
	@echo "  docker-compose- Run with Docker Compose"
	@echo "  docker-clean  - Clean Docker resources"
	@echo "  docker-logs   - View Docker logs"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/server main.go
	@echo "Build completed!"

# Run the application locally
run:
	@echo "Running application..."
	go run main.go

# Run tests
test:
	@echo "Running tests..."
	go test ./...

# Run tests with coverage
test-cover:
	@echo "Running tests with coverage..."
	go test -cover ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	@echo "Clean completed!"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t go-echo-api .

docker-run:
	@echo "Running with Docker..."
	docker run -d --name go-echo-api -p 8080:8080 go-echo-api

docker-compose:
	@echo "Running with Docker Compose..."
	docker-compose up -d

docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down 2>/dev/null || true
	docker rm -f go-echo-api 2>/dev/null || true
	docker rmi go-echo-api 2>/dev/null || true

docker-logs:
	@echo "Viewing Docker logs..."
	docker-compose logs -f app

# Development helpers
dev-setup:
	@echo "Setting up development environment..."
	cp config.yaml.example config.yaml
	@echo "Development setup completed!"

# Database helpers
db-migrate:
	@echo "Running database migrations..."
	psql -d users_db -f database/migrations/001_create_users_table.sql

# Quick start
quick-start: dev-setup docker-compose
	@echo "Application started! Visit http://localhost:8080" 