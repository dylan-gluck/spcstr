package state

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestNewStateManager(t *testing.T) {
	basePath := "/test/path"
	manager := NewStateManager(basePath)

	if manager.basePath != basePath {
		t.Errorf("NewStateManager() basePath = %q, want %q", manager.basePath, basePath)
	}

	if manager.timeout != DefaultTimeout {
		t.Errorf("NewStateManager() timeout = %v, want %v", manager.timeout, DefaultTimeout)
	}

	if manager.writer == nil {
		t.Error("NewStateManager() writer is nil")
	}
}

func TestNewStateManagerWithTimeout(t *testing.T) {
	basePath := "/test/path"
	customTimeout := 10 * time.Second
	manager := NewStateManagerWithTimeout(basePath, customTimeout)

	if manager.timeout != customTimeout {
		t.Errorf("NewStateManagerWithTimeout() timeout = %v, want %v", manager.timeout, customTimeout)
	}
}

func TestStateManager_InitializeState(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "manager_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewStateManager(tmpDir)

	tests := []struct {
		name      string
		sessionID string
		wantErr   bool
	}{
		{
			name:      "valid session ID",
			sessionID: "test_session_1",
			wantErr:   false,
		},
		{
			name:      "session ID with special characters",
			sessionID: "session-2024_01_01",
			wantErr:   false,
		},
		{
			name:      "empty session ID",
			sessionID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			state, err := manager.InitializeState(ctx, tt.sessionID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("InitializeState() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("InitializeState() unexpected error: %v", err)
				return
			}

			// Verify state properties
			if state.SessionID != tt.sessionID {
				t.Errorf("InitializeState() SessionID = %q, want %q", state.SessionID, tt.sessionID)
			}

			if !state.SessionActive {
				t.Error("InitializeState() SessionActive = false, want true")
			}

			if state.CreatedAt.IsZero() {
				t.Error("InitializeState() CreatedAt is zero")
			}

			if state.UpdatedAt.IsZero() {
				t.Error("InitializeState() UpdatedAt is zero")
			}

			// Verify file was created at correct path
			expectedPath := filepath.Join(tmpDir, "sessions", tt.sessionID, StateFileName)
			if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
				t.Errorf("InitializeState() file not created at expected path: %s", expectedPath)
			}

			// Verify JSON structure
			data, err := os.ReadFile(expectedPath)
			if err != nil {
				t.Errorf("Failed to read created state file: %v", err)
				return
			}

			var readState SessionState
			if err := json.Unmarshal(data, &readState); err != nil {
				t.Errorf("Created state file contains invalid JSON: %v", err)
			}
		})
	}
}

func TestStateManager_LoadState(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "manager_load_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewStateManager(tmpDir)
	ctx := context.Background()

	// Initialize a state first
	sessionID := "test_load_session"
	originalState, err := manager.InitializeState(ctx, sessionID)
	if err != nil {
		t.Fatalf("Failed to initialize state for test: %v", err)
	}

	tests := []struct {
		name      string
		sessionID string
		wantErr   bool
	}{
		{
			name:      "existing session",
			sessionID: sessionID,
			wantErr:   false,
		},
		{
			name:      "non-existent session",
			sessionID: "non_existent",
			wantErr:   true,
		},
		{
			name:      "empty session ID",
			sessionID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			loadedState, err := manager.LoadState(ctx, tt.sessionID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadState() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("LoadState() unexpected error: %v", err)
				return
			}

			// Verify loaded state matches original
			if loadedState.SessionID != originalState.SessionID {
				t.Errorf("LoadState() SessionID = %q, want %q", loadedState.SessionID, originalState.SessionID)
			}

			if loadedState.SessionActive != originalState.SessionActive {
				t.Errorf("LoadState() SessionActive = %v, want %v", loadedState.SessionActive, originalState.SessionActive)
			}
		})
	}
}

