package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LoadConfig loads configuration from a file
func LoadConfig(path string) (*Configuration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening config file: %w", err)
	}
	defer file.Close()

	return LoadConfigFromReader(file)
}

// LoadConfigFromReader loads configuration from an io.Reader
func LoadConfigFromReader(r io.Reader) (*Configuration, error) {
	var cfg Configuration
	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decoding config JSON: %w", err)
	}

	return &cfg, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(path string, cfg *Configuration) error {
	// Validate configuration before saving
	if err := ValidateConfig(cfg); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	// Write to temporary file first (atomic write)
	tempPath := path + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("creating temp config file: %w", err)
	}

	// Ensure temp file is removed on error
	defer func() {
		if err != nil {
			os.Remove(tempPath)
		}
	}()

	// Write configuration
	if err = SaveConfigToWriter(file, cfg); err != nil {
		file.Close()
		return err
	}

	// Close file before rename
	if err = file.Close(); err != nil {
		return fmt.Errorf("closing temp config file: %w", err)
	}

	// Set file permissions
	if err = os.Chmod(tempPath, 0644); err != nil {
		return fmt.Errorf("setting config file permissions: %w", err)
	}

	// Atomic rename
	if err = os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("renaming config file: %w", err)
	}

	return nil
}

// SaveConfigToWriter saves configuration to an io.Writer
func SaveConfigToWriter(w io.Writer, cfg *Configuration) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cfg); err != nil {
		return fmt.Errorf("encoding config to JSON: %w", err)
	}
	return nil
}

// ConfigExists checks if a configuration file exists
func ConfigExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// MergeConfigs merges multiple configurations
// Later configurations take precedence
func MergeConfigs(configs ...*Configuration) *Configuration {
	if len(configs) == 0 {
		return DefaultConfig()
	}

	result := configs[0]
	if result == nil {
		result = DefaultConfig()
	}

	for i := 1; i < len(configs); i++ {
		result.Merge(configs[i])
	}

	return result
}

// LoadOrCreateConfig loads config from path or creates default if not exists
func LoadOrCreateConfig(path string) (*Configuration, error) {
	if ConfigExists(path) {
		return LoadConfig(path)
	}

	// Create default config
	cfg := DefaultConfig()
	if err := SaveConfig(path, cfg); err != nil {
		return nil, fmt.Errorf("creating default config: %w", err)
	}

	return cfg, nil
}
