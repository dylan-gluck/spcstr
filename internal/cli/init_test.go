package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestInitCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantOut string
	}{
		{
			name:    "basic init",
			args:    []string{},
			wantErr: false,
			wantOut: "Initializing Spec⭐️",
		},
		{
			name:    "init with force flag",
			args:    []string{"--force"},
			wantErr: false,
			wantOut: "Initializing Spec⭐️",
		},
		{
			name:    "init with template flag",
			args:    []string{"--template", "default"},
			wantErr: false,
			wantOut: "Initializing Spec⭐️",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh command for each test to avoid state issues
			rootCmd := &cobra.Command{Use: "spcstr"}
			rootCmd.AddCommand(initCmd)
			rootCmd.SetArgs(append([]string{"init"}, tt.args...))

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

func TestInitCommandFlags(t *testing.T) {
	// Test that init command has expected flags
	tests := []struct {
		flagName  string
		shorthand string
		flagType  string
	}{
		{"force", "f", "bool"},
		{"template", "t", "string"},
	}

	for _, tt := range tests {
		t.Run("has_"+tt.flagName+"_flag", func(t *testing.T) {
			flag := initCmd.Flags().Lookup(tt.flagName)
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