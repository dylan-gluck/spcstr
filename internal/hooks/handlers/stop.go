package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dylan/spcstr/internal/state"
)

// StopParams defines the expected input for stop hook
type StopParams struct {
	SessionID string `json:"session_id"`
}

// StopHandler handles the stop hook
type StopHandler struct{}

// NewStopHandler creates a new StopHandler
func NewStopHandler() *StopHandler {
	return &StopHandler{}
}

// Name returns the hook name
func (h *StopHandler) Name() string {
	return "stop"
}

// Execute processes the stop hook
func (h *StopHandler) Execute(input []byte) error {
	var params StopParams
	if err := json.Unmarshal(input, &params); err != nil {
		return fmt.Errorf("failed to parse stop parameters: %w", err)
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

	// Set session as inactive
	ctx := context.Background()
	if err := stateManager.SetSessionActive(ctx, params.SessionID, false); err != nil {
		return fmt.Errorf("failed to set session inactive: %w", err)
	}

	return nil
}
