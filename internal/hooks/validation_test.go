package hooks

import (
	"testing"
)

// TestAllHooksRegistered validates that all 9 required hooks are registered
func TestAllHooksRegistered(t *testing.T) {
	requiredHooks := []string{
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

	// Verify each hook is registered
	for _, hookName := range requiredHooks {
		handler, exists := DefaultRegistry.GetHandler(hookName)
		if !exists {
			t.Errorf("Required hook '%s' is not registered", hookName)
			continue
		}

		if handler.Name() != hookName {
			t.Errorf("Handler name mismatch: expected '%s', got '%s'", hookName, handler.Name())
		}
	}

	// Verify we have exactly the expected number of hooks
	allHooks := DefaultRegistry.ListHooks()
	if len(allHooks) != len(requiredHooks) {
		t.Errorf("Expected %d hooks, got %d. Registered hooks: %v",
			len(requiredHooks), len(allHooks), allHooks)
	}
}

// TestJSONValidation tests that all handlers properly validate required JSON fields
func TestJSONValidation(t *testing.T) {
	tests := []struct {
		hookName    string
		validJSON   string
		invalidJSON string
	}{
		{
			hookName:    "session_start",
			validJSON:   `{"session_id": "test123", "source": "test"}`,
			invalidJSON: `{"session_id": ""}`,
		},
		{
			hookName:    "user_prompt_submit",
			validJSON:   `{"session_id": "test123", "prompt": "hello"}`,
			invalidJSON: `{"session_id": "test123"}`,
		},
		{
			hookName:    "pre_tool_use",
			validJSON:   `{"session_id": "test123", "tool_name": "Task"}`,
			invalidJSON: `{"session_id": "test123"}`,
		},
		{
			hookName:    "post_tool_use",
			validJSON:   `{"session_id": "test123", "tool_name": "Task"}`,
			invalidJSON: `{"session_id": ""}`,
		},
		{
			hookName:    "notification",
			validJSON:   `{"session_id": "test123", "message": "test"}`,
			invalidJSON: `{"session_id": "test123"}`,
		},
		{
			hookName:    "session_end",
			validJSON:   `{"session_id": "test123"}`,
			invalidJSON: `{"session_id": ""}`,
		},
		{
			hookName:    "stop",
			validJSON:   `{"session_id": "test123"}`,
			invalidJSON: `{"session_id": ""}`,
		},
		{
			hookName:    "subagent_stop",
			validJSON:   `{"session_id": "test123", "agent_name": "test_agent"}`,
			invalidJSON: `{"session_id": "test123"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.hookName, func(t *testing.T) {
			handler, exists := DefaultRegistry.GetHandler(tt.hookName)
			if !exists {
				t.Fatalf("Hook '%s' not found", tt.hookName)
			}

			// Test that malformed JSON fails appropriately
			err := handler.Execute([]byte("invalid json"))
			if err == nil {
				t.Error("Expected error for invalid JSON, got nil")
			}

			// Test that invalid field validation works
			err = handler.Execute([]byte(tt.invalidJSON))
			if err == nil {
				t.Error("Expected error for invalid fields, got nil")
			}

			// Note: We don't test valid JSON execution here as it would require
			// proper state setup and working directories
		})
	}
}

// TestHookNaming validates handler naming consistency
func TestHookNaming(t *testing.T) {
	allHooks := DefaultRegistry.ListHooks()

	for _, hookName := range allHooks {
		handler, _ := DefaultRegistry.GetHandler(hookName)

		if handler.Name() != hookName {
			t.Errorf("Handler for '%s' returns wrong name: '%s'", hookName, handler.Name())
		}

		// Validate snake_case naming convention
		for i, char := range hookName {
			if char >= 'A' && char <= 'Z' {
				t.Errorf("Hook name '%s' contains uppercase at position %d, should be snake_case", hookName, i)
			}
		}
	}
}

// TestExecutorValidation tests basic executor functionality
func TestExecutorValidation(t *testing.T) {
	// Test invalid project directory
	err := ExecuteHook("session_start", "/nonexistent/path", []byte(`{"session_id": "test"}`))
	if err == nil {
		t.Error("Expected error for nonexistent project directory")
	}

	// Test nonexistent hook
	tempDir := t.TempDir()
	err = ExecuteHook("nonexistent_hook", tempDir, []byte(`{"session_id": "test"}`))
	if err == nil {
		t.Error("Expected error for nonexistent hook")
	}
}
