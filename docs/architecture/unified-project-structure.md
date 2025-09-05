# Unified Project Structure

```
spcstr/
├── .github/                    # CI/CD workflows
│   └── workflows/
│       ├── ci.yaml            # Build and test pipeline
│       └── release.yaml       # Goreleaser workflow
├── cmd/                       # Application binaries
│   └── spcstr/               # Main binary entry point
│       └── main.go           # CLI bootstrap and Cobra setup
├── internal/                  # Internal packages (not importable)
│   ├── hooks/                # Hook command implementations
│   │   ├── handlers/         # Individual hook handlers
│   │   │   ├── session_start.go
│   │   │   ├── user_prompt_submit.go
│   │   │   ├── pre_tool_use.go
│   │   │   ├── post_tool_use.go
│   │   │   ├── notification.go
│   │   │   ├── pre_compact.go
│   │   │   ├── session_end.go
│   │   │   ├── stop.go
│   │   │   └── subagent_stop.go
│   │   ├── registry.go       # Hook registration system
│   │   └── executor.go       # Hook execution coordinator
│   ├── state/                # State management package
│   │   ├── manager.go        # State CRUD operations
│   │   ├── atomic.go         # Atomic file operations
│   │   ├── watcher.go        # File system monitoring
│   │   └── types.go          # State data structures
│   ├── tui/                  # TUI implementation
│   │   ├── app/              # Main TUI application
│   │   │   └── app.go        # Bubbletea app controller
│   │   ├── components/       # Reusable UI components
│   │   │   ├── header/       # Header bar
│   │   │   ├── footer/       # Status/keybind footer
│   │   │   ├── list/         # Generic list component
│   │   │   └── dashboard/    # Session dashboard
│   │   ├── views/            # Main view implementations
│   │   │   ├── plan/         # Plan view (document browser)
│   │   │   │   ├── plan.go   # Plan view controller
│   │   │   │   └── browser.go # Document browser logic
│   │   │   └── observe/      # Observe view (session monitor)
│   │   │       ├── observe.go # Observe view controller
│   │   │       └── dashboard.go # Dashboard rendering
│   │   ├── styles/           # Lipgloss styling
│   │   │   └── theme.go      # Color schemes and layouts
│   │   └── messages/         # Bubbletea messages
│   │       └── events.go     # Custom message types
│   ├── docs/                 # Document management
│   │   ├── scanner.go        # Document discovery
│   │   ├── indexer.go        # Document indexing
│   │   └── renderer.go       # Glamour markdown rendering
│   ├── config/               # Configuration management
│   │   ├── settings.go       # Application settings
│   │   ├── init.go           # Project initialization
│   │   └── paths.go          # Path management utilities
│   └── utils/                # Shared utilities
│       ├── filesystem.go     # File operation helpers
│       ├── json.go           # JSON processing utilities
│       └── terminal.go       # Terminal detection utilities
├── pkg/                      # Public API packages (if needed)
├── scripts/                  # Build and development scripts
│   ├── build.sh             # Local build script
│   ├── test.sh              # Testing script
│   └── install-hooks.sh     # Development hook setup
├── docs/                     # Project documentation
│   ├── prd.md               # Product Requirements Document
│   ├── architecture.md      # This document
│   ├── plan/                # Planning documents
│   └── vendor/              # External documentation
├── examples/                 # Usage examples
│   └── .spcstr/             # Example directory structure
├── .goreleaser.yaml         # Release configuration
├── go.mod                   # Go module definition
├── go.sum                   # Dependency checksums
├── Makefile                 # Build automation
└── README.md                # Project overview and usage
```