func TestStateManager_UpdateState(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "manager_update_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewStateManager(tmpDir)
	ctx := context.Background()

	// Initialize a state first
	sessionID := "test_update_session"
	_, err = manager.InitializeState(ctx, sessionID)
	if err != nil {
		t.Fatalf("Failed to initialize state for test: %v", err)
	}

	tests := []struct {
		name       string
		sessionID  string
		updateFunc func(*SessionState) error
		wantErr    bool
	}{
		{
			name:      "valid update",
			sessionID: sessionID,
			updateFunc: func(state *SessionState) error {
				state.SessionActive = false
				state.Agents = append(state.Agents, "test_agent")
				return nil
			},
			wantErr: false,
		},
		{
			name:      "update function error",
			sessionID: sessionID,
			updateFunc: func(state *SessionState) error {
				return &StateError{Code: "test_error", Message: "test error"}
			},
			wantErr: true,
		},
		{
			name:      "non-existent session",
			sessionID: "non_existent",
			updateFunc: func(state *SessionState) error {
				return nil
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalUpdatedAt := time.Time{}
			if !tt.wantErr && tt.sessionID == sessionID {
				// Get current UpdatedAt for comparison
				state, _ := manager.LoadState(ctx, sessionID)
				originalUpdatedAt = state.UpdatedAt
			}

			err := manager.UpdateState(ctx, tt.sessionID, tt.updateFunc)

			if tt.wantErr {
				if err == nil {
					t.Errorf("UpdateState() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("UpdateState() unexpected error: %v", err)
				return
			}

			// Verify update was applied
			updatedState, err := manager.LoadState(ctx, tt.sessionID)
			if err != nil {
				t.Errorf("Failed to load updated state: %v", err)
				return
			}

			// Verify UpdatedAt timestamp was updated
			if !updatedState.UpdatedAt.After(originalUpdatedAt) {
				t.Error("UpdateState() did not update the UpdatedAt timestamp")
			}
		})
	}
}

func TestStateManager_DeleteState(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "manager_delete_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewStateManager(tmpDir)
	ctx := context.Background()

	// Initialize a state first
	sessionID := "test_delete_session"
	_, err = manager.InitializeState(ctx, sessionID)
	if err != nil {
		t.Fatalf("Failed to initialize state for test: %v", err)
	}

	tests := []struct {
		name      string
		sessionID string
		wantErr   bool
	}{
		{
			name:      "existing session",
			sessionID: sessionID,
			wantErr:   false,
		},
		{
			name:      "non-existent session",
			sessionID: "non_existent",
			wantErr:   true,
		},
		{
			name:      "empty session ID",
			sessionID: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := manager.DeleteState(ctx, tt.sessionID)

			if tt.wantErr {
				if err == nil {
					t.Errorf("DeleteState() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("DeleteState() unexpected error: %v", err)
				return
			}

			// Verify file was deleted
			expectedPath := filepath.Join(tmpDir, "sessions", tt.sessionID, StateFileName)
			if _, err := os.Stat(expectedPath); !os.IsNotExist(err) {
				t.Errorf("DeleteState() file still exists: %s", expectedPath)
			}
		})
	}
}

func TestStateManager_ListSessions(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "manager_list_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewStateManager(tmpDir)
	ctx := context.Background()

	// Test empty directory
	sessions, err := manager.ListSessions(ctx)
	if err != nil {
		t.Errorf("ListSessions() unexpected error: %v", err)
	}
	if len(sessions) != 0 {
		t.Errorf("ListSessions() expected 0 sessions, got %d", len(sessions))
	}

	// Initialize some states
	sessionIDs := []string{"session1", "session2", "session3"}
	for _, sessionID := range sessionIDs {
		_, err := manager.InitializeState(ctx, sessionID)
		if err != nil {
			t.Fatalf("Failed to initialize state %s: %v", sessionID, err)
		}
	}

	// Test with existing sessions
	sessions, err = manager.ListSessions(ctx)
	if err != nil {
		t.Errorf("ListSessions() unexpected error: %v", err)
	}

	if len(sessions) != len(sessionIDs) {
		t.Errorf("ListSessions() expected %d sessions, got %d", len(sessionIDs), len(sessions))
	}

	// Verify all session IDs are present
	for _, expectedID := range sessionIDs {
		found := false
		for _, actualID := range sessions {
			if actualID == expectedID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ListSessions() missing session ID: %s", expectedID)
		}
	}
}

func TestStateManager_ConvenienceMethods(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "manager_convenience_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	manager := NewStateManager(tmpDir)
	ctx := context.Background()

	// Initialize a state
	sessionID := "convenience_test_session"
	_, err = manager.InitializeState(ctx, sessionID)
	if err != nil {
		t.Fatalf("Failed to initialize state: %v", err)
	}

	t.Run("AddAgent", func(t *testing.T) {
		agentName := "test_agent"
		err := manager.AddAgent(ctx, sessionID, agentName)
		if err != nil {
			t.Errorf("AddAgent() unexpected error: %v", err)
		}

		// Verify agent was added
		state, _ := manager.LoadState(ctx, sessionID)
		found := false
		for _, agent := range state.Agents {
			if agent == agentName {
				found = true
				break
			}
		}
		if !found {
			t.Error("AddAgent() agent not found in current agents")
		}

		// Verify agent was added to history
		historyFound := false
		for _, exec := range state.AgentsHistory {
			if exec.Name == agentName {
				historyFound = true
				break
			}
		}
		if !historyFound {
			t.Error("AddAgent() agent not found in agents history")
		}
	})

	t.Run("CompleteAgent", func(t *testing.T) {
		agentName := "test_agent"
		err := manager.CompleteAgent(ctx, sessionID, agentName)
		if err != nil {
			t.Errorf("CompleteAgent() unexpected error: %v", err)
		}

		// Verify agent was removed from current agents
		state, _ := manager.LoadState(ctx, sessionID)
		for _, agent := range state.Agents {
			if agent == agentName {
				t.Error("CompleteAgent() agent still in current agents")
			}
		}

		// Verify completion time was set in history
		completionFound := false
		for _, exec := range state.AgentsHistory {
			if exec.Name == agentName && exec.CompletedAt != nil {
				completionFound = true
				break
			}
		}
		if !completionFound {
			t.Error("CompleteAgent() completion time not set in history")
		}
	})

	t.Run("RecordError", func(t *testing.T) {
		err := manager.RecordError(ctx, sessionID, "test error", "test source", "high")
		if err != nil {
			t.Errorf("RecordError() unexpected error: %v", err)
		}

		state, _ := manager.LoadState(ctx, sessionID)
		if len(state.Errors) == 0 {
			t.Error("RecordError() no errors recorded")
		} else {
			lastError := state.Errors[len(state.Errors)-1]
			if lastError.Message != "test error" {
				t.Errorf("RecordError() message = %q, want %q", lastError.Message, "test error")
			}
		}
	})

	t.Run("RecordFileOperation", func(t *testing.T) {
		testCases := []struct {
			operation string
			filepath  string
		}{
			{"new", "/path/to/new/file.go"},
			{"edited", "/path/to/edited/file.go"},
			{"read", "/path/to/read/file.go"},
		}

		for _, tc := range testCases {
			err := manager.RecordFileOperation(ctx, sessionID, tc.operation, tc.filepath)
			if err != nil {
				t.Errorf("RecordFileOperation(%s) unexpected error: %v", tc.operation, err)
			}
		}

		state, _ := manager.LoadState(ctx, sessionID)
		if len(state.Files.New) == 0 || len(state.Files.Edited) == 0 || len(state.Files.Read) == 0 {
			t.Error("RecordFileOperation() file operations not recorded")
		}

		// Test invalid operation
		err := manager.RecordFileOperation(ctx, sessionID, "invalid", "/path")
		if err == nil {
			t.Error("RecordFileOperation() expected error for invalid operation")
		}
	})

	t.Run("IncrementToolUsage", func(t *testing.T) {
		toolName := "test_tool"

		// Increment multiple times
		for i := 0; i < 3; i++ {
			err := manager.IncrementToolUsage(ctx, sessionID, toolName)
			if err != nil {
				t.Errorf("IncrementToolUsage() unexpected error: %v", err)
			}
		}

		state, _ := manager.LoadState(ctx, sessionID)
		if count, exists := state.ToolsUsed[toolName]; !exists || count != 3 {
			t.Errorf("IncrementToolUsage() count = %d, want 3", count)
		}
	})

	t.Run("SetSessionActive", func(t *testing.T) {
		err := manager.SetSessionActive(ctx, sessionID, false)
		if err != nil {
			t.Errorf("SetSessionActive() unexpected error: %v", err)
		}

		state, _ := manager.LoadState(ctx, sessionID)
		if state.SessionActive {
			t.Error("SetSessionActive() session still active")
		}
	})
}

func TestJSONSchemaCompliance(t *testing.T) {
	// Test that our data structures marshal/unmarshal correctly with ISO8601 timestamps
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	completedAt := time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC)

	state := &SessionState{
		SessionID:     "schema_test",
		CreatedAt:     now,
		UpdatedAt:     now,
		SessionActive: true,
		Agents:        []string{"agent1"},
		AgentsHistory: []AgentExecution{
			{
				Name:        "agent1",
				StartedAt:   now,
				CompletedAt: &completedAt,
			},
		},
		Files: FileOperations{
			New:    []string{"file1.go"},
			Edited: []string{"file2.go"},
			Read:   []string{"file3.go"},
		},
		ToolsUsed: map[string]int{"read": 1},
		Errors: []ErrorEntry{
			{
				Timestamp: now,
				Message:   "test error",
				Source:    "test",
				Severity:  "low",
			},
		},
		Prompts: []PromptEntry{
			{
				Timestamp: now,
				Prompt:    "test prompt",
				Response:  "test response",
				ToolsUsed: []string{"read"},
			},
		},
		Notifications: []NotificationEntry{
			{
				Timestamp: now,
				Type:      "info",
				Message:   "test notification",
				Level:     "info",
			},
		},
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal state to JSON: %v", err)
	}

	// Verify ISO8601 timestamp format
	jsonString := string(jsonData)
	expectedTimestamp := "2024-01-01T12:00:00Z"
	if !strings.Contains(jsonString, expectedTimestamp) {
		t.Errorf("JSON does not contain expected ISO8601 timestamp %s", expectedTimestamp)
	}

	// Unmarshal back to struct
	var unmarshaledState SessionState
	if err := json.Unmarshal(jsonData, &unmarshaledState); err != nil {
		t.Fatalf("Failed to unmarshal JSON to state: %v", err)
	}

	// Verify data integrity
	if !reflect.DeepEqual(state.SessionID, unmarshaledState.SessionID) {
		t.Error("SessionID mismatch after marshal/unmarshal")
	}

	if !state.CreatedAt.Equal(unmarshaledState.CreatedAt) {
		t.Error("CreatedAt timestamp mismatch after marshal/unmarshal")
	}

	if !reflect.DeepEqual(state.Files, unmarshaledState.Files) {
		t.Error("Files mismatch after marshal/unmarshal")
	}

	if !reflect.DeepEqual(state.ToolsUsed, unmarshaledState.ToolsUsed) {
		t.Error("ToolsUsed mismatch after marshal/unmarshal")
	}
}
