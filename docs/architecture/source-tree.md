# Source Tree Structure

## spcstr Go CLI/TUI Implementation

The spcstr project follows this Go monorepo structure:

```
spcstr/                          # Project root
├── .github/                     # CI/CD workflows
│   └── workflows/
│       ├── ci.yaml             # Build and test pipeline
│       └── release.yaml        # Goreleaser workflow
├── cmd/                        # Application binaries
│   └── spcstr/                # Main binary entry point
│       └── main.go            # CLI bootstrap and Cobra setup
├── internal/                   # Internal packages (not importable)
│   ├── hooks/                 # Hook command implementations
│   │   ├── handlers/          # Individual hook handlers
│   │   │   ├── session_start.go
│   │   │   ├── user_prompt_submit.go
│   │   │   ├── pre_tool_use.go
│   │   │   ├── post_tool_use.go
│   │   │   ├── notification.go
│   │   │   ├── pre_compact.go
│   │   │   ├── session_end.go
│   │   │   ├── stop.go
│   │   │   └── subagent_stop.go
│   │   ├── registry.go        # Hook registration system
│   │   └── executor.go        # Hook execution coordinator
│   ├── state/                 # State management package
│   │   ├── manager.go         # State CRUD operations
│   │   ├── atomic.go          # Atomic file operations
│   │   ├── watcher.go         # File system monitoring
│   │   └── types.go           # State data structures
│   ├── tui/                   # TUI implementation
│   │   ├── app/               # Main TUI application
│   │   │   └── app.go         # Bubbletea app controller
│   │   ├── components/        # Reusable UI components
│   │   │   ├── header/        # Header bar
│   │   │   ├── footer/        # Status/keybind footer
│   │   │   ├── list/          # Generic list component
│   │   │   └── dashboard/     # Session dashboard
│   │   ├── views/             # Main view implementations
│   │   │   ├── plan/          # Plan view (document browser)
│   │   │   │   ├── plan.go    # Plan view controller
│   │   │   │   └── browser.go # Document browser logic
│   │   │   └── observe/       # Observe view (session monitor)
│   │   │       ├── observe.go # Observe view controller
│   │   │       └── dashboard.go # Dashboard rendering
│   │   ├── styles/            # Lipgloss styling
│   │   │   └── theme.go       # Color schemes and layouts
│   │   └── messages/          # Bubbletea messages
│   │       └── events.go      # Custom message types
│   ├── docs/                  # Document management
│   │   ├── scanner.go         # Document discovery
│   │   ├── indexer.go         # Document indexing
│   │   └── renderer.go        # Glamour markdown rendering
│   ├── config/                # Configuration management
│   │   ├── settings.go        # Application settings
│   │   ├── init.go            # Project initialization
│   │   └── paths.go           # Path management utilities
│   └── utils/                 # Shared utilities
│       ├── filesystem.go      # File operation helpers
│       ├── json.go            # JSON processing utilities
│       └── terminal.go        # Terminal detection utilities
├── pkg/                       # Public API packages (if needed)
├── tests/                     # Test files
│   ├── integration/           # Integration tests
│   ├── testdata/              # Test fixtures
│   └── manual/                # Manual testing procedures
├── scripts/                   # Build and development scripts
│   ├── build.sh              # Local build script
│   ├── test.sh               # Testing script
│   └── install-hooks.sh      # Development hook setup
├── examples/                  # Usage examples
│   └── .spcstr/              # Example directory structure
├── docs/                      # Project documentation (preserved)
├── .goreleaser.yaml          # Release configuration
├── go.mod                    # Go module definition
├── go.sum                    # Dependency checksums
├── Makefile                  # Build automation
├── README.md                 # Project overview and usage
└── LICENSE                   # MIT License
```

## Runtime Directory Structure

When spcstr is initialized in a project, it creates this directory structure:

```
{project_root}/
├── .spcstr/                    # spcstr configuration and data
│   ├── sessions/              # Session state directory
│   │   └── {session-id}/
│   │       └── state.json     # SessionState JSON
│   ├── logs/                  # Hook execution logs
│   │   ├── session_start.json      # Array of session start events
│   │   ├── user_prompt_submit.json # Array of prompt events
│   │   ├── pre_tool_use.json       # Array of tool invocation events
│   │   ├── post_tool_use.json      # Array of tool completion events
│   │   ├── notification.json       # Array of notification events
│   │   ├── pre_compact.json        # Array of compaction events
│   │   ├── session_end.json        # Array of session end events
│   │   ├── stop.json              # Array of stop events
│   │   └── subagent_stop.json     # Array of subagent stop events
│   └── config/                # Local spcstr configuration
│       └── settings.json      # Project-specific settings
└── docs/                      # Project documentation (for Plan view)
    ├── prd.md                # Product requirements
    ├── architecture.md       # Architecture document
    ├── epics/                # Epic documents
    └── stories/              # Story documents
```

## Package Dependencies

### Core Go Dependencies
- **Go Version:** 1.21+
- **CLI Framework:** github.com/spf13/cobra v1.8+
- **TUI Framework:** github.com/charmbracelet/bubbletea v0.25+
- **UI Styling:** github.com/charmbracelet/lipgloss v0.9+
- **Markdown Rendering:** github.com/charmbracelet/glamour v0.6+
- **File Watching:** github.com/fsnotify/fsnotify v1.7+
- **Standard Library:** json, os, filepath, time, context

### Development Dependencies
- **Build Automation:** github.com/goreleaser/goreleaser v1.21+
- **Testing:** Go standard testing package
- **Linting:** Various Go linting tools (configured in CI)

## Key Design Principles

1. **Single Binary Architecture:** All functionality embedded in one executable
2. **Privacy-First:** No network calls, all data remains local
3. **Atomic Operations:** State changes use temp file + rename pattern
4. **Real-time Updates:** File watching for immediate UI feedback
5. **Clean Architecture:** Clear separation between CLI, TUI, state, and hook layers

## Key Architectural Components

### Entry Points
- **`cmd/spcstr/main.go`** - Main binary entry point with Cobra CLI setup
- **`internal/tui/app/app.go`** - Bubbletea TUI application controller

### Core Systems
- **`internal/hooks/`** - Hook command implementations for Claude Code integration
- **`internal/state/`** - Atomic state management with JSON persistence
- **`internal/tui/`** - Terminal user interface with Bubbletea framework

### Key Features
- **Hook System:** Real-time Claude Code session tracking via executable hooks
- **TUI Interface:** Interactive terminal interface with Plan and Observe views
- **State Management:** Atomic JSON file operations for session persistence
- **Document Browser:** Markdown document discovery and rendering for Plan view

## Package Responsibilities

### `internal/hooks/`
Handles Claude Code integration via hook commands that receive JSON input and update session state.

### `internal/state/`
Manages atomic file operations for session state with filesystem-level atomicity using temp file + rename pattern.

### `internal/tui/`
Implements terminal user interface using Bubbletea with real-time file watching for immediate updates.

### `internal/docs/`
Provides document discovery and Glamour-based markdown rendering for the Plan view document browser.

## Runtime Integration

The spcstr binary integrates with Claude Code by being configured as hooks in the user's Claude Code settings, creating a seamless observability experience during development sessions.