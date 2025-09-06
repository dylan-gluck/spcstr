package integration

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestTUILaunch(t *testing.T) {
	// Skip TUI tests as they require TTY
	t.Skip("Skipping TUI test - requires TTY which is not available in test environment")

	// Test TUI launch without initialization
	t.Run("LaunchWithoutInit", func(t *testing.T) {
		tmpDir := t.TempDir()

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "./test_spcstr")
		cmd.Dir = tmpDir

		// The TUI should start and display init prompt
		// We'll just verify it starts without error
		err := cmd.Start()
		if err != nil {
			t.Fatalf("Failed to start TUI: %v", err)
		}

		// Give it a moment to initialize
		time.Sleep(100 * time.Millisecond)

		// Kill the process (since TUI runs indefinitely)
		cmd.Process.Kill()
	})

	// Test TUI launch with initialization
	t.Run("LaunchWithInit", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create .spcstr directory
		spcstrDir := filepath.Join(tmpDir, ".spcstr")
		if err := os.MkdirAll(filepath.Join(spcstrDir, "sessions"), 0755); err != nil {
			t.Fatalf("Failed to create .spcstr directory: %v", err)
		}
		if err := os.MkdirAll(filepath.Join(spcstrDir, "logs"), 0755); err != nil {
			t.Fatalf("Failed to create logs directory: %v", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		cmd := exec.CommandContext(ctx, "./test_spcstr")
		cmd.Dir = tmpDir

		// The TUI should start successfully
		err := cmd.Start()
		if err != nil {
			t.Fatalf("Failed to start TUI with init: %v", err)
		}

		// Give it a moment to initialize
		time.Sleep(100 * time.Millisecond)

		// Kill the process
		cmd.Process.Kill()
	})
}

func TestTUIInitCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary with absolute path
	binPath := filepath.Join(t.TempDir(), "test_spcstr")
	cmd := exec.Command("go", "build", "-o", binPath, "../../cmd/spcstr")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tmpDir := t.TempDir()

	// Run init command
	cmd = exec.Command(binPath, "init")
	cmd.Dir = tmpDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Init command failed: %v\nOutput: %s", err, output)
	}

	// Verify .spcstr directory was created
	spcstrDir := filepath.Join(tmpDir, ".spcstr")
	if _, err := os.Stat(spcstrDir); os.IsNotExist(err) {
		t.Error(".spcstr directory was not created")
	}

	// Verify subdirectories
	subdirs := []string{"sessions", "logs"}
	for _, subdir := range subdirs {
		path := filepath.Join(spcstrDir, subdir)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("%s directory was not created", subdir)
		}
	}
}

func TestTUIVersionCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Build the binary with absolute path
	binPath := filepath.Join(t.TempDir(), "test_spcstr")
	cmd := exec.Command("go", "build", "-o", binPath, "../../cmd/spcstr")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	// Run version command
	cmd = exec.Command(binPath, "version")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Version command failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)

	// Check for version output
	if !tuiContains(outputStr, "spcstr version") {
		t.Error("Version output should contain 'spcstr version'")
	}

	if !tuiContains(outputStr, "Git commit:") {
		t.Error("Version output should contain 'Git commit:'")
	}

	if !tuiContains(outputStr, "Built:") {
		t.Error("Version output should contain 'Built:'")
	}
}

func tuiContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
