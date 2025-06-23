.PHONY: build run test test-coverage test-integration clean deps migrate-up migrate-down docker-build docker-run

# Variables
APP_NAME=ecom-api
BUILD_DIR=bin
DOCKER_IMAGE=ecom-golang-clean-architecture

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) cmd/api/main.go

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run cmd/api/main.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	@go test -v -tags=integration ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Database migrations
migrate-up:
	@echo "Running database migrations..."
	@go run cmd/migrate/main.go up

migrate-down:
	@echo "Rolling back database migrations..."
	@go run cmd/migrate/main.go down

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(DOCKER_IMAGE)

# Development tools
dev-tools:
	@echo "Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

# Generate swagger docs
swagger:
	@echo "Generating swagger documentation..."
	@swag init -g cmd/api/main.go

# Live reload for development
dev:
	@echo "Starting development server with live reload..."
	@air

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golangci-lint run

# Security check
security:
	@echo "Running security check..."
	@gosec ./...

# All checks (format, lint, test)
check: fmt lint test

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build the application"
	@echo "  run            - Run the application"
	@echo "  deps           - Install dependencies"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  test-integration - Run integration tests"
	@echo "  clean          - Clean build artifacts"
	@echo "  migrate-up     - Run database migrations"
	@echo "  migrate-down   - Rollback database migrations"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  dev-tools      - Install development tools"
	@echo "  swagger        - Generate swagger documentation"
	@echo "  dev            - Start development server with live reload"
	@echo "  fmt            - Format code"
	@echo "  lint           - Lint code"
	@echo "  security       - Run security check"
	@echo "  check          - Run all checks (format, lint, test)"
	@echo "  help           - Show this help message"
