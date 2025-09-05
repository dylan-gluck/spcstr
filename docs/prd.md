# Spec⭐️ (spcstr) Product Requirements Document (PRD)

## Goals and Background Context

### Goals

- Enable real-time observability of Claude Code sessions through an intuitive TUI dashboard
- Automatically persist session state for analysis, replay, and recovery
- Provide seamless integration with BMad Method planning documents
- Support both planning and execution phases of development workflow
- Deliver sub-100ms UI responsiveness for all user interactions
- Create a single-binary tool installable via standard package managers
- Establish foundation for future multi-agent orchestration capabilities

### Background Context

Spec⭐️ addresses the critical gap in Claude Code development workflows where developers lack visibility into active AI-assisted coding sessions. As teams increasingly adopt Claude Code with the BMad Method, they need a purpose-built tool that provides real-time session monitoring, persistent state tracking, and integrated planning capabilities. The current landscape offers no terminal-native solutions that combine these capabilities, forcing developers to manually track session state and lose valuable context between sessions.

The solution leverages Claude Code's hook system to automatically capture session data without manual intervention, presenting it through a responsive TUI that developers can keep open alongside their development work. By maintaining session state in structured JSON format, Spec⭐️ enables session recovery, performance analysis, and workflow optimization that are impossible with current tooling.

### Change Log

| Date | Version | Description | Author |
|------|---------|-------------|--------|
| 2025-09-05 | 1.0 | Initial PRD creation | BMad Master |
| 2025-09-05 | 1.1 | Updated for Go-based hook system migration | Mary (Analyst) |

## Requirements

### Functional

1. FR1: The init command shall create all necessary directories, hook scripts, and configuration files with a single command
2. FR2: The system shall automatically track all Claude Code sessions when hooks are triggered
3. FR3: Session data shall be persisted to JSON files within 1 second of any state change
4. FR4: The TUI shall support keyboard-driven navigation with discoverable keybindings
5. FR5: Plan View shall recursively index and display all markdown planning documents
6. FR6: Plan View shall provide real-time markdown preview with syntax highlighting
7. FR7: Observe View shall display active and completed sessions in a selectable list
8. FR8: Session dashboard shall show agents, tasks, files, tools, and errors in organized sections
9. FR9: The system shall maintain session state across application restarts
10. FR10: Global keybindings shall allow view switching from any context
11. FR11: The TUI shall support terminal resize events without data loss
12. FR12: Session selection shall load and display full session details within 100ms
13. FR13: The system shall track file operations (new, edited, read) with full paths
14. FR14: Tool usage counters shall increment in real-time as hooks fire
15. FR15: Error tracking shall capture and display session errors with timestamps

### Non Functional

1. NFR1: UI response time shall not exceed 100ms for any user interaction
2. NFR2: Memory usage shall not exceed 10MB per tracked session
3. NFR3: CPU usage shall remain below 1% during idle periods
4. NFR4: The application shall compile to a single binary with no runtime dependencies
5. NFR5: Installation shall be achievable via standard package managers (brew, apt, yum)
6. NFR6: The system shall support 256-color terminal emulators
7. NFR7: Session data files shall use standard JSON format for interoperability
8. NFR8: The TUI shall maintain 60fps refresh rate during animations
9. NFR9: File system operations shall use platform-native path separators
10. NFR10: The application shall gracefully handle missing or corrupted session files

## User Interface Design Goals

### Overall UX Vision

Create a terminal-native experience that feels as responsive and intuitive as modern GUI applications while respecting terminal conventions and keyboard-driven workflows. The interface should provide information density without overwhelming users, using color and typography to create clear visual hierarchy.

### Key Interaction Paradigms

- Vim-style keybindings for navigation (hjkl movement, modes)
- Tab-based focus switching between panels
- Modal overlays for detailed views and settings
- Incremental search with fuzzy matching for document discovery
- Real-time updates without disrupting user focus
- Context-sensitive help available at any point

### Core Screens and Views

- **Plan View**: Split-panel layout with document tree and preview
- **Observe View**: Split-panel layout with session list and dashboard
- **Spec Mode**: Indexed PRD/Architecture documents with navigation
- **Workflow Mode**: Epic and story management interface
- **Config Mode**: Settings viewer for claude/bmad/spcstr configuration
- **Session Dashboard**: Multi-section layout with agents, tasks, files, metrics
- **Help Overlay**: Context-sensitive keybinding reference

### Accessibility: Terminal Native

Terminal accessibility through standard screen reader support, high contrast color themes, and keyboard-only navigation.

### Branding

Minimal, professional aesthetic with subtle use of color for status indication. Star emoji (⭐️) as visual identifier. Monospace typography throughout for alignment and readability.

### Target Device and Platforms: Terminal/Console

Cross-platform terminal application supporting Linux, macOS, and Windows (via WSL). Requires 80x24 minimum terminal size, optimized for standard developer terminal configurations.

## Technical Assumptions

### Repository Structure: Monorepo

Single repository containing all Go packages, hook scripts, and configuration templates.

### Service Architecture

