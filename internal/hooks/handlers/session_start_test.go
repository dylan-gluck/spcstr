package handlers

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dylan/spcstr/internal/state"
)

func TestSessionStartHandler(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Create .spcstr directory structure
	os.MkdirAll(".spcstr/sessions", 0755)

	handler := NewSessionStartHandler()

	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid session start",
			input:       `{"session_id": "test_session_123", "source": "startup"}`,
			expectError: false,
		},
		{
			name:        "missing session_id",
			input:       `{"source": "startup"}`,
			expectError: true,
			errorMsg:    "session_id is required",
		},
		{
			name:        "missing source",
			input:       `{"session_id": "test_session_123"}`,
			expectError: true,
			errorMsg:    "source is required",
		},
		{
			name:        "invalid JSON",
			input:       `{"session_id": "test_session_123", "source": "startup"`,
			expectError: true,
		},
		{
			name:        "empty session_id",
			input:       `{"session_id": "", "source": "startup"}`,
			expectError: true,
			errorMsg:    "session_id is required",
		},
		{
			name:        "empty source",
			input:       `{"session_id": "test_session_123", "source": ""}`,
			expectError: true,
			errorMsg:    "source is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.Execute([]byte(tt.input))

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, but got nil")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}

				// Verify state file was created
				statePath := ".spcstr/sessions/test_session_123/state.json"
				if _, err := os.Stat(statePath); os.IsNotExist(err) {
					t.Error("State file was not created")
					return
				}

				// Verify state content
				manager := state.NewStateManager(filepath.Join(tempDir, ".spcstr"))
				sessionState, err := manager.LoadState(context.Background(), "test_session_123")
				if err != nil {
					t.Errorf("Failed to load state: %v", err)
					return
				}

				if sessionState.SessionID != "test_session_123" {
					t.Errorf("Expected session ID 'test_session_123', got '%s'", sessionState.SessionID)
				}
				if !sessionState.SessionActive {
					t.Error("Expected session to be active")
				}
				if sessionState.CreatedAt.IsZero() {
					t.Error("Created time should not be zero")
				}
			}
		})
	}
}

func TestSessionStartHandlerName(t *testing.T) {
	handler := NewSessionStartHandler()
	if handler.Name() != "session_start" {
		t.Errorf("Expected handler name 'session_start', got '%s'", handler.Name())
	}
}

func TestNewSessionStartHandler(t *testing.T) {
	handler := NewSessionStartHandler()
	if handler == nil {
		t.Fatal("NewSessionStartHandler() returned nil")
	}
	// Handler doesn't have a stateManager field - it creates one during Execute
	if handler.Name() != "session_start" {
		t.Fatalf("Expected handler name 'session_start', got '%s'", handler.Name())
	}
}
