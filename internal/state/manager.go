package state

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// DefaultTimeout for file operations
	DefaultTimeout = 5 * time.Second
	// StateFileName is the standard name for session state files
	StateFileName = "state.json"
)

// StateManager provides CRUD operations for session state
type StateManager struct {
	writer   *AtomicWriter
	basePath string
	timeout  time.Duration
}

// NewStateManager creates a new StateManager with the specified base path
func NewStateManager(basePath string) *StateManager {
	return &StateManager{
		writer:   NewAtomicWriter(DefaultTimeout),
		basePath: basePath,
		timeout:  DefaultTimeout,
	}
}

// NewStateManagerWithTimeout creates a StateManager with custom timeout
func NewStateManagerWithTimeout(basePath string, timeout time.Duration) *StateManager {
	return &StateManager{
		writer:   NewAtomicWriter(timeout),
		basePath: basePath,
		timeout:  timeout,
	}
}

// InitializeState creates a new session state with proper directory structure
func (sm *StateManager) InitializeState(ctx context.Context, sessionID string) (*SessionState, error) {
	if sessionID == "" {
		return nil, &StateError{
			Code:    "invalid_session_id",
			Message: "session ID cannot be empty",
		}
	}

	// Create session state with current timestamp
	now := time.Now().UTC()
	state := &SessionState{
		SessionID:     sessionID,
		CreatedAt:     now,
		UpdatedAt:     now,
		SessionActive: true,
		Agents:        make([]string, 0),
		AgentsHistory: make([]AgentExecution, 0),
		Files: FileOperations{
			New:    make([]string, 0),
			Edited: make([]string, 0),
			Read:   make([]string, 0),
		},
		ToolsUsed:     make(map[string]int),
		Errors:        make([]ErrorEntry, 0),
		Prompts:       make([]PromptEntry, 0),
		Notifications: make([]NotificationEntry, 0),
	}

	// Create session directory and write initial state
	sessionPath := sm.getSessionPath(sessionID)
	if err := sm.writer.WriteJSON(ctx, sessionPath, state); err != nil {
		return nil, fmt.Errorf("failed to initialize state: %w", err)
	}

	return state, nil
}

// LoadState loads an existing session state from disk
func (sm *StateManager) LoadState(ctx context.Context, sessionID string) (*SessionState, error) {
	if sessionID == "" {
		return nil, &StateError{
			Code:    "invalid_session_id",
			Message: "session ID cannot be empty",
		}
	}

	sessionPath := sm.getSessionPath(sessionID)

	// Check if file exists
	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		return nil, &StateError{
			Code:    "session_not_found",
			Message: fmt.Sprintf("session %s does not exist", sessionID),
		}
	}

	// Create context with timeout for read operation
	readCtx, cancel := context.WithTimeout(ctx, sm.timeout)
	defer cancel()

	// Read file with context
	data, err := sm.readFileWithContext(readCtx, sessionPath)
	if err != nil {
		return nil, &FileError{
			Op:   "read_state_file",
			Path: sessionPath,
			Err:  err,
		}
	}

	// Unmarshal JSON data
	var state SessionState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, &FileError{
			Op:   "unmarshal_json",
			Path: sessionPath,
			Err:  err,
		}
	}

	return &state, nil
}

// UpdateState atomically updates an existing session state
func (sm *StateManager) UpdateState(ctx context.Context, sessionID string, updateFunc func(*SessionState) error) error {
	if sessionID == "" {
		return &StateError{
			Code:    "invalid_session_id",
			Message: "session ID cannot be empty",
		}
	}

	// Load current state
	state, err := sm.LoadState(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to load state for update: %w", err)
	}

	// Apply update function
	if err := updateFunc(state); err != nil {
		return fmt.Errorf("update function failed: %w", err)
	}

	// Update timestamp
	state.UpdatedAt = time.Now().UTC()

	// Write updated state atomically
	sessionPath := sm.getSessionPath(sessionID)
	if err := sm.writer.WriteJSON(ctx, sessionPath, state); err != nil {
		return fmt.Errorf("failed to write updated state: %w", err)
	}

	return nil
}

// DeleteState removes a session state file
func (sm *StateManager) DeleteState(ctx context.Context, sessionID string) error {
	if sessionID == "" {
		return &StateError{
			Code:    "invalid_session_id",
			Message: "session ID cannot be empty",
		}
	}

	sessionPath := sm.getSessionPath(sessionID)

	// Check if file exists
	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		return &StateError{
			Code:    "session_not_found",
			Message: fmt.Sprintf("session %s does not exist", sessionID),
		}
	}

	// Remove the state file
	if err := os.Remove(sessionPath); err != nil {
		return &FileError{
			Op:   "delete_state_file",
			Path: sessionPath,
			Err:  err,
		}
	}

	return nil
}

