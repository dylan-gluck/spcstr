package hooks

import (
	"fmt"
	"sync"
)

// HookHandler interface defines the contract for all hook handlers
type HookHandler interface {
	Name() string
	Execute(input []byte) error
}

// HookRegistry manages all registered hook handlers
type HookRegistry struct {
	handlers map[string]HookHandler
	mu       sync.RWMutex
}

// NewRegistry creates a new HookRegistry instance
func NewRegistry() *HookRegistry {
	return &HookRegistry{
		handlers: make(map[string]HookHandler),
	}
}

// Register adds a hook handler to the registry
func (r *HookRegistry) Register(handler HookHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[handler.Name()] = handler
}

// Execute runs the specified hook with the provided input
func (r *HookRegistry) Execute(name string, input []byte) error {
	r.mu.RLock()
	handler, exists := r.handlers[name]
	r.mu.RUnlock()

	if !exists {
		return fmt.Errorf("hook '%s' not found", name)
	}

	return handler.Execute(input)
}

// GetHandler retrieves a handler by name
func (r *HookRegistry) GetHandler(name string) (HookHandler, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, exists := r.handlers[name]
	return handler, exists
}

// ListHooks returns all registered hook names
func (r *HookRegistry) ListHooks() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	hooks := make([]string, 0, len(r.handlers))
	for name := range r.handlers {
		hooks = append(hooks, name)
	}
	return hooks
}

// DefaultRegistry is the global hook registry instance
var DefaultRegistry = NewRegistry()