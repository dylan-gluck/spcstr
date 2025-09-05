package handlers

import (
	"github.com/dylan-gluck/spcstr/internal/session"
	"log/slog"

	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

func HandleSessionEnd(projectRoot string, data map[string]interface{}) error {
	sessionID, _ := data["session_id"].(string)
	slog.Debug("session_end", "session_id", sessionID)

	state, err := session.LoadSessionState(projectRoot)
	if err != nil {
		slog.Error("failed to load state", "error", err)
		return nil
	}

	state.Status = "completed"
	state.LastUpdate = models.CurrentTimestamp()

	for i := range state.AgentsHistory {
		if state.AgentsHistory[i].EndedAt == "" {
			state.AgentsHistory[i].EndedAt = models.CurrentTimestamp()
		}
	}

	state.Agents = []string{}

	if err := session.SaveSessionState(projectRoot, state); err != nil {
		slog.Error("failed to save final state", "error", err)
	}

	return nil
}
