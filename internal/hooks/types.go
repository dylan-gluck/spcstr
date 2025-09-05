package hooks

import (
	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

// Type aliases for public models
type SessionState = models.SessionState
type AgentHistoryEntry = models.AgentHistoryEntry
type FileOperations = models.FileOperations
type ErrorEntry = models.ErrorEntry

// HookHandler interface for hook processors
type HookHandler interface {
	HandleHook(hookName string, data map[string]interface{}) error
}

// Expose utility functions from models
var CurrentTimestamp = models.CurrentTimestamp
var AddUniqueString = models.AddUniqueString
