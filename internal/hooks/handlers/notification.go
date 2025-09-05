package handlers

import (
	"github.com/dylan-gluck/spcstr/internal/session"
	"log/slog"

	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

func HandleNotification(projectRoot string, data map[string]interface{}) error {
	notificationType, _ := data["type"].(string)
	message, _ := data["message"].(string)

	slog.Debug("notification", "type", notificationType, "message", message)

	if notificationType == "error" {
		state, err := session.LoadSessionState(projectRoot)
		if err != nil {
			slog.Error("failed to load state", "error", err)
			return nil
		}

		errorEntry := models.ErrorEntry{
			Timestamp: models.CurrentTimestamp(),
			Hook:      "notification",
			Message:   message,
		}

		state.Errors = append(state.Errors, errorEntry)

		if err := session.SaveSessionState(projectRoot, state); err != nil {
			slog.Error("failed to save state", "error", err)
		}
	}

	return nil
}
