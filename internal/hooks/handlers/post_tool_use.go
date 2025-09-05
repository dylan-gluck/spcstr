package handlers

import (
	"github.com/dylan-gluck/spcstr/internal/session"
	"log/slog"
	"path/filepath"

	models "github.com/dylan-gluck/spcstr/pkg/hooks"
)

func HandlePostToolUse(projectRoot string, data map[string]interface{}) error {
	toolName, _ := data["tool_name"].(string)
	slog.Debug("post_tool_use", "tool", toolName)

	state, err := session.LoadSessionState(projectRoot)
	if err != nil {
		slog.Error("failed to load state", "error", err)
		return nil
	}

	if state.ToolsUsed == nil {
		state.ToolsUsed = make(map[string]int)
	}
	state.ToolsUsed[toolName]++

	if params, ok := data["parameters"].(map[string]interface{}); ok {
		switch toolName {
		case "Read":
			if filePath, ok := params["file_path"].(string); ok {
				absPath := makeAbsolutePath(projectRoot, filePath)
				state.Files.Read = models.AddUniqueString(state.Files.Read, absPath)
			}
		case "Write":
			if filePath, ok := params["file_path"].(string); ok {
				absPath := makeAbsolutePath(projectRoot, filePath)
				isNew := !fileExists(absPath)
				if isNew {
					state.Files.New = models.AddUniqueString(state.Files.New, absPath)
				} else {
					state.Files.Edited = models.AddUniqueString(state.Files.Edited, absPath)
				}
			}
		case "Edit", "MultiEdit":
			if filePath, ok := params["file_path"].(string); ok {
				absPath := makeAbsolutePath(projectRoot, filePath)
				state.Files.Edited = models.AddUniqueString(state.Files.Edited, absPath)
			}
		}
	}

	if err := session.SaveSessionState(projectRoot, state); err != nil {
		slog.Error("failed to save state", "error", err)
	}

	return nil
}

func makeAbsolutePath(projectRoot, filePath string) string {
	if filepath.IsAbs(filePath) {
		return filepath.Clean(filePath)
	}
	return filepath.Clean(filepath.Join(projectRoot, filePath))
}

func fileExists(path string) bool {
	_, err := filepath.EvalSymlinks(path)
	return err == nil
}
