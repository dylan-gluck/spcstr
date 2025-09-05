# API Specification

*No external API - This is a local CLI/TUI application. All communication occurs through:*

1. **Hook Command Interface:** Claude Code → `spcstr hook <name>` via stdin/stdout
2. **File System Interface:** JSON state files and log files
3. **Internal Go Interfaces:** Between TUI, state engine, and hook handlers

## Internal Command Interface

```go
// Hook command signature
type HookHandler func(input []byte) error

// State management interface
type StateManager interface {
    InitializeState(sessionID string) error
    LoadState(sessionID string) (*SessionState, error)
    UpdateState(sessionID string, updates StateUpdate) error
}

// TUI component interface
type TUIComponent interface {
    Update(msg tea.Msg) (tea.Model, tea.Cmd)
    View() string
}
```