// ListSessions returns all available session IDs
func (sm *StateManager) ListSessions(ctx context.Context) ([]string, error) {
	sessionsDir := filepath.Join(sm.basePath, "sessions")

	// Check if sessions directory exists
	if _, err := os.Stat(sessionsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	// Read directory entries
	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		return nil, &FileError{
			Op:   "read_sessions_directory",
			Path: sessionsDir,
			Err:  err,
		}
	}

	// Filter for directories containing state.json
	var sessions []string
	for _, entry := range entries {
		if entry.IsDir() {
			statePath := filepath.Join(sessionsDir, entry.Name(), StateFileName)
			if _, err := os.Stat(statePath); err == nil {
				sessions = append(sessions, entry.Name())
			}
		}
	}

	return sessions, nil
}

// getSessionPath returns the full path to a session's state file
func (sm *StateManager) getSessionPath(sessionID string) string {
	return filepath.Join(sm.basePath, "sessions", sessionID, StateFileName)
}

// readFileWithContext reads a file with context cancellation support
func (sm *StateManager) readFileWithContext(ctx context.Context, filename string) ([]byte, error) {
	// Create channel for read operation
	type result struct {
		data []byte
		err  error
	}
	done := make(chan result, 1)

	go func() {
		data, err := os.ReadFile(filename)
		done <- result{data: data, err: err}
	}()

	select {
	case res := <-done:
		return res.data, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Convenience functions for common update operations

// AddAgent adds an agent to the session state
func (sm *StateManager) AddAgent(ctx context.Context, sessionID, agentName string) error {
	return sm.UpdateState(ctx, sessionID, func(state *SessionState) error {
		// Add to current agents if not already present
		for _, agent := range state.Agents {
			if agent == agentName {
				return nil // Already present
			}
		}
		state.Agents = append(state.Agents, agentName)

		// Add to history
		state.AgentsHistory = append(state.AgentsHistory, AgentExecution{
			Name:      agentName,
			StartedAt: time.Now().UTC(),
		})

		return nil
	})
}

// CompleteAgent marks an agent as completed
func (sm *StateManager) CompleteAgent(ctx context.Context, sessionID, agentName string) error {
	return sm.UpdateState(ctx, sessionID, func(state *SessionState) error {
		// Remove from current agents
		for i, agent := range state.Agents {
			if agent == agentName {
				state.Agents = append(state.Agents[:i], state.Agents[i+1:]...)
				break
			}
		}

		// Update history with completion time
		now := time.Now().UTC()
		for i := range state.AgentsHistory {
			if state.AgentsHistory[i].Name == agentName && state.AgentsHistory[i].CompletedAt == nil {
				state.AgentsHistory[i].CompletedAt = &now
				break
			}
		}

		return nil
	})
}

// RecordError adds an error entry to the session state
func (sm *StateManager) RecordError(ctx context.Context, sessionID string, message, source, severity string) error {
	return sm.UpdateState(ctx, sessionID, func(state *SessionState) error {
		state.Errors = append(state.Errors, ErrorEntry{
			Timestamp: time.Now().UTC(),
			Message:   message,
			Source:    source,
			Severity:  severity,
		})
		return nil
	})
}

// RecordFileOperation adds file operations to the session state
func (sm *StateManager) RecordFileOperation(ctx context.Context, sessionID, operation, filepath string) error {
	if operation != "new" && operation != "edited" && operation != "read" {
		return &StateError{
			Code:    "invalid_operation",
			Message: fmt.Sprintf("invalid file operation: %s", operation),
		}
	}

	return sm.UpdateState(ctx, sessionID, func(state *SessionState) error {
		switch operation {
		case "new":
			state.Files.New = append(state.Files.New, filepath)
		case "edited":
			state.Files.Edited = append(state.Files.Edited, filepath)
		case "read":
			state.Files.Read = append(state.Files.Read, filepath)
		}
		return nil
	})
}

// IncrementToolUsage increments the usage counter for a tool
func (sm *StateManager) IncrementToolUsage(ctx context.Context, sessionID, toolName string) error {
	return sm.UpdateState(ctx, sessionID, func(state *SessionState) error {
		if state.ToolsUsed == nil {
			state.ToolsUsed = make(map[string]int)
		}
		state.ToolsUsed[toolName]++
		return nil
	})
}

// SetSessionActive sets the session active status
func (sm *StateManager) SetSessionActive(ctx context.Context, sessionID string, active bool) error {
	return sm.UpdateState(ctx, sessionID, func(state *SessionState) error {
		state.SessionActive = active
		return nil
	})
}
