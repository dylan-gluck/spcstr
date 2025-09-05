# REST API Spec

Based on the PRD and architecture, **Spec⭐️ does not include a REST API**. This is a terminal-based CLI/TUI application that:

- Operates entirely through command-line interface (Cobra CLI)
- Provides user interaction through terminal UI (Bubbletea)
- Integrates with Claude Code via file system hooks, not HTTP
- Uses file-based persistence, not network APIs

The application follows a different architectural pattern:
- **CLI Commands** instead of REST endpoints
- **File system events** instead of HTTP requests
- **TUI interactions** instead of web responses
- **Shell scripts** for external integration
