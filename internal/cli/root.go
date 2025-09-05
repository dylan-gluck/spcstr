package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "spcstr",
	Short: "Session persistence for Claude Code",
	Long: `Spec⭐️ (spec-star) is a session persistence tool for Claude Code.
It tracks and visualizes AI agent activities, managing project context
across multiple Claude sessions to maintain continuity and momentum.`,
	Version: "0.1.0",
	RunE: func(cmd *cobra.Command, args []string) error {
		// If no subcommand is provided, show help
		return cmd.Help()
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add commands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)

	// Global flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file path (default: .spcstr/config.json)")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}

// Version returns the current version
func Version() string {
	return rootCmd.Version
}
