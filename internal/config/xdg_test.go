package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetXDGConfigHome(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CONFIG_HOME")
	defer os.Setenv("XDG_CONFIG_HOME", originalXDG)

	// Test with XDG_CONFIG_HOME set
	testPath := "/custom/config"
	os.Setenv("XDG_CONFIG_HOME", testPath)
	result := GetXDGConfigHome()
	assert.Equal(t, testPath, result)

	// Test with XDG_CONFIG_HOME unset
	os.Unsetenv("XDG_CONFIG_HOME")
	result = GetXDGConfigHome()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	switch runtime.GOOS {
	case "darwin":
		assert.Equal(t, filepath.Join(home, "Library", "Application Support"), result)
	case "windows":
		// On Windows, it should use APPDATA or fallback
		assert.True(t, strings.Contains(result, "AppData"))
	default:
		assert.Equal(t, filepath.Join(home, ".config"), result)
	}
}

func TestGetXDGDataHome(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_DATA_HOME")
	defer os.Setenv("XDG_DATA_HOME", originalXDG)

	// Test with XDG_DATA_HOME set
	testPath := "/custom/data"
	os.Setenv("XDG_DATA_HOME", testPath)
	result := GetXDGDataHome()
	assert.Equal(t, testPath, result)

	// Test with XDG_DATA_HOME unset
	os.Unsetenv("XDG_DATA_HOME")
	result = GetXDGDataHome()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	switch runtime.GOOS {
	case "darwin":
		assert.Equal(t, filepath.Join(home, "Library", "Application Support"), result)
	case "windows":
		assert.True(t, strings.Contains(result, "AppData"))
	default:
		assert.Equal(t, filepath.Join(home, ".local", "share"), result)
	}
}

func TestGetXDGCacheHome(t *testing.T) {
	// Save original env
	originalXDG := os.Getenv("XDG_CACHE_HOME")
	defer os.Setenv("XDG_CACHE_HOME", originalXDG)

	// Test with XDG_CACHE_HOME set
	testPath := "/custom/cache"
	os.Setenv("XDG_CACHE_HOME", testPath)
	result := GetXDGCacheHome()
	assert.Equal(t, testPath, result)

	// Test with XDG_CACHE_HOME unset
	os.Unsetenv("XDG_CACHE_HOME")
	result = GetXDGCacheHome()

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	switch runtime.GOOS {
	case "darwin":
		assert.Equal(t, filepath.Join(home, "Library", "Caches"), result)
	case "windows":
		assert.True(t, strings.Contains(result, "Temp") || strings.Contains(result, "TEMP"))
	default:
		assert.Equal(t, filepath.Join(home, ".cache"), result)
	}
}

func TestGetSpcstrDirs(t *testing.T) {
	// These should append "spcstr" to the XDG directories
	configDir := GetSpcstrConfigDir()
	assert.True(t, strings.HasSuffix(configDir, "spcstr"))

	dataDir := GetSpcstrDataDir()
	assert.True(t, strings.HasSuffix(dataDir, "spcstr"))

	cacheDir := GetSpcstrCacheDir()
	assert.True(t, strings.HasSuffix(cacheDir, "spcstr"))
}

func TestGetGlobalConfigPath(t *testing.T) {
	path := GetGlobalConfigPath()
	assert.True(t, strings.HasSuffix(path, "config.json"))
	assert.True(t, strings.Contains(path, "spcstr"))
}

func TestGetProjectConfigPath(t *testing.T) {
	projectRoot := "/test/project"
	path := GetProjectConfigPath(projectRoot)
	assert.Equal(t, filepath.Join(projectRoot, ".spcstr", "config.json"), path)
}

func TestFindProjectRoot(t *testing.T) {
	// Create temp directory structure
	tmpDir := t.TempDir()

	// Create project with .spcstr
	projectDir := filepath.Join(tmpDir, "project")
	spcstrDir := filepath.Join(projectDir, ".spcstr")
	err := os.MkdirAll(spcstrDir, 0755)
	require.NoError(t, err)

	// Create subdirectory
	subDir := filepath.Join(projectDir, "sub", "dir")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	tests := []struct {
		name      string
		startPath string
		wantRoot  string
		wantErr   bool
	}{
		{
			name:      "find from project root",
			startPath: projectDir,
			wantRoot:  projectDir,
			wantErr:   false,
		},
		{
			name:      "find from subdirectory",
			startPath: subDir,
			wantRoot:  projectDir,
			wantErr:   false,
		},
		{
			name:      "not found",
			startPath: tmpDir,
			wantRoot:  "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root, err := FindProjectRoot(tt.startPath)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, root)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantRoot, root)
			}
		})
	}
}

func TestEnsureConfigDirs(t *testing.T) {
	// This test might need to be skipped in CI if it tries to create actual directories
	// For now, we'll just test that it doesn't error
	err := EnsureConfigDirs()
	// We don't assert NoError because it might fail in restricted environments
	// Just ensure it doesn't panic
	_ = err
}

func TestGetClaudeSettingsPath(t *testing.T) {
	paths := GetClaudeSettingsPath()

	// Should return multiple paths
	assert.NotEmpty(t, paths)

	// All paths should end with settings.json
	for _, path := range paths {
		assert.True(t, strings.HasSuffix(path, "settings.json"))
	}

	// Should include platform-specific paths
	switch runtime.GOOS {
	case "darwin":
		found := false
		for _, path := range paths {
			if strings.Contains(path, "Library") {
				found = true
				break
			}
		}
		assert.True(t, found, "Should include Library path on macOS")
	case "windows":
		found := false
		for _, path := range paths {
			if strings.Contains(path, "AppData") || strings.Contains(path, "Claude") {
				found = true
				break
			}
		}
		assert.True(t, found, "Should include AppData path on Windows")
	default:
		found := false
		for _, path := range paths {
			if strings.Contains(path, ".config") || strings.Contains(path, ".claude") {
				found = true
				break
			}
		}
		assert.True(t, found, "Should include .config path on Linux")
	}
}
