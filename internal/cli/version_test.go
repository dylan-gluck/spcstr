package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestVersionCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		wantOut []string
	}{
		{
			name:    "basic version",
			args:    []string{},
			wantErr: false,
			wantOut: []string{"Spec⭐️ version 0.1.0", "Build: development", "Go version: 1.21.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh command for each test
			rootCmd := &cobra.Command{Use: "spcstr"}
			rootCmd.AddCommand(versionCmd)
			rootCmd.SetArgs(append([]string{"version"}, tt.args...))

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

			// Check output contains all expected strings
			for _, want := range tt.wantOut {
				if !strings.Contains(output, want) {
					t.Errorf("Output doesn't contain expected string.\nWant substring: %s\nGot output: %s", want, output)
				}
			}
		})
	}
}