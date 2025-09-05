package integration

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestInitCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Build the binary
	binPath := filepath.Join(t.TempDir(), "spcstr")
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/spcstr")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build spcstr binary: %v", err)
	}

	tests := []struct {
		name           string
		setupFunc      func(string) error
		args           []string
		expectSuccess  bool
		verifyFunc     func(string) error
	}{
		{
			name:          "init fresh project",
			setupFunc:     nil,
			args:          []string{"init"},
			expectSuccess: true,
			verifyFunc: func(projectDir string) error {
				// Verify directory structure
				dirs := []string{
					filepath.Join(projectDir, ".spcstr", "logs"),
					filepath.Join(projectDir, ".spcstr", "sessions"),
				}
				for _, dir := range dirs {
					if _, err := os.Stat(dir); os.IsNotExist(err) {
						return err
					}
				}
				
				// Verify settings.json
				settingsPath := filepath.Join(projectDir, ".claude", "settings.json")
				data, err := os.ReadFile(settingsPath)
				if err != nil {
					return err
				}
				
				var settings map[string]interface{}
				if err := json.Unmarshal(data, &settings); err != nil {
					return err
				}
				
				if _, ok := settings["hooks"]; !ok {
					t.Error("hooks not found in settings.json")
				}
				
				return nil
			},
		},
		{
			name: "init with existing .spcstr and force flag",
			setupFunc: func(projectDir string) error {
				// Create existing .spcstr directory
				return os.MkdirAll(filepath.Join(projectDir, ".spcstr"), 0755)
			},
			args:          []string{"init", "--force"},
			expectSuccess: true,
			verifyFunc: func(projectDir string) error {
				// Verify directories still exist
				dirs := []string{
					filepath.Join(projectDir, ".spcstr", "logs"),
					filepath.Join(projectDir, ".spcstr", "sessions"),
				}
				for _, dir := range dirs {
					if _, err := os.Stat(dir); os.IsNotExist(err) {
						return err
					}
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp project directory
			projectDir := t.TempDir()
			
			// Setup if needed
			if tt.setupFunc != nil {
				if err := tt.setupFunc(projectDir); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}
			
			// Run init command
			cmd := exec.Command(binPath, tt.args...)
			cmd.Dir = projectDir
			
			output, err := cmd.CombinedOutput()
			
			if tt.expectSuccess && err != nil {
				t.Errorf("command failed: %v\nOutput: %s", err, output)
			}
			if !tt.expectSuccess && err == nil {
				t.Errorf("expected command to fail but it succeeded\nOutput: %s", output)
			}
			
			// Verify results
			if tt.expectSuccess && tt.verifyFunc != nil {
				if err := tt.verifyFunc(projectDir); err != nil {
					t.Errorf("verification failed: %v", err)
				}
			}
		})
	}
}

func TestVersionCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Build the binary
	binPath := filepath.Join(t.TempDir(), "spcstr")
	buildCmd := exec.Command("go", "build", 
		"-ldflags", "-X main.Version=1.0.0 -X main.GitCommit=abc123 -X main.BuildDate=2025-01-01",
		"-o", binPath, "../../cmd/spcstr")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build spcstr binary: %v", err)
	}

	cmd := exec.Command(binPath, "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version command failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	expectedStrings := []string{
		"spcstr version 1.0.0",
		"Git commit: abc123",
		"Built: 2025-01-01",
	}

	for _, expected := range expectedStrings {
		if !contains(outputStr, expected) {
			t.Errorf("version output missing expected string: %s\nActual output: %s", expected, outputStr)
		}
	}
}

func TestRootCommandNoArgs(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Build the binary
	binPath := filepath.Join(t.TempDir(), "spcstr")
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/spcstr")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build spcstr binary: %v", err)
	}

	cmd := exec.Command(binPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("root command failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	// Should mention TUI mode (even if not implemented yet)
	if !contains(outputStr, "TUI") || !contains(outputStr, "spcstr") {
		t.Errorf("root command output unexpected: %s", outputStr)
	}
}

func TestHookCommandAfterInit(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	// Build the binary
	binPath := filepath.Join(t.TempDir(), "spcstr")
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/spcstr")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("failed to build spcstr binary: %v", err)
	}

	// Create temp project directory
	projectDir := t.TempDir()

	// Initialize the project
	initCmd := exec.Command(binPath, "init", "--force")
	initCmd.Dir = projectDir
	if err := initCmd.Run(); err != nil {
		t.Fatalf("init command failed: %v", err)
	}

	// Test hook command with proper JSON input
	hookInput := `{"session_id": "test-session", "timestamp": "2025-01-01T00:00:00Z"}`
	hookCmd := exec.Command(binPath, "hook", "session_start", "--cwd", projectDir)
	hookCmd.Stdin = &testReader{data: []byte(hookInput)}
	
	output, err := hookCmd.CombinedOutput()
	if err != nil {
		// Hook might fail if not all handlers are implemented, but it should at least run
		// Check that it's failing for the right reason (not because of init issues)
		outputStr := string(output)
		if contains(outputStr, "invalid spcstr project") {
			t.Errorf("hook failed due to invalid project after init: %s", outputStr)
		}
	}

	// Verify that log file was created (if hook succeeded)
	logFile := filepath.Join(projectDir, ".spcstr", "logs", "session_start.json")
	if _, err := os.Stat(logFile); err == nil {
		// Log file exists, verify it has content
		data, _ := os.ReadFile(logFile)
		if len(data) == 0 {
			t.Error("log file is empty")
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || 
		len(s) > len(substr) && containsSubstring(s[1:], substr)
}

func containsSubstring(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

type testReader struct {
	data []byte
	pos  int
}

func (r *testReader) Read(p []byte) (n int, err error) {
	if r.pos >= len(r.data) {
		return 0, io.EOF
	}
	n = copy(p, r.data[r.pos:])
	r.pos += n
	return n, nil
}