package hooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator("/hooks", "/project")
	assert.NotNil(t, gen)
	assert.Equal(t, "/hooks", gen.hooksPath)
	assert.Equal(t, "/project", gen.projectPath)
}

func TestGenerator_GenerateHooks(t *testing.T) {
	// Create temp directories
	tmpDir := t.TempDir()
	hooksPath := filepath.Join(tmpDir, "hooks")
	projectPath := filepath.Join(tmpDir, "project")
	
	gen := NewGenerator(hooksPath, projectPath)
	
	// Generate hooks
	err := gen.GenerateHooks()
	assert.NoError(t, err)
	
	// Check that all hooks were created
	expectedHooks := []string{
		"common.sh",
		"pre-command.sh",
		"post-command.sh",
		"file-modified.sh",
		"session-end.sh",
	}
	
	for _, hook := range expectedHooks {
		hookPath := filepath.Join(hooksPath, hook)
		info, err := os.Stat(hookPath)
		require.NoError(t, err, "Hook %s should exist", hook)
		assert.False(t, info.IsDir(), "Hook %s should be a file", hook)
		
		// Check permissions (should be executable except common.sh)
		if hook != "common.sh" {
			assert.Equal(t, os.FileMode(0755), info.Mode().Perm(), 
				"Hook %s should be executable", hook)
		}
		
		// Read and verify content
		content, err := os.ReadFile(hookPath)
		require.NoError(t, err)
		contentStr := string(content)
		
		// All hooks should start with shebang
		assert.True(t, strings.HasPrefix(contentStr, "#!/bin/sh\n"), 
			"Hook %s should start with shebang", hook)
		
		// All hooks except common.sh should have set +e
		// (common.sh is sourced, not executed directly)
		assert.Contains(t, contentStr, "set +e", 
			"Hook %s should have 'set +e'", hook)
		
		// All hooks except common.sh should exit 0
		if hook != "common.sh" {
			assert.Contains(t, contentStr, "exit 0", 
				"Hook %s should exit 0", hook)
		}
	}
	
	// Verify hook-specific content
	t.Run("pre-command hook", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(hooksPath, "pre-command.sh"))
		require.NoError(t, err)
		contentStr := string(content)
		assert.Contains(t, contentStr, "Pre-command hook triggered")
		assert.Contains(t, contentStr, "CLAUDE_SESSION_ID")
		assert.Contains(t, contentStr, "CLAUDE_COMMAND")
	})
	
	t.Run("post-command hook", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(hooksPath, "post-command.sh"))
		require.NoError(t, err)
		contentStr := string(content)
		assert.Contains(t, contentStr, "Post-command hook triggered")
		assert.Contains(t, contentStr, "CLAUDE_EXIT_CODE")
	})
	
	t.Run("file-modified hook", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(hooksPath, "file-modified.sh"))
		require.NoError(t, err)
		contentStr := string(content)
		assert.Contains(t, contentStr, "File-modified hook triggered")
		assert.Contains(t, contentStr, "CLAUDE_MODIFIED_FILE")
	})
	
	t.Run("session-end hook", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(hooksPath, "session-end.sh"))
		require.NoError(t, err)
		contentStr := string(content)
		assert.Contains(t, contentStr, "Session-end hook triggered")
		assert.Contains(t, contentStr, "archive")
	})
	
	t.Run("common.sh", func(t *testing.T) {
		content, err := os.ReadFile(filepath.Join(hooksPath, "common.sh"))
		require.NoError(t, err)
		contentStr := string(content)
		assert.Contains(t, contentStr, "log_error")
		assert.Contains(t, contentStr, "log_info")
		assert.Contains(t, contentStr, "create_session_data")
	})
}

func TestGenerator_GenerateHooks_Overwrite(t *testing.T) {
	// Create temp directories
	tmpDir := t.TempDir()
	hooksPath := filepath.Join(tmpDir, "hooks")
	projectPath := filepath.Join(tmpDir, "project")
	
	gen := NewGenerator(hooksPath, projectPath)
	
	// Generate hooks first time
	err := gen.GenerateHooks()
	require.NoError(t, err)
	
	// Modify a hook
	hookPath := filepath.Join(hooksPath, "pre-command.sh")
	err = os.WriteFile(hookPath, []byte("modified content"), 0755)
	require.NoError(t, err)
	
	// Generate hooks again (should overwrite)
	err = gen.GenerateHooks()
	assert.NoError(t, err)
	
	// Check that hook was regenerated
	content, err := os.ReadFile(hookPath)
	require.NoError(t, err)
	assert.NotEqual(t, "modified content", string(content))
	assert.Contains(t, string(content), "Pre-command hook triggered")
}

func TestGenerator_writeFile(t *testing.T) {
	tmpDir := t.TempDir()
	gen := NewGenerator(tmpDir, tmpDir)
	
	// Test writing a simple template
	template := "Hello {{.ProjectPath}}"
	data := struct{ ProjectPath string }{ProjectPath: "/test/path"}
	filePath := filepath.Join(tmpDir, "test.txt")
	
	err := gen.writeFile(filePath, template, data)
	assert.NoError(t, err)
	
	// Verify content
	content, err := os.ReadFile(filePath)
	require.NoError(t, err)
	assert.Equal(t, "Hello /test/path", string(content))
	
	// Verify permissions
	info, err := os.Stat(filePath)
	require.NoError(t, err)
	assert.Equal(t, os.FileMode(0755), info.Mode().Perm())
}

func TestGenerator_writeFile_InvalidTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	gen := NewGenerator(tmpDir, tmpDir)
	
	// Test with invalid template
	template := "{{.InvalidSyntax"
	filePath := filepath.Join(tmpDir, "test.txt")
	
	err := gen.writeFile(filePath, template, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "parsing template")
	
	// File should not exist
	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}