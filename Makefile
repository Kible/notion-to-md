# note: call scripts from /scripts

# Define Go binary path explicitly
GOBIN := $(shell go env GOPATH)/bin

# Define tools with full paths
GENQLIENT := $(GOBIN)/genqlient
GQLGEN := $(GOBIN)/gqlgen
CFF := $(GOBIN)/cff
OAPI_CODEGEN := $(GOBIN)/oapi-codegen

# Docker Compose command
COMPOSE := docker compose

# App name and binary location
APP_NAME := app
BIN_DIR := ./bin
BIN := $(BIN_DIR)/$(APP_NAME)

.PHONY: all build clean test run lint docker help tools

# Default target
all: lint test build

# Build the application
build:
	@echo "Building..."
	@go build -o $(BIN) ./cmd/app

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)
	@go clean

# Run the application
run:
	@echo "Running app..."
	@go run ./cmd/app

# Lint the code
lint:
	@echo "Linting..."
	@./scripts/lint.sh

# Format the code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Build docker image
docker:
	@echo "Building Docker image..."
	@docker build -t github.com/Kible/kible-backend:latest .

# Check for updates in dependencies
deps-check:
	@echo "Checking for dependency updates..."
	@./scripts/check-deps.sh

# Update dependencies
deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# Generate mocks for testing
mocks:
	@echo "Generating mocks..."
	@./scripts/generate-mocks.sh

# Run security scan
security:
	@echo "Running security scan..."
	@./scripts/security-scan.sh

# Install required tools
tools:
	@echo "Installing required tools..."
	@./scripts/install-tools.sh

# Show help
help:
	@echo "Available targets:"
	@echo "  all          - Run lint, test, and build"
	@echo "  build        - Build the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Remove build artifacts"
	@echo "  run          - Run the application"
	@echo "  lint         - Run linters"
	@echo "  fmt          - Format code"
	@echo "  docker       - Build Docker image"
	@echo "  deps-check   - Check for dependency updates"
	@echo "  deps-update  - Update dependencies"
	@echo "  mocks        - Generate test mocks"
	@echo "  security     - Run security scanning"
	@echo "  tools        - Install required development tools"
	@echo "  help         - Show this help message"