Monolithic Go binary with modular internal package structure. Event-driven architecture internally with channels for real-time updates. File-based storage using OS file system for persistence.

### Testing Requirements

Comprehensive unit tests for core logic, integration tests for file operations and hook processing, and manual testing utilities for TUI interactions.

### Additional Technical Assumptions and Requests

- Use Go 1.21+ for improved performance and standard library features
- Leverage Bubbletea for TUI framework with proven stability
- Implement Cobra for CLI command structure and help generation
- Use standard library JSON package for session state serialization
- Integrated Go-based hook system via `spcstr hook` command for <10ms execution
- Atomic file operations using temp file + rename pattern for data integrity
- Avoid CGO dependencies to ensure easy cross-compilation
- Configuration follows XDG Base Directory specification
- Support for both project-local and global configuration

## Epic List

- **Epic 1: Foundation & Core Infrastructure**: Establish project structure, implement init command, and create basic TUI shell with view switching
- **Epic 2: Session Tracking & Persistence**: Implement Go-based hook handlers, session state management with proper data model, and atomic JSON persistence layer
- **Epic 3: Plan View Implementation**: Build document indexing, markdown preview, and navigation features for planning documents
- **Epic 4: Observe View & Dashboard**: Create session list, implement dashboard layout, and enable real-time session monitoring

## Epic 1: Foundation & Core Infrastructure

Establish the foundational Go project structure with Cobra CLI framework and Bubbletea TUI. Implement the init command that sets up hook scripts and configuration. Create the basic TUI shell with global keybindings and view switching capabilities.

### Story 1.1: Project Setup and CLI Structure

As a developer,
I want a well-structured Go project with CLI commands,
so that I can easily extend and maintain the codebase.

#### Acceptance Criteria
1. Go module initialized with appropriate package structure
2. Cobra CLI framework integrated with root command
3: Main entry point established with proper error handling
4: Makefile created for common development tasks (build, test, install)
5: Basic README with project description and development setup

### Story 1.2: Init Command Implementation

As a user,
I want to initialize spcstr in my project with a single command,
so that all necessary hooks and configuration are automatically set up.

#### Acceptance Criteria
1: `spcstr init` creates .spcstr/ directory structure
2: Hook scripts are written to .spcstr/hooks/ directory
3: Claude settings.json is updated with hook configurations
4: Command detects existing configuration and prompts before overwriting
5: Success message confirms initialization with next steps

### Story 1.3: Go Hook System Implementation

As a developer,
I want an integrated Go-based hook system,
so that Claude Code events are processed with type safety and performance.

#### Acceptance Criteria
1: `spcstr hook <hook-name>` command processes all Claude hook events
2: Hook execution completes within 10ms including file I/O
3: Session state uses proper data model with categorized file operations
4: Atomic writes prevent data corruption during concurrent operations
5: Claude settings.json updated to use Go hook commands
6: Existing shell hooks from Story 1.2 are detected and migrated
7: Tool usage tracking with accurate counts per tool type
8: Agent history maintained with timestamps for all executions

### Story 1.4: Basic TUI Shell with View Management

As a user,
I want a responsive TUI with keyboard navigation between views,
so that I can switch between planning and observability modes.

#### Acceptance Criteria
1: TUI launches and displays initial view within 500ms
2: Global keybindings (p, o, q) work from any context
3: View transitions are smooth without flicker
4: Terminal resize events are handled gracefully
5: Help text shows available keybindings

### Story 1.5: Configuration Management

As a user,
I want spcstr to respect both project and global configuration,
so that I can customize behavior per project or globally.

#### Acceptance Criteria
1: Project config (.spcstr/settings.json) takes precedence over global
2: Global config location follows XDG specification
3: Default configuration is created if none exists
4: Configuration changes take effect without restart
5: Invalid configuration shows clear error messages

## Epic 2: Session Tracking & Persistence

Implement the core session tracking functionality with integrated Go hook handlers that capture Claude Code events. Create the atomic JSON-based persistence layer with proper data model that maintains session state across restarts. Ensure real-time updates flow from hooks to the TUI with sub-100ms latency.

### Story 2.1: Hook Event Processing Implementation

As a developer,
I want Go-based hook handlers that process all Claude Code events,
so that session data is accurately tracked with type safety.

#### Acceptance Criteria
1: Go handlers implemented for all 9 Claude hooks (pre-tool-use, post-tool-use, etc.)
2: Handlers validate and create session directories atomically
3: Session ID generation uses format sess_{uuid} consistently
4: Handlers never block Claude (exit 0 for success, 2 for blocking)
5: Performance validated at <10ms per hook execution
6: Safety checks implemented for dangerous operations (rm -rf, .env access)

### Story 2.2: Session State Management

As a system,
I want to maintain session state with proper data model in structured JSON,
so that sessions can be persisted and recovered accurately.

#### Acceptance Criteria
1: Session state includes: agents, agents_history, categorized files (new/edited/read), tools_used map, errors array
2: State updates use atomic writes (temp file + rename)
3: File-based locking prevents concurrent corruption
4: Session files remain human-readable JSON
5: Modified flag tracks dirty state without persistence
6: AgentHistoryEntry tracks start/end times for all agents

