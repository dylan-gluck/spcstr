package handlers

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dylan/spcstr/internal/state"
)

func TestUserPromptSubmitHandler(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Create .spcstr directory structure
	os.MkdirAll(".spcstr/sessions", 0755)

	// First create a session state
	manager := state.NewStateManager(filepath.Join(tempDir, ".spcstr"))
	_, err := manager.InitializeState(context.Background(), "test_session_123")
	if err != nil {
		t.Fatalf("Failed to initialize state: %v", err)
	}

	handler := NewUserPromptSubmitHandler()

	tests := []struct {
		name        string
		input       string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid prompt submission",
			input:       `{"session_id": "test_session_123", "prompt": "Hello, Claude!", "timestamp": "2023-01-01T10:00:00Z"}`,
			expectError: false,
		},
		{
			name:        "prompt without timestamp",
			input:       `{"session_id": "test_session_123", "prompt": "Hello, Claude!"}`,
			expectError: false,
		},
		{
			name:        "missing session_id",
			input:       `{"prompt": "Hello, Claude!"}`,
			expectError: true,
			errorMsg:    "session_id is required",
		},
		{
			name:        "missing prompt",
			input:       `{"session_id": "test_session_123"}`,
			expectError: true,
			errorMsg:    "prompt is required",
		},
		{
			name:        "empty session_id",
			input:       `{"session_id": "", "prompt": "Hello, Claude!"}`,
			expectError: true,
			errorMsg:    "session_id is required",
		},
		{
			name:        "empty prompt",
			input:       `{"session_id": "test_session_123", "prompt": ""}`,
			expectError: true,
			errorMsg:    "prompt is required",
		},
		{
			name:        "invalid timestamp format",
			input:       `{"session_id": "test_session_123", "prompt": "Hello, Claude!", "timestamp": "invalid-timestamp"}`,
			expectError: false, // Should use current time instead
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

				// Verify prompt was added to session state
				updatedState, err := manager.LoadState(context.Background(), "test_session_123")
				if err != nil {
					t.Errorf("Failed to load updated state: %v", err)
					return
				}

				if len(updatedState.Prompts) == 0 {
					t.Error("No prompts found in session state")
					return
				}

				lastPrompt := updatedState.Prompts[len(updatedState.Prompts)-1]
				if lastPrompt.Prompt != "Hello, Claude!" {
					t.Errorf("Expected prompt content 'Hello, Claude!', got '%s'", lastPrompt.Prompt)
				}
				if lastPrompt.Timestamp.IsZero() {
					t.Error("Prompt timestamp should not be zero")
				}
			}
		})
	}
}

func TestUserPromptSubmitHandlerName(t *testing.T) {
	handler := NewUserPromptSubmitHandler()
	if handler.Name() != "user_prompt_submit" {
		t.Errorf("Expected handler name 'user_prompt_submit', got '%s'", handler.Name())
	}
}

func TestUserPromptSubmitHandlerTimestampParsing(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	oldDir, _ := os.Getwd()
	os.Chdir(tempDir)
	defer os.Chdir(oldDir)

	// Create .spcstr directory structure
	os.MkdirAll(".spcstr/sessions", 0755)

	// First create a session state
	manager := state.NewStateManager(filepath.Join(tempDir, ".spcstr"))
	_, err := manager.InitializeState(context.Background(), "timestamp_test_123")
	if err != nil {
		t.Fatalf("Failed to initialize state: %v", err)
	}

	handler := NewUserPromptSubmitHandler()

	// Test with valid timestamp
	input := `{"session_id": "timestamp_test_123", "prompt": "Test prompt", "timestamp": "2023-06-15T14:30:00Z"}`
	err = handler.Execute([]byte(input))
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	// Verify the timestamp was parsed correctly
	updatedState, err := manager.LoadState(context.Background(), "timestamp_test_123")
	if err != nil {
		t.Fatalf("Failed to load updated state: %v", err)
	}

	if len(updatedState.Prompts) == 0 {
		t.Fatal("No prompts found in session state")
	}

	expectedTime, _ := time.Parse(time.RFC3339, "2023-06-15T14:30:00Z")
	actualTime := updatedState.Prompts[0].Timestamp

	if !actualTime.Equal(expectedTime) {
		t.Errorf("Expected timestamp %v, got %v", expectedTime, actualTime)
	}
}