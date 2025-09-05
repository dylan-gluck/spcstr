package hooks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ClaudeSettings represents Claude's settings.json structure
type ClaudeSettings struct {
	Hooks map[string]interface{} `json:"hooks,omitempty"`
	// Other fields are preserved but not typed
	Other map[string]interface{} `json:"-"`
}

// MarshalJSON custom marshaler to preserve unknown fields
func (cs *ClaudeSettings) MarshalJSON() ([]byte, error) {
	// Start with the Other fields
	result := make(map[string]interface{})
	for k, v := range cs.Other {
		result[k] = v
	}
	
	// Add hooks if present
	if len(cs.Hooks) > 0 {
		result["hooks"] = cs.Hooks
	}
	
	return json.Marshal(result)
}

// UnmarshalJSON custom unmarshaler to preserve unknown fields
func (cs *ClaudeSettings) UnmarshalJSON(data []byte) error {
	// First unmarshal into a map
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	
	// Extract hooks if present
	if hooks, ok := raw["hooks"]; ok {
		if hooksMap, ok := hooks.(map[string]interface{}); ok {
			cs.Hooks = hooksMap
		} else {
			cs.Hooks = make(map[string]interface{})
		}
		delete(raw, "hooks")
	} else {
		cs.Hooks = make(map[string]interface{})
	}
	
	// Store remaining fields
	cs.Other = raw
	return nil
}

// ClaudeUpdater handles updating Claude's settings.json
type ClaudeUpdater struct{}

// NewClaudeUpdater creates a new Claude settings updater
func NewClaudeUpdater() *ClaudeUpdater {
	return &ClaudeUpdater{}
}

// UpdateClaudeSettings updates Claude's settings.json with hook configurations
func (cu *ClaudeUpdater) UpdateClaudeSettings(hooksPath string, force bool) error {
	// Extract project directory from hooks path
	// hooksPath is like "/path/to/project/.spcstr/hooks"
	projectDir := filepath.Dir(filepath.Dir(hooksPath))
	
	// Find Claude settings file in project directory
	settingsPath, err := cu.findClaudeSettings(projectDir)
	if err != nil {
		return fmt.Errorf("finding Claude settings: %w", err)
	}
	
	if settingsPath == "" {
		return fmt.Errorf("Claude settings.json not found")
	}

	// Load existing settings
	settings, err := cu.loadSettings(settingsPath)
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	// Check for existing hooks
	if !force && len(settings.Hooks) > 0 {
		// Check if any spcstr hooks already exist
		for key := range settings.Hooks {
			if key == "PreToolUse" || key == "PostToolUse" || key == "Notification" || key == "SessionEnd" {
				return fmt.Errorf("hooks already configured in Claude settings (use --force to overwrite)")
			}
		}
	}

	// Create backup
	backupPath := settingsPath + ".backup." + time.Now().Format("20060102_150405")
	if err := cu.createBackup(settingsPath, backupPath); err != nil {
		return fmt.Errorf("creating backup: %w", err)
	}

	// Update hooks
	if settings.Hooks == nil {
		settings.Hooks = make(map[string]interface{})
	}
	
	// Create proper Claude Code hook format using $CLAUDE_PROJECT_DIR
	settings.Hooks["PreToolUse"] = []map[string]interface{}{
		{
			"matcher": "*",
			"hooks": []map[string]interface{}{
				{
					"type":    "command",
					"command": "$CLAUDE_PROJECT_DIR/.spcstr/hooks/pre-command.sh",
				},
			},
		},
	}
	
	settings.Hooks["PostToolUse"] = []map[string]interface{}{
		{
			"matcher": "*", 
			"hooks": []map[string]interface{}{
				{
					"type":    "command",
					"command": "$CLAUDE_PROJECT_DIR/.spcstr/hooks/post-command.sh",
				},
			},
		},
	}
	
	settings.Hooks["Notification"] = []map[string]interface{}{
		{
			"hooks": []map[string]interface{}{
				{
					"type":    "command",
					"command": "$CLAUDE_PROJECT_DIR/.spcstr/hooks/file-modified.sh",
				},
			},
		},
	}
	
	settings.Hooks["SessionEnd"] = []map[string]interface{}{
		{
			"hooks": []map[string]interface{}{
				{
					"type":    "command",
					"command": "$CLAUDE_PROJECT_DIR/.spcstr/hooks/session-end.sh",
				},
			},
		},
	}

	// Save updated settings
	if err := cu.saveSettings(settingsPath, settings); err != nil {
		// Try to restore backup on failure
		_ = cu.restoreBackup(backupPath, settingsPath)
		return fmt.Errorf("saving settings: %w", err)
	}

	return nil
}

