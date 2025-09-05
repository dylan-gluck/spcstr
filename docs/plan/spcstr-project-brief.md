# Project Brief: Spec⭐️

## Executive Summary

**spcstr** (pronounced "spec star") is a Terminal User Interface (TUI) application providing multi-agent observability and orchestration for Claude Code sessions integrated with the BMad Method. The system captures real-time session state through a hook-based architecture, enabling developers to monitor agent activities, track file operations, and manage spec-driven development workflows through an intuitive terminal interface. By installing a single binary and initializing hooks in their project, users gain comprehensive visibility into AI-assisted development sessions with automatic state tracking and rich observability dashboards.

## Problem Statement

Modern AI-assisted development with Claude Code and multi-agent systems lacks visibility into session state, agent orchestration, and workflow progress. Developers currently experience:

- **Zero visibility** into active agent states, task progress, and system resource usage during Claude Code sessions
- **No persistent tracking** of file operations, tool invocations, or error states across development sessions
- **Manual overhead** in managing planning documents, PRDs, architecture specs, and workflow definitions
- **Lost context** between sessions with no historical record of agent activities or decision paths

The impact is significant: developers waste time recreating context, debugging invisible agent failures, and manually tracking progress across complex multi-step workflows. Current solutions require manual instrumentation or provide only partial visibility, failing to capture the full lifecycle of AI-assisted development sessions.

## Proposed Solution

spcstr provides a comprehensive observability framework through a single-binary architecture:

1. **Embedded Hook System**: Hook commands integrated as Cobra subcommands (`spcstr hook <name>`) that respond to Claude Code's lifecycle events to capture session state in real-time
2. **TUI Application**: Rich terminal interface for viewing session data, managing planning documents, and monitoring agent orchestration

The solution uniquely combines automatic state capture with zero configuration overhead - users simply run `spcstr init` to enable comprehensive tracking. Unlike existing solutions, spcstr maintains complete session history in structured JSON, enabling both real-time monitoring and historical analysis while remaining fully local and privacy-preserving.

## Target Users

### Primary User Segment: AI-Native Developers

- **Profile**: Software engineers actively using Claude Code for development, typically working on complex projects with multiple agents
- **Current workflow**: Iterative development with Claude Code, frequent agent switching, multi-file operations
- **Pain points**: Lack of visibility into agent decisions, lost context between sessions, difficulty tracking file changes
- **Goals**: Maintain development velocity while gaining insight into AI-assisted workflows

### Secondary User Segment: Engineering Teams

- **Profile**: Development teams adopting AI-assisted development practices, needing standardization and observability
- **Current workflow**: Shared Claude Code usage across team members, collaborative spec-driven development
- **Pain points**: No team-wide visibility standards, inconsistent planning documentation, difficulty onboarding new members
- **Goals**: Establish consistent AI development practices with proper observability and documentation

## Goals & Success Metrics

### Business Objectives
- Achieve 80% reduction in context recreation time between Claude Code sessions within 30 days of adoption
- Enable 100% automatic capture of session state without manual instrumentation
- Support 10+ concurrent agent tracking with <100ms UI latency

### User Success Metrics
- Time to initialize project: <30 seconds
- Session state capture accuracy: 100% of Claude Code events
- Dashboard refresh rate: Real-time (<500ms latency)

### Key Performance Indicators (KPIs)
- **Adoption Rate**: Number of projects with spcstr initialized / total Claude Code projects
- **Session Completeness**: Percentage of sessions with full state capture (start to end)
- **User Engagement**: Average time spent in TUI per session

## MVP Scope

### Core Features (Must Have)

- **Binary Installation**: Single executable installable via package managers (brew, apt, pacman)
- **Project Initialization**: `spcstr init` command creating `.spcstr/` directory structure and configuring Claude Code hooks
- **Hook Commands**: Complete implementation of 9 Claude Code hooks as embedded subcommands following exact state management specification
- **Session State Tracking**: Real-time capture and persistence of session data to `.spcstr/sessions/{session-id}/state.json`
- **TUI Plan View**: Document browser for PRDs, architecture specs with markdown preview
- **TUI Observe View**: Session list and detail dashboard showing agents, tasks, files, and activity
- **Atomic State Updates**: File-system safe state mutations with temp file + rename pattern

### Out of Scope for MVP
- Cloud synchronization or remote storage
- Custom agent definitions or orchestration rules
- IDE integrations beyond Claude Code
- Web-based dashboard
- Team collaboration features
- Custom theming beyond light/dark modes

### MVP Success Criteria

