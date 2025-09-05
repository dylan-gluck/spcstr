package config

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfigFromReader(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Configuration
		wantErr bool
	}{
		{
			name: "valid config",
			input: `{
				"version": "1.0.0",
				"scope": "project",
				"paths": {
					"hooks": ".spcstr/hooks",
					"sessions": ".spcstr/sessions",
					"docs": ["docs/"]
				},
				"ui": {
					"theme": "dark",
					"refreshInterval": 500
				},
				"session": {
					"retentionDays": 30,
					"autoArchive": true,
					"maxActiveSession": 10
				}
			}`,
			want:    DefaultConfig(),
			wantErr: false,
		},
		{
			name: "minimal config",
			input: `{
				"version": "1.0.0",
				"scope": "global",
				"paths": {
					"hooks": "hooks",
					"sessions": "sessions",
					"docs": ["docs"]
				},
				"ui": {
					"theme": "light",
					"refreshInterval": 1000
				},
				"session": {
					"retentionDays": 7,
					"autoArchive": false,
					"maxActiveSession": 5
				}
			}`,
			want: &Configuration{
				Version: "1.0.0",
				Scope:   "global",
				Paths: PathConfig{
					Hooks:    "hooks",
					Sessions: "sessions",
					Docs:     []string{"docs"},
				},
				UI: UIConfig{
					Theme:           "light",
					RefreshInterval: 1000,
				},
				Session: SessionConfig{
					RetentionDays:    7,
					AutoArchive:      false,
					MaxActiveSession: 5,
				},
			},
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{invalid json}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   `{}`,
			want:    &Configuration{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			got, err := LoadConfigFromReader(reader)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestSaveConfigToWriter(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Configuration
		want string
	}{
		{
			name: "default config",
			cfg:  DefaultConfig(),
			want: `{
  "version": "1.0.0",
  "scope": "project",
  "paths": {
    "hooks": ".spcstr/hooks",
    "sessions": ".spcstr/sessions",
    "docs": [
      "docs/"
    ]
  },
  "ui": {
    "theme": "dark",
    "refreshInterval": 500
  },
  "session": {
    "retentionDays": 30,
    "autoArchive": true,
    "maxActiveSession": 10
  }
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := SaveConfigToWriter(&buf, tt.cfg)

			assert.NoError(t, err)
			assert.JSONEq(t, tt.want, buf.String())
		})
	}
}

func TestLoadConfig(t *testing.T) {
	// Create temp directory for test
	tmpDir := t.TempDir()

	// Test loading existing config
	configPath := filepath.Join(tmpDir, "config.json")
	cfg := DefaultConfig()
	err := SaveConfig(configPath, cfg)
	require.NoError(t, err)

	loaded, err := LoadConfig(configPath)
	assert.NoError(t, err)
	assert.Equal(t, cfg, loaded)

	// Test loading non-existent file
	nonExistentPath := filepath.Join(tmpDir, "nonexistent.json")
	_, err = LoadConfig(nonExistentPath)
	assert.Error(t, err)
}

func TestSaveConfig(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		cfg     *Configuration
		wantErr bool
	}{
		{
			name:    "save valid config",
			path:    filepath.Join(tmpDir, "config.json"),
			cfg:     DefaultConfig(),
			wantErr: false,
		},
		{
			name:    "save to new directory",
			path:    filepath.Join(tmpDir, "subdir", "config.json"),
			cfg:     DefaultConfig(),
			wantErr: false,
		},
		{
			name:    "save invalid config",
			path:    filepath.Join(tmpDir, "invalid.json"),
			cfg:     &Configuration{}, // Invalid - missing required fields
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SaveConfig(tt.path, tt.cfg)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// Verify file exists and has correct permissions
				info, err := os.Stat(tt.path)
				assert.NoError(t, err)
				assert.Equal(t, os.FileMode(0644), info.Mode().Perm())

				// Verify content
				loaded, err := LoadConfig(tt.path)
				assert.NoError(t, err)
				assert.Equal(t, tt.cfg, loaded)
			}
		})
	}
}

func TestConfigExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Test non-existent file
	nonExistentPath := filepath.Join(tmpDir, "nonexistent.json")
	assert.False(t, ConfigExists(nonExistentPath))

	// Create a file
	existingPath := filepath.Join(tmpDir, "existing.json")
	err := SaveConfig(existingPath, DefaultConfig())
	require.NoError(t, err)

	// Test existing file
	assert.True(t, ConfigExists(existingPath))
}

func TestMergeConfigs(t *testing.T) {
	tests := []struct {
		name     string
		configs  []*Configuration
		expected *Configuration
	}{
		{
			name:     "no configs",
			configs:  []*Configuration{},
			expected: DefaultConfig(),
		},
		{
			name:     "single config",
			configs:  []*Configuration{DefaultConfig()},
			expected: DefaultConfig(),
		},
		{
			name: "nil first config",
			configs: []*Configuration{
				nil,
				&Configuration{
					UI: UIConfig{Theme: "light"},
				},
			},
			expected: func() *Configuration {
				cfg := DefaultConfig()
				cfg.UI.Theme = "light"
				return cfg
			}(),
		},
		{
			name: "multiple configs",
			configs: []*Configuration{
				DefaultConfig(),
				&Configuration{
					UI: UIConfig{Theme: "light"},
				},
				&Configuration{
					Session: SessionConfig{
						RetentionDays: 60,
						AutoArchive:   true, // Explicitly set to match expected
					},
				},
			},
			expected: func() *Configuration {
				cfg := DefaultConfig()
				cfg.UI.Theme = "light"
				cfg.Session.RetentionDays = 60
				cfg.Session.AutoArchive = true
				return cfg
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MergeConfigs(tt.configs...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoadOrCreateConfig(t *testing.T) {
	tmpDir := t.TempDir()

	// Test creating new config
	newPath := filepath.Join(tmpDir, "new.json")
	cfg, err := LoadOrCreateConfig(newPath)
	assert.NoError(t, err)
	assert.Equal(t, DefaultConfig(), cfg)
	assert.True(t, ConfigExists(newPath))

	// Test loading existing config
	cfg.UI.Theme = "light"
	err = SaveConfig(newPath, cfg)
	require.NoError(t, err)

	loaded, err := LoadOrCreateConfig(newPath)
	assert.NoError(t, err)
	assert.Equal(t, cfg, loaded)
}
