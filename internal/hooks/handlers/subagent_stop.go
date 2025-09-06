package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dylan/spcstr/internal/state"
)

// SubagentStopParams defines the expected input for subagent_stop hook
type SubagentStopParams struct {
	SessionID string `json:"session_id"`
	AgentName string `json:"agent_name"`
}

// SubagentStopHandler handles the subagent_stop hook
type SubagentStopHandler struct{}

// NewSubagentStopHandler creates a new SubagentStopHandler
func NewSubagentStopHandler() *SubagentStopHandler {
	return &SubagentStopHandler{}
}

// Name returns the hook name
func (h *SubagentStopHandler) Name() string {
	return "subagent_stop"
}

// Execute processes the subagent_stop hook
func (h *SubagentStopHandler) Execute(input []byte) error {
	var params SubagentStopParams
	if err := json.Unmarshal(input, &params); err != nil {
		return fmt.Errorf("failed to parse subagent_stop parameters: %w", err)
	}

	// Validate required fields
	if params.SessionID == "" {
		return fmt.Errorf("session_id is required")
	}
	
	// If agent_name is not provided, use a default value
	// This handles cases where Claude Code doesn't provide the agent_name parameter
	if params.AgentName == "" {
		params.AgentName = "claude"
	}

	// Create StateManager using current working directory (after --cwd change)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	stateManager := state.NewStateManager(filepath.Join(cwd, ".spcstr"))

	// Complete the specified agent
	ctx := context.Background()
	if err := stateManager.CompleteAgent(ctx, params.SessionID, params.AgentName); err != nil {
		return fmt.Errorf("failed to complete agent: %w", err)
	}

	return nil
}
