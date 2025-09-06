package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dylan/spcstr/internal/state"
)

// SessionEndParams defines the expected input for session_end hook
type SessionEndParams struct {
	SessionID string `json:"session_id"`
}

// SessionEndHandler handles the session_end hook
type SessionEndHandler struct{}

// NewSessionEndHandler creates a new SessionEndHandler
func NewSessionEndHandler() *SessionEndHandler {
	return &SessionEndHandler{}
}

// Name returns the hook name
func (h *SessionEndHandler) Name() string {
	return "session_end"
}

// Execute processes the session_end hook
func (h *SessionEndHandler) Execute(input []byte) error {
	var params SessionEndParams
	if err := json.Unmarshal(input, &params); err != nil {
		return fmt.Errorf("failed to parse session_end parameters: %w", err)
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
