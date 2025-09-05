package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRunCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantOut string
	}{
		{
			name:    "basic run",
			args:    []string{},
			wantErr: false,
			wantOut: "Launching Spec⭐️ TUI",
		},
		{
			name:    "run with session flag",
			args:    []string{"--session", "sess_12345"},
			wantErr: false,
			wantOut: "Launching Spec⭐️ TUI",
		},
		{
			name:    "run with plan flag",
			args:    []string{"--plan"},
			wantErr: false,
			wantOut: "Launching Spec⭐️ TUI",
		},
		{
			name:    "run with observe flag",
			args:    []string{"--observe"},
			wantErr: false,
			wantOut: "Launching Spec⭐️ TUI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh command for each test
			rootCmd := &cobra.Command{Use: "spcstr"}
			rootCmd.AddCommand(runCmd)
			rootCmd.SetArgs(append([]string{"run"}, tt.args...))

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

func TestRunCommandFlags(t *testing.T) {
	// Test that run command has expected flags
	tests := []struct {
		flagName  string
		shorthand string
		flagType  string
	}{
		{"session", "s", "string"},
		{"plan", "p", "bool"},
		{"observe", "o", "bool"},
	}

	for _, tt := range tests {
		t.Run("has_"+tt.flagName+"_flag", func(t *testing.T) {
			flag := runCmd.Flags().Lookup(tt.flagName)
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