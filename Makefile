# Scripts Tool Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Main binary
BINARY_NAME=scripts
BINARY_PATH=./$(BINARY_NAME)

# Test directories
TEST_DIRS=./tests

# Build the project
build:
	$(GOBUILD) -o $(BINARY_NAME) -v .

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Run all tests
test:
	$(GOTEST) -v $(TEST_DIRS)

# Run unit tests (config, script management, compilation logic)
test-unit:
	$(GOTEST) -v $(TEST_DIRS) -run "Test(Config|IsExecutable|MakeExecutable|CreateTestScript|AddScript|ScriptRemoval|BinaryRemoval|ScriptDirectoryStructure|InvalidScriptOperations)"

# Run integration tests (CLI commands)
test-integration:
	$(GOTEST) -v $(TEST_DIRS) -run "Test(CLI|Compile)"

# Run tests with coverage
test-coverage:
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic $(TEST_DIRS)
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run tests in short mode (skip long-running tests)
test-short:
	$(GOTEST) -short -v $(TEST_DIRS)

# Run specific test package
test-package:
	@echo "Usage: make test-package PACKAGE=./tests/unit/config_test.go"

# Lint the code (optional - skips if golangci-lint not available)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not found - skipping lint check"; \
	fi

# Format code
fmt:
	$(GOCMD) fmt ./...

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Install development tools
install-tools:
	$(GOCMD) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run all checks (format, lint, test)
check: fmt lint test

# Build and test
all: deps build check

# Show help
help:
	@echo "Scripts Tool - Available targets:"
	@echo "  build          - Build the binary"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-short     - Run tests in short mode"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  deps           - Download and tidy dependencies"
	@echo "  install-tools  - Install development tools"
	@echo "  check          - Run format, lint, and tests"
	@echo "  all            - Full build and test pipeline"
	@echo "  help           - Show this help"

# Development workflow
dev: clean deps fmt lint test build

# CI/CD pipeline
ci: deps fmt lint test-coverage build

.PHONY: build clean test test-unit test-integration test-coverage test-short lint fmt deps install-tools check all help dev ci
