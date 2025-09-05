package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestInitCommand(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		setupFunc  func(string) error
		wantErr    bool
		wantOut    string
		notWantOut string
	}{
		{
			name: "basic init in clean directory",
			args: []string{},
			setupFunc: func(dir string) error {
				// Clean directory, no setup needed
				return nil
			},
			wantErr: false,
			wantOut: "Successfully initialized spcstr",
		},
		{
			name: "init with existing .spcstr",
			args: []string{},
			setupFunc: func(dir string) error {
				// Create .spcstr directory
				return os.MkdirAll(filepath.Join(dir, ".spcstr"), 0755)
			},
			wantErr:    true,
			wantOut:    ".spcstr directory already exists",
			notWantOut: "Successfully initialized",
		},
		{
			name: "init with force flag and existing .spcstr",
			args: []string{"--force"},
			setupFunc: func(dir string) error {
				// Create .spcstr directory
				return os.MkdirAll(filepath.Join(dir, ".spcstr"), 0755)
			},
			wantErr: false,
			wantOut: "Successfully initialized spcstr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory for each test
			tmpDir := t.TempDir()
			
			// Change to temp directory
			oldDir, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}
			defer os.Chdir(oldDir)
			
			if err := os.Chdir(tmpDir); err != nil {
				t.Fatal(err)
			}
			
			// Run setup function if provided
			if tt.setupFunc != nil {
				if err := tt.setupFunc(tmpDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}
			
			// Create a fresh command for each test to avoid state issues
			rootCmd := &cobra.Command{Use: "spcstr"}
			rootCmd.AddCommand(initCmd)
			rootCmd.SetArgs(append([]string{"init"}, tt.args...))

			// Capture output
			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)

			err = rootCmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			output := buf.String()
			if tt.wantOut != "" && !strings.Contains(output, tt.wantOut) {
				t.Errorf("Output does not contain expected text.\nWant: %s\nGot: %s", tt.wantOut, output)
			}
			
			if tt.notWantOut != "" && strings.Contains(output, tt.notWantOut) {
				t.Errorf("Output contains unexpected text.\nDon't want: %s\nGot: %s", tt.notWantOut, output)
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