# Tech Stack

## Technology Stack Table

| Category | Technology | Version | Purpose | Rationale |
|----------|------------|---------|---------|-----------|
| Primary Language | Go | 1.21+ | CLI/TUI application development | Memory safety, cross-platform binaries, excellent CLI ecosystem |
| TUI Framework | Bubbletea | v0.25+ | Terminal user interface | Industry standard for Go TUI apps with excellent event handling |
| UI Styling | Lipgloss | v0.9+ | Terminal styling and layout | Seamless integration with Bubbletea for consistent visual design |
| Markdown Rendering | Glamour | v0.6+ | Document display with syntax highlighting | Rich markdown rendering in terminal environments |
| CLI Framework | Cobra | v1.8+ | Command structure and hook subcommands | Standard Go CLI framework with excellent subcommand support |
| File Watching | fsnotify | v1.7+ | Real-time file system monitoring | Cross-platform file watching for live TUI updates |
| JSON Processing | Standard Library | Go 1.21+ | State serialization and parsing | Built-in JSON support eliminates external dependencies |
| Atomic Operations | Standard Library | Go 1.21+ | Safe concurrent file operations | Native filesystem atomicity through temp file + rename |
| Testing Framework | Go Testing | Go 1.21+ | Unit and integration testing | Built-in testing with table-driven test patterns |
| Build System | Standard Go Build | Go 1.21+ | Binary compilation | Native Go build tools with cross-compilation support |
| Release Automation | Goreleaser | v1.21+ | Multi-platform binary distribution | Automated releases to package managers and GitHub |
| Version Control | Git | 2.40+ | Source code management | Standard version control with GitHub integration |
