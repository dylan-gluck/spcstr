# Coding Standards

## Critical Go Rules

- **Single Binary Rule:** All functionality must be accessible through the main spcstr binary via subcommands
- **Atomic Operations:** Always use temp file + rename pattern for state modifications, never direct writes
- **Error Propagation:** Hook handlers must return appropriate exit codes (0=success, 2=block operation)
- **JSON Schema Compliance:** All state operations must conform to exact schema from hooks-state-management.md
- **File Path Safety:** Always use filepath.Join() and validate paths to prevent directory traversal
- **Context Timeouts:** Use context.WithTimeout for all file operations to prevent hanging
- **Resource Cleanup:** Always defer file.Close() and handle cleanup in error paths
- **Hook Isolation:** Hook command execution must not modify global state or affect TUI operation

## Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Types | PascalCase | `SessionState`, `HookHandler` |
| Functions | PascalCase (exported), camelCase (internal) | `NewManager()`, `loadState()` |
| Constants | UPPER_SNAKE_CASE | `DEFAULT_TIMEOUT`, `STATE_FILE_NAME` |
| File Names | snake_case | `session_start.go`, `state_manager.go` |
| Package Names | lowercase | `hooks`, `state`, `tui` |
| Hook Commands | snake_case | `session_start`, `pre_tool_use` |
| JSON Fields | snake_case | `session_id`, `created_at` |
