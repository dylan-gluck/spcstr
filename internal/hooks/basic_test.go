package hooks

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/dylan/spcstr/internal/state"
)

// TestBasicHookExecution tests basic hook workflow without complex state management
func TestBasicHookExecution(t *testing.T) {
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

	// Initialize all handlers
	InitializeHandlers()

	sessionID := "basic_test_session"
	
	// Test session_start
	t.Run("session_start_basic", func(t *testing.T) {
		input := `{
			"session_id": "` + sessionID + `",
			"source": "basic_test"
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
		
		// Verify we can load the state
		oldDir, _ := os.Getwd()
		os.Chdir(projectDir)
		defer os.Chdir(oldDir)
		
		manager := state.NewStateManager(".spcstr")
		sessionState, err := manager.LoadState(context.Background(), sessionID)
		if err != nil {
			t.Errorf("Failed to load created state: %v", err)
		} else {
			if sessionState.SessionID != sessionID {
				t.Errorf("Expected session ID '%s', got '%s'", sessionID, sessionState.SessionID)
			}
			if !sessionState.SessionActive {
				t.Error("Session should be active")
			}
		}
	})

	// Test session_end
	t.Run("session_end_basic", func(t *testing.T) {
		input := `{
			"session_id": "` + sessionID + `"
		}`
		
		err := ExecuteHook("session_end", projectDir, []byte(input))
		if err != nil {
			t.Fatalf("session_end hook failed: %v", err)
		}
		
		// Verify session is now inactive
		oldDir, _ := os.Getwd()
		os.Chdir(projectDir)
		defer os.Chdir(oldDir)
		
		manager := state.NewStateManager(".spcstr")
		sessionState, err := manager.LoadState(context.Background(), sessionID)
		if err != nil {
			t.Errorf("Failed to load state after session_end: %v", err)
		} else {
			if sessionState.SessionActive {
				t.Error("Session should be inactive after session_end")
			}
		}
	})

	// Test log files were created
	t.Run("verify_log_files", func(t *testing.T) {
		expectedLogs := []string{
			"session_start.json",
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

func TestHookRegistry(t *testing.T) {
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