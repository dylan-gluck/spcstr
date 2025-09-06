package hooks

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dylan/spcstr/internal/state"
)

// TestHookWorkflowIntegration tests the complete hook workflow
func TestHookWorkflowIntegration(t *testing.T) {
	// Create temporary project structure
	projectDir := t.TempDir()
	spcstrDir := filepath.Join(projectDir, ".spcstr")
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

	// Initialize all handlers (this should happen automatically via init())
	InitializeHandlers()

	sessionID := "integration_test_session"

	// Test 1: Session Start
	t.Run("session_start", func(t *testing.T) {
		input := `{
			"session_id": "` + sessionID + `",
			"source": "integration_test"
		}`

		err := ExecuteHook("session_start", projectDir, []byte(input))
		if err != nil {
			t.Fatalf("session_start hook failed: %v", err)
		}

		// Verify session state was created
		statePath := filepath.Join(projectDir, ".spcstr", "sessions", sessionID, "state.json")
		if _, err := os.Stat(statePath); os.IsNotExist(err) {
			t.Error("Session state file was not created")
		}
	})

	// Test 2: User Prompt Submit
	t.Run("user_prompt_submit", func(t *testing.T) {
		input := `{
			"session_id": "` + sessionID + `",
			"prompt": "Hello, this is an integration test!",
			"timestamp": "2023-06-15T14:30:00Z"
		}`

		err := ExecuteHook("user_prompt_submit", projectDir, []byte(input))
		if err != nil {
			t.Fatalf("user_prompt_submit hook failed: %v", err)
		}
	})

	// Test 3: Pre Tool Use
	t.Run("pre_tool_use", func(t *testing.T) {
		input := `{
			"session_id": "` + sessionID + `",
			"tool_name": "Task",
			"tool_input": {
				"description": "Test task",
				"prompt": "Test prompt",
				"subagent_type": "test_agent"
			}
		}`

		err := ExecuteHook("pre_tool_use", projectDir, []byte(input))
		if err != nil {
			t.Fatalf("pre_tool_use hook failed: %v", err)
		}
	})

	// Test 4: Post Tool Use for file operations
	t.Run("post_tool_use_write", func(t *testing.T) {
		input := `{
			"session_id": "` + sessionID + `",
			"tool_name": "Write",
			"tool_input": {
				"file_path": "test.go"
			},
			"tool_response": {
				"filePath": "test.go",
				"type": "create"
			}
		}`

		err := ExecuteHook("post_tool_use", projectDir, []byte(input))
		if err != nil {
			t.Fatalf("post_tool_use hook failed: %v", err)
		}
	})

	// Test 4b: Post Tool Use for Edit
	t.Run("post_tool_use_edit", func(t *testing.T) {
		input := `{
			"session_id": "` + sessionID + `",
			"tool_name": "Edit",
			"tool_input": {
				"file_path": "main.go"
			},
			"tool_response": {
				"filePath": "main.go",
				"type": "edit"
			}
		}`

		err := ExecuteHook("post_tool_use", projectDir, []byte(input))
		if err != nil {
			t.Fatalf("post_tool_use hook failed: %v", err)
		}
	})

	// Test 4c: Post Tool Use for Read
	t.Run("post_tool_use_read", func(t *testing.T) {
		input := `{
			"session_id": "` + sessionID + `",
			"tool_name": "Read",
			"tool_input": {
				"file_path": "config.json"
			}
		}`

		err := ExecuteHook("post_tool_use", projectDir, []byte(input))
		if err != nil {
			t.Fatalf("post_tool_use hook failed: %v", err)
		}
	})

	// Test 4d: Post Tool Use for Task completion
	t.Run("post_tool_use_task", func(t *testing.T) {
		input := `{
			"session_id": "` + sessionID + `",
			"tool_name": "Task",
			"tool_response": {}
		}`

		err := ExecuteHook("post_tool_use", projectDir, []byte(input))
		if err != nil {
			t.Fatalf("post_tool_use hook failed: %v", err)
		}
	})

	// Test 5: Notification
	t.Run("notification", func(t *testing.T) {
		input := `{
			"session_id": "` + sessionID + `",
			"message": "Integration test notification",
			"level": "info",
			"timestamp": "2023-06-15T14:35:00Z"
		}`

		err := ExecuteHook("notification", projectDir, []byte(input))
		if err != nil {
			t.Fatalf("notification hook failed: %v", err)
		}
	})

	// Test 6: Session End
	t.Run("session_end", func(t *testing.T) {
		input := `{
			"session_id": "` + sessionID + `"
		}`

		err := ExecuteHook("session_end", projectDir, []byte(input))
		if err != nil {
			t.Fatalf("session_end hook failed: %v", err)
		}
	})

	// Verify final state
	t.Run("verify_final_state", func(t *testing.T) {
		// Change to project directory to load state correctly
		oldDir, _ := os.Getwd()
		os.Chdir(projectDir)
		defer os.Chdir(oldDir)

		manager := state.NewStateManager(".spcstr")
		finalState, err := manager.LoadState(context.Background(), sessionID)
		if err != nil {
			t.Fatalf("Failed to load final state: %v", err)
		}

		// Verify session is inactive
		if finalState.SessionActive {
			t.Error("Session should be inactive after session_end")
		}

		// Verify prompt was added
		if len(finalState.Prompts) != 1 {
			t.Errorf("Expected 1 prompt, got %d", len(finalState.Prompts))
		} else if finalState.Prompts[0].Prompt != "Hello, this is an integration test!" {
			t.Errorf("Unexpected prompt content: %s", finalState.Prompts[0].Prompt)
		}

		// Verify notification was added
		if len(finalState.Notifications) != 1 {
			t.Errorf("Expected 1 notification, got %d", len(finalState.Notifications))
		} else if finalState.Notifications[0].Message != "Integration test notification" {
			t.Errorf("Unexpected notification message: %s", finalState.Notifications[0].Message)
		}

		// Verify file operations were tracked
		if len(finalState.Files.New) != 1 || finalState.Files.New[0] != "test.go" {
			t.Errorf("Unexpected files created: %v", finalState.Files.New)
		}
		if len(finalState.Files.Edited) != 1 || finalState.Files.Edited[0] != "main.go" {
			t.Errorf("Unexpected files edited: %v", finalState.Files.Edited)
		}
		if len(finalState.Files.Read) != 1 || finalState.Files.Read[0] != "config.json" {
			t.Errorf("Unexpected files read: %v", finalState.Files.Read)
		}

		// Verify agent was moved to history
		if len(finalState.Agents) != 0 {
			t.Errorf("Expected 0 active agents, got %d", len(finalState.Agents))
		}
		if len(finalState.AgentsHistory) != 1 {
			t.Errorf("Expected 1 agent in history, got %d", len(finalState.AgentsHistory))
		} else if finalState.AgentsHistory[0].Name != "test_agent" {
			t.Errorf("Unexpected agent name in history: %s", finalState.AgentsHistory[0].Name)
		}
	})

	// Verify all log files were created
	t.Run("verify_log_files", func(t *testing.T) {
		expectedLogs := []string{
			"session_start.json",
			"user_prompt_submit.json",
			"pre_tool_use.json",
			"post_tool_use.json",
			"notification.json",
			"session_end.json",
		}

		for _, logFile := range expectedLogs {
			logPath := filepath.Join(logsDir, logFile)
			if _, err := os.Stat(logPath); os.IsNotExist(err) {
				t.Errorf("Log file %s was not created", logFile)
			}
		}
	})
}

func TestHookCommandRouting(t *testing.T) {
	// Test that all expected hooks are registered
	expectedHooks := []string{
		"session_start",
		"user_prompt_submit",
		"pre_tool_use",
		"post_tool_use",
		"notification",
		"pre_compact",
		"session_end",
		"stop",
		"subagent_stop",
	}

	for _, hookName := range expectedHooks {
		t.Run(hookName, func(t *testing.T) {
			handler, exists := DefaultRegistry.GetHandler(hookName)
			if !exists {
				t.Errorf("Hook '%s' is not registered", hookName)
				return
			}

			if handler.Name() != hookName {
				t.Errorf("Handler name mismatch: expected '%s', got '%s'", hookName, handler.Name())
			}
		})
	}

	// Verify total number of registered hooks
	allHooks := DefaultRegistry.ListHooks()
	if len(allHooks) != len(expectedHooks) {
		t.Errorf("Expected %d registered hooks, got %d", len(expectedHooks), len(allHooks))
	}
}
