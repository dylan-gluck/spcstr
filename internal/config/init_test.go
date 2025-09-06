package config

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDirectoryStructure(t *testing.T) {
	tests := []struct {
		name        string
		projectRoot string
		expectError bool
	}{
		{
			name:        "create directories in temp dir",
			projectRoot: t.TempDir(),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := createDirectoryStructure(ctx, tt.projectRoot)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify directories were created
			if !tt.expectError {
				logsDir := filepath.Join(tt.projectRoot, ".spcstr", "logs")
				sessionsDir := filepath.Join(tt.projectRoot, ".spcstr", "sessions")

				if !dirExists(logsDir) {
					t.Errorf("logs directory not created: %s", logsDir)
				}
				if !dirExists(sessionsDir) {
					t.Errorf("sessions directory not created: %s", sessionsDir)
				}
			}
		})
	}
}

func TestConfigureClaudeHooks(t *testing.T) {
	tests := []struct {
		name             string
		existingSettings map[string]interface{}
		expectError      bool
	}{
		{
			name:             "new settings file",
			existingSettings: nil,
			expectError:      false,
		},
		{
			name:             "existing empty settings",
			existingSettings: map[string]interface{}{},
			expectError:      false,
		},
		{
			name: "existing settings with other data",
			existingSettings: map[string]interface{}{
				"otherKey": "otherValue",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectRoot := t.TempDir()
			ctx := context.Background()

			// Create existing settings if specified
			if tt.existingSettings != nil {
				claudeDir := filepath.Join(projectRoot, ".claude")
				os.MkdirAll(claudeDir, 0755)

				settingsPath := filepath.Join(claudeDir, "settings.json")
				data, _ := json.MarshalIndent(tt.existingSettings, "", "  ")
				os.WriteFile(settingsPath, data, 0644)
			}

			err := configureClaudeHooks(ctx, projectRoot)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify settings file was created and has correct hooks
			if !tt.expectError {
				settingsPath := filepath.Join(projectRoot, ".claude", "settings.json")

				data, err := os.ReadFile(settingsPath)
				if err != nil {
					t.Fatalf("failed to read settings.json: %v", err)
				}

				var settings map[string]interface{}
				if err := json.Unmarshal(data, &settings); err != nil {
					t.Fatalf("failed to parse settings.json: %v", err)
				}

				hooks, ok := settings["hooks"].(map[string]interface{})
				if !ok {
					t.Fatal("hooks not found in settings or wrong type")
				}

				// Check for PascalCase event names and proper structure
				expectedHooks := map[string]string{
					"SessionStart":     "session_start",
					"UserPromptSubmit": "user_prompt_submit",
					"PreToolUse":       "pre_tool_use",
					"PostToolUse":      "post_tool_use",
					"Notification":     "notification",
					"PreCompact":       "pre_compact",
					"SessionEnd":       "session_end",
					"Stop":             "stop",
					"SubagentStop":     "subagent_stop",
				}

				for eventName, commandName := range expectedHooks {
					hookArray, exists := hooks[eventName]
					if !exists {
						t.Errorf("event %s not found in settings", eventName)
						continue
					}

					// Verify it's an array
					arr, ok := hookArray.([]interface{})
					if !ok || len(arr) == 0 {
						t.Errorf("event %s is not an array or is empty", eventName)
						continue
					}

					// Check first element structure
					firstConfig, ok := arr[0].(map[string]interface{})
					if !ok {
						t.Errorf("event %s first element is not a map", eventName)
						continue
					}

					// Check for hooks array
					hooksArr, ok := firstConfig["hooks"].([]interface{})
					if !ok || len(hooksArr) == 0 {
						t.Errorf("event %s does not have hooks array", eventName)
						continue
					}

					// Check first hook structure
					firstHook, ok := hooksArr[0].(map[string]interface{})
					if !ok {
						t.Errorf("event %s first hook is not a map", eventName)
						continue
					}

					// Verify type and command fields
					if hookType, ok := firstHook["type"].(string); !ok || hookType != "command" {
						t.Errorf("event %s hook does not have type 'command'", eventName)
					}

					if command, ok := firstHook["command"].(string); !ok {
						t.Errorf("event %s hook does not have command field", eventName)
					} else {
						expectedCmd := `spcstr hook ` + commandName + ` --cwd="${CLAUDE_PROJECT_DIR}"`
						if command != expectedCmd {
							t.Errorf("event %s command mismatch: got %q, want %q", eventName, command, expectedCmd)
						}
					}

					// Check matcher field for events that need it
					needsMatcher := []string{"PreToolUse", "PostToolUse", "PreCompact"}
					for _, name := range needsMatcher {
						if eventName == name {
							if matcher, exists := firstConfig["matcher"]; !exists || matcher != "*" {
								t.Errorf("event %s should have matcher='*'", eventName)
							}
						}
					}
				}

				// Check that existing data is preserved
				if tt.existingSettings != nil {
					for key := range tt.existingSettings {
						if key != "hooks" {
							if _, exists := settings[key]; !exists {
								t.Errorf("existing key %s was not preserved", key)
							}
						}
					}
				}
			}
		})
	}
}

