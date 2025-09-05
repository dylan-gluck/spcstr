package handlers

import (
	"fmt"
	"github.com/dylan-gluck/spcstr/internal/session"
	"log/slog"
	"strings"
)

func HandlePreToolUse(projectRoot string, data map[string]interface{}) error {
	toolName, _ := data["tool_name"].(string)
	slog.Debug("pre_tool_use", "tool", toolName)

	if toolName == "Bash" {
		if params, ok := data["parameters"].(map[string]interface{}); ok {
			if command, ok := params["command"].(string); ok {
				if isDangerousCommand(command) {
					slog.Warn("blocking dangerous command", "command", command)
					return fmt.Errorf("blocked: dangerous command attempted: %s", command)
				}
			}
		}
	}

	state, err := session.LoadSessionState(projectRoot)
	if err != nil {
		slog.Error("failed to load state", "error", err)
		return nil
	}

	if state.ToolsUsed == nil {
		state.ToolsUsed = make(map[string]int)
	}

	if err := session.SaveSessionState(projectRoot, state); err != nil {
		slog.Error("failed to save state", "error", err)
	}

	return nil
}

func isDangerousCommand(command string) bool {
	dangerous := []string{
		"rm -rf /",
		"dd if=/dev/zero",
		"mkfs",
		"format",
		"> /dev/sda",
		"chmod -R 777 /",
		"curl | sh",
		"wget | sh",
		"eval",
	}

	cmd := strings.ToLower(command)
	for _, d := range dangerous {
		if strings.Contains(cmd, strings.ToLower(d)) {
			return true
		}
	}

	return false
}
