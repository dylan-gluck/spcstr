
# Spec⭐️

## Multi Session Observability TUI for Claude Code

An opinionated orchestration framework, designed to work with [Claude-code](https://github.com/anthropics/claude-code) & [BMad Method](https://github.com/bmad-code-org/BMAD-METHOD)

## `spcstr`

Single executable, installed via package manager.

```
brew install spcstr
```

### Init Command

```bash
# Initialize project hooks, scripts, settings

spcstr init

# Start the TUI in the current directory

spcstr
```

*Once initialized, all claude-code session data will be automatically logged to the project sessions folder: `.spcstr/sessions/{session-id}/state.json`*

The session state json file is updated by `sh` scripts which are triggered by claude-code hooks.[^1] On *init* the scripts and hooks are added to the project.

**Init process:**
1. Create `.spcstr/` dir
2. Write hook scripts to `.spcstr/hooks/*.sh`
3. Add hook settings to `.claude/settings.json`

## TUI

The `spcstr` TUI has different *views*, toggled with keystrokes. Each view may have multiple *modes* or pages within the view.

Global key-binds are available at all times. These are used to change views, exit program.

```markdown
Global key-binds:
[p] -> Change view: *Plan*
[o] -> Change view: *Observe*
[q] -> Quit
```

### Plan View

Main dashboard to view and create planning documents. Documents are recursively indexed, `fzf`.

**Layout**

- Left column: list of planning docs in project
- Right column: selected document markdown preview / edit

**Key-binds**

```markdown
[s]  -> Change mode: *Spec*
[w]  -> Change mode: *Workflow*
[c]  -> Change mode: *Config*

[tab]     -> Focus column (toggle)
[up/down] -> Cursor / selection movement
```

*Additional commands for toggling edit mode when a document is focused, vim key-binds*

*Additional commands for triggering BMAD workflows / starting new sessions from a document*

##### SPEC Mode

Dashboard that indexes and tracks all project docs:

```
PRD                     → docs/prd.md
PRDs Shared             → docs/prd/*.md
Architecture            → docs/architecture.md
Architecture Shared     → docs/architecture/*.md
```

##### Workflow Mode

Workflow builder (epics, stories):

```
Sharded Epics    → docs/epics/
Sharded Stories  → docs/stories/
```

##### Config Mode

View active config settings for `claude`, `bmad`, `spcstr`

```
Claude config     ->  .claude/
BMAD config       ->  .bmad-core/
Spcstr config     ->  .spcstr/settings.json
```

### Observe View

Session observability dashboard

**Layout**

- Left column: list of active / completed sessions
- Right column: selected session details

Each session tracks agents, tasks, files changes, shared context.

**Key-binds**

```markdown
[up/down]  -> Cursor / selection movement
[enter]    -> Select session
[^k]       -> Kill session
[^r]       -> Resume session
```

##### Observability Mode

The selected session data should be rendered in a dashboard style layout with tables and charts.

**Example Dashboard Layout**

```
╭─── ORCHESTRATION DASHBOARD ─────────────────────────────────────────────╮
  │ Session: {session_id}          Runtime: {runtime}          Mode: {mode} │
  ├──────────────────────────────────────────────────────────────────────────┤
  │ TEAMS                              │ ACTIVE AGENTS                      │
  │ ├─ Engineering ({eng_count})       │ • {active_agent_1} [{status}]     │
  │ │  └─ {eng_status}                 │ • {active_agent_2} [{status}]     │
  │ ├─ Product ({prod_count})          │ • {active_agent_3} [{status}]     │
  │ │  └─ {prod_status}                │ • {active_agent_4} [{status}]     │
  │ ├─ QA ({qa_count})                 │ • {active_agent_5} [{status}]     │
  │ │  └─ {qa_status}                  │                                    │
  │ └─ DevOps ({devops_count})         │ Load: {system_load}                │
  │    └─ {devops_status}              │ Tasks: {task_queue_length}         │
  ├──────────────────────────────────────────────────────────────────────────┤
  │ SPRINT: {sprint_name}                           VELOCITY: {velocity}    │
  │ ┌────────────────────────────────────────────────────────────────────┐  │
  │ │ Todo ({todo_count})  │ In Progress ({wip})  │ Done ({done_count})  │  │
  │ │ {todo_tasks}          │ {wip_tasks}          │ {done_tasks}         │  │
  │ └────────────────────────────────────────────────────────────────────┘  │
  │ Blockers: {blocker_count}  |  At Risk: {at_risk_count}                  │
  ├──────────────────────────────────────────────────────────────────────────┤
  │ RECENT ACTIVITY                                                          │
  │ {timestamp_1} │ {event_1}                                               │
  │ {timestamp_2} │ {event_2}                                               │
  │ {timestamp_3} │ {event_3}                                               │
  ╰──────────────────────────────────────────────────────────────────────────╯

  Commands: /help | /team <name> | /sprint | /metrics | @<agent> | #<sprint>
  > {command_prompt}
```

**Example State JSON**

The following json is real session data from a previous iteration of this system:

 `.spcstr/sessions/6edc20cd-6688-4f7f-af04-e2b1cbcc123b/state.json`

```json
{
  "session_id": "6edc20cd-6688-4f7f-af04-e2b1cbcc123b",
  "created_at": "2025-08-28T15:25:35.255593Z",
  "updated_at": "2025-08-28T15:54:30.166887Z",
  "source": "startup",
  "workflow": "",
  "agents": [],
  "agents_history": [
    {
      "name": "engineering-tests",
      "started_at": "2025-08-28T15:51:20.902890Z"
    }
  ],
  "files": {
    "new": [
      "/Users/dylan/Workspace/projects/atspro-bun/migrations/006_update_subscription_tiers.sql",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/subscription.remote.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/components/subscription/subscription-badge.svelte",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/components/rate-limit-toast.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/scripts/fix-subscription-migration.js",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/__tests__/rate-limit.test.ts"
    ],
    "edited": [
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/rate-limit.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/thoughts/shared/plans/subscription-tier-implementation.md",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/utils.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/routes/(app)/app/+layout.svelte",
      "/Users/dylan/Workspace/projects/atspro-bun/src/routes/(app)/app/settings/+page.svelte",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/document.remote.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/resume.remote.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/job.remote.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/migrations/006_update_subscription_tiers.sql",
      "/Users/dylan/Workspace/projects/atspro-bun/scripts/fix-subscription-migration.js",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/components/subscription/subscription-badge.svelte"
    ],
    "read": [
      "/Users/dylan/Workspace/projects/atspro-bun/thoughts/dylan/features/subscription-logic.md",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/rate-limit.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/migrations/004_add_subscription_tiers.sql",
      "/Users/dylan/Workspace/projects/atspro-bun/src/routes/(app)/app/settings/+page.svelte",
      "/Users/dylan/Workspace/projects/atspro-bun/src/routes/(app)/app/+layout.svelte",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/utils.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/db/index.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/components/ui/badge/badge.svelte",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/resume.remote.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/document.remote.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/job.remote.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/scripts/migrate.js",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/rate-limit.remote.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/__tests__/rate-limit.test.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/__tests__/job.remote.test.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/__tests__/resume.remote.test.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/services/subscription.remote.ts",
      "/Users/dylan/Workspace/projects/atspro-bun/src/lib/components/subscription/subscription-badge.svelte",
      "/Users/dylan/Workspace/projects/atspro-bun/package.json"
    ]
  },
  "tools_used": {
    "TodoWrite": 19,
    "Read": 22,
    "Write": 6,
    "Edit": 25,
    "Glob": 7,
    "LS": 2,
    "Bash": 25,
    "MultiEdit": 1,
    "Grep": 7,
    "mcp__postgres__query": 6,
    "Task": 1
  },
  "errors": []
}
```

The dashboard must display this data in organized, human readable tables charts and lists that update in real time when the json changes.

*Toggle between dashboard view & streaming JSON output for a specific hook or the stdout of the main claude session* [^2][^3]

---

## Architecture

The Spec⭐️ system is comprised of two main components:

- The TUI program that the user directly interacts with
- The hook scripts that are triggered during claude lifecycle

#### TUI Development Stack: Go

- [Cobra](https://github.com/spf13/cobra)
- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)
- [JSON/v2](https://pkg.go.dev/encoding/json/v2)

#### Hook Scripts

Lightweight `sh` scripts that log session data & update shared state. Scripts are called by claude-code hook events. There is one script for each hook.

**On Init**
- The scripts are written to the project: `.spcstr/hooks/*.sh`
- The hooks settings are added to `.claude/settings.json`

### Spcstr Config

Project config `.spcstr/settings.json` takes precedence, default config at `$XDG_CONFIG_HOME/spcstr/settings.json`

**Settings**
- Key-binds
- Theme

---
**Footnotes**

[^1]: Claude-code hooks: https://docs.anthropic.com/en/docs/claude-code/hooks
[^2]: Claude-code SDK: https://docs.anthropic.com/en/docs/claude-code/sdk
[^3]: Claude-code CLI flags: https://docs.anthropic.com/en/docs/claude-code/cli-reference#cli-flags
