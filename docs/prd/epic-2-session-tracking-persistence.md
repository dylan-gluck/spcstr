# Epic 2: Session Tracking & Persistence

Implement the core session tracking functionality with hook scripts that capture Claude Code events. Create the JSON-based persistence layer that maintains session state across restarts. Ensure real-time updates flow from hooks to the TUI.

## Story 2.1: Hook Script Implementation

As a developer,
I want hook scripts that capture all Claude Code events,
so that session data is automatically tracked.

### Acceptance Criteria
1: Shell scripts created for all critical Claude hooks
2: Scripts validate and create session directories as needed
3: Session ID generation is unique and consistent
4: Scripts handle errors gracefully without blocking Claude
5: Minimal performance impact (<10ms per hook execution)

## Story 2.2: Session State Management

As a system,
I want to maintain session state in structured JSON,
so that sessions can be persisted and recovered.

### Acceptance Criteria
1: Session JSON schema includes all required fields
2: State updates are atomic to prevent corruption
3: Concurrent modifications are handled safely
4: Session files are human-readable and editable
5: Old sessions are archived after configurable duration

## Story 2.3: Real-time Data Pipeline

As a user,
I want session updates to appear immediately in the TUI,
so that I can monitor activity in real-time.

### Acceptance Criteria
1: File watchers detect session changes within 100ms
2: Updates flow to UI without blocking user interaction
3: Batch updates are processed efficiently
4: Memory usage remains bounded with many updates
5: Recovery from missed updates is automatic

## Story 2.4: Session File Operations Tracking

As a user,
I want to see all files that were created, edited, or read in a session,
so that I understand the scope of changes.

### Acceptance Criteria
1: File operations are categorized correctly (new/edited/read)
2: Full absolute paths are captured for all files
3: Duplicate entries are prevented within categories
4: File lists are sorted alphabetically for readability
5: Relative paths can be displayed optionally
