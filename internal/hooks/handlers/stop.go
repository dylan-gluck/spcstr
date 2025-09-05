package handlers

import (
	"github.com/dylan-gluck/spcstr/internal/session"
	"log/slog"

	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

func HandleStop(projectRoot string, data map[string]interface{}) error {
	reason, _ := data["reason"].(string)
	slog.Debug("stop", "reason", reason)

	state, err := session.LoadSessionState(projectRoot)
	if err != nil {
		slog.Error("failed to load state", "error", err)
		return nil
	}

	if reason != "" {
		errorEntry := models.ErrorEntry{
			Timestamp: models.CurrentTimestamp(),
			Hook:      "stop",
			Message:   "Session stopped: " + reason,
		}
		state.Errors = append(state.Errors, errorEntry)
	}

	state.Status = "stopped"
	state.LastUpdate = models.CurrentTimestamp()

	if err := session.SaveSessionState(projectRoot, state); err != nil {
		slog.Error("failed to save state", "error", err)
	}

	return nil
}
