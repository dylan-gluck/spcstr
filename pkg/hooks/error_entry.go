package hooks

// ErrorEntry represents a structured error log entry
type ErrorEntry struct {
	Timestamp string `json:"timestamp"`
	Hook      string `json:"hook"`
	Message   string `json:"message"`
}
