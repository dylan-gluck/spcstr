package config

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidateConfig validates a configuration
func ValidateConfig(cfg *Configuration) error {
	if cfg == nil {
		return fmt.Errorf("configuration is nil")
	}

	// Validate version
	if cfg.Version == "" {
		return fmt.Errorf("version is required")
	}
	if !isValidVersion(cfg.Version) {
		return fmt.Errorf("invalid version format: %s", cfg.Version)
	}

	// Validate scope
	if cfg.Scope != "project" && cfg.Scope != "global" {
		return fmt.Errorf("scope must be either 'project' or 'global', got: %s", cfg.Scope)
	}

	// Validate paths
	if err := validatePaths(&cfg.Paths); err != nil {
		return fmt.Errorf("invalid paths: %w", err)
	}

	// Validate UI
	if err := validateUI(&cfg.UI); err != nil {
		return fmt.Errorf("invalid UI config: %w", err)
	}

	// Validate Session
	if err := validateSession(&cfg.Session); err != nil {
		return fmt.Errorf("invalid session config: %w", err)
	}

	return nil
}

func isValidVersion(version string) bool {
	// Simple semantic version validation (x.y.z)
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return false
	}

	for _, part := range parts {
		if part == "" {
			return false
		}
		// Check if it's a valid number
		for _, ch := range part {
			if ch < '0' || ch > '9' {
				return false
			}
		}
	}
	return true
}

func validatePaths(paths *PathConfig) error {
	if paths.Hooks == "" {
		return fmt.Errorf("hooks path is required")
	}
	if !isValidPath(paths.Hooks) {
		return fmt.Errorf("invalid hooks path: %s", paths.Hooks)
	}

	if paths.Sessions == "" {
		return fmt.Errorf("sessions path is required")
	}
	if !isValidPath(paths.Sessions) {
		return fmt.Errorf("invalid sessions path: %s", paths.Sessions)
	}

	if len(paths.Docs) == 0 {
		return fmt.Errorf("at least one docs path is required")
	}
	for _, doc := range paths.Docs {
		if !isValidPath(doc) {
			return fmt.Errorf("invalid docs path: %s", doc)
		}
	}

	return nil
}

func isValidPath(path string) bool {
	// Basic path validation
	if path == "" {
		return false
	}

	// Check for dangerous path patterns
	if strings.Contains(path, "..") {
		return false // No parent directory traversal
	}

	// Clean the path and check if it's valid
	cleaned := filepath.Clean(path)
	// Allow original path, cleaned path, or paths with trailing slash
	return cleaned == path || cleaned+"/" == path || path == cleaned || "./"+cleaned == path
}

func validateUI(ui *UIConfig) error {
	// Validate theme
	validThemes := []string{"dark", "light", "auto"}
	if !contains(validThemes, ui.Theme) {
		return fmt.Errorf("invalid theme: %s (must be one of: %s)", ui.Theme, strings.Join(validThemes, ", "))
	}

	// Validate refresh interval
	if ui.RefreshInterval < 100 {
		return fmt.Errorf("refresh interval must be at least 100ms, got: %d", ui.RefreshInterval)
	}
	if ui.RefreshInterval > 10000 {
		return fmt.Errorf("refresh interval must be at most 10000ms, got: %d", ui.RefreshInterval)
	}

	return nil
}

func validateSession(session *SessionConfig) error {
	// Validate retention days
	if session.RetentionDays < 0 {
		return fmt.Errorf("retention days cannot be negative, got: %d", session.RetentionDays)
	}
	if session.RetentionDays > 365 {
		return fmt.Errorf("retention days cannot exceed 365, got: %d", session.RetentionDays)
	}

	// Validate max active sessions
	if session.MaxActiveSession < 1 {
		return fmt.Errorf("max active sessions must be at least 1, got: %d", session.MaxActiveSession)
	}
	if session.MaxActiveSession > 100 {
		return fmt.Errorf("max active sessions cannot exceed 100, got: %d", session.MaxActiveSession)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// IsValidConfigPath checks if a path is valid for configuration file
func IsValidConfigPath(path string) error {
	if path == "" {
		return fmt.Errorf("config path cannot be empty")
	}

	// Must end with .json
	if !strings.HasSuffix(path, ".json") {
		return fmt.Errorf("config file must have .json extension")
	}

	// Check for directory traversal
	if strings.Contains(path, "..") {
		return fmt.Errorf("config path cannot contain parent directory references")
	}

	return nil
}
