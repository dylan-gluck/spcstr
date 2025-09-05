package hooks

import (
	"encoding/json"
	"testing"
	"time"
)

func TestSessionStateSerialization(t *testing.T) {
	state := &SessionState{
		SessionID:   "sess_test_123",
		Source:      "test",
		ProjectPath: "/test/path",
		Timestamp:   "2024-01-01T00:00:00Z",
		LastUpdate:  "2024-01-01T00:01:00Z",
		Status:      "active",
		Agents:      []string{"agent1", "agent2"},
		AgentsHistory: []AgentHistoryEntry{
			{
				Name:      "agent1",
				StartedAt: "2024-01-01T00:00:00Z",
				EndedAt:   "2024-01-01T00:00:30Z",
			},
		},
		Files: FileOperations{
			New:    []string{"/file1"},
			Edited: []string{"/file2"},
			Read:   []string{"/file3"},
		},
		ToolsUsed: map[string]int{
			"Read":  5,
			"Write": 3,
		},
		Errors: []ErrorEntry{
			{
				Timestamp: "2024-01-01T00:00:45Z",
				Hook:      "test_hook",
				Message:   "test error",
			},
		},
		Modified: true,
	}

	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal SessionState: %v", err)
	}

	var decoded SessionState
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal SessionState: %v", err)
	}

	if decoded.SessionID != state.SessionID {
		t.Errorf("SessionID mismatch: got %s, want %s", decoded.SessionID, state.SessionID)
	}

	if len(decoded.Agents) != len(state.Agents) {
		t.Errorf("Agents count mismatch: got %d, want %d", len(decoded.Agents), len(state.Agents))
	}

	if decoded.ToolsUsed["Read"] != state.ToolsUsed["Read"] {
		t.Errorf("ToolsUsed[Read] mismatch: got %d, want %d", decoded.ToolsUsed["Read"], state.ToolsUsed["Read"])
	}
}

func TestCurrentTimestamp(t *testing.T) {
	ts := CurrentTimestamp()

	_, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t.Errorf("CurrentTimestamp returned invalid RFC3339 format: %s, error: %v", ts, err)
	}
}

func TestAddUniqueString(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		str      string
		expected []string
	}{
		{
			name:     "add to empty slice",
			slice:    []string{},
			str:      "test",
			expected: []string{"test"},
		},
		{
			name:     "add unique string",
			slice:    []string{"a", "b"},
			str:      "c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "add duplicate string",
			slice:    []string{"a", "b", "c"},
			str:      "b",
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddUniqueString(tt.slice, tt.str)

			if len(result) != len(tt.expected) {
				t.Errorf("length mismatch: got %d, want %d", len(result), len(tt.expected))
				return
			}

			for i, v := range result {
				if v != tt.expected[i] {
					t.Errorf("element %d mismatch: got %s, want %s", i, v, tt.expected[i])
				}
			}
		})
	}
}
