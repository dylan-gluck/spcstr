package hooks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestLoadSessionState(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("new session", func(t *testing.T) {
		state, err := LoadSessionState(tmpDir)
		if err != nil {
			t.Fatalf("Failed to load new session state: %v", err)
		}

		if state.SessionID == "" {
			t.Error("SessionID should not be empty for new session")
		}

		if state.Status != "active" {
			t.Errorf("Status should be 'active', got %s", state.Status)
		}

		if state.Source != "startup" {
			t.Errorf("Source should be 'startup', got %s", state.Source)
		}

		if state.ToolsUsed == nil {
			t.Error("ToolsUsed should be initialized")
		}
	})

	t.Run("existing session", func(t *testing.T) {
		existingState := &SessionState{
			SessionID:     "sess_existing_123",
			Source:        "test",
			ProjectPath:   tmpDir,
			Timestamp:     CurrentTimestamp(),
			LastUpdate:    CurrentTimestamp(),
			Status:        "active",
			Agents:        []string{"agent1"},
			AgentsHistory: []AgentHistoryEntry{},
			Files: FileOperations{
				New:    []string{"/file1"},
				Edited: []string{},
				Read:   []string{},
			},
			ToolsUsed: map[string]int{"Read": 1},
			Errors:    []ErrorEntry{},
			Modified:  true,
		}

		sessionFile := filepath.Join(tmpDir, ".spcstr", "session_state.json")
		os.MkdirAll(filepath.Dir(sessionFile), 0755)

		data, _ := json.MarshalIndent(existingState, "", "  ")
		os.WriteFile(sessionFile, data, 0644)

		loadedState, err := LoadSessionState(tmpDir)
		if err != nil {
			t.Fatalf("Failed to load existing session state: %v", err)
		}

		if loadedState.SessionID != "sess_existing_123" {
			t.Errorf("SessionID mismatch: got %s, want sess_existing_123", loadedState.SessionID)
		}

		if len(loadedState.Files.New) != 1 || loadedState.Files.New[0] != "/file1" {
			t.Error("Files.New not loaded correctly")
		}

		if loadedState.ToolsUsed["Read"] != 1 {
			t.Errorf("ToolsUsed[Read] mismatch: got %d, want 1", loadedState.ToolsUsed["Read"])
		}
	})
}

func TestSaveSessionState(t *testing.T) {
	tmpDir := t.TempDir()

	state := &SessionState{
		SessionID:     "sess_save_123",
		Source:        "test",
		ProjectPath:   tmpDir,
		Timestamp:     CurrentTimestamp(),
		LastUpdate:    "",
		Status:        "active",
		Agents:        []string{},
		AgentsHistory: []AgentHistoryEntry{},
		Files: FileOperations{
			New:    []string{"/file1"},
			Edited: []string{"/file2"},
			Read:   []string{"/file3"},
		},
		ToolsUsed: map[string]int{"Write": 2},
		Errors:    []ErrorEntry{},
		Modified:  false,
	}

	err := SaveSessionState(tmpDir, state)
	if err != nil {
		t.Fatalf("Failed to save session state: %v", err)
	}

	if state.LastUpdate == "" {
		t.Error("LastUpdate should be set after save")
	}

	if !state.Modified {
		t.Error("Modified should be true after save")
	}

	sessionFile := filepath.Join(tmpDir, ".spcstr", "session_state.json")
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		t.Error("Session file was not created")
	}

	data, err := os.ReadFile(sessionFile)
	if err != nil {
		t.Fatalf("Failed to read saved session file: %v", err)
	}

	var loaded SessionState
	if err := json.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("Failed to unmarshal saved session: %v", err)
	}

	if loaded.SessionID != "sess_save_123" {
		t.Errorf("Saved SessionID mismatch: got %s, want sess_save_123", loaded.SessionID)
	}

	if len(loaded.Files.Edited) != 1 || loaded.Files.Edited[0] != "/file2" {
		t.Error("Files.Edited not saved correctly")
	}
}

func TestConcurrentAccess(t *testing.T) {
	tmpDir := t.TempDir()

	initialState := &SessionState{
		SessionID:     "sess_concurrent_123",
		Source:        "test",
		ProjectPath:   tmpDir,
		Timestamp:     CurrentTimestamp(),
		LastUpdate:    CurrentTimestamp(),
		Status:        "active",
		Agents:        []string{},
		AgentsHistory: []AgentHistoryEntry{},
		Files: FileOperations{
			New:    []string{},
			Edited: []string{},
			Read:   []string{},
		},
		ToolsUsed: make(map[string]int),
		Errors:    []ErrorEntry{},
		Modified:  false,
	}

	err := SaveSessionState(tmpDir, initialState)
	if err != nil {
		t.Fatalf("Failed to save initial state: %v", err)
	}

	var wg sync.WaitGroup
	numGoroutines := 10

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			state, err := LoadSessionState(tmpDir)
			if err != nil {
				t.Errorf("Goroutine %d: Failed to load state: %v", id, err)
				return
			}

			toolName := "Tool" + string(rune('A'+id))
			if state.ToolsUsed == nil {
				state.ToolsUsed = make(map[string]int)
			}
			state.ToolsUsed[toolName]++

			err = SaveSessionState(tmpDir, state)
			if err != nil {
				t.Errorf("Goroutine %d: Failed to save state: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	finalState, err := LoadSessionState(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load final state: %v", err)
	}

	if len(finalState.ToolsUsed) == 0 {
		t.Error("No tools were recorded in concurrent access test")
	}
}
