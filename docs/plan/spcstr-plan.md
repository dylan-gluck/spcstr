
# Spec⭐️

## Multi Session Observability TUI for Claude Code

An opinionated orchestration and observibility framework, designed to work with [Claude-code](https://github.com/anthropics/claude-code) & [BMad Method](https://github.com/bmad-code-org/BMAD-METHOD)

*Pronounced spec star, always written as spcstr in lowercase.*

## `spcstr`

Spcstr is a TUI for multi-agent observibility and spec driven development (SDD).

User's can install the spcstr binary via package manager (brew, pacman, apt, etc).

```
brew install spcstr
```

Once installed, user can initialize spcstr in a project's root directory by running `spcstr init`. This creates a `.spcstr/` directory in their project and adds the hook commands to the project's claude settings `.claude/settings.json`.


### Init Command

```bash
# Initialize project hooks, scripts, settings
spcstr init
```

**Init process:**
1. Create `.spcstr/` directory
2. Create `.spcstr/{logs,sessions,hooks}` directory
3. Write hook executibles to `.spcstr/hooks/*`
4. Add hook settings to `.claude/settings.json`

The session state json file is updated by small executibles in `.spcstr/hooks/*` which are triggered by claude-code hooks.[^1] On *spcstr init* the executibles are added to the project and local claude hook settings are updated.

*Once initialized, all claude-code session data will be automatically logged to the project sessions folder: `.spcstr/sessions/{session-id}/state.json`*


#### State Management & Hook Executibles

A core component of `spcstr` is a state management system that is directly integrated with claude-code lifecycle events, "hooks".

For each hook, there will be a matching executible in the `.spcstr/hooks/*` directory. The source-code for these hooks should be written in the same Go monolith and compiled with the main program at build. They will be written to the project hooks directory on init.

The specifications for each hook and the full state management system are documented in `docs/plan/hooks-state-management.md`. **This spec must be followed exactly.**

The latest claude-code hooks documentation may be used for additional reference when writing the architecture plan & prd documents.
- `docs/vendor/cc-hooks-guide.md`
- `docs/vendor/cc-hooks-reference.md`


## TUI

Once installed, a user can run `spcstr` in a project directory to launch the main TUI application.

*If initialization has not been completed, the user will be prompted to complete it.*

The `spcstr` TUI has different *views*, toggled with keystrokes. Each view may have multiple *modes* or pages within the view.

Global key-binds are available at all times. These are used to change views, exit program.

```markdown
Global key-binds:
[p] -> Change view: *Plan*
[o] -> Change view: *Observe*
[q] -> Quit
```

### Plan View

Main dashboard to view and create planning documents. Documents are recursively indexed, `fzf`, `jq`.

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

### Observe View

Session observability dashboard

**Layout**

- Left column: list of active / completed sessions
- Right column: selected session details

**Key-binds**

```markdown
[up/down]  -> Cursor / selection movement
[enter]    -> Select session
```

##### Observability Mode

The selected session data should be rendered in a dashboard style layout with tables and charts. The following design is just an example of layout, the final data output must reflect the state management system data model.

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
  │ EPIC: {epic_name}                           STORY: {story_id}    │
  │ ┌────────────────────────────────────────────────────────────────────┐  │
  │ │ Todo ({todo_count})  │ In Progress ({wip})  │ Done ({done_count})  │  │
  │ │ {todo_tasks}          │ {wip_tasks}          │ {done_tasks}         │  │
  │ └────────────────────────────────────────────────────────────────────┘  │
  │ Errors: {errors_count}  |  At Risk: {at_risk_count}                  │
  ├──────────────────────────────────────────────────────────────────────────┤
  │ RECENT ACTIVITY                                                          │
  │ {timestamp_1} │ {event_1}                                               │
  │ {timestamp_2} │ {event_2}                                               │
  │ {timestamp_3} │ {event_3}                                               │
  ╰──────────────────────────────────────────────────────────────────────────╯
```

The dashboard must display this data in organized, human readable tables charts and lists that update in real time when the json changes.

---

## Architecture

The Spec⭐️ system is comprised of two main components:

- The TUI program that the user directly interacts with
- The hook executibles that are triggered during claude lifecycle

### TUI Development Stack: Go

- [Cobra](https://github.com/spf13/cobra)
- [Bubbletea](https://github.com/charmbracelet/bubbletea)
- [Lipgloss](https://github.com/charmbracelet/lipgloss)

### Spcstr Config

Project config is located at `.spcstr/settings.json`.
Global user config is located at `~/.spcstr/settings.json`.

**Settings**
- Key-bind overrides
- Theme (light|dark|custom)
- Auto-initialize

---
**Footnotes**

[^1]: Claude-code hooks: https://docs.anthropic.com/en/docs/claude-code/hooks
