# Spec⭐️ (spcstr) Makefile
.PHONY: build test install lint clean help

# Variables
BINARY_NAME := spcstr
BINARY_PATH := bin/$(BINARY_NAME)
MAIN_PATH := cmd/spcstr/main.go
INSTALL_PATH := $(GOPATH)/bin

# Default target
help:
	@echo "Available targets:"
	@echo "  build    - Build the binary to bin/spcstr"
	@echo "  test     - Run all tests with coverage"
	@echo "  install  - Install binary to $$GOPATH/bin"
	@echo "  lint     - Run golangci-lint"
	@echo "  clean    - Remove build artifacts"
	@echo "  help     - Show this help message"

# Build the binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p bin
	@go build -o $(BINARY_PATH) $(MAIN_PATH)
	@echo "Binary built at $(BINARY_PATH)"

# Run tests with coverage
test:
	@echo "Running tests with coverage..."
	@go test -v -cover -coverprofile=coverage.out ./...
	@echo "Coverage report generated at coverage.out"
	@echo "Run 'go tool cover -html=coverage.out' to view detailed coverage"

# Install to $GOPATH/bin
install: build
	@echo "Installing to $(INSTALL_PATH)..."
	@cp $(BINARY_PATH) $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed to $(INSTALL_PATH)"

# Run linter
lint:
	@echo "Running golangci-lint..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Install with:"; \
		echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.0"; \
		exit 1; \
	fi

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@rm -f coverage.out
	@echo "Clean complete"

# Run the application
run: build
	@$(BINARY_PATH)