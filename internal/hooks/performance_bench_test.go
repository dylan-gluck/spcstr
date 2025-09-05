package hooks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func BenchmarkHookExecution(b *testing.B) {
	tmpDir := b.TempDir()
	projectRoot := tmpDir

	handler := NewHandler(projectRoot)

	testData := map[string]interface{}{
		"session_id": "sess_test_123",
		"tool_name":  "Read",
		"parameters": map[string]interface{}{
			"file_path": "/test/file.txt",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.HandleHook("post_tool_use", testData)
	}
}

func BenchmarkSessionStateLoad(b *testing.B) {
	tmpDir := b.TempDir()

	state := &SessionState{
		SessionID:     "sess_bench_123",
		Source:        "test",
		ProjectPath:   tmpDir,
		Timestamp:     CurrentTimestamp(),
		LastUpdate:    CurrentTimestamp(),
		Status:        "active",
		Agents:        []string{"agent1", "agent2"},
		AgentsHistory: []AgentHistoryEntry{},
		Files: FileOperations{
			New:    []string{"/file1", "/file2"},
			Edited: []string{"/file3", "/file4"},
			Read:   []string{"/file5", "/file6"},
		},
		ToolsUsed: map[string]int{"Read": 10, "Write": 5},
		Errors:    []ErrorEntry{},
		Modified:  true,
	}

	sessionFile := filepath.Join(tmpDir, ".spcstr", "session_state.json")
	os.MkdirAll(filepath.Dir(sessionFile), 0755)

	data, _ := json.MarshalIndent(state, "", "  ")
	os.WriteFile(sessionFile, data, 0644)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadSessionState(tmpDir)
	}
}

func BenchmarkSessionStateSave(b *testing.B) {
	tmpDir := b.TempDir()

	state := &SessionState{
		SessionID:     "sess_bench_123",
		Source:        "test",
		ProjectPath:   tmpDir,
		Timestamp:     CurrentTimestamp(),
		LastUpdate:    CurrentTimestamp(),
		Status:        "active",
		Agents:        []string{"agent1", "agent2"},
		AgentsHistory: []AgentHistoryEntry{},
		Files: FileOperations{
			New:    []string{"/file1", "/file2"},
			Edited: []string{"/file3", "/file4"},
			Read:   []string{"/file5", "/file6"},
		},
		ToolsUsed: map[string]int{"Read": 10, "Write": 5},
		Errors:    []ErrorEntry{},
		Modified:  true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SaveSessionState(tmpDir, state)
	}
}

func TestHookPerformanceUnder10ms(t *testing.T) {
	tmpDir := t.TempDir()
	projectRoot := tmpDir

	handler := NewHandler(projectRoot)

	testCases := []struct {
		name string
		hook string
		data map[string]interface{}
	}{
		{
			name: "pre_tool_use",
			hook: "pre_tool_use",
			data: map[string]interface{}{
				"tool_name": "Bash",
				"parameters": map[string]interface{}{
					"command": "ls -la",
				},
			},
		},
		{
			name: "post_tool_use",
			hook: "post_tool_use",
			data: map[string]interface{}{
				"tool_name": "Read",
				"parameters": map[string]interface{}{
					"file_path": "/test/file.txt",
				},
			},
		},
		{
			name: "session_start",
			hook: "session_start",
			data: map[string]interface{}{
				"session_id": "sess_test_123",
				"source":     "startup",
			},
		},
		{
			name: "session_end",
			hook: "session_end",
			data: map[string]interface{}{
				"session_id": "sess_test_123",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			err := handler.HandleHook(tc.hook, tc.data)
			duration := time.Since(start)

			if err != nil {
				t.Errorf("hook %s failed: %v", tc.hook, err)
			}

			if duration > 10*time.Millisecond {
				t.Errorf("hook %s took %v, exceeds 10ms requirement", tc.hook, duration)
			}
		})
	}
}
