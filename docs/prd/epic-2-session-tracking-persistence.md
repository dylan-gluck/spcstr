# Epic 2: Session Tracking & Persistence

Implement the core session tracking functionality with integrated Go hook handlers that capture Claude Code events. Create the atomic JSON-based persistence layer with proper data model that maintains session state across restarts. Ensure real-time updates flow from hooks to the TUI with sub-100ms latency.

## Story 2.1: Hook Event Processing Implementation

As a developer,
I want Go-based hook handlers that process all Claude Code events,
so that session data is accurately tracked with type safety.

### Acceptance Criteria
1: Go handlers implemented for all 9 Claude hooks (pre-tool-use, post-tool-use, etc.)
2: Handlers validate and create session directories atomically
3: Session ID generation uses format sess_{uuid} consistently
4: Handlers never block Claude (exit 0 for success, 2 for blocking)
5: Performance validated at <10ms per hook execution
6: Safety checks implemented for dangerous operations (rm -rf, .env access)

## Story 2.2: Session State Management

As a system,
I want to maintain session state with proper data model in structured JSON,
so that sessions can be persisted and recovered accurately.

### Acceptance Criteria
1: Session state includes: agents, agents_history, categorized files (new/edited/read), tools_used map, errors array
2: State updates use atomic writes (temp file + rename)
3: File-based locking prevents concurrent corruption
4: Session files remain human-readable JSON
5: Modified flag tracks dirty state without persistence
6: AgentHistoryEntry tracks start/end times for all agents

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
1: File operations categorized into FileOperations struct (new/edited/read arrays)
2: Full absolute paths captured from tool usage
3: Duplicate entries prevented within each category
4: Tool usage counts maintained in tools_used map
5: Error tracking with timestamp, hook name, and message
