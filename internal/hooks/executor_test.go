package hooks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestExecuteHook(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	spcstrDir := filepath.Join(tempDir, ".spcstr")
	sessionsDir := filepath.Join(spcstrDir, "sessions")
	logsDir := filepath.Join(spcstrDir, "logs")
	
	err := os.MkdirAll(sessionsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create sessions directory: %v", err)
	}
	
	err = os.MkdirAll(logsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}
	
	// Register a test handler
	testRegistry := NewRegistry()
	handler := &mockHandler{name: "test_hook", execErr: nil}
	testRegistry.Register(handler)
	
	// Temporarily replace default registry
	oldRegistry := DefaultRegistry
	DefaultRegistry = testRegistry
	defer func() {
		DefaultRegistry = oldRegistry
	}()
	
	tests := []struct {
		name        string
		projectDir  string
		hookName    string
		input       []byte
		expectError bool
	}{
		{
			name:        "valid project and hook",
			projectDir:  tempDir,
			hookName:    "test_hook",
			input:       []byte(`{"session_id": "test123"}`),
			expectError: false,
		},
		{
			name:        "invalid project directory",
			projectDir:  "/nonexistent/path",
			hookName:    "test_hook",
			input:       []byte(`{"session_id": "test123"}`),
			expectError: true,
		},
		{
			name:        "nonexistent hook",
			projectDir:  tempDir,
			hookName:    "nonexistent_hook",
			input:       []byte(`{"session_id": "test123"}`),
			expectError: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ExecuteHook(tt.hookName, tt.projectDir, tt.input)
			
			if tt.expectError {
				if err == nil {
					t.Error("Expected error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

func TestIsValidSpcstrProject(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() string
		expected bool
	}{
		{
			name: "valid project structure",
			setup: func() string {
				tempDir := t.TempDir()
				os.MkdirAll(filepath.Join(tempDir, ".spcstr", "sessions"), 0755)
				os.MkdirAll(filepath.Join(tempDir, ".spcstr", "logs"), 0755)
				return tempDir
			},
			expected: true,
		},
		{
			name: "missing sessions directory",
			setup: func() string {
				tempDir := t.TempDir()
				os.MkdirAll(filepath.Join(tempDir, ".spcstr", "logs"), 0755)
				return tempDir
			},
			expected: false,
		},
		{
			name: "missing logs directory",
			setup: func() string {
				tempDir := t.TempDir()
				os.MkdirAll(filepath.Join(tempDir, ".spcstr", "sessions"), 0755)
				return tempDir
			},
			expected: false,
		},
		{
			name: "missing .spcstr directory",
			setup: func() string {
				return t.TempDir()
			},
			expected: false,
		},
		{
			name: "nonexistent directory",
			setup: func() string {
				return "/nonexistent/path"
			},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectDir := tt.setup()
			result := isValidSpcstrProject(projectDir)
			
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for directory: %s", tt.expected, result, projectDir)
			}
		})
	}
}

func TestDirExists(t *testing.T) {
	// Test existing directory
	tempDir := t.TempDir()
	if !dirExists(tempDir) {
		t.Error("dirExists returned false for existing directory")
	}
	
	// Test non-existent directory
	if dirExists("/nonexistent/path") {
		t.Error("dirExists returned true for non-existent directory")
	}
	
	// Test file (not directory)
	tempFile := filepath.Join(tempDir, "testfile")
	os.WriteFile(tempFile, []byte("test"), 0644)
	if dirExists(tempFile) {
		t.Error("dirExists returned true for file")
	}
}

func TestExecuteHookWithLogging(t *testing.T) {
	// Create temporary directory structure
	tempDir := t.TempDir()
	spcstrDir := filepath.Join(tempDir, ".spcstr")
	sessionsDir := filepath.Join(spcstrDir, "sessions")
	logsDir := filepath.Join(spcstrDir, "logs")
	
	err := os.MkdirAll(sessionsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create sessions directory: %v", err)
	}
	
	err = os.MkdirAll(logsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create logs directory: %v", err)
	}
	
	// Change to temp directory for test
	oldDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)
	
	// Register a test handler
	testRegistry := NewRegistry()
	handler := &mockHandler{name: "logging_test_hook", execErr: nil}
	testRegistry.Register(handler)
	
	// Temporarily replace default registry
	oldRegistry := DefaultRegistry
	DefaultRegistry = testRegistry
	defer func() {
		DefaultRegistry = oldRegistry
	}()
	
	// Execute hook
	input := []byte(`{"session_id": "log_test_session"}`)
	err = ExecuteHook("logging_test_hook", tempDir, input)
	if err != nil {
		t.Fatalf("Hook execution failed: %v", err)
	}
	
	// Check if log file was created
	logPath := filepath.Join(logsDir, "logging_test_hook.json")
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Error("Log file was not created")
		return
	}
	
	// Check log file contents
	logData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	
	var events []HookEvent
	err = json.Unmarshal(logData, &events)
	if err != nil {
		t.Fatalf("Failed to parse log file: %v", err)
	}
	
	if len(events) != 1 {
		t.Errorf("Expected 1 log event, got %d", len(events))
		return
	}
	
	event := events[0]
	if event.SessionID != "log_test_session" {
		t.Errorf("Expected session ID 'log_test_session', got '%s'", event.SessionID)
	}
	if event.HookName != "logging_test_hook" {
		t.Errorf("Expected hook name 'logging_test_hook', got '%s'", event.HookName)
	}
	if !event.Success {
		t.Error("Expected success to be true")
	}
}