The MVP succeeds when a developer can install spcstr, run `spcstr init` in their project, and immediately see real-time session data in the TUI during their next Claude Code session, with 100% of hooks firing correctly and state persisting across session restarts.

## Post-MVP Vision

### Phase 2 Features

- **Advanced Analytics**: Historical trend analysis, agent performance metrics, error pattern detection
- **Workflow Automation**: Custom workflow definitions, automated agent orchestration based on patterns
- **Team Features**: Shared session visibility, team dashboards, collaborative planning documents
- **Export Capabilities**: Session replay, audit logs, compliance reporting

### Long-term Vision

Transform spcstr into the definitive observability platform for AI-assisted development, supporting multiple AI coding assistants beyond Claude Code, with enterprise features for compliance, security, and team governance.

### Expansion Opportunities

- Integration with other AI coding tools (GitHub Copilot, Cursor, etc.)
- SaaS offering for team analytics and insights
- Marketplace for custom workflow templates and agent configurations
- API for programmatic access to session data

## Technical Considerations

### Platform Requirements
- **Target Platforms**: macOS (primary), Linux, Windows (WSL)
- **Browser/OS Support**: Terminal emulators with 256-color support
- **Performance Requirements**: <100ms UI response time, <10MB memory footprint for hook commands

### Technology Preferences
- **Frontend**: Bubbletea TUI framework with Lipgloss styling
- **Backend**: Pure Go for both TUI and hook commands
- **Database**: File-based JSON storage in `.spcstr/` directory
- **Hosting/Infrastructure**: Local-only, no cloud dependencies

### Architecture Considerations
- **Repository Structure**: Monolithic Go module with cmd/ for binaries and internal/ for shared code
- **Service Architecture**: Single binary with subcommands via Cobra, hook functionality embedded as subcommands
- **Integration Requirements**: Claude Code hooks API v1, file system permissions for project directories
- **Security/Compliance**: Local-only data storage, no network calls, no sensitive data capture

## Constraints & Assumptions

### Constraints
- **Budget**: Open-source project, no paid dependencies
- **Timeline**: MVP within 2-3 days
- **Resources**: Single developer initially, community contributions post-launch
- **Technical**: Must work within Claude Code hook limitations, file system based state only

### Key Assumptions
- Claude Code hook API remains stable and backward compatible
- Users have appropriate file system permissions in project directories
- Terminal environments support standard ANSI escape codes
- Go 1.21+ is acceptable as minimum version requirement
- JSON state files under 10MB perform adequately

## Risks & Open Questions

### Key Risks
- **Hook API Changes**: Claude Code modifying hook interface could break integration
- **Performance at Scale**: Large session states (>1000 operations) may impact TUI responsiveness
- **Cross-platform Compatibility**: Terminal rendering differences across OS platforms

### Open Questions
- How to handle hook failures gracefully without blocking Claude Code operations?
- Should state files be automatically pruned after certain age/size?
- What's the optimal update frequency for TUI dashboard refresh?

## Appendices

### A. Research Summary

**Claude Code Hook System Analysis:**
- 9 distinct lifecycle hooks available
- JSON input/output with exit codes for control flow
- Hooks can block operations with exit code 2
- Maximum execution time constraints per hook type

**TUI Framework Evaluation:**
- Bubbletea provides robust event handling and state management
- Lipgloss enables consistent cross-platform styling
- Cobra offers industry-standard CLI patterns

### B. Stakeholder Input

Initial user feedback indicates strong demand for:
- Real-time agent visibility
- Historical session replay
- Error tracking and debugging
- Integration with existing development workflows

### C. References

- Claude Code Hooks Documentation: `docs/vendor/cc-hooks-*.md`
- Bubbletea Framework: https://github.com/charmbracelet/bubbletea
- BMad Method: https://github.com/bmad-code-org/BMAD-METHOD
- Project Specifications: `docs/plan/spcstr-plan.md`, `docs/plan/hooks-state-management.md`

## Next Steps

### Immediate Actions
1. Set up Go monorepo structure with Cobra and Bubbletea dependencies
2. Implement core state management module with atomic write operations
3. Create hook command implementations with JSON parsing
4. Build TUI scaffold with Plan and Observe views
5. Implement `spcstr init` command with hook configuration logic
6. Test end-to-end with sample Claude Code session
7. Package binaries for brew formula creation

### PM Handoff

This Project Brief provides the full context for spcstr. Please start in 'PRD Generation Mode', review the brief thoroughly to work with the user to create the PRD section by section as the template indicates, asking for any necessary clarification or suggesting improvements. Key areas requiring detailed specification include hook error handling strategies, TUI keyboard navigation patterns, and state file rotation policies.
