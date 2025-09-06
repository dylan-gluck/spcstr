
❯ bin/spcstr --help
Spec⭐️ (spec-star) is a session persistence tool for Claude Code.
It tracks and visualizes AI agent activities, managing project context
across multiple Claude sessions to maintain continuity and momentum.

Usage:
  spcstr [flags]
  spcstr [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      Manage spcstr configuration
  help        Help about any command
  init        Initialize spcstr for current project
  run         Run the spcstr TUI
  version     Print version information

Flags:
  -c, --config string   config file path (default: .spcstr/config.json)
  -h, --help            help for spcstr
  -v, --verbose         verbose output
      --version         version for spcstr

Use "spcstr [command] --help" for more information about a command.


---


❯ bin/spcstr init --force
Creating .spcstr directory structure...
Generating configuration...
Configuring Claude Code hooks...

✅ Successfully initialized spcstr!

Created:
  📁 .spcstr/              - Project configuration directory
  📄 .spcstr/config.json   - Configuration file
  🪝 Claude hooks          - Integrated Go hook handlers
  💾 .spcstr/sessions/     - Session data storage
  📦 .spcstr/cache/        - Application cache
  📝 .spcstr/logs/         - Log files

Next steps:
  1. Review configuration in .spcstr/config.json
  2. Run 'spcstr run' to start the TUI
  3. Start a new Claude Code session to activate hooks
