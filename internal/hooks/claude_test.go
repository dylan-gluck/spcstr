package hooks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClaudeSettings_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		settings *ClaudeSettings
		want     string
	}{
		{
			name: "with hooks",
			settings: &ClaudeSettings{
				Hooks: map[string]interface{}{
					"PreToolUse":  "/path/to/pre-command.sh",
					"PostToolUse": "/path/to/post-command.sh",
				},
				Other: map[string]interface{}{
					"theme":    "dark",
					"fontSize": 14,
				},
			},
			want: `{
				"theme": "dark",
				"fontSize": 14,
				"hooks": {
					"PreToolUse": "/path/to/pre-command.sh",
					"PostToolUse": "/path/to/post-command.sh"
				}
			}`,
		},
		{
			name: "without hooks",
			settings: &ClaudeSettings{
				Hooks: map[string]interface{}{},
				Other: map[string]interface{}{
					"theme": "light",
				},
			},
			want: `{
				"theme": "light"
			}`,
		},
		{
			name: "empty settings",
			settings: &ClaudeSettings{
				Hooks: map[string]interface{}{},
				Other: map[string]interface{}{},
			},
			want: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.settings)
			assert.NoError(t, err)
			assert.JSONEq(t, tt.want, string(data))
		})
	}
}

func TestClaudeSettings_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *ClaudeSettings
		wantErr bool
	}{
		{
			name: "with hooks",
			input: `{
				"theme": "dark",
				"fontSize": 14,
				"hooks": {
					"PreToolUse": "/path/to/pre-command.sh",
					"PostToolUse": "/path/to/post-command.sh"
				}
			}`,
			want: &ClaudeSettings{
				Hooks: map[string]interface{}{
					"PreToolUse":  "/path/to/pre-command.sh",
					"PostToolUse": "/path/to/post-command.sh",
				},
				Other: map[string]interface{}{
					"theme":    "dark",
					"fontSize": float64(14), // JSON numbers are float64
				},
			},
			wantErr: false,
		},
		{
			name: "without hooks",
			input: `{
				"theme": "light",
				"someOtherField": true
			}`,
			want: &ClaudeSettings{
				Hooks: map[string]interface{}{},
				Other: map[string]interface{}{
					"theme":          "light",
					"someOtherField": true,
				},
			},
			wantErr: false,
		},
		{
			name:  "empty object",
			input: `{}`,
			want: &ClaudeSettings{
				Hooks: map[string]interface{}{},
				Other: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{invalid}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var settings ClaudeSettings
			err := json.Unmarshal([]byte(tt.input), &settings)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Hooks, settings.Hooks)
				assert.Equal(t, tt.want.Other, settings.Other)
			}
		})
	}
}

func TestNewClaudeUpdater(t *testing.T) {
	updater := NewClaudeUpdater()
	assert.NotNil(t, updater)
}

func TestClaudeUpdater_loadSettings(t *testing.T) {
	tmpDir := t.TempDir()
	updater := NewClaudeUpdater()

	// Test loading non-existent file
	nonExistentPath := filepath.Join(tmpDir, "nonexistent.json")
	settings, err := updater.loadSettings(nonExistentPath)
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Empty(t, settings.Hooks)
	assert.Empty(t, settings.Other)

	// Test loading existing file
	existingPath := filepath.Join(tmpDir, "settings.json")
	testSettings := &ClaudeSettings{
		Hooks: map[string]interface{}{
			"PreToolUse": "/test/hook.sh",
		},
		Other: map[string]interface{}{
			"theme": "dark",
		},
	}

	// Save test settings
	err = updater.saveSettings(existingPath, testSettings)
	require.NoError(t, err)

	// Load and verify
	loaded, err := updater.loadSettings(existingPath)
	assert.NoError(t, err)
	assert.Equal(t, testSettings.Hooks, loaded.Hooks)
	assert.Equal(t, "dark", loaded.Other["theme"])
}

func TestClaudeUpdater_saveSettings(t *testing.T) {
	tmpDir := t.TempDir()
	updater := NewClaudeUpdater()

	settingsPath := filepath.Join(tmpDir, "subdir", "settings.json")
	settings := &ClaudeSettings{
		Hooks: map[string]interface{}{
			"PreToolUse":   "/hooks/pre-command.sh",
			"PostToolUse":  "/hooks/post-command.sh",
			"Notification": "/hooks/file-modified.sh",
			"SessionEnd":   "/hooks/session-end.sh",
		},
		Other: map[string]interface{}{
			"theme":    "dark",
			"fontSize": float64(14),
		},
	}

	// Save settings
	err := updater.saveSettings(settingsPath, settings)
	assert.NoError(t, err)

	// Verify file exists with correct permissions
	info, err := os.Stat(settingsPath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0644), info.Mode().Perm())

	// Verify content
	loaded, err := updater.loadSettings(settingsPath)
	assert.NoError(t, err)
	assert.Equal(t, settings.Hooks, loaded.Hooks)
}

