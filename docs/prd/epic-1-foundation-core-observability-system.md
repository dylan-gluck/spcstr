# Epic 1 - Foundation & Core Observability System

**Goal:** Establish the complete spcstr system with project initialization, embedded hook commands for Claude Code session tracking, and a TUI providing real-time observability into agent activities, file operations, and task progress. This epic delivers the entire MVP as a single, cohesive binary.

## Story 1.1: Project Structure and Build System

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

## Story 1.2: State Management Package

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

## Story 1.3: Hook Command Implementation

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

## Story 1.4: CLI and Init Command

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

## Story 1.5: TUI Foundation and Navigation

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

## Story 1.6: Plan View Implementation

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

## Story 1.7: Observe View Implementation

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

## Story 1.8: Integration Testing and Polish

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
