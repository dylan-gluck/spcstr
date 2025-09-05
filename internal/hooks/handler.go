package hooks

import (
	"log/slog"

	"github.com/dylan-gluck/spcstr/internal/hooks/handlers"
)

// Handler is the main hook event dispatcher
type Handler struct {
	projectRoot string
	sessionDir  string
}

// NewHandler creates a new hook handler
func NewHandler(projectRoot string) *Handler {
	return &Handler{
		projectRoot: projectRoot,
	}
}

// HandleHook is the main entry point for processing hook events
func (h *Handler) HandleHook(hookName string, data map[string]interface{}) error {
	slog.Debug("handling hook", "hook", hookName)

	var err error
	switch hookName {
	case "pre_tool_use":
		err = handlers.HandlePreToolUse(h.projectRoot, data)
	case "post_tool_use":
		err = handlers.HandlePostToolUse(h.projectRoot, data)
	case "session_start":
		err = handlers.HandleSessionStart(h.projectRoot, data)
	case "session_end":
		err = handlers.HandleSessionEnd(h.projectRoot, data)
	case "user_prompt_submit":
		err = handlers.HandleUserPrompt(h.projectRoot, data)
	case "notification":
		err = handlers.HandleNotification(h.projectRoot, data)
	case "stop":
		err = handlers.HandleStop(h.projectRoot, data)
	case "subagent_start":
		err = handlers.HandleSubagentStart(h.projectRoot, data)
	case "subagent_stop":
		err = handlers.HandleSubagentStop(h.projectRoot, data)
	case "pre_compact":
		err = handlers.HandlePreCompact(h.projectRoot, data)
	default:
		slog.Warn("unknown hook", "hook", hookName)
		return nil
	}

	if err != nil {
		LogStructuredError(hookName, err, data)
		return err
	}

	return nil
}

func (h *Handler) getSessionDir() string {
	if h.sessionDir != "" {
		return h.sessionDir
	}
	return h.projectRoot
}
