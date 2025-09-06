package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dylan/spcstr/internal/state"
)

// UserPromptSubmitParams defines the expected input for user_prompt_submit hook
type UserPromptSubmitParams struct {
	SessionID string `json:"session_id"`
	Prompt    string `json:"prompt"`
	Timestamp string `json:"timestamp"`
}

// UserPromptSubmitHandler handles the user_prompt_submit hook
type UserPromptSubmitHandler struct{}

// NewUserPromptSubmitHandler creates a new UserPromptSubmitHandler
func NewUserPromptSubmitHandler() *UserPromptSubmitHandler {
	return &UserPromptSubmitHandler{}
}

// Name returns the hook name
func (h *UserPromptSubmitHandler) Name() string {
	return "user_prompt_submit"
}

// Execute processes the user_prompt_submit hook
func (h *UserPromptSubmitHandler) Execute(input []byte) error {
	var params UserPromptSubmitParams
	if err := json.Unmarshal(input, &params); err != nil {
		return fmt.Errorf("failed to parse user_prompt_submit parameters: %w", err)
	}

	// Validate required fields
	if params.SessionID == "" {
		return fmt.Errorf("session_id is required")
	}
	if params.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}

	// Parse timestamp or use current time
	var promptTime time.Time
	if params.Timestamp != "" {
		var err error
		promptTime, err = time.Parse(time.RFC3339, params.Timestamp)
		if err != nil {
			promptTime = time.Now()
		}
	} else {
		promptTime = time.Now()
	}

	// Create StateManager using current working directory (after --cwd change)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	stateManager := state.NewStateManager(filepath.Join(cwd, ".spcstr"))

	// Add prompt to session state using UpdateState
	ctx := context.Background()
	err = stateManager.UpdateState(ctx, params.SessionID, func(sessionState *state.SessionState) error {
		// Create new prompt entry
		newPrompt := state.PromptEntry{
			Timestamp: promptTime,
			Prompt:    params.Prompt,
			Response:  "", // Will be filled by response hooks if needed
			ToolsUsed: []string{},
		}

		sessionState.Prompts = append(sessionState.Prompts, newPrompt)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to add prompt to session: %w", err)
	}

	return nil
}
