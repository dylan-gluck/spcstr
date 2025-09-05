# Tech Stack

This is the **DEFINITIVE** technology selection section. These choices will be referenced by all other documentation and AI agents.

## Cloud Infrastructure
- **Provider:** N/A - Local application
- **Key Services:** File system only
- **Deployment Regions:** N/A - Distributed as binary

## Technology Stack Table

| Category | Technology | Version | Purpose | Rationale |
|----------|------------|---------|---------|-----------|
| **Language** | Go | 1.21.0 | Primary development language | Latest stable with improved performance, required by PRD |
| **CLI Framework** | Cobra | 1.8.0 | Command-line interface | Industry standard, explicitly required in PRD |
| **TUI Framework** | Bubbletea | 0.25.0 | Terminal UI | Proven stability, explicitly required in PRD |
| **TUI Components** | Bubbles | 0.17.1 | UI components for Bubbletea | Official component library for Bubbletea |
| **TUI Styling** | Lipgloss | 0.9.1 | Terminal styling | Official styling library for Bubbletea |
| **Markdown** | Glamour | 0.6.0 | Markdown rendering in terminal | Integrates with Bubbletea ecosystem |
| **Shell Scripts** | POSIX sh | N/A | Hook scripts | Maximum compatibility as required |
| **File Watching** | fsnotify | 1.7.0 | File system events | Production-proven, cross-platform |
| **Testing** | Go testing | stdlib | Unit tests | Standard Go testing package |
| **Testing Mock** | testify | 1.8.4 | Test assertions and mocks | Most popular Go testing toolkit |
| **Linting** | golangci-lint | 1.55.0 | Code quality | Comprehensive Go linter |
| **Build Tool** | Make | N/A | Build automation | Standard, mentioned in PRD stories |
| **Module Mgmt** | Go Modules | 1.21.0 | Dependency management | Standard Go dependency system |
| **JSON** | encoding/json | stdlib | Session persistence | Explicitly required, standard library |
| **Config Format** | JSON | N/A | Configuration files | Consistency with session format |
| **Fuzzy Search** | fuzzy | 1.0.0 | Document search in Plan View | Lightweight fuzzy matching |
