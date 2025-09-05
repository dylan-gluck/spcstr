package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dylan/spcstr/internal/state"
)

// SessionStartParams defines the expected input for session_start hook
type SessionStartParams struct {
	SessionID string `json:"session_id"`
	Source    string `json:"source"`
}

// SessionStartHandler handles the session_start hook
type SessionStartHandler struct{}

// NewSessionStartHandler creates a new SessionStartHandler
func NewSessionStartHandler() *SessionStartHandler {
	return &SessionStartHandler{}
}

// Name returns the hook name
func (h *SessionStartHandler) Name() string {
	return "session_start"
}

// Execute processes the session_start hook
func (h *SessionStartHandler) Execute(input []byte) error {
	var params SessionStartParams
	if err := json.Unmarshal(input, &params); err != nil {
		return fmt.Errorf("failed to parse session_start parameters: %w", err)
	}

	// Validate required fields
	if params.SessionID == "" {
		return fmt.Errorf("session_id is required")
	}
	if params.Source == "" {
		return fmt.Errorf("source is required")
	}

	// Create StateManager using current working directory (after --cwd change)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}
	
	stateManager := state.NewStateManager(filepath.Join(cwd, ".spcstr"))

	// Initialize state using StateManager
	ctx := context.Background()
	if _, err := stateManager.InitializeState(ctx, params.SessionID); err != nil {
		return fmt.Errorf("failed to initialize session state: %w", err)
	}

	return nil
}