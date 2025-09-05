package hooks

// SessionState represents the complete session tracking structure
type SessionState struct {
	SessionID     string              `json:"session_id"`
	Source        string              `json:"source"`
	ProjectPath   string              `json:"project_path"`
	Timestamp     string              `json:"timestamp"`
	LastUpdate    string              `json:"last_update"`
	Status        string              `json:"status"`
	Agents        []string            `json:"agents"`
	AgentsHistory []AgentHistoryEntry `json:"agents_history"`
	Files         FileOperations      `json:"files"`
	ToolsUsed     map[string]int      `json:"tools_used"`
	Errors        []ErrorEntry        `json:"errors"`
	Modified      bool                `json:"modified"`
}