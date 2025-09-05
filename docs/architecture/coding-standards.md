# Coding Standards

**⚠️ IMPORTANT:** These standards are MANDATORY for AI agents generating code for this project. We'll keep this minimal and focused on critical project-specific rules only.

## Core Standards
- **Languages & Runtimes:** Go 1.21.0 (no other languages except POSIX shell for hooks)
- **Style & Linting:** golangci-lint with .golangci.yml configuration (gofmt, goimports, govet enabled)
- **Test Organization:** Tests in same package with _test.go suffix, integration tests in test/integration/

## Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Packages | lowercase, no underscores | `session`, `indexer` |
| Interfaces | PascalCase with -er suffix | `SessionManager`, `Renderer` |
| Structs | PascalCase | `SessionData`, `Config` |
| Functions/Methods | PascalCase (exported), camelCase (private) | `LoadSession()`, `parseJSON()` |
| Constants | PascalCase | `DefaultTimeout`, `MaxRetries` |
| Files | snake_case.go | `session_manager.go` |

## Critical Rules

- **Never use fmt.Print/Println in production code - use log/slog:** All output must go through structured logging
- **All file operations must use atomic writes (temp + rename):** Prevents data corruption on crash
- **Hook commands must complete within 10ms:** Never block Claude Code operations
- **Hook exit codes: 0 for success, 2 for blocking:** Enables safety checks for dangerous operations
- **All errors must be wrapped with context using %w:** Enables proper error tracing
- **Session IDs must use format sess_{uuid}:** Consistent identification across system
- **Never panic in library code - return errors:** Only main() can panic on fatal errors
- **All public types must have godoc comments:** Documentation is mandatory for exported items
- **Use context.Context for cancellation in long operations:** Enables graceful shutdown

## Go-Specific Guidelines

- **Interfaces belong in consumer package, not provider:** Define interfaces where they're used
- **Prefer small interfaces (1-3 methods):** Enables better composition and testing
- **Use channels for coordination, mutexes for state:** Don't communicate by sharing memory
- **Embed types rather than inherit:** Go doesn't have inheritance, use composition
