package hooks

import "github.com/dylan/spcstr/internal/hooks/handlers"

// InitializeHandlers registers all hook handlers with the default registry
func InitializeHandlers() {
	DefaultRegistry.Register(handlers.NewSessionStartHandler())
	DefaultRegistry.Register(handlers.NewUserPromptSubmitHandler())
	DefaultRegistry.Register(handlers.NewPreToolUseHandler())
	DefaultRegistry.Register(handlers.NewPostToolUseHandler())
	DefaultRegistry.Register(handlers.NewNotificationHandler())
	DefaultRegistry.Register(handlers.NewPreCompactHandler())
	DefaultRegistry.Register(handlers.NewSessionEndHandler())
	DefaultRegistry.Register(handlers.NewStopHandler())
	DefaultRegistry.Register(handlers.NewSubagentStopHandler())
}

func init() {
	InitializeHandlers()
}