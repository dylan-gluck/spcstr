# User Interface Design Goals

## Overall UX Vision
Minimalist, keyboard-driven TUI providing instant visibility into Claude Code sessions with zero mouse interaction. Focus on information density and real-time updates with beautiful markdown rendering via Glamour.

## Key Interaction Paradigms
- Pure keyboard navigation with single-key commands
- Two-pane layouts (list/detail pattern) for both Plan and Observe views
- Context-aware footer showing available keybinds that update per view/mode
- Immediate feedback on all actions (no loading states for local operations)

## Core Screens and Views

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

## Accessibility: None
MVP focuses on core functionality. Accessibility features deferred to post-MVP.

## Branding
Terminal-native aesthetic using Lipgloss styling with Glamour for rich markdown rendering. Light/dark theme support based on terminal colors. No custom branding or logos.

## Target Device and Platforms: Terminal/CLI Only
- macOS Terminal/iTerm2
- Linux terminal emulators
- Windows Terminal with WSL
- Minimum 80x24 terminal size, optimized for 120x40
