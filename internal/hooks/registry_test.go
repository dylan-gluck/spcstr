package hooks

import (
	"errors"
	"testing"
)

// mockHandler implements HookHandler for testing
type mockHandler struct {
	name    string
	execErr error
}

func (m *mockHandler) Name() string {
	return m.name
}

func (m *mockHandler) Execute(input []byte) error {
	return m.execErr
}

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()
	if registry == nil {
		t.Fatal("NewRegistry() returned nil")
	}
	if registry.handlers == nil {
		t.Fatal("Registry handlers map is nil")
	}
}

func TestRegistryRegister(t *testing.T) {
	registry := NewRegistry()
	handler := &mockHandler{name: "test_hook"}

	registry.Register(handler)

	retrieved, exists := registry.GetHandler("test_hook")
	if !exists {
		t.Fatal("Handler was not registered")
	}
	if retrieved.Name() != "test_hook" {
		t.Errorf("Expected handler name 'test_hook', got '%s'", retrieved.Name())
	}
}

func TestRegistryExecute(t *testing.T) {
	registry := NewRegistry()

	tests := []struct {
		name           string
		handler        *mockHandler
		hookName       string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:        "successful execution",
			handler:     &mockHandler{name: "success_hook", execErr: nil},
			hookName:    "success_hook",
			expectError: false,
		},
		{
			name:        "handler execution error",
			handler:     &mockHandler{name: "error_hook", execErr: errors.New("execution failed")},
			hookName:    "error_hook",
			expectError: true,
		},
		{
			name:           "hook not found",
			handler:        nil,
			hookName:       "nonexistent_hook",
			expectError:    true,
			expectedErrMsg: "hook 'nonexistent_hook' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.handler != nil {
				registry.Register(tt.handler)
			}

			err := registry.Execute(tt.hookName, []byte("{}"))

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, but got nil")
				}
				if tt.expectedErrMsg != "" && err.Error() != tt.expectedErrMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

func TestRegistryGetHandler(t *testing.T) {
	registry := NewRegistry()
	handler := &mockHandler{name: "test_handler"}

	// Test getting non-existent handler
	_, exists := registry.GetHandler("nonexistent")
	if exists {
		t.Error("Expected handler to not exist")
	}

	// Register and test getting existing handler
	registry.Register(handler)
	retrieved, exists := registry.GetHandler("test_handler")
	if !exists {
		t.Error("Expected handler to exist")
	}
	if retrieved.Name() != "test_handler" {
		t.Errorf("Expected handler name 'test_handler', got '%s'", retrieved.Name())
	}
}

func TestRegistryListHooks(t *testing.T) {
	registry := NewRegistry()

	// Test empty registry
	hooks := registry.ListHooks()
	if len(hooks) != 0 {
		t.Errorf("Expected 0 hooks, got %d", len(hooks))
	}

	// Register handlers
	handler1 := &mockHandler{name: "hook1"}
	handler2 := &mockHandler{name: "hook2"}
	registry.Register(handler1)
	registry.Register(handler2)

	hooks = registry.ListHooks()
	if len(hooks) != 2 {
		t.Errorf("Expected 2 hooks, got %d", len(hooks))
	}

	// Check that both hooks are present
	hookMap := make(map[string]bool)
	for _, hook := range hooks {
		hookMap[hook] = true
	}

	if !hookMap["hook1"] {
		t.Error("hook1 not found in list")
	}
	if !hookMap["hook2"] {
		t.Error("hook2 not found in list")
	}
}

func TestRegistryConcurrency(t *testing.T) {
	registry := NewRegistry()

	// Test concurrent registration and execution
	done := make(chan bool, 2)

	go func() {
		handler := &mockHandler{name: "concurrent_hook", execErr: nil}
		registry.Register(handler)
		done <- true
	}()

	go func() {
		// This might fail if handler isn't registered yet, but shouldn't crash
		registry.Execute("concurrent_hook", []byte("{}"))
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Verify the handler was registered
	_, exists := registry.GetHandler("concurrent_hook")
	if !exists {
		t.Error("Handler was not registered in concurrent test")
	}
}
