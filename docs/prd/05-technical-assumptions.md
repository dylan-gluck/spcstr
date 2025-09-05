# Technical Assumptions

## Repository Structure: Monorepo

Single repository containing all Go packages, hook scripts, and configuration templates.

## Service Architecture

Monolithic Go binary with modular internal package structure. Event-driven architecture internally with channels for real-time updates. File-based storage using OS file system for persistence.

## Testing Requirements

Comprehensive unit tests for core logic, integration tests for file operations and hook processing, and manual testing utilities for TUI interactions.

## Additional Technical Assumptions and Requests

- Use Go 1.21+ for improved performance and standard library features
- Leverage Bubbletea for TUI framework with proven stability
- Implement Cobra for CLI command structure and help generation
- Use standard library JSON package for session state serialization
- Shell scripts must be POSIX-compliant for maximum compatibility
- Avoid CGO dependencies to ensure easy cross-compilation
- Configuration follows XDG Base Directory specification
- Support for both project-local and global configuration