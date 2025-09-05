package handlers

import (
	"github.com/dylan-gluck/spcstr/internal/session"
	"log/slog"

	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

func HandlePreCompact(projectRoot string, data map[string]interface{}) error {
	slog.Debug("pre_compact invoked")

	state, err := session.LoadSessionState(projectRoot)
	if err != nil {
		slog.Error("failed to load state", "error", err)
		return nil
	}

	state.LastUpdate = models.CurrentTimestamp()

	if err := session.SaveSessionState(projectRoot, state); err != nil {
		slog.Error("failed to save state", "error", err)
	}

	return nil
}
