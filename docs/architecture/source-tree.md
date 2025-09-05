# Source Tree

```
spcstr/
├── cmd/
│   └── spcstr/
│       ├── main.go                 # Application entry point
│       └── hook.go                 # Hook subcommand entry
├── internal/
│   ├── cli/
│   │   ├── root.go                 # Root command setup
│   │   ├── init.go                 # Init command implementation
│   │   ├── run.go                  # Run command implementation
│   │   ├── config.go               # Config command implementation
│   │   └── version.go              # Version command
│   ├── config/
│   │   ├── config.go               # Configuration structures
│   │   ├── loader.go               # Config loading logic
│   │   ├── validator.go            # Config validation
│   │   └── xdg.go                  # XDG directory handling
│   ├── hooks/
│   │   ├── handler.go              # Main hook event dispatcher
│   │   ├── state.go                # Session state management
│   │   ├── types.go                # Hook data structures
│   │   ├── persistence.go          # Atomic file I/O operations
│   │   ├── claude.go               # Claude settings updater
│   │   └── handlers/               # Individual hook handlers
│   │       ├── pre_tool_use.go    # Pre-tool safety checks
│   │       ├── post_tool_use.go   # Tool tracking & file ops
│   │       ├── session_start.go   # Session initialization
│   │       ├── session_end.go     # Session finalization
│   │       ├── user_prompt.go     # Prompt filtering
│   │       ├── notification.go    # Notification logging
│   │       ├── stop.go            # Stop handling
│   │       ├── subagent_stop.go   # Agent tracking
│   │       └── pre_compact.go     # Compaction prep
│   ├── session/
│   │   ├── session.go              # Session data structures
│   │   ├── manager.go              # Session lifecycle management
│   │   ├── persistence.go          # JSON read/write operations
│   │   └── watcher.go              # File system monitoring
│   ├── events/
│   │   ├── bus.go                  # Event bus implementation
│   │   ├── types.go                # Event type definitions
│   │   └── subscriber.go           # Subscription management
│   ├── tui/
│   │   ├── app.go                  # Main TUI application
│   │   ├── keys.go                 # Global keybinding definitions
│   │   ├── styles.go               # Lipgloss style definitions
│   │   └── components/
│   │       ├── help.go             # Help overlay component
│   │       ├── status.go           # Status bar component
│   │       └── input.go            # Search input component
│   ├── views/
│   │   ├── plan/
│   │   │   ├── view.go             # Plan view implementation
│   │   │   ├── document_list.go    # Document list component
│   │   │   ├── preview.go          # Markdown preview pane
│   │   │   └── modes.go            # Spec/Workflow/Config modes
│   │   └── observe/
│   │       ├── view.go             # Observe view implementation
│   │       ├── session_list.go     # Session list component
│   │       └── dashboard.go        # Dashboard layout
│   ├── dashboard/
│   │   ├── renderer.go             # Dashboard rendering engine
│   │   ├── sections.go             # Section definitions
│   │   ├── agents.go               # Agent status section
│   │   ├── tasks.go                # Task progress section
│   │   ├── files.go                # File operations section
│   │   ├── tools.go                # Tool usage metrics
│   │   └── errors.go               # Error log section
│   ├── indexer/
│   │   ├── indexer.go              # Document indexing engine
│   │   ├── scanner.go              # Directory scanner
│   │   ├── classifier.go           # Document type detection
│   │   ├── metadata.go             # Metadata extraction
│   │   └── search.go               # Fuzzy search implementation
│   └── utils/
│       ├── logger.go               # Logging utilities
│       ├── json.go                 # JSON helpers
│       ├── filepath.go             # Path manipulation
│       └── terminal.go             # Terminal utilities
├── pkg/
│   ├── models/
│   │   ├── session.go              # Public session types
│   │   └── document.go             # Document model
│   └── hooks/
│       ├── session_state.go        # SessionState structure
│       ├── file_operations.go      # FileOperations model
│       ├── agent_history.go        # AgentHistoryEntry model
│       └── error_entry.go          # ErrorEntry model
├── scripts/
│   ├── install.sh                  # Installation script
│   └── release.sh                  # Release build script
├── test/
│   ├── integration/
│   │   ├── cli_test.go             # CLI integration tests
│   │   ├── session_test.go         # Session management tests
│   │   └── tui_test.go             # TUI interaction tests
│   ├── fixtures/
│   │   ├── sessions/               # Test session data
│   │   └── docs/                   # Test documents
│   └── mocks/
│       ├── filesystem.go           # File system mocks
│       └── events.go               # Event bus mocks
├── docs/
│   ├── prd.md                      # Product Requirements Document
│   ├── architecture.md             # This architecture document
│   ├── frontend-spec.md            # Frontend specification
│   └── api/                        # Generated API documentation
├── .spcstr/                         # Default project config location
│   ├── config.json                 # Project configuration
│   ├── sessions/                   # Session data directories
│   └── cache/                      # Application cache
├── go.mod                           # Go module definition
├── go.sum                           # Dependency checksums
├── Makefile                         # Build automation
├── .golangci.yml                   # Linter configuration
├── .gitignore                      # Git ignore rules
├── LICENSE                          # MIT License
└── README.md                        # Project documentation
```
