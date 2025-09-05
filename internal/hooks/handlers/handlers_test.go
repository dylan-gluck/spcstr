package handlers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dylan-gluck/spcstr/internal/session"
	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

func setupTestEnvironment(t *testing.T) string {
	tmpDir := t.TempDir()
	sessionDir := filepath.Join(tmpDir, ".spcstr")
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		t.Fatalf("Failed to create session dir: %v", err)
	}

	// Create initial session state
	state := &models.SessionState{
		SessionID:     "sess_test_" + t.Name(),
		Source:        "test",
		ProjectPath:   tmpDir,
		Timestamp:     models.CurrentTimestamp(),
		LastUpdate:    models.CurrentTimestamp(),
		Status:        "active",
		Agents:        []string{},
		AgentsHistory: []models.AgentHistoryEntry{},
		Files: models.FileOperations{
			New:    []string{},
			Edited: []string{},
			Read:   []string{},
		},
		ToolsUsed: make(map[string]int),
		Errors:    []models.ErrorEntry{},
		Modified:  false,
	}

	if err := session.SaveSessionState(tmpDir, state); err != nil {
		t.Fatalf("Failed to save initial state: %v", err)
	}

	return tmpDir
}

func TestHandlePreToolUse(t *testing.T) {
	tmpDir := setupTestEnvironment(t)

	tests := []struct {
		name      string
		data      map[string]interface{}
		wantError bool
	}{
		{
			name: "safe_command",
			data: map[string]interface{}{
				"tool_name": "Bash",
				"parameters": map[string]interface{}{
					"command": "ls -la",
				},
			},
			wantError: false,
		},
		{
			name: "dangerous_rm_command",
			data: map[string]interface{}{
				"tool_name": "Bash",
				"parameters": map[string]interface{}{
					"command": "rm -rf /",
				},
			},
			wantError: true, // Should block
		},
		{
			name: "dangerous_fork_bomb",
			data: map[string]interface{}{
				"tool_name": "Bash",
				"parameters": map[string]interface{}{
					"command": ":(){ :|:& };:",
				},
			},
			wantError: true, // Should block
		},
		{
			name: "non_bash_tool",
			data: map[string]interface{}{
				"tool_name": "Read",
				"parameters": map[string]interface{}{
					"file_path": "/some/file.txt",
				},
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := HandlePreToolUse(tmpDir, tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("HandlePreToolUse() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestHandlePostToolUse(t *testing.T) {
	tests := []struct {
		name       string
		data       map[string]interface{}
		wantError  bool
		checkState func(*testing.T, *models.SessionState)
	}{
		{
			name: "track_read_tool",
			data: map[string]interface{}{
				"tool_name": "Read",
				"parameters": map[string]interface{}{
					"file_path": "/test/file.txt",
				},
			},
			checkState: func(t *testing.T, state *models.SessionState) {
				if state.ToolsUsed["Read"] != 1 {
					t.Errorf("ToolsUsed[Read] = %d, want 1", state.ToolsUsed["Read"])
				}
				if len(state.Files.Read) != 1 || state.Files.Read[0] != "/test/file.txt" {
					t.Errorf("Files.Read = %v, want [/test/file.txt]", state.Files.Read)
				}
			},
		},
		{
			name: "track_write_tool",
			data: map[string]interface{}{
				"tool_name": "Write",
				"parameters": map[string]interface{}{
					"file_path": "/test/new.txt",
				},
			},
			checkState: func(t *testing.T, state *models.SessionState) {
				if state.ToolsUsed["Write"] != 1 {
					t.Errorf("ToolsUsed[Write] = %d, want 1", state.ToolsUsed["Write"])
				}
				if len(state.Files.New) != 1 || state.Files.New[0] != "/test/new.txt" {
					t.Errorf("Files.New = %v, want [/test/new.txt]", state.Files.New)
				}
			},
		},
		{
			name: "track_edit_tool",
			data: map[string]interface{}{
				"tool_name": "Edit",
				"parameters": map[string]interface{}{
					"file_path": "/test/edit.txt",
				},
			},
			checkState: func(t *testing.T, state *models.SessionState) {
				if state.ToolsUsed["Edit"] != 1 {
					t.Errorf("ToolsUsed[Edit] = %d, want 1", state.ToolsUsed["Edit"])
				}
				if len(state.Files.Edited) != 1 || state.Files.Edited[0] != "/test/edit.txt" {
					t.Errorf("Files.Edited = %v, want [/test/edit.txt]", state.Files.Edited)
				}
			},
		},
		{
			name: "track_bash_tool",
			data: map[string]interface{}{
				"tool_name": "Bash",
				"parameters": map[string]interface{}{
					"command": "echo test",
				},
			},
			checkState: func(t *testing.T, state *models.SessionState) {
				if state.ToolsUsed["Bash"] != 1 {
					t.Errorf("ToolsUsed[Bash] = %d, want 1", state.ToolsUsed["Bash"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset state for each test
			tmpDir := setupTestEnvironment(t)

			err := HandlePostToolUse(tmpDir, tt.data)
			if err != nil {
				t.Errorf("HandlePostToolUse() error = %v", err)
			}

			// Load state and verify
			state, err := session.LoadSessionState(tmpDir)
			if err != nil {
				t.Fatalf("Failed to load state: %v", err)
			}

			if tt.checkState != nil {
				tt.checkState(t, state)
			}
		})
	}
}

func TestHandleSessionStart(t *testing.T) {
	tmpDir := setupTestEnvironment(t)

	data := map[string]interface{}{
		"source": "startup",
	}

	err := HandleSessionStart(tmpDir, data)
	if err != nil {
		t.Errorf("HandleSessionStart() error = %v", err)
	}

	// Verify state
	state, err := session.LoadSessionState(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	if state.Status != "active" {
		t.Errorf("Status = %s, want active", state.Status)
	}

	if state.Source != "startup" {
		t.Errorf("Source = %s, want startup", state.Source)
	}
}

func TestHandleSessionEnd(t *testing.T) {
	tmpDir := setupTestEnvironment(t)

	err := HandleSessionEnd(tmpDir, nil)
	if err != nil {
		t.Errorf("HandleSessionEnd() error = %v", err)
	}

	// Verify state
	state, err := session.LoadSessionState(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	if state.Status != "completed" {
		t.Errorf("Status = %s, want completed", state.Status)
	}
}

func TestHandleUserPrompt(t *testing.T) {
	tmpDir := setupTestEnvironment(t)

	data := map[string]interface{}{
		"text": "test prompt",
	}

	err := HandleUserPrompt(tmpDir, data)
	if err != nil {
		t.Errorf("HandleUserPrompt() error = %v", err)
	}
}

func TestHandleNotification(t *testing.T) {
	tmpDir := setupTestEnvironment(t)

	data := map[string]interface{}{
		"type":    "info",
		"message": "test notification",
	}

	err := HandleNotification(tmpDir, data)
	if err != nil {
		t.Errorf("HandleNotification() error = %v", err)
	}
}

func TestHandleStop(t *testing.T) {
	tmpDir := setupTestEnvironment(t)

	data := map[string]interface{}{
		"reason": "user_requested",
	}

	err := HandleStop(tmpDir, data)
	if err != nil {
		t.Errorf("HandleStop() error = %v", err)
	}

	// Verify state
	state, err := session.LoadSessionState(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	if state.Status != "stopped" {
		t.Errorf("Status = %s, want stopped", state.Status)
	}
}

func TestHandleSubagentStart(t *testing.T) {
	tmpDir := setupTestEnvironment(t)

	data := map[string]interface{}{
		"agent_name": "test_agent",
	}

	err := HandleSubagentStart(tmpDir, data)
	if err != nil {
		t.Errorf("HandleSubagentStart() error = %v", err)
	}

	// Verify state
	state, err := session.LoadSessionState(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	if len(state.Agents) != 1 || state.Agents[0] != "test_agent" {
		t.Errorf("Agents = %v, want [test_agent]", state.Agents)
	}

	if len(state.AgentsHistory) != 1 {
		t.Errorf("AgentsHistory length = %d, want 1", len(state.AgentsHistory))
	}
}

func TestHandleSubagentStop(t *testing.T) {
	tmpDir := setupTestEnvironment(t)

	// First start an agent
	startData := map[string]interface{}{
		"agent_name": "test_agent",
	}
	HandleSubagentStart(tmpDir, startData)

	// Then stop it
	stopData := map[string]interface{}{
		"agent_name": "test_agent",
	}

	err := HandleSubagentStop(tmpDir, stopData)
	if err != nil {
		t.Errorf("HandleSubagentStop() error = %v", err)
	}

	// Verify state
	state, err := session.LoadSessionState(tmpDir)
	if err != nil {
		t.Fatalf("Failed to load state: %v", err)
	}

	if len(state.Agents) != 0 {
		t.Errorf("Agents = %v, want []", state.Agents)
	}

	if len(state.AgentsHistory) != 1 {
		t.Errorf("AgentsHistory length = %d, want 1", len(state.AgentsHistory))
	}

	// Check that end time is set
	if state.AgentsHistory[0].EndedAt == "" {
		t.Error("AgentHistoryEntry.EndedAt should be set")
	}
}

func TestHandlePreCompact(t *testing.T) {
	tmpDir := setupTestEnvironment(t)

	err := HandlePreCompact(tmpDir, nil)
	if err != nil {
		t.Errorf("HandlePreCompact() error = %v", err)
	}
}
