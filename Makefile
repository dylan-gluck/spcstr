# Makefile for spcstr Go CLI/TUI application

# Build variables
BINARY_NAME=spcstr
BUILD_DIR=bin
MAIN_PATH=cmd/spcstr

# Go build flags
GOOS_LOCAL := $(shell go env GOOS)
GOARCH_LOCAL := $(shell go env GOARCH)

# Default target
.PHONY: all
all: clean build

# Build for local platform
.PHONY: build
build:
	@echo "Building $(BINARY_NAME) for $(GOOS_LOCAL)/$(GOARCH_LOCAL)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for macOS
.PHONY: build-darwin
build-darwin:
	@echo "Building $(BINARY_NAME) for macOS..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./$(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./$(MAIN_PATH)
	@echo "macOS builds complete"

# Build for Linux
.PHONY: build-linux
build-linux:
	@echo "Building $(BINARY_NAME) for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(MAIN_PATH)
	@GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./$(MAIN_PATH)
	@echo "Linux builds complete"

# Build for all platforms
.PHONY: build-all
build-all: build-darwin build-linux

# Test the project
.PHONY: test
test:
	@echo "Running tests..."
	@go test ./...

# Lint the project (requires golangci-lint to be installed)
.PHONY: lint
lint:
	@echo "Running linters..."
	@golangci-lint run

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Show build info
.PHONY: info
info:
	@echo "Build Information:"
	@echo "  Binary Name: $(BINARY_NAME)"
	@echo "  Build Directory: $(BUILD_DIR)"
	@echo "  Main Path: $(MAIN_PATH)"
	@echo "  Local OS/Arch: $(GOOS_LOCAL)/$(GOARCH_LOCAL)"
	@echo "  Go Version: $(shell go version)"

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Clean and build for local platform"
	@echo "  build        - Build for local platform"
	@echo "  build-darwin - Build for macOS (amd64 and arm64)"
	@echo "  build-linux  - Build for Linux (amd64 and arm64)"
	@echo "  build-all    - Build for all supported platforms"
	@echo "  test         - Run all tests"
	@echo "  lint         - Run linters (requires golangci-lint)"
	@echo "  clean        - Remove build artifacts"
	@echo "  deps         - Install and tidy dependencies"
	@echo "  info         - Show build information"
	@echo "  help         - Show this help message"