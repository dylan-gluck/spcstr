# Spec⭐️ (spcstr)

[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Spec⭐️ (spcstr) is a session persistence tool for Claude Code that tracks and visualizes AI agent activities. It manages project context across multiple Claude sessions to maintain continuity and momentum in your development workflow.

## Installation

### From Source

Requirements:
- Go 1.21.0 or later
- Make (optional, for using Makefile targets)

Clone the repository and build:

```bash
git clone https://github.com/dylan-gluck/spcstr.git
cd spcstr
make build
```

Install to your `$GOPATH/bin`:

```bash
make install
```

Or install directly with Go:

```bash
go install github.com/dylan-gluck/spcstr/cmd/spcstr@latest
```

## Development Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/dylan-gluck/spcstr.git
   cd spcstr
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the project:
   ```bash
   make build
   ```

4. Run tests:
   ```bash
   make test
   ```

5. Run linting:
   ```bash
   make lint
   ```

## Usage Examples

### Initialize Spec⭐️ in your project

```bash
spcstr init
```

This sets up Spec⭐️ configuration and hooks for your project.

### Launch the TUI

```bash
spcstr run
```

Opens the terminal user interface where you can:
- View Claude Code session history (Observe view)
- Browse project documentation (Plan view)
- Monitor agent activities in real-time

### View version information

```bash
spcstr version
```

### Manage configuration

```bash
# List all configuration values
spcstr config list

# Get a specific configuration value
spcstr config get session.path

# Set a configuration value
spcstr config set ui.theme dark
```

### Command-line flags

Global flags available for all commands:
- `-c, --config`: Specify custom config file path
- `-v, --verbose`: Enable verbose output
- `-h, --help`: Show help for any command

## Project Structure

```
spcstr/
├── cmd/spcstr/     # Application entry point
├── internal/       # Internal packages
│   ├── cli/        # Command-line interface
│   ├── config/     # Configuration management
│   ├── hooks/      # Claude Code hook integration
│   ├── session/    # Session persistence
│   ├── tui/        # Terminal UI components
│   └── views/      # TUI views (Plan & Observe)
├── docs/           # Documentation
└── test/           # Integration tests
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
