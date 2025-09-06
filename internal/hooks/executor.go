package hooks

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ExecuteHook executes a hook in the context of a project directory
func ExecuteHook(hookName string, projectDir string, input []byte) error {
	// 1. Validate project directory
	if !isValidSpcstrProject(projectDir) {
		return fmt.Errorf("invalid spcstr project directory: %s", projectDir)
	}

	// 2. Change to project directory
	oldDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}

	err = os.Chdir(projectDir)
	if err != nil {
		return fmt.Errorf("failed to change to project directory '%s': %w", projectDir, err)
	}
	defer func() {
		os.Chdir(oldDir)
	}()

	// 3. Parse input to get session ID for logging
	var inputData map[string]interface{}
	sessionID := ""
	if err := json.Unmarshal(input, &inputData); err == nil {
		if id, ok := inputData["session_id"].(string); ok {
			sessionID = id
		}
	}

	// 4. Execute hook in project context
	err = DefaultRegistry.Execute(hookName, input)

	// 5. Log the event
	success := err == nil
	if logErr := DefaultLogger.LogEvent(sessionID, hookName, inputData, success); logErr != nil {
		// Don't fail the hook execution due to logging issues, but print a warning
		fmt.Fprintf(os.Stderr, "Warning: Failed to log hook event: %v\n", logErr)
	}

	return err
}

// isValidSpcstrProject checks if the directory contains valid .spcstr structure
func isValidSpcstrProject(dir string) bool {
	sessionsPath := filepath.Join(dir, ".spcstr", "sessions")
	logsPath := filepath.Join(dir, ".spcstr", "logs")

	return dirExists(sessionsPath) && dirExists(logsPath)
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
