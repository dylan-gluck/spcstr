package state

import (
	"time"
)

// SessionState represents the complete state of a Claude Code session
type SessionState struct {
	SessionID     string              `json:"session_id"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	SessionActive bool                `json:"session_active"`
	Agents        []string            `json:"agents"`
	AgentsHistory []AgentExecution    `json:"agents_history"`
	Files         FileOperations      `json:"files"`
	ToolsUsed     map[string]int      `json:"tools_used"`
	Errors        []ErrorEntry        `json:"errors"`
	Prompts       []PromptEntry       `json:"prompts"`
	Notifications []NotificationEntry `json:"notifications"`
	Todos         TodoState           `json:"todos"`
}

// AgentExecution tracks the execution lifecycle of an agent
type AgentExecution struct {
	Name        string     `json:"name"`
	StartedAt   time.Time  `json:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// FileOperations tracks all file operations during a session
type FileOperations struct {
	New    []string `json:"new"`
	Edited []string `json:"edited"`
	Read   []string `json:"read"`
}

// ErrorEntry represents an error that occurred during the session
type ErrorEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	Severity  string    `json:"severity"`
}

// PromptEntry tracks user prompts and responses
type PromptEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Prompt    string    `json:"prompt"`
	Response  string    `json:"response"`
	ToolsUsed []string  `json:"tools_used"`
}

// NotificationEntry tracks system notifications
type NotificationEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
}

// TodoState tracks todo items and their status
type TodoState struct {
	Total       int        `json:"total"`
	Pending     int        `json:"pending"`
	InProgress  int        `json:"in_progress"`
	Completed   int        `json:"completed"`
	Recent      []TodoItem `json:"recent"`
	LastUpdated string     `json:"last_updated"`
}

// TodoItem represents a single todo item
type TodoItem struct {
	Content    string `json:"content"`
	Status     string `json:"status"`
	ActiveForm string `json:"activeForm"`
}
