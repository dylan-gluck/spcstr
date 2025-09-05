package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "1.0.0", cfg.Version)
	assert.Equal(t, "project", cfg.Scope)
	assert.Equal(t, ".spcstr/hooks", cfg.Paths.Hooks)
	assert.Equal(t, ".spcstr/sessions", cfg.Paths.Sessions)
	assert.Equal(t, []string{"docs/"}, cfg.Paths.Docs)
	assert.Equal(t, "dark", cfg.UI.Theme)
	assert.Equal(t, 500, cfg.UI.RefreshInterval)
	assert.Equal(t, 30, cfg.Session.RetentionDays)
	assert.True(t, cfg.Session.AutoArchive)
	assert.Equal(t, 10, cfg.Session.MaxActiveSession)
}

func TestGlobalConfig(t *testing.T) {
	cfg := GlobalConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, "global", cfg.Scope)
	// Other fields should match defaults
	assert.Equal(t, "1.0.0", cfg.Version)
}

func TestConfiguration_Merge(t *testing.T) {
	tests := []struct {
		name     string
		base     *Configuration
		other    *Configuration
		expected *Configuration
	}{
		{
			name:     "merge with nil",
			base:     DefaultConfig(),
			other:    nil,
			expected: DefaultConfig(),
		},
		{
			name: "merge paths",
			base: DefaultConfig(),
			other: &Configuration{
				Paths: PathConfig{
					Hooks:    "/custom/hooks",
					Sessions: "/custom/sessions",
					Docs:     []string{"/custom/docs"},
				},
				Session: SessionConfig{
					RetentionDays:    1, // Need to set something to indicate Session was populated
					AutoArchive:      false,
					MaxActiveSession: 1,
				},
			},
			expected: func() *Configuration {
				cfg := DefaultConfig()
				cfg.Paths.Hooks = "/custom/hooks"
				cfg.Paths.Sessions = "/custom/sessions"
				cfg.Paths.Docs = []string{"/custom/docs"}
				cfg.Session.RetentionDays = 1
				cfg.Session.AutoArchive = false
				cfg.Session.MaxActiveSession = 1
				return cfg
			}(),
		},
		{
			name: "merge UI config",
			base: DefaultConfig(),
			other: &Configuration{
				UI: UIConfig{
					Theme:           "light",
					RefreshInterval: 1000,
				},
				Session: SessionConfig{
					RetentionDays:    1, // Need to set something to indicate Session was populated
					AutoArchive:      false,
					MaxActiveSession: 1,
				},
			},
			expected: func() *Configuration {
				cfg := DefaultConfig()
				cfg.UI.Theme = "light"
				cfg.UI.RefreshInterval = 1000
				cfg.Session.RetentionDays = 1
				cfg.Session.AutoArchive = false
				cfg.Session.MaxActiveSession = 1
				return cfg
			}(),
		},
		{
			name: "merge session config",
			base: DefaultConfig(),
			other: &Configuration{
				Session: SessionConfig{
					RetentionDays:    60,
					AutoArchive:      false,
					MaxActiveSession: 20,
				},
			},
			expected: func() *Configuration {
				cfg := DefaultConfig()
				cfg.Session.RetentionDays = 60
				cfg.Session.AutoArchive = false
				cfg.Session.MaxActiveSession = 20
				return cfg
			}(),
		},
		{
			name: "partial merge",
			base: DefaultConfig(),
			other: &Configuration{
				UI: UIConfig{
					Theme: "light",
					// RefreshInterval not set, should keep original
				},
				// Session not populated, so AutoArchive should remain unchanged
			},
			expected: func() *Configuration {
				cfg := DefaultConfig()
				cfg.UI.Theme = "light"
				// RefreshInterval remains 500
				// AutoArchive remains true since Session was not populated in other
				return cfg
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.base.Merge(tt.other)
			assert.Equal(t, tt.expected, tt.base)
		})
	}
}

func TestConfiguration_IsExpired(t *testing.T) {
	cfg := &Configuration{
		Session: SessionConfig{
			RetentionDays: 30,
		},
	}

	tests := []struct {
		name        string
		sessionTime time.Time
		retention   int
		expected    bool
	}{
		{
			name:        "not expired - recent session",
			sessionTime: time.Now().Add(-24 * time.Hour),
			retention:   30,
			expected:    false,
		},
		{
			name:        "expired - old session",
			sessionTime: time.Now().Add(-31 * 24 * time.Hour),
			retention:   30,
			expected:    true,
		},
		{
			name:        "not expired - exactly at limit",
			sessionTime: time.Now().Add(-30 * 24 * time.Hour).Add(time.Hour), // Add an hour to ensure not expired
			retention:   30,
			expected:    false,
		},
		{
			name:        "no expiration - zero retention",
			sessionTime: time.Now().Add(-365 * 24 * time.Hour),
			retention:   0,
			expected:    false,
		},
		{
			name:        "no expiration - negative retention",
			sessionTime: time.Now().Add(-365 * 24 * time.Hour),
			retention:   -1,
			expected:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg.Session.RetentionDays = tt.retention
			result := cfg.IsExpired(tt.sessionTime)
			assert.Equal(t, tt.expected, result)
		})
	}
}
