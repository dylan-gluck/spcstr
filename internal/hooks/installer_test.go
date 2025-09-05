package hooks

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewInstaller(t *testing.T) {
	installer := NewInstaller("/hooks", "/project")
	assert.NotNil(t, installer)
	assert.Equal(t, "/hooks", installer.hooksPath)
	assert.Equal(t, "/project", installer.projectPath)
}

func TestInstaller_InstallHooks(t *testing.T) {
	// Create temp directories
	tmpDir := t.TempDir()
	hooksPath := filepath.Join(tmpDir, "hooks")
	projectPath := filepath.Join(tmpDir, "project")
	
	installer := NewInstaller(hooksPath, projectPath)
	
	// Install hooks
	err := installer.InstallHooks()
	assert.NoError(t, err)
	
	// Verify all hooks are installed
	err = installer.VerifyHooks()
	assert.NoError(t, err)
}

func TestInstaller_VerifyHooks(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	hooksPath := filepath.Join(tmpDir, "hooks")
	projectPath := filepath.Join(tmpDir, "project")
	
	installer := NewInstaller(hooksPath, projectPath)
	
	// Verify before installation should fail
	err := installer.VerifyHooks()
	assert.Error(t, err)
	
	// Install hooks
	err = installer.InstallHooks()
	require.NoError(t, err)
	
	// Verify after installation should succeed
	err = installer.VerifyHooks()
	assert.NoError(t, err)
	
	// Remove a hook and verify should fail
	err = os.Remove(filepath.Join(hooksPath, "pre-command.sh"))
	require.NoError(t, err)
	
	err = installer.VerifyHooks()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pre-command.sh not found")
}

func TestInstaller_VerifyHooks_Permissions(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	hooksPath := filepath.Join(tmpDir, "hooks")
	projectPath := filepath.Join(tmpDir, "project")
	
	installer := NewInstaller(hooksPath, projectPath)
	
	// Install hooks
	err := installer.InstallHooks()
	require.NoError(t, err)
	
	// Make a hook non-executable
	hookPath := filepath.Join(hooksPath, "pre-command.sh")
	err = os.Chmod(hookPath, 0644)
	require.NoError(t, err)
	
	// Verify should fail
	err = installer.VerifyHooks()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not executable")
}

func TestInstaller_VerifyHooks_Directory(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	hooksPath := filepath.Join(tmpDir, "hooks")
	projectPath := filepath.Join(tmpDir, "project")
	
	installer := NewInstaller(hooksPath, projectPath)
	
	// Create hooks directory
	err := os.MkdirAll(hooksPath, 0755)
	require.NoError(t, err)
	
	// Create a directory instead of a file
	dirPath := filepath.Join(hooksPath, "pre-command.sh")
	err = os.MkdirAll(dirPath, 0755)
	require.NoError(t, err)
	
	// Verify should fail
	err = installer.VerifyHooks()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is a directory")
}

func TestInstaller_UninstallHooks(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	hooksPath := filepath.Join(tmpDir, "hooks")
	projectPath := filepath.Join(tmpDir, "project")
	
	installer := NewInstaller(hooksPath, projectPath)
	
	// Install hooks
	err := installer.InstallHooks()
	require.NoError(t, err)
	
	// Verify hooks exist
	err = installer.VerifyHooks()
	require.NoError(t, err)
	
	// Uninstall hooks
	err = installer.UninstallHooks()
	assert.NoError(t, err)
	
	// Verify hooks don't exist
	err = installer.VerifyHooks()
	assert.Error(t, err)
}

func TestInstaller_UninstallHooks_NonExistent(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	hooksPath := filepath.Join(tmpDir, "hooks")
	projectPath := filepath.Join(tmpDir, "project")
	
	installer := NewInstaller(hooksPath, projectPath)
	
	// Uninstall without installing first
	err := installer.UninstallHooks()
	// Should not error for non-existent files
	assert.NoError(t, err)
}

func TestInstaller_GetHookPath(t *testing.T) {
	installer := NewInstaller("/hooks", "/project")
	
	path := installer.GetHookPath("pre-command.sh")
	assert.Equal(t, "/hooks/pre-command.sh", path)
	
	path = installer.GetHookPath("custom.sh")
	assert.Equal(t, "/hooks/custom.sh", path)
}

func TestInstaller_ListHooks(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	hooksPath := filepath.Join(tmpDir, "hooks")
	projectPath := filepath.Join(tmpDir, "project")
	
	installer := NewInstaller(hooksPath, projectPath)
	
	// List before installation
	hooks, err := installer.ListHooks()
	assert.NoError(t, err)
	assert.Empty(t, hooks)
	
	// Install hooks
	err = installer.InstallHooks()
	require.NoError(t, err)
	
	// List after installation
	hooks, err = installer.ListHooks()
	assert.NoError(t, err)
	assert.Len(t, hooks, 5)
	
	// Check that all expected hooks are listed
	hookMap := make(map[string]bool)
	for _, hook := range hooks {
		hookMap[hook] = true
	}
	
	assert.True(t, hookMap["common.sh"])
	assert.True(t, hookMap["pre-command.sh"])
	assert.True(t, hookMap["post-command.sh"])
	assert.True(t, hookMap["file-modified.sh"])
	assert.True(t, hookMap["session-end.sh"])
}

func TestInstaller_ListHooks_WithExtraFiles(t *testing.T) {
	// Create temp directory
	tmpDir := t.TempDir()
	hooksPath := filepath.Join(tmpDir, "hooks")
	projectPath := filepath.Join(tmpDir, "project")
	
	installer := NewInstaller(hooksPath, projectPath)
	
	// Install hooks
	err := installer.InstallHooks()
	require.NoError(t, err)
	
	// Add extra files
	err = os.WriteFile(filepath.Join(hooksPath, "custom.sh"), []byte("custom"), 0755)
	require.NoError(t, err)
	
	err = os.WriteFile(filepath.Join(hooksPath, "README.md"), []byte("readme"), 0644)
	require.NoError(t, err)
	
	err = os.MkdirAll(filepath.Join(hooksPath, "subdir"), 0755)
	require.NoError(t, err)
	
	// List hooks
	hooks, err := installer.ListHooks()
	assert.NoError(t, err)
	
	// Should include .sh files but not others
	assert.Contains(t, hooks, "custom.sh")
	assert.NotContains(t, hooks, "README.md")
	assert.NotContains(t, hooks, "subdir")
}