// RemoveHooks removes spcstr hooks from Claude settings
func (cu *ClaudeUpdater) RemoveHooks() error {
	// Get current working directory as project directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}
	
	settingsPath, err := cu.findClaudeSettings(cwd)
	if err != nil {
		return fmt.Errorf("finding Claude settings: %w", err)
	}
	
	if settingsPath == "" {
		return nil // No settings file, nothing to remove
	}

	settings, err := cu.loadSettings(settingsPath)
	if err != nil {
		return fmt.Errorf("loading settings: %w", err)
	}

	// Remove spcstr hooks
	delete(settings.Hooks, "PreToolUse")
	delete(settings.Hooks, "PostToolUse")
	delete(settings.Hooks, "Notification")
	delete(settings.Hooks, "SessionEnd")

	// Save updated settings
	return cu.saveSettings(settingsPath, settings)
}

// findClaudeSettings finds the Claude settings.json file in the project directory
func (cu *ClaudeUpdater) findClaudeSettings(projectDir string) (string, error) {
	// Look for .claude/settings.json in the project directory
	settingsPath := filepath.Join(projectDir, ".claude", "settings.json")
	
	// Check if it exists
	if _, err := os.Stat(settingsPath); err == nil {
		return settingsPath, nil
	}
	
	// Try to create it if it doesn't exist
	claudeDir := filepath.Join(projectDir, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return "", fmt.Errorf("creating .claude directory: %w", err)
	}
	
	// Create empty settings file
	settings := &ClaudeSettings{
		Hooks: make(map[string]interface{}),
		Other: make(map[string]interface{}),
	}
	if err := cu.saveSettings(settingsPath, settings); err != nil {
		return "", fmt.Errorf("creating settings file: %w", err)
	}
	
	return settingsPath, nil
}

// loadSettings loads Claude settings from file
func (cu *ClaudeUpdater) loadSettings(path string) (*ClaudeSettings, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty settings if file doesn't exist
			return &ClaudeSettings{
				Hooks: make(map[string]interface{}),
				Other: make(map[string]interface{}),
			}, nil
		}
		return nil, err
	}
	defer file.Close()

	var settings ClaudeSettings
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&settings); err != nil {
		return nil, err
	}

	if settings.Other == nil {
		settings.Other = make(map[string]interface{})
	}
	if settings.Hooks == nil {
		settings.Hooks = make(map[string]interface{})
	}

	return &settings, nil
}

// saveSettings saves Claude settings to file
func (cu *ClaudeUpdater) saveSettings(path string, settings *ClaudeSettings) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write to temporary file first
	tempPath := path + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(settings); err != nil {
		file.Close()
		os.Remove(tempPath)
		return err
	}

	if err := file.Close(); err != nil {
		os.Remove(tempPath)
		return err
	}

	// Set permissions
	if err := os.Chmod(tempPath, 0644); err != nil {
		os.Remove(tempPath)
		return err
	}

	// Atomic rename
	return os.Rename(tempPath, path)
}

// createBackup creates a backup of the settings file
func (cu *ClaudeUpdater) createBackup(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file to backup
		}
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy content
	buf := make([]byte, 8192)
	for {
		n, err := srcFile.Read(buf)
		if n > 0 {
			if _, werr := dstFile.Write(buf[:n]); werr != nil {
				return werr
			}
		}
		if err != nil {
			if err == os.ErrClosed || err.Error() == "EOF" {
				break
			}
			return err
		}
	}

	return nil
}

// restoreBackup restores a backup file
func (cu *ClaudeUpdater) restoreBackup(backup, target string) error {
	return os.Rename(backup, target)
}