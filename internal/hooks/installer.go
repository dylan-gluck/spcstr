package hooks

import (
	"fmt"
	"os"
	"path/filepath"
)

// Installer handles hook installation
type Installer struct {
	hooksPath   string
	projectPath string
}

// NewInstaller creates a new hook installer
func NewInstaller(hooksPath, projectPath string) *Installer {
	return &Installer{
		hooksPath:   hooksPath,
		projectPath: projectPath,
	}
}

// InstallHooks installs hooks in the target directory
func (i *Installer) InstallHooks() error {
	// For now, InstallHooks just ensures the hooks are generated
	// The actual installation happens via the Generator
	// This separation allows for future enhancements like:
	// - Installing to different locations
	// - Symlinking instead of copying
	// - Version checking before installation
	
	generator := NewGenerator(i.hooksPath, i.projectPath)
	return generator.GenerateHooks()
}

// VerifyHooks checks if all required hooks are installed
func (i *Installer) VerifyHooks() error {
	requiredHooks := []string{
		"pre-command.sh",
		"post-command.sh",
		"file-modified.sh",
		"session-end.sh",
		"common.sh",
	}

	for _, hook := range requiredHooks {
		hookPath := filepath.Join(i.hooksPath, hook)
		info, err := os.Stat(hookPath)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("hook %s not found", hook)
			}
			return fmt.Errorf("checking hook %s: %w", hook, err)
		}

		// Verify it's a file and executable (except common.sh)
		if info.IsDir() {
			return fmt.Errorf("%s is a directory, expected file", hook)
		}

		if hook != "common.sh" {
			// Check if executable
			if info.Mode().Perm()&0111 == 0 {
				return fmt.Errorf("hook %s is not executable", hook)
			}
		}
	}

	return nil
}

// UninstallHooks removes installed hooks
func (i *Installer) UninstallHooks() error {
	// Remove all hook files
	hooks := []string{
		"pre-command.sh",
		"post-command.sh",
		"file-modified.sh",
		"session-end.sh",
		"common.sh",
	}

	var lastErr error
	for _, hook := range hooks {
		hookPath := filepath.Join(i.hooksPath, hook)
		if err := os.Remove(hookPath); err != nil && !os.IsNotExist(err) {
			lastErr = fmt.Errorf("removing %s: %w", hook, err)
		}
	}

	return lastErr
}

// GetHookPath returns the full path to a specific hook
func (i *Installer) GetHookPath(hookName string) string {
	return filepath.Join(i.hooksPath, hookName)
}

// ListHooks returns a list of installed hooks
func (i *Installer) ListHooks() ([]string, error) {
	entries, err := os.ReadDir(i.hooksPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, fmt.Errorf("reading hooks directory: %w", err)
	}

	var hooks []string
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".sh" {
			hooks = append(hooks, entry.Name())
		}
	}

	return hooks, nil
}