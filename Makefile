.PHONY: all build clean test lint run help build-all

# Variables
BINARY_NAME=superscan
BIN_DIR=bin
VERSION ?= dev
GO=go
GOLINT=golangci-lint
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
LDFLAGS=-X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)

# Default target
all: clean test build

# Build the application
build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BIN_DIR)
	@$(GO) build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY_NAME)

# Build for all platforms
build-all:
	@echo "Building $(BINARY_NAME) for all platforms..."
	@mkdir -p $(BIN_DIR)
	@for platform in darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64; do \
		os=$${platform%/*}; \
		arch=$${platform#*/}; \
		output="$(BIN_DIR)/$(BINARY_NAME)-$$os-$$arch"; \
		if [ "$$os" = "windows" ]; then \
			output="$$output.exe"; \
		fi; \
		echo "Building for $$os/$$arch..."; \
		GOOS=$$os GOARCH=$$arch $(GO) build -ldflags "$(LDFLAGS)" -o $$output; \
	done

# Run tests
test:
	@echo "Running tests..."
	@$(GO) test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(GO) test -v -coverprofile=coverage.out ./...
	@$(GO) tool cover -html=coverage.out

# Run linter
lint:
	@echo "Running linter..."
	@$(GOLINT) run

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@rm -f coverage.out

# Run the application
run:
	@echo "Running $(BINARY_NAME)..."
	@$(GO) run main.go

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@$(GO) mod download
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Show help
help:
	@echo "Available targets:"
	@echo "  all            - Clean, test, and build the application"
	@echo "  build          - Build the application for current platform"
	@echo "  build-all      - Build the application for all supported platforms"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  lint           - Run linter"
	@echo "  clean          - Remove build artifacts"
	@echo "  run            - Run the application"
	@echo "  deps           - Install dependencies"
	@echo "  help           - Show this help message"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION        - Set version (default: dev)"
	@echo "  GOOS           - Set target OS (default: current OS)"
	@echo "  GOARCH         - Set target architecture (default: current arch)" 