func TestWriteSettingsAtomic(t *testing.T) {
	tests := []struct {
		name        string
		settings    map[string]interface{}
		expectError bool
	}{
		{
			name: "write simple settings",
			settings: map[string]interface{}{
				"key": "value",
			},
			expectError: false,
		},
		{
			name: "write complex nested settings",
			settings: map[string]interface{}{
				"hooks": map[string]string{
					"session_start": "command",
				},
				"other": "data",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			settingsPath := filepath.Join(tempDir, "settings.json")
			ctx := context.Background()

			err := writeSettingsAtomic(ctx, settingsPath, tt.settings)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify file was written correctly
			if !tt.expectError {
				data, err := os.ReadFile(settingsPath)
				if err != nil {
					t.Fatalf("failed to read written file: %v", err)
				}

				var readSettings map[string]interface{}
				if err := json.Unmarshal(data, &readSettings); err != nil {
					t.Fatalf("failed to parse written file: %v", err)
				}

				// Simple check - just verify keys match
				for key := range tt.settings {
					if _, exists := readSettings[key]; !exists {
						t.Errorf("key %s not found in written settings", key)
					}
				}
			}
		})
	}
}

func TestDirExists(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "file.txt")
	os.WriteFile(tempFile, []byte("test"), 0644)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "existing directory",
			path:     tempDir,
			expected: true,
		},
		{
			name:     "non-existent path",
			path:     filepath.Join(tempDir, "nonexistent"),
			expected: false,
		},
		{
			name:     "file not directory",
			path:     tempFile,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := dirExists(tt.path)
			if result != tt.expected {
				t.Errorf("dirExists(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestFileExists(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "file.txt")
	os.WriteFile(tempFile, []byte("test"), 0644)

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "existing file",
			path:     tempFile,
			expected: true,
		},
		{
			name:     "non-existent path",
			path:     filepath.Join(tempDir, "nonexistent.txt"),
			expected: false,
		},
		{
			name:     "directory not file",
			path:     tempDir,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fileExists(tt.path)
			if result != tt.expected {
				t.Errorf("fileExists(%s) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestInitializeProject(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(string) error
		force        bool
		expectError  bool
		expectPrompt bool
	}{
		{
			name:         "fresh project initialization",
			setupFunc:    nil,
			force:        false,
			expectError:  false,
			expectPrompt: false,
		},
		{
			name: "existing .spcstr with force flag",
			setupFunc: func(dir string) error {
				return os.MkdirAll(filepath.Join(dir, ".spcstr"), 0755)
			},
			force:        true,
			expectError:  false,
			expectPrompt: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory and change to it
			tempDir := t.TempDir()
			oldWd, _ := os.Getwd()
			defer os.Chdir(oldWd)
			os.Chdir(tempDir)

			// Setup if needed
			if tt.setupFunc != nil {
				if err := tt.setupFunc(tempDir); err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}

			// Run initialization
			err := InitializeProject(tt.force)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Verify structure was created
			if !tt.expectError {
				// Check .spcstr directories
				if !dirExists(filepath.Join(tempDir, ".spcstr", "logs")) {
					t.Error(".spcstr/logs directory not created")
				}
				if !dirExists(filepath.Join(tempDir, ".spcstr", "sessions")) {
					t.Error(".spcstr/sessions directory not created")
				}

				// Check settings.json
				settingsPath := filepath.Join(tempDir, ".claude", "settings.json")
				if !fileExists(settingsPath) {
					t.Error(".claude/settings.json not created")
				}
			}
		})
	}
}
