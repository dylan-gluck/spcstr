# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.0.1] - 2025-09-06

### Added
- Initial release of spcstr observability framework
- TUI with two primary views:
  - **Observe View**: Real-time session monitoring with dashboard displaying agents, files, tools, and tasks
  - **Plan View**: Browse PRDs, architecture docs, and workflows with rich markdown rendering
- Hook system for Claude Code integration:
  - Session lifecycle tracking (start/end)
  - User prompt capture
  - Tool usage monitoring (pre/post)
  - Agent activity tracking
  - File operation tracking (new/edited/read)
  - Task/TODO state management
- State persistence in `.spcstr/sessions/` directory
- Automatic session discovery and loading
- Real-time updates via file watching (fsnotify)
- Keyboard navigation and shortcuts
- Session activity indicators (active/inactive)
- Chronological activity feed with prompts and notifications

### Technical Features
- Built with Go and Bubbletea/Lipgloss TUI framework
- JSON-based state management
- Hook-based architecture for minimal Claude Code coupling
- Zero network calls - all data stored locally
- Responsive layout with pane focus management

### Known Issues
- Unit test coverage incomplete for observe view
- File watcher error recovery could be improved
- Session list caching not implemented

### Documentation
- Comprehensive architecture documentation
- Story-driven development records
- QA gate process established
- Testing strategy defined

[0.0.1]: https://github.com/dylan/spcstr/releases/tag/v0.0.1