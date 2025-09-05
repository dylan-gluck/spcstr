package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestExecute(t *testing.T) {
	// Save original rootCmd and restore after test
	originalCmd := rootCmd
	defer func() { rootCmd = originalCmd }()

	// Execute should not return error for help
	rootCmd = &cobra.Command{
		Use:   "test",
		Short: "test command",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	err := Execute()
	if err != nil {
		t.Errorf("Execute() returned error: %v", err)
	}
}

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantOut string
	}{
		{
			name:    "no args shows help",
			args:    []string{},
			wantErr: false,
			wantOut: "Session persistence for Claude Code",
		},
		{
			name:    "help flag",
			args:    []string{"--help"},
			wantErr: false,
			wantOut: "Session persistence for Claude Code",
		},
		{
			name:    "version flag",
			args:    []string{"--version"},
			wantErr: false,
			wantOut: "spcstr version 0.1.0",
		},
		{
			name:    "invalid flag",
			args:    []string{"--invalid"},
			wantErr: true,
			wantOut: "unknown flag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for each test
			cmd := &cobra.Command{
				Use:     "spcstr",
				Short:   "Session persistence for Claude Code",
				Version: "0.1.0",
				RunE: func(cmd *cobra.Command, args []string) error {
					return cmd.Help()
				},
			}

			// Add commands
			cmd.AddCommand(initCmd)
			cmd.AddCommand(runCmd)
			cmd.AddCommand(configCmd)
			cmd.AddCommand(versionCmd)

			// Add flags
			cmd.PersistentFlags().StringP("config", "c", "", "config file path")
			cmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			output := buf.String()

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check output contains expected string
			if !strings.Contains(output, tt.wantOut) {
				t.Errorf("Output doesn't contain expected string.\nWant substring: %s\nGot output: %s", tt.wantOut, output)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	expected := "0.1.0"
	got := Version()
	if got != expected {
		t.Errorf("Version() = %v, want %v", got, expected)
	}
}

func TestSubcommands(t *testing.T) {
	// Test that all expected subcommands are registered
	expectedCommands := []string{"init", "run", "config", "version"}

	for _, cmdName := range expectedCommands {
		t.Run("has_"+cmdName+"_command", func(t *testing.T) {
			found := false
			for _, cmd := range rootCmd.Commands() {
				if cmd.Use == cmdName || strings.HasPrefix(cmd.Use, cmdName+" ") {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Command %s not found in root command", cmdName)
			}
		})
	}
}

func TestGlobalFlags(t *testing.T) {
	// Test that global flags are properly configured
	tests := []struct {
		flagName  string
		shorthand string
		flagType  string
	}{
		{"config", "c", "string"},
		{"verbose", "v", "bool"},
	}

	for _, tt := range tests {
		t.Run("has_"+tt.flagName+"_flag", func(t *testing.T) {
			flag := rootCmd.PersistentFlags().Lookup(tt.flagName)
			if flag == nil {
				t.Errorf("Flag %s not found", tt.flagName)
				return
			}

			// Check shorthand
			if flag.Shorthand != tt.shorthand {
				t.Errorf("Flag %s shorthand = %s, want %s", tt.flagName, flag.Shorthand, tt.shorthand)
			}

			// Check type
			if flag.Value.Type() != tt.flagType {
				t.Errorf("Flag %s type = %s, want %s", tt.flagName, flag.Value.Type(), tt.flagType)
			}
		})
	}
}
