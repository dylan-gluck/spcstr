package config

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const DEFAULT_TIMEOUT = 30 * time.Second

// InitializeProject initializes a project for spcstr usage
func InitializeProject(force bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()

	// Get current working directory
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Check if .spcstr already exists
	spcstrDir := filepath.Join(projectRoot, ".spcstr")
	if dirExists(spcstrDir) && !force {
		fmt.Printf("Directory .spcstr already exists in %s\n", projectRoot)
		
		// Prompt for confirmation
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Do you want to reinitialize? This will preserve existing data. (y/N): ")
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read user input: %w", err)
		}
		
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Initialization cancelled.")
			return nil
		}
	}

	// Create directory structure
	if err := createDirectoryStructure(ctx, projectRoot); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	// Configure Claude Code hooks
	if err := configureClaudeHooks(ctx, projectRoot); err != nil {
		return fmt.Errorf("failed to configure Claude Code hooks: %w", err)
	}

	fmt.Printf("✓ Successfully initialized spcstr in %s\n", projectRoot)
	fmt.Println("✓ Created .spcstr/logs and .spcstr/sessions directories")
	fmt.Println("✓ Configured Claude Code hooks in .claude/settings.json")
	fmt.Println("\nYour project is now ready for Claude Code session tracking!")
	
	return nil
}

// createDirectoryStructure creates the .spcstr directory structure
func createDirectoryStructure(ctx context.Context, projectRoot string) error {
	// Check context before operations
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Create .spcstr/logs directory
	logsDir := filepath.Join(projectRoot, ".spcstr", "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create .spcstr/sessions directory
	sessionsDir := filepath.Join(projectRoot, ".spcstr", "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	return nil
}

// configureClaudeHooks configures hooks in .claude/settings.json
func configureClaudeHooks(ctx context.Context, projectRoot string) error {
	// Check context before operations
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Ensure .claude directory exists
	claudeDir := filepath.Join(projectRoot, ".claude")
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .claude directory: %w", err)
	}

	settingsPath := filepath.Join(claudeDir, "settings.json")
	
	// Read existing settings if file exists
	var settings map[string]interface{}
	if fileExists(settingsPath) {
		data, err := os.ReadFile(settingsPath)
		if err != nil {
			return fmt.Errorf("failed to read existing settings.json: %w", err)
		}
		
		if len(data) > 0 {
			if err := json.Unmarshal(data, &settings); err != nil {
				return fmt.Errorf("failed to parse existing settings.json: %w", err)
			}
		}
	}
	
	if settings == nil {
		settings = make(map[string]interface{})
	}

	// Define hook commands
	hooks := map[string]string{
		"session_start":      `spcstr hook session_start --cwd="${CLAUDE_PROJECT_DIR}"`,
		"user_prompt_submit": `spcstr hook user_prompt_submit --cwd="${CLAUDE_PROJECT_DIR}"`,
		"pre_tool_use":       `spcstr hook pre_tool_use --cwd="${CLAUDE_PROJECT_DIR}"`,
		"post_tool_use":      `spcstr hook post_tool_use --cwd="${CLAUDE_PROJECT_DIR}"`,
		"notification":       `spcstr hook notification --cwd="${CLAUDE_PROJECT_DIR}"`,
		"pre_compact":        `spcstr hook pre_compact --cwd="${CLAUDE_PROJECT_DIR}"`,
		"session_end":        `spcstr hook session_end --cwd="${CLAUDE_PROJECT_DIR}"`,
		"stop":               `spcstr hook stop --cwd="${CLAUDE_PROJECT_DIR}"`,
		"subagent_stop":      `spcstr hook subagent_stop --cwd="${CLAUDE_PROJECT_DIR}"`,
	}

	// Set hooks in settings
	settings["hooks"] = hooks

	// Write settings using atomic operation (temp file + rename)
	if err := writeSettingsAtomic(ctx, settingsPath, settings); err != nil {
		return fmt.Errorf("failed to write settings.json: %w", err)
	}

	return nil
}

// writeSettingsAtomic writes settings.json using atomic file operation
func writeSettingsAtomic(ctx context.Context, path string, settings map[string]interface{}) error {
	// Check context before operations
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Marshal settings with indentation
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Create temp file
	tempFile, err := os.CreateTemp(filepath.Dir(path), ".settings-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()
	
	// Ensure cleanup
	defer func() {
		tempFile.Close()
		os.Remove(tempPath) // Clean up temp file if it still exists
	}()

	// Write data to temp file
	if _, err := tempFile.Write(data); err != nil {
		return fmt.Errorf("failed to write to temp file: %w", err)
	}
	
	// Add newline at end of file
	if _, err := tempFile.Write([]byte("\n")); err != nil {
		return fmt.Errorf("failed to write newline: %w", err)
	}

	// Sync to disk
	if err := tempFile.Sync(); err != nil {
		return fmt.Errorf("failed to sync temp file: %w", err)
	}
	
	// Close temp file before rename
	if err := tempFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("failed to rename temp file to settings.json: %w", err)
	}

	return nil
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}