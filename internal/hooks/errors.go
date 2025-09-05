package hooks

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/dylan-gluck/spcstr/internal/session"
	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

var (
	ErrBlockedOperation = errors.New("operation blocked by safety check")
	ErrSessionNotFound  = errors.New("session state not found")
	ErrInvalidHookData  = errors.New("invalid hook data format")
	ErrCorruptedState   = errors.New("session state is corrupted")
)

type RecoverableError struct {
	Err       error
	Recovered bool
	Action    string
}

func (r RecoverableError) Error() string {
	if r.Recovered {
		return fmt.Sprintf("%v (recovered: %s)", r.Err, r.Action)
	}
	return r.Err.Error()
}

func RecoverCorruptedState(projectRoot string) error {
	slog.Warn("attempting to recover corrupted session state")

	sessionFile := filepath.Join(projectRoot, ".spcstr", "session_state.json")
	backupFile := sessionFile + ".corrupted." + time.Now().Format("20060102_150405")

	if err := os.Rename(sessionFile, backupFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to backup corrupted state: %w", err)
	}

	newState := &models.SessionState{
		SessionID:     session.GenerateSessionID(),
		Source:        "recovery",
		ProjectPath:   projectRoot,
		Timestamp:     models.CurrentTimestamp(),
		LastUpdate:    models.CurrentTimestamp(),
		Status:        "active",
		Agents:        []string{},
		AgentsHistory: []models.AgentHistoryEntry{},
		Files: models.FileOperations{
			New:    []string{},
			Edited: []string{},
			Read:   []string{},
		},
		ToolsUsed: make(map[string]int),
		Errors: []models.ErrorEntry{
			{
				Timestamp: models.CurrentTimestamp(),
				Hook:      "recovery",
				Message:   fmt.Sprintf("Session recovered from corrupted state (backup: %s)", backupFile),
			},
		},
		Modified: true,
	}

	if err := session.SaveSessionState(projectRoot, newState); err != nil {
		return fmt.Errorf("failed to create recovery state: %w", err)
	}

	slog.Info("successfully recovered session state", "backup", backupFile)
	return nil
}

func LogStructuredError(hook string, err error, data map[string]interface{}) {
	attrs := []any{
		"hook", hook,
		"error", err.Error(),
	}

	if sessionID, ok := data["session_id"].(string); ok {
		attrs = append(attrs, "session_id", sessionID)
	}

	if toolName, ok := data["tool_name"].(string); ok {
		attrs = append(attrs, "tool", toolName)
	}

	slog.Error("hook error", attrs...)
}

func IsRecoverable(err error) bool {
	if err == nil {
		return true
	}

	var recErr RecoverableError
	if errors.As(err, &recErr) {
		return recErr.Recovered
	}

	return !errors.Is(err, ErrBlockedOperation) &&
		!errors.Is(err, ErrCorruptedState)
}
