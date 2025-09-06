package handlers

import (
	"encoding/json"
	"fmt"
)

// PreCompactParams defines the expected input for pre_compact hook
type PreCompactParams struct {
	SessionID string `json:"session_id"`
}

// PreCompactHandler handles the pre_compact hook
type PreCompactHandler struct{}

// NewPreCompactHandler creates a new PreCompactHandler
func NewPreCompactHandler() *PreCompactHandler {
	return &PreCompactHandler{}
}

// Name returns the hook name
func (h *PreCompactHandler) Name() string {
	return "pre_compact"
}

// Execute processes the pre_compact hook
func (h *PreCompactHandler) Execute(input []byte) error {
	var params PreCompactParams
	if err := json.Unmarshal(input, &params); err != nil {
		return fmt.Errorf("failed to parse pre_compact parameters: %w", err)
	}

	// Validate required fields
	if params.SessionID == "" {
		return fmt.Errorf("session_id is required")
	}

	// This hook is mainly for observability, no state changes required
	// Could be used for compaction metrics in the future

	return nil
}
