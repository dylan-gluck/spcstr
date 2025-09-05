package hooks

// FileOperations tracks categorized file operations
type FileOperations struct {
	New    []string `json:"new"`
	Edited []string `json:"edited"`
	Read   []string `json:"read"`
}
