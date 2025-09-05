# Requirements

## Functional

1. FR1: The init command shall create all necessary directories, hook scripts, and configuration files with a single command
2. FR2: The system shall automatically track all Claude Code sessions when hooks are triggered
3. FR3: Session data shall be persisted to JSON files within 1 second of any state change
4. FR4: The TUI shall support keyboard-driven navigation with discoverable keybindings
5. FR5: Plan View shall recursively index and display all markdown planning documents
6. FR6: Plan View shall provide real-time markdown preview with syntax highlighting
7. FR7: Observe View shall display active and completed sessions in a selectable list
8. FR8: Session dashboard shall show agents, tasks, files, tools, and errors in organized sections
9. FR9: The system shall maintain session state across application restarts
10. FR10: Global keybindings shall allow view switching from any context
11. FR11: The TUI shall support terminal resize events without data loss
12. FR12: Session selection shall load and display full session details within 100ms
13. FR13: The system shall track file operations (new, edited, read) with full paths
14. FR14: Tool usage counters shall increment in real-time as hooks fire
15. FR15: Error tracking shall capture and display session errors with timestamps

## Non Functional

1. NFR1: UI response time shall not exceed 100ms for any user interaction
2. NFR2: Memory usage shall not exceed 10MB per tracked session
3. NFR3: CPU usage shall remain below 1% during idle periods
4. NFR4: The application shall compile to a single binary with no runtime dependencies
5. NFR5: Installation shall be achievable via standard package managers (brew, apt, yum)
6. NFR6: The system shall support 256-color terminal emulators
7. NFR7: Session data files shall use standard JSON format for interoperability
8. NFR8: The TUI shall maintain 60fps refresh rate during animations
9. NFR9: File system operations shall use platform-native path separators
10. NFR10: The application shall gracefully handle missing or corrupted session files
