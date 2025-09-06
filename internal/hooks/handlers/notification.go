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

// NotificationParams defines the expected input for notification hook
type NotificationParams struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
}

// NotificationHandler handles the notification hook
type NotificationHandler struct{}

// NewNotificationHandler creates a new NotificationHandler
func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{}
}

// Name returns the hook name
func (h *NotificationHandler) Name() string {
	return "notification"
}

// Execute processes the notification hook
func (h *NotificationHandler) Execute(input []byte) error {
	var params NotificationParams
	if err := json.Unmarshal(input, &params); err != nil {
		return fmt.Errorf("failed to parse notification parameters: %w", err)
	}

	// Validate required fields
	if params.SessionID == "" {
		return fmt.Errorf("session_id is required")
	}
	if params.Message == "" {
		return fmt.Errorf("message is required")
	}

	// Parse timestamp or use current time
	var notificationTime time.Time
	if params.Timestamp != "" {
		var err error
		notificationTime, err = time.Parse(time.RFC3339, params.Timestamp)
		if err != nil {
			notificationTime = time.Now()
		}
	} else {
		notificationTime = time.Now()
	}

	// Set default level if not provided
	level := params.Level
	if level == "" {
		level = "info"
	}

	// Create StateManager using current working directory (after --cwd change)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	stateManager := state.NewStateManager(filepath.Join(cwd, ".spcstr"))

	// Add notification to session state using UpdateState
	ctx := context.Background()
	err = stateManager.UpdateState(ctx, params.SessionID, func(sessionState *state.SessionState) error {
		// Create new notification entry
		newNotification := state.NotificationEntry{
			Timestamp: notificationTime,
			Type:      "hook",
			Message:   params.Message,
			Level:     level,
		}

		sessionState.Notifications = append(sessionState.Notifications, newNotification)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to add notification to session: %w", err)
	}

	return nil
}
