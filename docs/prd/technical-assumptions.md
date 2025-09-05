# Technical Assumptions

## Repository Structure: Monorepo
Single Go module containing all components - TUI application and hook executables as separate binaries. Clean separation with `cmd/` for binaries and `internal/` for shared code.

## Service Architecture
**Single binary with embedded hook functionality.** Main `spcstr` binary provides CLI, TUI, and hook execution via Cobra subcommands. Hook logic implemented as separate internal packages but exposed through `spcstr hook <hook_name>` subcommands. No separate binaries to manage.

**Cobra Command Structure:**
```
spcstr                           # Root command - launches TUI
├── init                         # Initialize project with hook settings
├── version                      # Display version information  
├── config                       # Manage configuration
├── tui                         # Explicitly launch TUI
└── hook                        # Execute hook functionality (called by Claude Code)
    ├── session_start
    ├── user_prompt_submit
    ├── pre_tool_use
    ├── post_tool_use
    ├── notification
    ├── pre_compact
    ├── session_end
    ├── stop
    └── subagent_stop
```

## Testing Requirements
**Minimal testing for MVP.** Focus on manual testing of critical paths: init command, hook execution, state management, TUI navigation. Unit tests only for state management atomic write operations. Full testing pyramid deferred to post-MVP.

## Additional Technical Assumptions and Requests

- **Go 1.21+** minimum version for modern stdlib features
- **Embedded hook functionality** in main binary via Cobra subcommands
- **Build process** compiles single binary with all hook commands
- **Bubbletea** for TUI framework with event-driven architecture
- **Lipgloss** for consistent cross-platform terminal styling
- **Glamour** for markdown rendering with syntax highlighting
- **Cobra** for CLI structure and subcommands
- **Shared internal packages** for state management used by both TUI and hooks
- **JSON** for all data persistence (state.json, logs, settings)
- **File watching** via fsnotify for real-time TUI updates
- **No database** - filesystem only with `.spcstr/` directory
- **No network calls** - completely offline, privacy-preserving
- **Platform-specific builds** via goreleaser creating distribution packages