### Story 2.3: Real-time Data Pipeline

As a user,
I want session updates to appear immediately in the TUI,
so that I can monitor activity in real-time.

#### Acceptance Criteria
1: File watchers detect session changes within 100ms
2: Updates flow to UI without blocking user interaction
3: Batch updates are processed efficiently
4: Memory usage remains bounded with many updates
5: Recovery from missed updates is automatic

### Story 2.4: Session File Operations Tracking

As a user,
I want to see all files that were created, edited, or read in a session,
so that I understand the scope of changes.

#### Acceptance Criteria
1: File operations categorized into FileOperations struct (new/edited/read arrays)
2: Full absolute paths captured from tool usage
3: Duplicate entries prevented within each category
4: Tool usage counts maintained in tools_used map
5: Error tracking with timestamp, hook name, and message

## Epic 3: Plan View Implementation

Build the Plan View that provides access to planning documents with multiple modes (Spec, Workflow, Config). Implement document indexing with search capabilities and markdown preview with syntax highlighting.

### Story 3.1: Document Indexing and Discovery

As a user,
I want spcstr to automatically find all planning documents,
so that I can quickly navigate to any document.

#### Acceptance Criteria
1: Recursive search finds all markdown files in configured paths
2: Documents are organized by type (PRD, Architecture, Epic, Story)
3: Search supports fuzzy matching on file names and content
4: Index updates automatically when files change
5: Performance remains fast with hundreds of documents

### Story 3.2: Markdown Preview with Syntax Highlighting

As a user,
I want to preview markdown documents with proper formatting,
so that I can read planning documents without leaving the TUI.

#### Acceptance Criteria
1: Markdown renders with headers, lists, and emphasis
2: Code blocks have syntax highlighting
3: Links are visually distinguished
4: Preview updates in real-time during edits
5: Large documents scroll smoothly

### Story 3.3: Plan View Mode Navigation

As a user,
I want to switch between Spec, Workflow, and Config modes,
so that I can access different types of planning information.

#### Acceptance Criteria
1: Mode switching via keyboard shortcuts (s, w, c)
2: Each mode shows relevant document categories
3: Mode state persists across view switches
4: Visual indicator shows current mode
5: Mode-specific keybindings are available

### Story 3.4: Document Navigation and Selection

As a user,
I want to navigate and select documents efficiently,
so that I can quickly access the information I need.

#### Acceptance Criteria
1: Arrow keys navigate document list
2: Tab switches focus between list and preview
3: Enter opens document for detailed view
4: Search narrows list in real-time
5: Recently accessed documents appear first

## Epic 4: Observe View & Dashboard

Create the Observe View with session list and detailed dashboard. Implement the orchestration dashboard layout with organized sections for agents, tasks, files, and metrics. Enable session management capabilities.

### Story 4.1: Session List Implementation

As a user,
I want to see all active and completed sessions,
so that I can select and monitor any session.

#### Acceptance Criteria
1: Sessions display with ID, status, and duration
2: Active sessions appear at top with visual indicator
3: List updates automatically as sessions change
4: Sessions can be filtered by status or date
5: Completed sessions show completion time

### Story 4.2: Orchestration Dashboard Layout

As a user,
I want a comprehensive dashboard showing session details,
so that I understand what's happening in each session.

#### Acceptance Criteria
1: Dashboard sections for agents, agents_history, files, tools, errors
2: Agents show current active agents and complete agent history with timestamps
3: Files displayed in categories: new, edited, read
4: Tool usage shows as map with tool name and execution count
5: Errors section displays timestamp, hook name, and message
6: Layout adjusts to terminal size intelligently

### Story 4.3: Real-time Dashboard Updates

As a user,
I want the dashboard to update in real-time,
so that I can monitor ongoing activity.

#### Acceptance Criteria
1: Updates appear within 1 second of hook trigger
2: Animations indicate changing values
3: Historical data remains visible with timestamps
4: Performance remains smooth with rapid updates
5: User can pause updates to examine details

### Story 4.4: Session Management Controls

As a user,
I want to manage sessions from the TUI,
so that I can kill or resume sessions as needed.

#### Acceptance Criteria
1: Ctrl+K kills selected session with confirmation
2: Ctrl+R resumes stopped session if possible
3: Session details can be exported to file
4: Clear error messages for invalid operations
5: Audit log tracks session management actions

## Checklist Results Report

*To be populated after PM checklist execution*

## Next Steps

### UX Expert Prompt

Please create the front-end specification for Spec⭐️ using the provided PRD. Focus on terminal-based UI patterns, ASCII art layouts for the dashboard, and detailed keybinding maps for all interactions. The specification should guide implementation of a Bubbletea-based TUI.

### Architect Prompt

Please create the fullstack architecture document for Spec⭐️ using the provided PRD. Focus on Go package structure, integration with Claude Code hooks, file-based session persistence patterns, and the event-driven architecture for real-time updates. Consider scalability for large projects with many sessions.