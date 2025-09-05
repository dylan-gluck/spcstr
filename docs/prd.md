# spcstr Product Requirements Document (PRD)

## Goals and Background Context

### Goals
- Ship working MVP with `spcstr init` command that installs hooks and creates `.spcstr/` directory structure
- Implement exact state management system tracking session data to `.spcstr/sessions/{session-id}/state.json`
- Deliver TUI with Plan view (spec document browser) and Observe view (session dashboard)
- Enable 100% automatic Claude Code session tracking via 9 hook executables
- Provide real-time observability into agents, tasks, files, and tool usage
- Build foundation for future iteration with clean Go monorepo architecture

### Background Context

spcstr provides multi-agent observability for Claude Code sessions integrated with the BMad Method. By installing a single binary and running `spcstr init`, developers gain comprehensive visibility into AI-assisted development sessions with automatic state tracking through a hook-based architecture. The system captures real-time session state, tracks file operations, monitors agent activities, and manages spec-driven development workflows through an intuitive TUI.

### Change Log

| Date | Version | Description | Author |
|------|---------|-------------|--------|
| 2025-09-05 | v1.0 | Initial PRD creation focused on MVP | John (PM) |

## Requirements

### Functional Requirements

**FR1:** Binary installation via package managers (brew, apt, pacman) with single `spcstr` executable

**FR2:** `spcstr init` command creates `.spcstr/{logs,sessions,hooks}` directory structure and writes 9 hook executables

**FR3:** Hook executables implement exact state management spec from `docs/plan/hooks-state-management.md`:
- session_start: Initialize state.json with session structure
- user_prompt_submit: Track user prompts in prompts array
- pre_tool_use: Track tool invocations, manage agents for Task tool
- post_tool_use: Track tool completions, file operations (Write/Edit/Read), update agents_history
- notification: Track notification messages
- pre_compact: Log context compaction events
- session_end: Set session_active to false with termination reason
- stop: Handle session termination
- subagent_stop: Log sub-agent termination

**FR4:** Atomic state updates using temp file + rename pattern for filesystem safety

