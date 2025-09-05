package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Setup before tests
	code := m.Run()
	// Cleanup after tests
	os.Exit(code)
}

func TestVersionCommand(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "version variable has default",
			want: "dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Version != tt.want {
				t.Errorf("Version = %v, want %v", Version, tt.want)
			}
		})
	}
}

func TestRootCommand(t *testing.T) {
	if rootCmd == nil {
		t.Error("rootCmd should not be nil")
	}

	if rootCmd.Use != "spcstr" {
		t.Errorf("rootCmd.Use = %v, want %v", rootCmd.Use, "spcstr")
	}

	if rootCmd.Version != Version {
		t.Errorf("rootCmd.Version = %v, want %v", rootCmd.Version, Version)
	}
}

func TestVersionCommandExists(t *testing.T) {
	versionCommand := rootCmd.Commands()
	found := false

	for _, cmd := range versionCommand {
		if cmd.Use == "version" {
			found = true
			break
		}
	}

	if !found {
		t.Error("version command should be registered with root command")
	}
}
