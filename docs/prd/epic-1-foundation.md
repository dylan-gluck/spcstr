# Epic 1: Foundation & Core Infrastructure

Establish the foundational Go project structure with Cobra CLI framework and Bubbletea TUI. Implement the init command that sets up hook scripts and configuration. Create the basic TUI shell with global keybindings and view switching capabilities.

## Story 1.1: Project Setup and CLI Structure

As a developer,
I want a well-structured Go project with CLI commands,
so that I can easily extend and maintain the codebase.

### Acceptance Criteria
1. Go module initialized with appropriate package structure
2. Cobra CLI framework integrated with root command
3: Main entry point established with proper error handling
4: Makefile created for common development tasks (build, test, install)
5: Basic README with project description and development setup

## Story 1.2: Init Command Implementation

As a user,
I want to initialize spcstr in my project with a single command,
so that all necessary hooks and configuration are automatically set up.

### Acceptance Criteria
1: `spcstr init` creates .spcstr/ directory structure
2: Hook scripts are written to .spcstr/hooks/ directory
3: Claude settings.json is updated with hook configurations
4: Command detects existing configuration and prompts before overwriting
5: Success message confirms initialization with next steps

## Story 1.3: Basic TUI Shell with View Management

As a user,
I want a responsive TUI with keyboard navigation between views,
so that I can switch between planning and observability modes.

### Acceptance Criteria
1: TUI launches and displays initial view within 500ms
2: Global keybindings (p, o, q) work from any context
3: View transitions are smooth without flicker
4: Terminal resize events are handled gracefully
5: Help text shows available keybindings

## Story 1.4: Configuration Management

As a user,
I want spcstr to respect both project and global configuration,
so that I can customize behavior per project or globally.

### Acceptance Criteria
1: Project config (.spcstr/settings.json) takes precedence over global
2: Global config location follows XDG specification
3: Default configuration is created if none exists
4: Configuration changes take effect without restart
5: Invalid configuration shows clear error messages