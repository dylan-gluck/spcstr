# Requirements

## Functional Requirements

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

## Non-Functional Requirements

**NFR1:** TUI response time <100ms for view switching and navigation

**NFR2:** Hook execution must not block Claude Code operations (except when returning exit code 2)

**NFR3:** State file updates must be atomic to prevent corruption

**NFR4:** Support 256-color terminal emulators on macOS, Linux, Windows (WSL)

**NFR5:** Memory footprint <10MB for hook executables

**NFR6:** Single Go monorepo with embedded hook binaries compiled at build time
