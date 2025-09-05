package hooks

// AgentHistoryEntry tracks agent execution history
type AgentHistoryEntry struct {
	Name      string `json:"name"`
	StartedAt string `json:"started_at"`
	EndedAt   string `json:"ended_at,omitempty"`
}