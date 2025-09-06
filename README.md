# spcstr

Multi-agent observability for Claude Code sessions. Track agent activities, file operations, and task progress in real-time through an intuitive TUI.

![Planning Dashboard](docs/screenshot.png)

**Version:** 0.0.1 (POC)

## Features

- **Automatic Session Tracking** - Transparent hook-based monitoring of Claude Code sessions
- **Real-time Dashboard** - Live view of agents, tasks, files, and tool usage
- **Plan Browser** - Navigate PRDs, architecture docs, and workflows with rich markdown rendering
- **Zero Configuration** - Single `spcstr init` command sets up everything
- **Privacy First** - All data stored locally, no network calls

## Installation

### From Source
```bash
git clone https://github.com/dylan-gluck/spcstr.git
cd spcstr
make build
# Binary available at bin/spcstr
# Add to $PATH for hooks to function
```

## Quick Start

1. Initialize spcstr in your project:
```bash
spcstr init
```

2. Launch the TUI:
```bash
spcstr
```

3. Navigate with keyboard shortcuts:
- `p` - Plan view (browse documentation)
- `o` - Observe view (session dashboard)
- `q` - Quit

## How It Works

Spec⭐️ integrates with Claude Code through a hook system that captures session events:

1. **Initialization** - `spcstr init` configures your project with hook executables
2. **State Tracking** - Hooks capture prompts, tool usage, file operations, and agent activities
3. **Persistence** - Session data stored in `.spcstr/sessions/{session-id}/state.json`
4. **Visualization** - TUI provides real-time and historical session analysis

## Project Structure

```
.spcstr/
├── sessions/          # Session state files
│   └── {session-id}/
│       └── state.json
├── logs/             # Hook execution logs
└── hooks/            # Hook configuration
└── settings.json     # Configuration settings
```

## Development Setup

### Prerequisites
- Go 1.21 or higher
- Make

### Building
```bash
make build        # Build binary to bin/
make test         # Run tests
make clean        # Clean build artifacts
```

### Project Layout
```
cmd/spcstr/       # Main application entry point
internal/
├── config/       # Configuration management
├── hooks/        # Hook command implementations
├── state/        # State management and persistence
└── tui/          # Terminal UI components
```

## Hook System

Spec⭐️ tracks Claude Code sessions through 9 hooks:

- `session_start` - Initialize new session
- `user_prompt_submit` - Capture user prompts
- `pre_tool_use` - Track tool invocations
- `post_tool_use` - Record tool completions
- `notification` - Log notifications
- `pre_compact` - Monitor context compaction
- `session_end` - Mark session complete
- `stop` - Handle termination
- `subagent_stop` - Track sub-agent lifecycle

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Standards
- Follow Go conventions and idioms
- Add tests for new functionality
- Update documentation as needed
- Keep commits atomic and well-described

## License

MIT License - see [LICENSE](LICENSE) file for details

## Support

- Documentation: [docs/](docs/)
- Report issues: [GitHub Issues](https://github.com/dylan/spcstr/issues)
- Community: [Discussions](https://github.com/dylan/spcstr/discussions)
