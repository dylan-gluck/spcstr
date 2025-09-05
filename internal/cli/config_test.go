package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestConfigCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantOut string
	}{
		{
			name:    "basic config",
			args:    []string{},
			wantErr: false,
			wantOut: "Configuration management",
		},
		{
			name:    "config get",
			args:    []string{"get", "session.path"},
			wantErr: false,
			wantOut: "Getting config value for: session.path",
		},
		{
			name:    "config set",
			args:    []string{"set", "ui.theme", "dark"},
			wantErr: false,
			wantOut: "Setting config: ui.theme = dark",
		},
		{
			name:    "config list",
			args:    []string{"list"},
			wantErr: false,
			wantOut: "Listing all configuration values",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh command for each test
			rootCmd := &cobra.Command{Use: "spcstr"}
			rootCmd.AddCommand(configCmd)
			rootCmd.SetArgs(append([]string{"config"}, tt.args...))

			// Capture output
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)

			err := rootCmd.Execute()
			output := buf.String()

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check output
			if !strings.Contains(output, tt.wantOut) {
				t.Errorf("Output doesn't contain expected string.\nWant substring: %s\nGot output: %s", tt.wantOut, output)
			}
		})
	}
}

func TestConfigSubcommands(t *testing.T) {
	// Test that config subcommands exist
	expectedSubcommands := []string{"get", "set", "list"}

	for _, cmdName := range expectedSubcommands {
		t.Run("has_"+cmdName+"_subcommand", func(t *testing.T) {
			found := false
			for _, cmd := range configCmd.Commands() {
				if cmd.Use == cmdName || strings.HasPrefix(cmd.Use, cmdName+" ") {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Subcommand %s not found in config command", cmdName)
			}
		})
	}
}
