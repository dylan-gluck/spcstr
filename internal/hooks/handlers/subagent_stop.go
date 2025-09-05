package handlers

import (
	"github.com/dylan-gluck/spcstr/internal/session"
	"log/slog"

	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

func HandleSubagentStop(projectRoot string, data map[string]interface{}) error {
	agentName, _ := data["agent_name"].(string)
	agentType, _ := data["agent_type"].(string)

	slog.Debug("subagent_stop", "agent_name", agentName, "agent_type", agentType)

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

	for i := range state.AgentsHistory {
		if state.AgentsHistory[i].Name == agentName && state.AgentsHistory[i].EndedAt == "" {
			state.AgentsHistory[i].EndedAt = models.CurrentTimestamp()
			break
		}
	}

	newAgents := []string{}
	for _, a := range state.Agents {
		if a != agentName {
			newAgents = append(newAgents, a)
		}
	}
	state.Agents = newAgents

	if err := session.SaveSessionState(projectRoot, state); err != nil {
		slog.Error("failed to save state", "error", err)
	}

	return nil
}
