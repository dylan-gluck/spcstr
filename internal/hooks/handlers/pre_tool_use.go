package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dylan/spcstr/internal/state"
)

// PreToolUseParams defines the expected input for pre_tool_use hook
type PreToolUseParams struct {
	SessionID string `json:"session_id"`
	ToolName  string `json:"tool_name"`
	AgentName string `json:"agent_name,omitempty"`
}

// PreToolUseHandler handles the pre_tool_use hook
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
	var params PreToolUseParams
	if err := json.Unmarshal(input, &params); err != nil {
		return fmt.Errorf("failed to parse pre_tool_use parameters: %w", err)
	}

	// Validate required fields
	if params.SessionID == "" {
		return fmt.Errorf("session_id is required")
	}
	if params.ToolName == "" {
		return fmt.Errorf("tool_name is required")
	}

	// Create StateManager using current working directory (after --cwd change)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	
	stateManager := state.NewStateManager(filepath.Join(cwd, ".spcstr"))

	ctx := context.Background()

	// Increment tool usage
	if err := stateManager.IncrementToolUsage(ctx, params.SessionID, params.ToolName); err != nil {
		return fmt.Errorf("failed to increment tool usage: %w", err)
	}

	// Handle Task tool invocations (agent management)
	if params.ToolName == "Task" && params.AgentName != "" {
		// Add agent to active agents
		if err := stateManager.AddAgent(ctx, params.SessionID, params.AgentName); err != nil {
			return fmt.Errorf("failed to add agent to session: %w", err)
		}
	}

	return nil
}