package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dylan-gluck/spcstr/internal/hooks"
)

func TestHookSystemEndToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_PATH", tmpDir)
	defer os.Unsetenv("CLAUDE_PROJECT_PATH")

	tests := []struct {
		name     string
		hook     string
		input    map[string]interface{}
		validate func(t *testing.T, tmpDir string)
	}{
		{
			name: "session_start creates state",
			hook: "session_start",
			input: map[string]interface{}{
				"session_id":   "sess_test_001",
				"source":       "startup",
				"project_path": tmpDir,
			},
			validate: func(t *testing.T, tmpDir string) {
				state, err := hooks.LoadSessionState(tmpDir)
				if err != nil {
					t.Errorf("Failed to load state after session_start: %v", err)
					return
				}
				if state.SessionID != "sess_test_001" {
					t.Errorf("SessionID mismatch: got %s, want sess_test_001", state.SessionID)
				}
				if state.Status != "active" {
					t.Errorf("Status should be active, got %s", state.Status)
				}
			},
		},
		{
			name: "post_tool_use tracks tools",
			hook: "post_tool_use",
			input: map[string]interface{}{
				"tool_name": "Read",
				"parameters": map[string]interface{}{
					"file_path": "/test/file.txt",
				},
			},
			validate: func(t *testing.T, tmpDir string) {
				state, err := hooks.LoadSessionState(tmpDir)
				if err != nil {
					t.Errorf("Failed to load state: %v", err)
					return
				}
				if state.ToolsUsed["Read"] != 1 {
					t.Errorf("Read tool count mismatch: got %d, want 1", state.ToolsUsed["Read"])
				}
				if len(state.Files.Read) != 1 {
					t.Errorf("Files.Read should have 1 entry, got %d", len(state.Files.Read))
				}
			},
		},
		{
			name: "pre_tool_use blocks dangerous commands",
			hook: "pre_tool_use",
			input: map[string]interface{}{
				"tool_name": "Bash",
				"parameters": map[string]interface{}{
					"command": "rm -rf /",
				},
			},
			validate: func(t *testing.T, tmpDir string) {
				// This should be blocked, validate by checking handler response
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := hooks.NewHandler(tmpDir)
			err := handler.HandleHook(tc.hook, tc.input)

			if tc.hook == "pre_tool_use" && strings.Contains(tc.input["parameters"].(map[string]interface{})["command"].(string), "rm -rf") {
				if err == nil {
					t.Error("Dangerous command should have been blocked")
				}
			} else if err != nil {
				t.Errorf("Hook %s failed: %v", tc.hook, err)
			}

			if tc.validate != nil {
				tc.validate(t, tmpDir)
			}
		})
	}
}

func TestHookCommandLineInterface(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_PATH", tmpDir)
	defer os.Unsetenv("CLAUDE_PROJECT_PATH")

	testData := map[string]interface{}{
		"session_id": "sess_cli_test",
		"source":     "startup",
	}

	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("Failed to marshal test data: %v", err)
	}

	cmd := exec.Command("spcstr", "hook", "session_start")
	cmd.Stdin = strings.NewReader(string(jsonData))
	cmd.Env = append(os.Environ(), "CLAUDE_PROJECT_PATH="+tmpDir)

	output, err := cmd.CombinedOutput()
	if err != nil {
		// Command might not be in PATH during tests, skip if not found
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 127 {
			t.Skip("spcstr command not found in PATH")
		}
		t.Errorf("Command failed: %v\nOutput: %s", err, output)
	}

	// Verify session was created
	sessionFile := filepath.Join(tmpDir, ".spcstr", "session_state.json")
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		t.Error("Session file was not created by CLI command")
	}
}

func TestHookPerformanceRequirements(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	handler := hooks.NewHandler(tmpDir)

	// Initialize session
	handler.HandleHook("session_start", map[string]interface{}{
		"session_id": "sess_perf_test",
	})

	// Simulate realistic hook sequence
	hooks := []struct {
		name string
		data map[string]interface{}
	}{
		{
			"pre_tool_use",
			map[string]interface{}{
				"tool_name": "Read",
				"parameters": map[string]interface{}{
					"file_path": "/test/file.txt",
				},
			},
		},
		{
			"post_tool_use",
			map[string]interface{}{
				"tool_name": "Read",
				"parameters": map[string]interface{}{
					"file_path": "/test/file.txt",
				},
			},
		},
		{
			"subagent_start",
			map[string]interface{}{
				"agent_name": "test_agent",
			},
		},
		{
			"subagent_stop",
			map[string]interface{}{
				"agent_name": "test_agent",
			},
		},
	}

	for _, h := range hooks {
		start := time.Now()
		err := handler.HandleHook(h.name, h.data)
		duration := time.Since(start)

		if err != nil && !strings.Contains(err.Error(), "blocked") {
			t.Errorf("Hook %s failed: %v", h.name, err)
		}

		if duration > 10*time.Millisecond {
			t.Errorf("Hook %s took %v, exceeds 10ms requirement", h.name, duration)
		}
	}
}

func TestSessionLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	handler := hooks.NewHandler(tmpDir)

	// Start session
	err := handler.HandleHook("session_start", map[string]interface{}{
		"session_id":   "sess_lifecycle",
		"source":       "startup",
		"project_path": tmpDir,
	})
	if err != nil {
		t.Fatalf("session_start failed: %v", err)
	}

	// Perform some operations
	operations := []struct {
		hook string
		data map[string]interface{}
	}{
		{
			"post_tool_use",
			map[string]interface{}{
				"tool_name": "Write",
				"parameters": map[string]interface{}{
					"file_path": "/new/file.txt",
				},
			},
		},
		{
			"post_tool_use",
			map[string]interface{}{
				"tool_name": "Edit",
				"parameters": map[string]interface{}{
					"file_path": "/existing/file.txt",
				},
			},
		},
		{
			"notification",
			map[string]interface{}{
				"type":    "error",
				"message": "Test error message",
			},
		},
	}

	for _, op := range operations {
		if err := handler.HandleHook(op.hook, op.data); err != nil {
			t.Errorf("%s failed: %v", op.hook, err)
		}
	}

	// End session
	err = handler.HandleHook("session_end", map[string]interface{}{
		"session_id": "sess_lifecycle",
	})
	if err != nil {
		t.Fatalf("session_end failed: %v", err)
	}

	// Verify final state
	state, err := hooks.LoadSessionState(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load final state: %v", err)
	}

	if state.Status != "completed" {
		t.Errorf("Final status should be 'completed', got %s", state.Status)
	}

	if state.ToolsUsed["Write"] != 1 || state.ToolsUsed["Edit"] != 1 {
		t.Error("Tool counts not tracked correctly")
	}

	if len(state.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(state.Errors))
	}

	if len(state.Files.New) != 1 || len(state.Files.Edited) != 1 {
		t.Error("File operations not tracked correctly")
	}
}
