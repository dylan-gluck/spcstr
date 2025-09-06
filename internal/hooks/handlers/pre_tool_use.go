package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dylan/spcstr/internal/hooks/events"
	"github.com/dylan/spcstr/internal/state"
)

// PreToolUseHandler handles the pre_tool_use hook
// Reference .spcstr/logs/pre_tool_use.json for event structure
// Updates .spcstr/sessions/{session-id}/state.json with tool usage and agent info
type PreToolUseHandler struct{}

// NewPreToolUseHandler creates a new PreToolUseHandler
func NewPreToolUseHandler() *PreToolUseHandler {
	return &PreToolUseHandler{}
}

// Name returns the hook name
func (h *PreToolUseHandler) Name() string {
	return "pre_tool_use"
}

// Execute processes the pre_tool_use hook
func (h *PreToolUseHandler) Execute(input []byte) error {
	var event events.ClaudeEvent
	if err := json.Unmarshal(input, &event); err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	if event.SessionID == "" || event.ToolName == "" {
		return fmt.Errorf("missing required fields")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	stateManager := state.NewStateManager(filepath.Join(cwd, ".spcstr"))
	ctx := context.Background()

	// Always increment tool usage
	if err := stateManager.IncrementToolUsage(ctx, event.SessionID, event.ToolName); err != nil {
		return fmt.Errorf("failed to increment tool usage: %w", err)
	}

	// Handle Task tool to extract agent info
	if event.ToolName == "Task" {
		var taskInput events.TaskInput
		if err := json.Unmarshal(event.ToolInput, &taskInput); err == nil && taskInput.SubagentType != "" {
			if err := stateManager.AddAgent(ctx, event.SessionID, taskInput.SubagentType); err != nil {
				return fmt.Errorf("failed to add agent: %w", err)
			}
		}
	}

	// Never block operations - let Claude handle permissions
	return nil
}