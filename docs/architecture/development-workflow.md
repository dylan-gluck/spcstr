# Development Workflow

## Local Development Setup

### Prerequisites
```bash
# Install Go 1.21+
go version  # Verify Go installation

# Install development tools
go install github.com/goreleaser/goreleaser@latest
```

### Initial Setup
```bash
# Clone and setup
git clone https://github.com/username/spcstr.git
cd spcstr

# Install dependencies
go mod tidy

# Build binary
make build

# Run tests
make test

# Install locally for development
make install
```

### Development Commands
```bash
# Build all binaries
make build

# Run TUI directly
go run cmd/spcstr/main.go

# Test hook execution
echo '{"session_id": "test_session", "source": "startup"}' | go run cmd/spcstr/main.go hook session_start

# Run tests with coverage
make test-coverage

# Run linting
make lint
```

## Environment Configuration

### Required Environment Variables
```bash
# Development (.env.local - optional)
SPCSTR_DEBUG=true                    # Enable debug logging
SPCSTR_CONFIG_PATH=./dev-config     # Custom config path for development

# Production (system environment)
# No environment variables required - all configuration via .spcstr/ directory

# CI/CD
GITHUB_TOKEN=ghp_xxx                 # For release automation
GO_VERSION=1.21                      # Go version for builds
```
