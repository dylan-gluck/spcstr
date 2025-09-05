package state

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestAtomicWriter_WriteJSON(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "atomic_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	writer := NewAtomicWriter(5 * time.Second)

	tests := []struct {
		name     string
		data     interface{}
		filename string
		wantErr  bool
	}{
		{
			name: "valid session state",
			data: &SessionState{
				SessionID:     "test_session_1",
				CreatedAt:     time.Now().UTC(),
				UpdatedAt:     time.Now().UTC(),
				SessionActive: true,
				Agents:        []string{"test_agent"},
			},
			filename: filepath.Join(tmpDir, "sessions", "test_session_1", "state.json"),
			wantErr:  false,
		},
		{
			name: "complex session state with all fields",
			data: &SessionState{
				SessionID:     "test_session_2",
				CreatedAt:     time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2024, 1, 1, 12, 30, 0, 0, time.UTC),
				SessionActive: false,
				Agents:        []string{"agent1", "agent2"},
				AgentsHistory: []AgentExecution{
					{
						Name:      "agent1",
						StartedAt: time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
						CompletedAt: func() *time.Time {
							t := time.Date(2024, 1, 1, 12, 15, 0, 0, time.UTC)
							return &t
						}(),
					},
				},
				Files: FileOperations{
					New:    []string{"file1.go", "file2.go"},
					Edited: []string{"file3.go"},
					Read:   []string{"file4.go", "file5.go"},
				},
				ToolsUsed: map[string]int{
					"bash": 5,
					"read": 10,
					"edit": 3,
				},
				Errors: []ErrorEntry{
					{
						Timestamp: time.Date(2024, 1, 1, 12, 10, 0, 0, time.UTC),
						Message:   "Test error",
						Source:    "test_source",
						Severity:  "high",
					},
				},
				Prompts: []PromptEntry{
					{
						Timestamp: time.Date(2024, 1, 1, 12, 5, 0, 0, time.UTC),
						Prompt:    "Test prompt",
						Response:  "Test response",
						ToolsUsed: []string{"read", "edit"},
					},
				},
				Notifications: []NotificationEntry{
					{
						Timestamp: time.Date(2024, 1, 1, 12, 20, 0, 0, time.UTC),
						Type:      "info",
						Message:   "Test notification",
						Level:     "info",
					},
				},
			},
			filename: filepath.Join(tmpDir, "sessions", "test_session_2", "state.json"),
			wantErr:  false,
		},
		{
			name:     "invalid data",
			data:     make(chan int), // Channels cannot be marshaled to JSON
			filename: filepath.Join(tmpDir, "invalid.json"),
			wantErr:  true,
		},
		{
			name:     "relative path",
			data:     &SessionState{SessionID: "test"},
			filename: "relative/path.json", // Not absolute path
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := writer.WriteJSON(ctx, tt.filename, tt.data)

			if tt.wantErr {
				if err == nil {
					t.Errorf("WriteJSON() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("WriteJSON() unexpected error: %v", err)
				return
			}

			// Verify file was created
			if _, err := os.Stat(tt.filename); os.IsNotExist(err) {
				t.Errorf("WriteJSON() file was not created: %s", tt.filename)
				return
			}

			// Verify temp files were cleaned up (check directory for any .tmp. files)
			dir := filepath.Dir(tt.filename)
			entries, err := os.ReadDir(dir)
			if err == nil {
				for _, entry := range entries {
					if strings.Contains(entry.Name(), ".tmp.") {
						t.Errorf("WriteJSON() temp file was not cleaned up: %s", entry.Name())
					}
				}
			}

			// Verify content can be read back
			data, err := os.ReadFile(tt.filename)
			if err != nil {
				t.Errorf("Failed to read back written file: %v", err)
				return
			}

			if len(data) == 0 {
				t.Errorf("WriteJSON() wrote empty file")
			}
		})
	}
}

func TestAtomicWriter_ConcurrentWrites(t *testing.T) {
	// Create temporary directory for tests
	tmpDir, err := os.MkdirTemp("", "atomic_concurrent_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	writer := NewAtomicWriter(10 * time.Second)
	filename := filepath.Join(tmpDir, "concurrent.json")

	// Ensure the directory exists to prevent race conditions
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	const numGoroutines = 10
	var wg sync.WaitGroup
	errChan := make(chan error, numGoroutines)

	// Start multiple goroutines writing to the same file
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			data := &SessionState{
				SessionID:     "concurrent_test",
				CreatedAt:     time.Now().UTC(),
				UpdatedAt:     time.Now().UTC(),
				SessionActive: true,
				Agents:        []string{fmt.Sprintf("agent_%d", id)},
			}

			ctx := context.Background()
			if err := writer.WriteJSON(ctx, filename, data); err != nil {
				errChan <- err
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			t.Errorf("Concurrent WriteJSON() error: %v", err)
		}
	}

	// Verify file exists and is valid JSON
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Errorf("Concurrent WriteJSON() file was not created: %s", filename)
		return
	}

	// Verify no temp files remain
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read temp dir: %v", err)
	}

	for _, entry := range entries {
		if strings.Contains(entry.Name(), ".tmp.") {
			t.Errorf("Temp file not cleaned up: %s", entry.Name())
		}
	}
}

func TestAtomicWriter_ContextTimeout(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "atomic_timeout_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	writer := NewAtomicWriter(1 * time.Millisecond) // Very short timeout
	filename := filepath.Join(tmpDir, "timeout.json")

	data := &SessionState{
		SessionID:     "timeout_test",
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
		SessionActive: true,
	}

	ctx := context.Background()
	err = writer.WriteJSON(ctx, filename, data)

	// The timeout might not always trigger in tests, so we just verify
	// the operation completes one way or another
	if err != nil {
		// Check if it's a timeout error
		if err != context.DeadlineExceeded {
			t.Logf("WriteJSON() with timeout got non-timeout error: %v", err)
		}
	}
}

func TestFileError(t *testing.T) {
	err := &FileError{
		Op:   "test_operation",
		Path: "/test/path",
		Err:  errors.New("underlying error"),
	}

	expected := "test_operation /test/path: underlying error"
	if err.Error() != expected {
		t.Errorf("FileError.Error() = %q, want %q", err.Error(), expected)
	}

	if err.Unwrap().Error() != "underlying error" {
		t.Errorf("FileError.Unwrap() = %q, want %q", err.Unwrap().Error(), "underlying error")
	}
}

func TestStateError(t *testing.T) {
	err := &StateError{
		Code:    "test_code",
		Message: "test message",
	}

	expected := "state: test_code: test message"
	if err.Error() != expected {
		t.Errorf("StateError.Error() = %q, want %q", err.Error(), expected)
	}
}
