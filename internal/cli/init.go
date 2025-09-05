package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dylan-gluck/spcstr/internal/config"
	"github.com/dylan-gluck/spcstr/internal/hooks"
	"github.com/spf13/cobra"
)

var (
	forceInit bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize spcstr for current project",
	Long: `Initialize spcstr for the current project by setting up the .spcstr directory,
creating configuration files, generating hook scripts, and updating Claude Code settings.

This command will:
  - Create the .spcstr directory structure
  - Generate default configuration in .spcstr/config.json
  - Create hook scripts in .spcstr/hooks/
  - Update Claude Code's settings.json with hook configurations

Use --force to overwrite existing configuration.`,
	RunE: runInit,
}

func init() {
	initCmd.Flags().BoolVarP(&forceInit, "force", "f", false, "Force initialization, overwriting existing configuration")
}

func runInit(cmd *cobra.Command, args []string) error {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	projectPath := filepath.Join(cwd, ".spcstr")

	// Check for existing configuration
	if !forceInit {
		if _, err := os.Stat(projectPath); err == nil {
			fmt.Fprintln(cmd.OutOrStdout(), "Error: .spcstr directory already exists in this project.")
			fmt.Fprintln(cmd.OutOrStdout(), "Use --force to overwrite existing configuration.")
			return fmt.Errorf("configuration already exists")
		}
	}

	// Create .spcstr directory structure
	fmt.Fprintln(cmd.OutOrStdout(), "Creating .spcstr directory structure...")
	if err := createDirectoryStructure(projectPath); err != nil {
		return fmt.Errorf("creating directory structure: %w", err)
	}

	// Create default configuration
	fmt.Fprintln(cmd.OutOrStdout(), "Generating configuration...")
	cfg := config.DefaultConfig()
	cfg.Scope = "project"

	configPath := filepath.Join(projectPath, "config.json")
	if err := config.SaveConfig(configPath, cfg); err != nil {
		return fmt.Errorf("saving configuration: %w", err)
	}

	// Generate hook scripts
	fmt.Fprintln(cmd.OutOrStdout(), "Generating hook scripts...")
	hooksPath := filepath.Join(projectPath, "hooks")
	generator := hooks.NewGenerator(hooksPath, projectPath)
	if err := generator.GenerateHooks(); err != nil {
		return fmt.Errorf("generating hooks: %w", err)
	}

	// Update Claude settings
	fmt.Fprintln(cmd.OutOrStdout(), "Updating Claude Code settings...")
	updater := hooks.NewClaudeUpdater()
	if err := updater.UpdateClaudeSettings(hooksPath, forceInit); err != nil {
		// Non-fatal error - Claude settings update is optional
		fmt.Fprintf(cmd.OutOrStdout(), "Warning: Could not update Claude settings: %v\n", err)
		fmt.Fprintln(cmd.OutOrStdout(), "You may need to manually configure hooks in Claude Code settings.")
	}

	// Success message
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "✅ Successfully initialized spcstr!")
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "Created:")
	fmt.Fprintf(cmd.OutOrStdout(), "  📁 .spcstr/              - Project configuration directory\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  📄 .spcstr/config.json   - Configuration file\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  🪝 .spcstr/hooks/        - Hook scripts\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  💾 .spcstr/sessions/     - Session data storage\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  📦 .spcstr/cache/        - Application cache\n")
	fmt.Fprintf(cmd.OutOrStdout(), "  📝 .spcstr/logs/         - Log files\n")
	fmt.Fprintln(cmd.OutOrStdout(), "")
	fmt.Fprintln(cmd.OutOrStdout(), "Next steps:")
	fmt.Fprintln(cmd.OutOrStdout(), "  1. Review configuration in .spcstr/config.json")
	fmt.Fprintln(cmd.OutOrStdout(), "  2. Run 'spcstr run' to start the TUI")
	fmt.Fprintln(cmd.OutOrStdout(), "  3. Start a new Claude Code session to activate hooks")

	return nil
}

func createDirectoryStructure(basePath string) error {
	// Define directory structure
	dirs := []string{
		basePath,
		filepath.Join(basePath, "sessions"),
		filepath.Join(basePath, "sessions", "active"),
		filepath.Join(basePath, "sessions", "archive"),
		filepath.Join(basePath, "hooks"),
		filepath.Join(basePath, "cache"),
		filepath.Join(basePath, "logs"),
	}

	// Create directories with appropriate permissions
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	return nil
}