**FR5:** TUI Plan view displays indexed planning documents:
- PRD (docs/prd.md and docs/prd/*.md)
- Architecture (docs/architecture.md and docs/architecture/*.md)
- Workflow mode for epics/stories

**FR6:** TUI Observe view displays session list and real-time dashboard showing:
- Active/completed sessions
- Current agents and agents_history
- Todo/In Progress/Done task counts
- Files created/edited/read
- Tool usage statistics
- Recent activity feed

**FR7:** Global keybinds: [p] Plan view, [o] Observe view, [q] Quit

**FR8:** Session state persisted in JSON format at `.spcstr/sessions/{session-id}/state.json` following exact schema

**FR9:** Event logging to `.spcstr/logs/{hook_name}.json` in append-only JSON array format

### Non-Functional Requirements

**NFR1:** TUI response time <100ms for view switching and navigation

**NFR2:** Hook execution must not block Claude Code operations (except when returning exit code 2)

**NFR3:** State file updates must be atomic to prevent corruption

**NFR4:** Support 256-color terminal emulators on macOS, Linux, Windows (WSL)

**NFR5:** Memory footprint <10MB for hook executables

**NFR6:** Single Go monorepo with embedded hook binaries compiled at build time

## User Interface Design Goals

### Overall UX Vision
Minimalist, keyboard-driven TUI providing instant visibility into Claude Code sessions with zero mouse interaction. Focus on information density and real-time updates with beautiful markdown rendering via Glamour.

### Key Interaction Paradigms
- Pure keyboard navigation with single-key commands
- Two-pane layouts (list/detail pattern) for both Plan and Observe views
- Context-aware footer showing available keybinds that update per view/mode
- Immediate feedback on all actions (no loading states for local operations)

### Core Screens and Views

**Main TUI Screen**
- Header bar showing current view, session status
- Content area switching between Plan and Observe views
- **Footer/status bar displaying context-aware keybinds that change based on current view**

**Plan View**
- Left pane: Document tree (PRD, Architecture, Epics, Stories)
- Right pane: **Glamour-rendered markdown preview** with syntax highlighting and formatting
- Tab to switch focus between panes
- Footer shows: [p]lan [o]bserve [q]uit [tab] switch-pane [↑↓] navigate [s]pec [w]orkflow [c]onfig

**Observe View**
- Left pane: Session list showing ID, status (active/completed), timestamp
- Right pane: Dashboard with agents, tasks, files, tools data from state.json
- Real-time updates when state.json changes
- Footer shows: [p]lan [o]bserve [q]uit [↑↓] navigate [enter] select

### Accessibility: None
MVP focuses on core functionality. Accessibility features deferred to post-MVP.

### Branding
Terminal-native aesthetic using Lipgloss styling with Glamour for rich markdown rendering. Light/dark theme support based on terminal colors. No custom branding or logos.

### Target Device and Platforms: Terminal/CLI Only
- macOS Terminal/iTerm2
- Linux terminal emulators
- Windows Terminal with WSL
- Minimum 80x24 terminal size, optimized for 120x40

## Technical Assumptions

### Repository Structure: Monorepo
Single Go module containing all components - TUI application and hook executables as separate binaries. Clean separation with `cmd/` for binaries and `internal/` for shared code.

### Service Architecture
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

### Testing Requirements
**Minimal testing for MVP.** Focus on manual testing of critical paths: init command, hook execution, state management, TUI navigation. Unit tests only for state management atomic write operations. Full testing pyramid deferred to post-MVP.

### Additional Technical Assumptions and Requests

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

## Epic List

**Epic 1: Foundation & Core Observability System**
Establish complete spcstr system with initialization, hook-based state management, and TUI for session observability

## Epic 1 - Foundation & Core Observability System

**Goal:** Establish the complete spcstr system with project initialization, embedded hook commands for Claude Code session tracking, and a TUI providing real-time observability into agent activities, file operations, and task progress. This epic delivers the entire MVP as a single, cohesive binary.

### Story 1.1: Project Structure and Build System

As a developer,
I want to set up the Go monorepo with proper structure and build configuration,
so that I can compile a single spcstr binary with embedded hook functionality.

**Acceptance Criteria:**
1. Go module initialized with go.mod at project root
2. Directory structure created: `cmd/spcstr/`, `internal/hooks/`, `internal/state/`, `internal/tui/`
3. Single Makefile or build script compiles one spcstr binary
4. Build process runs without errors on macOS and Linux
5. All Cobra, Bubbletea, Lipgloss, and Glamour dependencies properly imported
6. Binary size remains reasonable (<50MB) with embedded functionality

### Story 1.2: State Management Package

As a developer,
I want to implement the shared state management library,
so that both TUI and hook commands can read/write session state atomically.

**Acceptance Criteria:**
1. `internal/state/` package implements InitializeState, LoadState, UpdateState, AtomicWrite functions
2. State operations use temp file + rename pattern for atomic writes
3. JSON marshaling/unmarshaling follows exact schema from hooks-state-management.md
4. State files created at `.spcstr/sessions/{session-id}/state.json`
5. All timestamp fields use ISO8601 format
6. Unit tests verify atomic write behavior

### Story 1.3: Hook Command Implementation

As a developer,
I want to implement all 9 hook commands as Cobra subcommands,
so that Claude Code can trigger session state tracking via `spcstr hook <name>`.

**Acceptance Criteria:**
1. All 9 hooks implemented as `spcstr hook <hook_name>` subcommands
2. Each hook reads JSON from stdin and returns appropriate exit codes
3. `spcstr hook session_start` creates new state.json with initial structure
4. `spcstr hook user_prompt_submit` appends to prompts array
5. `spcstr hook pre_tool_use` manages agents array for Task tool invocations
6. `spcstr hook post_tool_use` tracks file operations and updates agents_history
7. `spcstr hook notification` appends to notifications array
8. `spcstr hook session_end` and `spcstr hook stop` set session_active to false
9. All hooks log to `.spcstr/logs/{hook_name}.json` in append-only format
10. Hook execution completes within Claude Code timeout constraints

### Story 1.4: CLI and Init Command

As a user,
I want to run `spcstr init` to configure my project for Claude Code integration,
so that sessions are automatically tracked via hook commands.

**Acceptance Criteria:**
1. `spcstr init` creates `.spcstr/{logs,sessions}` directory structure
2. Hook settings added to `.claude/settings.json` with `spcstr hook <name>` commands
3. Commands include `--cwd={$CLAUDE_PROJECT_DIR}` parameter
4. Command detects existing `.spcstr/` and prompts for confirmation
5. `--force` flag reinitializes without prompting
6. Success message confirms initialization complete
7. `spcstr version` displays version information
8. Root `spcstr` command without args launches TUI

### Story 1.5: TUI Foundation and Navigation

As a user,
I want to launch the TUI and navigate between views,
so that I can access planning documents and session data.

**Acceptance Criteria:**
1. `spcstr` launches TUI using Bubbletea framework
2. Header shows current view and session status
3. Footer displays context-aware keybinds that update per view
4. [p] switches to Plan view, [o] switches to Observe view, [q] quits
5. TUI detects if project not initialized and prompts to run init
6. View switching occurs in <100ms
7. Terminal resize handled gracefully

### Story 1.6: Plan View Implementation

As a user,
I want to browse planning documents in the TUI,
so that I can review PRDs, architecture docs, and workflows.

**Acceptance Criteria:**
1. Left pane shows document tree: PRD, Architecture, Epics, Stories
2. Documents indexed from configured paths (docs/prd.md, docs/architecture.md, etc.)
3. Right pane displays Glamour-rendered markdown with syntax highlighting
4. [tab] switches focus between panes
5. Arrow keys navigate document list
6. [s]pec, [w]orkflow, [c]onfig modes available (can be minimal for MVP)
7. Markdown rendering updates when selecting different documents

### Story 1.7: Observe View Implementation

As a user,
I want to view real-time session data in the TUI,
so that I can monitor agent activities and file operations.

**Acceptance Criteria:**
1. Left pane lists all sessions from `.spcstr/sessions/` with ID and status
2. Sessions marked as active/completed based on session_active field
3. Right pane shows dashboard for selected session with data from state.json
4. Dashboard displays: current agents, agents_history, files (new/edited/read), tools_used
5. File watching via fsnotify triggers dashboard refresh when state.json changes
6. Recent activity shows prompts and notifications in chronological order
7. Task counts displayed if TODO data available in state

### Story 1.8: Integration Testing and Polish

As a developer,
I want to test the complete system end-to-end,
so that all components work together seamlessly.

**Acceptance Criteria:**
1. Manual test: `spcstr init` successfully configures a fresh project
2. Manual test: Claude Code session triggers hook commands and updates state.json
3. Manual test: TUI displays real-time updates during active session
4. Manual test: Plan view correctly renders markdown documents
5. Hook commands execute without blocking Claude Code operations
6. Error messages are helpful when issues occur
7. README.md documents installation and usage

## Checklist Results Report

## PRD & EPIC VALIDATION SUMMARY

### Executive Summary

**Overall PRD Completeness:** 85%  
**MVP Scope Appropriateness:** Just Right  
**Readiness for Architecture Phase:** Ready  
**Most Critical Concerns:** Minor gaps in error handling and performance monitoring

### Category Analysis Table

| Category                         | Status  | Critical Issues |
| -------------------------------- | ------- | --------------- |
| 1. Problem Definition & Context  | PASS    | None |
| 2. MVP Scope Definition          | PASS    | None |
| 3. User Experience Requirements  | PASS    | None |
| 4. Functional Requirements       | PASS    | None |
| 5. Non-Functional Requirements   | PARTIAL | Limited error handling specificity |
| 6. Epic & Story Structure        | PASS    | None |
| 7. Technical Guidance            | PASS    | None |
| 8. Cross-Functional Requirements | PARTIAL | Monitoring approach minimal |
| 9. Clarity & Communication       | PASS    | None |

### Top Issues by Priority

**HIGH:**
- Error handling requirements could be more specific (NFR4 references but doesn't detail recovery strategies)
- Performance monitoring approach specified but minimal detail

**MEDIUM:**
- Data retention policies not addressed (logs will grow indefinitely)
- Support requirements not documented (though appropriate for open-source MVP)

**LOW:**
- Future integration points could be more explicit
- Configuration management deferred but mentioned

### MVP Scope Assessment

**Scope is appropriately minimal:**
- Single binary with embedded hooks eliminates distribution complexity
- Two-view TUI focuses on core observability needs
- State management system is well-defined and bounded
- No feature creep evident

**No missing essential features identified**

**Complexity is manageable:**
- File watching for real-time updates is standard Go practice
- JSON-based state management is simple and reliable
- Hook system follows Claude Code specifications exactly

**Timeline appears realistic** for focused 2-3 day sprint

### Technical Readiness

**Technical constraints are clear:**
- Go 1.21+ requirement specified
- Single binary architecture defined
- Exact dependency list provided (Cobra, Bubbletea, Lipgloss, Glamour)

**Technical risks identified and mitigated:**
- Atomic writes specified for state safety
- Hook timeout constraints acknowledged
- File watching performance considered

**Areas for architect investigation:** None blocking, all implementation details

### Recommendations

1. **Address error handling specifics:** Define what happens when hooks fail or TUI encounters corrupted state files
2. **Add basic monitoring:** Specify how to detect if hooks are working correctly
3. **Consider log rotation:** Add basic policy to prevent log files from growing indefinitely

### Final Decision

**READY FOR ARCHITECT**: The PRD and epics are comprehensive, properly structured, and ready for architectural design. The identified gaps are minor and won't block implementation planning.

## Next Steps

### UX Expert Prompt

The spcstr PRD is complete and ready for UX design. Please create the information architecture and interaction design for the TUI, focusing on the two-view system (Plan and Observe) with keyboard navigation patterns. Pay special attention to the real-time dashboard layout for session observability and the document browser experience for planning materials.

### Architect Prompt

The spcstr PRD is complete and validated. Please create the technical architecture document specifying the Go monorepo structure, state management implementation, hook command system, and TUI framework integration. Focus on the single-binary approach with embedded hook functionality, atomic state operations, and real-time file watching for dashboard updates. The exact state schema and hook specifications are documented in the referenced planning documents.