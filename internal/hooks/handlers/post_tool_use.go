package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dylan/spcstr/internal/hooks/events"
	"github.com/dylan/spcstr/internal/state"
)

// PostToolUseHandler handles the post_tool_use hook
// Reference .spcstr/logs/post_tool_use.json for event structure
// Updates .spcstr/sessions/{session-id}/state.json with extracted data
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
	var event events.ClaudeEvent
	if err := json.Unmarshal(input, &event); err != nil {
		return fmt.Errorf("failed to parse event: %w", err)
	}

	if event.SessionID == "" || event.ToolName == "" {
		return fmt.Errorf("missing required fields")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	stateManager := state.NewStateManager(filepath.Join(cwd, ".spcstr"))
	ctx := context.Background()

	switch event.ToolName {
	case "TodoWrite":
		var todoInput events.TodoWriteInput
		if err := json.Unmarshal(event.ToolInput, &todoInput); err == nil {
			todoState := state.TodoState{
				Total:       len(todoInput.Todos),
				Pending:     0,
				InProgress:  0,
				Completed:   0,
				Recent:      []state.TodoItem{},
				LastUpdated: time.Now().Format(time.RFC3339),
			}

			for i, todo := range todoInput.Todos {
				switch todo.Status {
				case "pending":
					todoState.Pending++
				case "in_progress":
					todoState.InProgress++
				case "completed":
					todoState.Completed++
				}

				if i < 5 {
					todoState.Recent = append(todoState.Recent, state.TodoItem{
						Content:    todo.Content,
						Status:     todo.Status,
						ActiveForm: todo.ActiveForm,
					})
				}
			}

			if err := stateManager.UpdateTodos(ctx, event.SessionID, todoState); err != nil {
				return fmt.Errorf("failed to update todos: %w", err)
			}
		}

	case "Write", "Edit":
		var fileInput events.FileOperationInput
		var fileResponse events.FileOperationResponse

		if err := json.Unmarshal(event.ToolInput, &fileInput); err == nil && fileInput.FilePath != "" {
			if err := json.Unmarshal(event.ToolResponse, &fileResponse); err == nil {
				opType := "edited"
				if fileResponse.Type == "create" {
					opType = "new"
				}
				stateManager.RecordFileOperation(ctx, event.SessionID, opType, fileInput.FilePath)
			}
		}

	case "MultiEdit":
		var multiEditInput events.MultiEditInput
		var fileResponse events.FileOperationResponse

		if err := json.Unmarshal(event.ToolInput, &multiEditInput); err == nil && multiEditInput.FilePath != "" {
			if err := json.Unmarshal(event.ToolResponse, &fileResponse); err == nil {
				opType := "edited"
				if fileResponse.Type == "create" {
					opType = "new"
				}
				stateManager.RecordFileOperation(ctx, event.SessionID, opType, multiEditInput.FilePath)
			}
		}

	case "Read":
		var fileInput events.FileOperationInput
		if err := json.Unmarshal(event.ToolInput, &fileInput); err == nil && fileInput.FilePath != "" {
			stateManager.RecordFileOperation(ctx, event.SessionID, "read", fileInput.FilePath)
		}

	case "Task":
		// Move agent from active to history when Task completes
		sessionState, err := stateManager.GetSessionState(ctx, event.SessionID)
		if err == nil && len(sessionState.Agents) > 0 {
			// Complete the first active agent (most recently added)
			stateManager.CompleteAgent(ctx, event.SessionID, sessionState.Agents[0])
		}
	}

	return nil
}