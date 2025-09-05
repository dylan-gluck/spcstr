package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// GetXDGConfigHome returns the XDG config home directory
func GetXDGConfigHome() string {
	// Check XDG_CONFIG_HOME environment variable
	if xdgHome := os.Getenv("XDG_CONFIG_HOME"); xdgHome != "" {
		return xdgHome
	}

	// Fall back to default based on OS
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home cannot be determined
		return ".config"
	}

	switch runtime.GOOS {
	case "darwin":
		// macOS uses ~/Library/Application Support by convention
		return filepath.Join(home, "Library", "Application Support")
	case "windows":
		// Windows uses %APPDATA%
		if appData := os.Getenv("APPDATA"); appData != "" {
			return appData
		}
		return filepath.Join(home, "AppData", "Roaming")
	default:
		// Linux and other Unix-like systems use ~/.config
		return filepath.Join(home, ".config")
	}
}

// GetXDGDataHome returns the XDG data home directory
func GetXDGDataHome() string {
	// Check XDG_DATA_HOME environment variable
	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
		return xdgData
	}

	// Fall back to default based on OS
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home cannot be determined
		return ".local/share"
	}

	switch runtime.GOOS {
	case "darwin":
		// macOS uses ~/Library/Application Support for data as well
		return filepath.Join(home, "Library", "Application Support")
	case "windows":
		// Windows uses %LOCALAPPDATA%
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return localAppData
		}
		return filepath.Join(home, "AppData", "Local")
	default:
		// Linux and other Unix-like systems use ~/.local/share
		return filepath.Join(home, ".local", "share")
	}
}

// GetXDGCacheHome returns the XDG cache home directory
func GetXDGCacheHome() string {
	// Check XDG_CACHE_HOME environment variable
	if xdgCache := os.Getenv("XDG_CACHE_HOME"); xdgCache != "" {
		return xdgCache
	}

	// Fall back to default based on OS
	home, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home cannot be determined
		return ".cache"
	}

	switch runtime.GOOS {
	case "darwin":
		// macOS uses ~/Library/Caches
		return filepath.Join(home, "Library", "Caches")
	case "windows":
		// Windows uses %TEMP% or %LOCALAPPDATA%\Temp
		if temp := os.Getenv("TEMP"); temp != "" {
			return temp
		}
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			return filepath.Join(localAppData, "Temp")
		}
		return filepath.Join(home, "AppData", "Local", "Temp")
	default:
		// Linux and other Unix-like systems use ~/.cache
		return filepath.Join(home, ".cache")
	}
}

// GetSpcstrConfigDir returns the spcstr-specific config directory
func GetSpcstrConfigDir() string {
	return filepath.Join(GetXDGConfigHome(), "spcstr")
}

// GetSpcstrDataDir returns the spcstr-specific data directory
func GetSpcstrDataDir() string {
	return filepath.Join(GetXDGDataHome(), "spcstr")
}

// GetSpcstrCacheDir returns the spcstr-specific cache directory
func GetSpcstrCacheDir() string {
	return filepath.Join(GetXDGCacheHome(), "spcstr")
}

// GetGlobalConfigPath returns the path to the global spcstr config file
func GetGlobalConfigPath() string {
	return filepath.Join(GetSpcstrConfigDir(), "config.json")
}

// GetProjectConfigPath returns the path to the project config file
func GetProjectConfigPath(projectRoot string) string {
	return filepath.Join(projectRoot, ".spcstr", "config.json")
}

// FindProjectRoot searches for .spcstr directory in current or parent directories
func FindProjectRoot(startPath string) (string, error) {
	current := startPath
	
	for {
		// Check if .spcstr exists in current directory
		spcstrPath := filepath.Join(current, ".spcstr")
		if info, err := os.Stat(spcstrPath); err == nil && info.IsDir() {
			return current, nil
		}

		// Move to parent directory
		parent := filepath.Dir(current)
		
		// Check if we've reached the root
		if parent == current {
			break
		}
		
		current = parent
	}

	return "", fmt.Errorf("no .spcstr directory found in current or parent directories")
}

// EnsureConfigDirs creates all necessary configuration directories
func EnsureConfigDirs() error {
	dirs := []string{
		GetSpcstrConfigDir(),
		GetSpcstrDataDir(),
		GetSpcstrCacheDir(),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("creating directory %s: %w", dir, err)
		}
	}

	return nil
}

// GetClaudeSettingsPath returns possible paths for Claude's settings.json
func GetClaudeSettingsPath() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{}
	}

	paths := []string{}

	switch runtime.GOOS {
	case "darwin":
		// macOS paths for Claude
		paths = append(paths,
			filepath.Join(home, "Library", "Application Support", "Claude", "settings.json"),
			filepath.Join(home, ".claude", "settings.json"),
		)
	case "windows":
		// Windows paths for Claude
		if appData := os.Getenv("APPDATA"); appData != "" {
			paths = append(paths, filepath.Join(appData, "Claude", "settings.json"))
		}
		paths = append(paths,
			filepath.Join(home, "AppData", "Roaming", "Claude", "settings.json"),
			filepath.Join(home, ".claude", "settings.json"),
		)
	default:
		// Linux and other Unix-like systems
		paths = append(paths,
			filepath.Join(home, ".config", "claude", "settings.json"),
			filepath.Join(home, ".claude", "settings.json"),
		)
	}

	return paths
}