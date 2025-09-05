package handlers

import (
	"fmt"
	"github.com/dylan-gluck/spcstr/internal/session"
	"log/slog"
	"time"

	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

func HandleSessionStart(projectRoot string, data map[string]interface{}) error {
	sessionID, _ := data["session_id"].(string)
	source, _ := data["source"].(string)
	projectPath, _ := data["project_path"].(string)

	if sessionID == "" {
		sessionID = fmt.Sprintf("sess_%d", time.Now().UnixNano())
	}
	if source == "" {
		source = "startup"
	}
	if projectPath == "" {
		projectPath = projectRoot
	}

	slog.Debug("session_start", "session_id", sessionID, "source", source)

	state := &models.SessionState{
		SessionID:     sessionID,
		Source:        source,
		ProjectPath:   projectPath,
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
		Errors:    []models.ErrorEntry{},
		Modified:  false,
	}

	if err := session.SaveSessionState(projectRoot, state); err != nil {
		slog.Error("failed to save initial state", "error", err)
		return nil
	}

	return nil
}
