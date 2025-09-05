# Components

## CLI Command Handler

**Responsibility:** Process command-line arguments and route to appropriate subsystems

**Key Interfaces:**
- Cobra root command with subcommands (init, version, config, tui, hook)
- Hook subcommand routing to individual hook handlers
- Configuration management for `.spcstr/` setup

**Dependencies:** Cobra framework, State Manager, Hook System

**Technology Stack:** Cobra CLI framework with embedded hook logic

## State Engine

**Responsibility:** Atomic state management for session tracking and persistence

**Key Interfaces:**
- SessionState CRUD operations with atomic write guarantees
- JSON serialization/deserialization with validation
- File system operations with error handling

**Dependencies:** Standard library (os, json), atomic file operations

**Technology Stack:** Go standard library with custom atomic write implementation

## TUI Application

**Responsibility:** Interactive terminal interface with real-time updates

**Key Interfaces:**
- Bubbletea model/update/view pattern implementation
- Component switching between Plan and Observe views
- Keyboard event handling and navigation

**Dependencies:** Bubbletea, Lipgloss, Glamour, File Watcher, State Engine

**Technology Stack:** Bubbletea framework with Lipgloss styling and Glamour rendering

## Hook System

**Responsibility:** Claude Code integration via executable hook commands

**Key Interfaces:**
- Standard input JSON parsing for hook parameters
- Exit code management (0=success, 2=block operation)
- State updates triggered by Claude Code events

**Dependencies:** State Engine, JSON parsing, Standard I/O

**Technology Stack:** Go standard library with JSON processing

## File Watcher

**Responsibility:** Real-time file system monitoring for TUI updates

**Key Interfaces:**
- fsnotify integration for `.spcstr/` directory monitoring
- Event filtering for relevant state file changes
- Bubbletea command generation for UI updates

**Dependencies:** fsnotify library, TUI Application

**Technology Stack:** fsnotify cross-platform file watching

## Document Engine

**Responsibility:** Markdown document discovery and rendering for Plan view

**Key Interfaces:**
- File system scanning for markdown documents
- Glamour markdown rendering with syntax highlighting
- Document indexing and caching

**Dependencies:** Glamour renderer, file system operations

**Technology Stack:** Glamour markdown rendering with file system integration
