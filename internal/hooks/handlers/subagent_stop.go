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

// SubagentStopHandler handles the subagent_stop hook
// Reference .spcstr/logs/subagent_stop.json for event structure
// Note: agent_name field doesn't exist in actual events
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
	var params events.SubagentStopParams
	if err := json.Unmarshal(input, &params); err != nil {
		return fmt.Errorf("failed to parse subagent_stop parameters: %w", err)
	}

	// Validate required fields
	if params.SessionID == "" {
		return fmt.Errorf("session_id is required")
	}

	// Create StateManager using current working directory (after --cwd change)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	stateManager := state.NewStateManager(filepath.Join(cwd, ".spcstr"))
	ctx := context.Background()

	// Get current state to find active agent
	sessionState, err := stateManager.GetSessionState(ctx, params.SessionID)
	if err != nil {
		// Session might not exist yet, that's ok
		return nil
	}

	// Complete the most recently added agent if any
	if len(sessionState.Agents) > 0 {
		if err := stateManager.CompleteAgent(ctx, params.SessionID, sessionState.Agents[0]); err != nil {
			// Don't fail if we can't complete the agent
			// This is an observability tool, not a blocker
			return nil
		}
	}

	return nil
}