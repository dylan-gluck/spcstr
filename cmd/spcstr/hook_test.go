package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHook(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_PATH", tmpDir)
	defer os.Unsetenv("CLAUDE_PROJECT_PATH")

	tests := []struct {
		name     string
		hookName string
		input    map[string]interface{}
		wantCode int
	}{
		{
			name:     "session_start",
			hookName: "session_start",
			input: map[string]interface{}{
				"session_id": "test_123",
				"source":     "startup",
			},
			wantCode: 0,
		},
		{
			name:     "blocked command",
			hookName: "pre_tool_use",
			input: map[string]interface{}{
				"tool_name": "Bash",
				"parameters": map[string]interface{}{
					"command": "rm -rf /",
				},
			},
			wantCode: 2,
		},
		{
			name:     "unknown hook",
			hookName: "unknown_hook",
			input:    map[string]interface{}{},
			wantCode: 0,
		},
		{
			name:     "empty input",
			hookName: "notification",
			input:    nil,
			wantCode: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a pipe to simulate stdin
			r, w, _ := os.Pipe()
			oldStdin := os.Stdin
			os.Stdin = r
			defer func() {
				os.Stdin = oldStdin
				r.Close()
			}()

			// Write test data to the pipe
			go func() {
				defer w.Close()
				if tt.input != nil {
					jsonData, _ := json.Marshal(tt.input)
					w.Write(jsonData)
				}
			}()

			code := runHook(tt.hookName)
			if code != tt.wantCode {
				t.Errorf("runHook(%s) = %d, want %d", tt.hookName, code, tt.wantCode)
			}
		})
	}
}

func TestIsBlockingError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"blocked error", errors.New("operation blocked"), true},
		{"dangerous error", errors.New("dangerous operation detected"), true},
		{"prohibited error", errors.New("prohibited action"), true},
		{"normal error", errors.New("file not found"), false},
		{"nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBlockingError(tt.err)
			if got != tt.want {
				t.Errorf("isBlockingError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestHookLogging(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("CLAUDE_PROJECT_PATH", tmpDir)
	defer os.Unsetenv("CLAUDE_PROJECT_PATH")

	// Create a pipe to simulate stdin
	r, w, _ := os.Pipe()
	oldStdin := os.Stdin
	os.Stdin = r
	defer func() {
		os.Stdin = oldStdin
		r.Close()
	}()

	// Write test data to the pipe
	go func() {
		defer w.Close()
		input := map[string]interface{}{
			"session_id": "test_logging",
			"source":     "test",
		}
		jsonData, _ := json.Marshal(input)
		w.Write(jsonData)
	}()

	code := runHook("session_start")
	if code != 0 {
		t.Errorf("Expected exit code 0, got %d", code)
	}

	logFile := filepath.Join(tmpDir, ".spcstr", "logs", "debug.log")
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("Debug log file was not created")
	}

	logContent, err := os.ReadFile(logFile)
	if err != nil {
		t.Errorf("Failed to read log file: %v", err)
	}

	if !strings.Contains(string(logContent), "hook invoked") {
		t.Error("Log file doesn't contain expected 'hook invoked' message")
	}

	if !strings.Contains(string(logContent), "session_start") {
		t.Error("Log file doesn't contain hook name")
	}
}

// Helper to create a reader that implements io.Reader for testing stdin
type stdinReader struct {
	*strings.Reader
}

func (sr *stdinReader) Fd() uintptr {
	return 0
}

func newStdinReader(s string) io.Reader {
	return &stdinReader{strings.NewReader(s)}
}