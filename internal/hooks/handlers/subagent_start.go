package handlers

import (
	"github.com/dylan-gluck/spcstr/internal/session"
	"log/slog"

	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

func HandleSubagentStart(projectRoot string, data map[string]interface{}) error {
	agentName, _ := data["agent_name"].(string)
	agentType, _ := data["agent_type"].(string)

	slog.Debug("subagent_start", "agent_name", agentName, "agent_type", agentType)

	state, err := session.LoadSessionState(projectRoot)
	if err != nil {
		slog.Error("failed to load state", "error", err)
		return nil
	}

	if agentName == "" {
		agentName = agentType
	}
	if agentName == "" {
		agentName = "unknown"
	}

	state.Agents = models.AddUniqueString(state.Agents, agentName)

	historyEntry := models.AgentHistoryEntry{
		Name:      agentName,
		StartedAt: models.CurrentTimestamp(),
	}
	state.AgentsHistory = append(state.AgentsHistory, historyEntry)

	if err := session.SaveSessionState(projectRoot, state); err != nil {
		slog.Error("failed to save state", "error", err)
	}

	return nil
}
