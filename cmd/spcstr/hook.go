package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/dylan-gluck/spcstr/internal/hooks"
)

func runHook(hookName string) int {
	projectRoot := os.Getenv("CLAUDE_PROJECT_PATH")
	if projectRoot == "" {
		projectRoot = "."
	}

	logFile := filepath.Join(projectRoot, ".spcstr", "logs", "debug.log")
	if err := os.MkdirAll(filepath.Dir(logFile), 0755); err != nil {
		return 0
	}

	lf, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		defer lf.Close()
		logger := slog.New(slog.NewTextHandler(lf, &slog.HandlerOptions{Level: slog.LevelDebug}))
		slog.SetDefault(logger)
	}

	slog.Debug("hook invoked", "hook", hookName, "project_root", projectRoot)

	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		slog.Error("failed to read stdin", "error", err)
		return 0
	}

	var hookData map[string]interface{}
	if len(stdin) > 0 && strings.TrimSpace(string(stdin)) != "" {
		if err := json.Unmarshal(stdin, &hookData); err != nil {
			slog.Error("failed to parse hook data", "error", err)
			return 0
		}
	}

	handler := hooks.NewHandler(projectRoot)
	if err := handler.HandleHook(hookName, hookData); err != nil {
		slog.Error("hook handler error", "hook", hookName, "error", err)
		if hookName == "pre_tool_use" {
			if isBlockingError(err) {
				fmt.Fprintf(os.Stderr, "Hook blocked: %v\n", err)
				return 2
			}
		}
		return 0
	}

	slog.Debug("hook completed", "hook", hookName)
	return 0
}

func isBlockingError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "blocked") ||
		strings.Contains(err.Error(), "dangerous") ||
		strings.Contains(err.Error(), "prohibited")
}
