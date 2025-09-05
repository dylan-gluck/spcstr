# Goals and Background Context

## Goals
- Ship working MVP with `spcstr init` command that installs hooks and creates `.spcstr/` directory structure
- Implement exact state management system tracking session data to `.spcstr/sessions/{session-id}/state.json`
- Deliver TUI with Plan view (spec document browser) and Observe view (session dashboard)
- Enable 100% automatic Claude Code session tracking via 9 hook executables
- Provide real-time observability into agents, tasks, files, and tool usage
- Build foundation for future iteration with clean Go monorepo architecture

## Background Context

spcstr provides multi-agent observability for Claude Code sessions integrated with the BMad Method. By installing a single binary and running `spcstr init`, developers gain comprehensive visibility into AI-assisted development sessions with automatic state tracking through a hook-based architecture. The system captures real-time session state, tracks file operations, monitors agent activities, and manages spec-driven development workflows through an intuitive TUI.

## Change Log

| Date | Version | Description | Author |
|------|---------|-------------|--------|
| 2025-09-05 | v1.0 | Initial PRD creation focused on MVP | John (PM) |