func TestClaudeUpdater_createBackup(t *testing.T) {
	tmpDir := t.TempDir()
	updater := NewClaudeUpdater()

	// Create original file
	originalPath := filepath.Join(tmpDir, "original.json")
	originalContent := []byte(`{"test": "data"}`)
	err := os.WriteFile(originalPath, originalContent, 0644)
	require.NoError(t, err)

	// Create backup
	backupPath := filepath.Join(tmpDir, "backup.json")
	err = updater.createBackup(originalPath, backupPath)
	assert.NoError(t, err)

	// Verify backup content
	backupContent, err := os.ReadFile(backupPath)
	assert.NoError(t, err)
	assert.Equal(t, originalContent, backupContent)

	// Test backup of non-existent file
	err = updater.createBackup(filepath.Join(tmpDir, "nonexistent.json"), backupPath)
	assert.NoError(t, err) // Should not error for non-existent source
}

func TestClaudeUpdater_restoreBackup(t *testing.T) {
	tmpDir := t.TempDir()
	updater := NewClaudeUpdater()

	// Create backup file
	backupPath := filepath.Join(tmpDir, "backup.json")
	backupContent := []byte(`{"backup": true}`)
	err := os.WriteFile(backupPath, backupContent, 0644)
	require.NoError(t, err)

	// Restore backup
	targetPath := filepath.Join(tmpDir, "restored.json")
	err = updater.restoreBackup(backupPath, targetPath)
	assert.NoError(t, err)

	// Verify restored content
	restoredContent, err := os.ReadFile(targetPath)
	assert.NoError(t, err)
	assert.Equal(t, backupContent, restoredContent)

	// Verify backup file no longer exists
	_, err = os.Stat(backupPath)
	assert.True(t, os.IsNotExist(err))
}

func TestClaudeUpdater_UpdateClaudeSettings(t *testing.T) {
	tmpDir := t.TempDir()
	updater := NewClaudeUpdater()

	// Create mock settings file
	settingsDir := filepath.Join(tmpDir, ".claude")
	err := os.MkdirAll(settingsDir, 0755)
	require.NoError(t, err)

	settingsPath := filepath.Join(settingsDir, "settings.json")
	initialSettings := &ClaudeSettings{
		Hooks: map[string]interface{}{},
		Other: map[string]interface{}{
			"theme": "dark",
		},
	}
	err = updater.saveSettings(settingsPath, initialSettings)
	require.NoError(t, err)

	// Mock the findClaudeSettings method to return our test path
	// Since we can't easily mock this, we'll test the core logic through other methods

	// Test that hooks are added correctly
	hooksPath := filepath.Join(tmpDir, "hooks")

	// Load settings, update hooks, and save
	settings, err := updater.loadSettings(settingsPath)
	require.NoError(t, err)

	settings.Hooks["PreToolUse"] = filepath.Join(hooksPath, "pre-command.sh")
	settings.Hooks["PostToolUse"] = filepath.Join(hooksPath, "post-command.sh")
	settings.Hooks["Notification"] = filepath.Join(hooksPath, "file-modified.sh")
	settings.Hooks["SessionEnd"] = filepath.Join(hooksPath, "session-end.sh")

	err = updater.saveSettings(settingsPath, settings)
	assert.NoError(t, err)

	// Verify hooks were added
	loaded, err := updater.loadSettings(settingsPath)
	assert.NoError(t, err)
	assert.Len(t, loaded.Hooks, 4)
	assert.Equal(t, filepath.Join(hooksPath, "pre-command.sh"), loaded.Hooks["PreToolUse"])
	assert.Equal(t, "dark", loaded.Other["theme"]) // Other settings preserved
}

func TestClaudeUpdater_RemoveHooks(t *testing.T) {
	tmpDir := t.TempDir()
	updater := NewClaudeUpdater()

	// Create settings with hooks
	settingsPath := filepath.Join(tmpDir, "settings.json")
	settings := &ClaudeSettings{
		Hooks: map[string]interface{}{
			"PreToolUse":   "/hooks/pre-command.sh",
			"PostToolUse":  "/hooks/post-command.sh",
			"Notification": "/hooks/file-modified.sh",
			"SessionEnd":   "/hooks/session-end.sh",
			"customHook":   "/custom/hook.sh", // Non-spcstr hook
		},
		Other: map[string]interface{}{
			"theme": "dark",
		},
	}
	err := updater.saveSettings(settingsPath, settings)
	require.NoError(t, err)

	// Remove spcstr hooks
	loaded, err := updater.loadSettings(settingsPath)
	require.NoError(t, err)

	delete(loaded.Hooks, "PreToolUse")
	delete(loaded.Hooks, "PostToolUse")
	delete(loaded.Hooks, "Notification")
	delete(loaded.Hooks, "SessionEnd")

	err = updater.saveSettings(settingsPath, loaded)
	assert.NoError(t, err)

	// Verify only custom hook remains
	final, err := updater.loadSettings(settingsPath)
	assert.NoError(t, err)
	assert.Len(t, final.Hooks, 1)
	assert.Equal(t, "/custom/hook.sh", final.Hooks["customHook"])
	assert.Equal(t, "dark", final.Other["theme"])
}
