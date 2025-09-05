package hooks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// HookEvent represents a logged hook event
type HookEvent struct {
	Timestamp time.Time   `json:"timestamp"`
	SessionID string      `json:"session_id"`
	HookName  string      `json:"hook_name"`
	InputData interface{} `json:"input_data"`
	Success   bool        `json:"success"`
}

// HookLogger handles logging of hook events
type HookLogger struct {
	mu sync.Mutex
}

// NewHookLogger creates a new HookLogger instance
func NewHookLogger() *HookLogger {
	return &HookLogger{}
}

// LogEvent logs a hook event to the appropriate log file
func (l *HookLogger) LogEvent(sessionID, hookName string, inputData interface{}, success bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	event := HookEvent{
		Timestamp: time.Now(),
		SessionID: sessionID,
		HookName:  hookName,
		InputData: inputData,
		Success:   success,
	}

	logPath := filepath.Join(".spcstr", "logs", fmt.Sprintf("%s.json", hookName))
	
	// Create logs directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	return l.appendToLogFile(logPath, event)
}

// appendToLogFile appends an event to the log file in a thread-safe manner
func (l *HookLogger) appendToLogFile(logPath string, event HookEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create temporary file for atomic write
	tempFile := logPath + ".tmp"
	
	// Read existing events
	var existingEvents []HookEvent
	if data, err := os.ReadFile(logPath); err == nil {
		if err := json.Unmarshal(data, &existingEvents); err != nil {
			return fmt.Errorf("failed to parse existing log file: %w", err)
		}
	}

	// Append new event
	existingEvents = append(existingEvents, event)

	// Marshal to JSON
	data, err := json.MarshalIndent(existingEvents, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal log events: %w", err)
	}

	// Write to temp file with timeout
	done := make(chan error, 1)
	go func() {
		done <- os.WriteFile(tempFile, data, 0644)
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to write temp log file: %w", err)
		}
	case <-ctx.Done():
		return fmt.Errorf("log file write timeout")
	}

	// Atomic rename
	if err := os.Rename(tempFile, logPath); err != nil {
		os.Remove(tempFile) // Cleanup temp file on error
		return fmt.Errorf("failed to rename temp log file: %w", err)
	}

	return nil
}

// DefaultLogger is the global hook logger instance
var DefaultLogger = NewHookLogger()