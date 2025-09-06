package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dylan/spcstr/internal/state"
)

// PostToolUseParams defines the expected input for post_tool_use hook
type PostToolUseParams struct {
	SessionID    string   `json:"session_id"`
	ToolName     string   `json:"tool_name"`
	AgentName    string   `json:"agent_name,omitempty"`
	FilesCreated []string `json:"files_created,omitempty"`
	FilesEdited  []string `json:"files_edited,omitempty"`
	FilesRead    []string `json:"files_read,omitempty"`
}

// PostToolUseHandler handles the post_tool_use hook
type PostToolUseHandler struct{}

// NewPostToolUseHandler creates a new PostToolUseHandler
func NewPostToolUseHandler() *PostToolUseHandler {
	return &PostToolUseHandler{}
}

// Name returns the hook name
func (h *PostToolUseHandler) Name() string {
	return "post_tool_use"
}

// Execute processes the post_tool_use hook
func (h *PostToolUseHandler) Execute(input []byte) error {
	var params PostToolUseParams
	if err := json.Unmarshal(input, &params); err != nil {
		return fmt.Errorf("failed to parse post_tool_use parameters: %w", err)
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

	// Track file operations
	for _, file := range params.FilesCreated {
		if err := stateManager.RecordFileOperation(ctx, params.SessionID, "new", file); err != nil {
			return fmt.Errorf("failed to record file creation: %w", err)
		}
	}
	for _, file := range params.FilesEdited {
		if err := stateManager.RecordFileOperation(ctx, params.SessionID, "edited", file); err != nil {
			return fmt.Errorf("failed to record file edit: %w", err)
		}
	}
	for _, file := range params.FilesRead {
		if err := stateManager.RecordFileOperation(ctx, params.SessionID, "read", file); err != nil {
			return fmt.Errorf("failed to record file read: %w", err)
		}
	}

	// Handle Task tool completion (agent management)
	if params.ToolName == "Task" && params.AgentName != "" {
		// Move agent from active to history
		if err := stateManager.CompleteAgent(ctx, params.SessionID, params.AgentName); err != nil {
			return fmt.Errorf("failed to complete agent: %w", err)
		}
	}

	return nil
}
