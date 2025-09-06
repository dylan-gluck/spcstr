package events

import "encoding/json"

// ClaudeEvent represents the full structure of Claude hook events
type ClaudeEvent struct {
	SessionID       string          `json:"session_id"`
	HookEventName   string          `json:"hook_event_name"`
	ToolName        string          `json:"tool_name"`
	ToolInput       json.RawMessage `json:"tool_input"`
	ToolResponse    json.RawMessage `json:"tool_response"`
	PermissionMode  string          `json:"permission_mode"`
	TranscriptPath  string          `json:"transcript_path"`
	CWD            string          `json:"cwd"`
}

// TodoWriteInput for TodoWrite tool
type TodoWriteInput struct {
	Todos []TodoItem `json:"todos"`
}

type TodoItem struct {
	Content    string `json:"content"`
	Status     string `json:"status"`
	ActiveForm string `json:"activeForm"`
}

// FileOperationInput for Write/Edit/Read tools
type FileOperationInput struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content,omitempty"`
}

// MultiEditInput for MultiEdit tool
type MultiEditInput struct {
	FilePath string `json:"file_path"`
	Edits    []EditOperation `json:"edits"`
}

type EditOperation struct {
	OldString   string `json:"old_string"`
	NewString   string `json:"new_string"`
	ReplaceAll  bool   `json:"replace_all"`
}

// FileOperationResponse for file tools
type FileOperationResponse struct {
	FilePath string `json:"filePath"`
	Type     string `json:"type"` // "create" or "edit"
}

// TaskInput for Task tool
type TaskInput struct {
	Description   string `json:"description"`
	Prompt       string `json:"prompt"`
	SubagentType string `json:"subagent_type"`
}

// BashInput for Bash tool
type BashInput struct {
	Command         string `json:"command"`
	Description     string `json:"description"`
	RunInBackground bool   `json:"run_in_background"`
	Timeout         int    `json:"timeout"`
}

// NotificationParams for notification events
type NotificationParams struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
	Type      string `json:"type"`
}

// UserPromptSubmitParams for user prompt events
type UserPromptSubmitParams struct {
	SessionID string `json:"session_id"`
	Prompt    string `json:"prompt"`
}

// SessionStartParams for session start events
type SessionStartParams struct {
	SessionID string `json:"session_id"`
	CWD       string `json:"cwd"`
}

// SessionEndParams for session end events
type SessionEndParams struct {
	SessionID string `json:"session_id"`
}

// StopParams for stop events
type StopParams struct {
	SessionID string `json:"session_id"`
	Reason    string `json:"reason"`
}

// SubagentStopParams for subagent stop events
type SubagentStopParams struct {
	SessionID string `json:"session_id"`
	// Note: agent_name field doesn't exist in actual events
}

// PreCompactParams for pre-compact events
type PreCompactParams struct {
	SessionID string `json:"session_id"`
}