package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/dylan/spcstr/internal/hooks"
)

const version = "1.0.0"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "spcstr",
	Short:   "spcstr - a CLI/TUI tool for Claude Code session observability",
	Long:    `spcstr provides real-time observability for Claude Code sessions through hook integration and an interactive terminal interface.`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("spcstr v%s\n", version)
		fmt.Println("Use 'spcstr --help' for more information.")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of spcstr",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("spcstr v%s\n", version)
	},
}

var hookCmd = &cobra.Command{
	Use:   "hook [hook_name]",
	Short: "Execute a Claude Code hook command",
	Long:  `Execute a Claude Code hook command with JSON input from stdin`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		hookName := args[0]
		
		// Get working directory from flag
		cwdFlag, _ := cmd.Flags().GetString("cwd")
		workingDir := cwdFlag
		if workingDir == "" {
			var err error
			workingDir, err = os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}
		}
		
		// Convert to absolute path
		absPath, err := filepath.Abs(workingDir)
		if err != nil {
			return fmt.Errorf("failed to resolve absolute path: %w", err)
		}
		
		// Read JSON input from stdin
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		
		// Execute the hook
		err = hooks.ExecuteHook(hookName, absPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Hook execution failed: %v\n", err)
			os.Exit(2) // Block operation exit code
		}
		
		return nil
	},
}

func init() {
	hookCmd.Flags().StringP("cwd", "c", "", "Working directory for hook execution (project root)")
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(hookCmd)
}
