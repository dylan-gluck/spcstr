package handlers

import (
	"encoding/json"
	"fmt"
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

	// Stop hook is called when Claude finishes a response turn
	// The session remains active - it only becomes inactive on session_end
	// For now, we don't need to do anything here, but we keep the handler
	// to acknowledge the hook was received successfully
	
	return nil
}
