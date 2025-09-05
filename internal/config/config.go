package config

import (
	"time"
)

// Configuration represents the spcstr configuration structure
type Configuration struct {
	Version string        `json:"version"`
	Scope   string        `json:"scope"`
	Paths   PathConfig    `json:"paths"`
	UI      UIConfig      `json:"ui"`
	Session SessionConfig `json:"session"`
}

// PathConfig defines path-related configuration
type PathConfig struct {
	Hooks    string   `json:"hooks"`
	Sessions string   `json:"sessions"`
	Docs     []string `json:"docs"`
}

// UIConfig defines UI-related configuration
type UIConfig struct {
	Theme           string `json:"theme"`
	RefreshInterval int    `json:"refreshInterval"`
}

// SessionConfig defines session-related configuration
type SessionConfig struct {
	RetentionDays    int  `json:"retentionDays"`
	AutoArchive      bool `json:"autoArchive"`
	MaxActiveSession int  `json:"maxActiveSession"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Configuration {
	return &Configuration{
		Version: "1.0.0",
		Scope:   "project",
		Paths: PathConfig{
			Hooks:    ".spcstr/hooks",
			Sessions: ".spcstr/sessions",
			Docs:     []string{"docs/"},
		},
		UI: UIConfig{
			Theme:           "dark",
			RefreshInterval: 500,
		},
		Session: SessionConfig{
			RetentionDays:    30,
			AutoArchive:      true,
			MaxActiveSession: 10,
		},
	}
}

// GlobalConfig returns the default global configuration
func GlobalConfig() *Configuration {
	cfg := DefaultConfig()
	cfg.Scope = "global"
	return cfg
}

// Merge merges another configuration into this one
// The provided configuration takes precedence
func (c *Configuration) Merge(other *Configuration) {
	if other == nil {
		return
	}

	// Version is not merged, it's always from the base config

	// Merge Paths
	if other.Paths.Hooks != "" {
		c.Paths.Hooks = other.Paths.Hooks
	}
	if other.Paths.Sessions != "" {
		c.Paths.Sessions = other.Paths.Sessions
	}
	if len(other.Paths.Docs) > 0 {
		c.Paths.Docs = other.Paths.Docs
	}

	// Merge UI
	if other.UI.Theme != "" {
		c.UI.Theme = other.UI.Theme
	}
	if other.UI.RefreshInterval != 0 {
		c.UI.RefreshInterval = other.UI.RefreshInterval
	}

	// Merge Session
	if other.Session.RetentionDays != 0 {
		c.Session.RetentionDays = other.Session.RetentionDays
	}
	// Only update AutoArchive if other configuration explicitly has Session config
	// Check if the other Session struct was actually populated
	if other.Session.RetentionDays != 0 || other.Session.MaxActiveSession != 0 {
		c.Session.AutoArchive = other.Session.AutoArchive
	}
	if other.Session.MaxActiveSession != 0 {
		c.Session.MaxActiveSession = other.Session.MaxActiveSession
	}
}

// IsExpired checks if a session should be archived based on retention policy
func (c *Configuration) IsExpired(sessionTime time.Time) bool {
	if c.Session.RetentionDays <= 0 {
		return false // No expiration if retention days is 0 or negative
	}

	expirationTime := sessionTime.Add(time.Duration(c.Session.RetentionDays) * 24 * time.Hour)
	return time.Now().After(expirationTime)
}
