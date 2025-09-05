package main

import (
	"log/slog"
	"os"

	"github.com/dylan-gluck/spcstr/internal/cli"
)

func main() {
	// Check if this is a hook invocation
	if len(os.Args) >= 3 && os.Args[1] == "hook" {
		exitCode := runHook(os.Args[2])
		os.Exit(exitCode)
	}

	// Set up structured logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Setup panic recovery (only in main)
	defer func() {
		if r := recover(); r != nil {
			slog.Error("fatal error occurred", "error", r)
			os.Exit(1)
		}
	}()

	// Execute root command
	if err := cli.Execute(); err != nil {
		slog.Error("command execution failed", "error", err)
		os.Exit(1)
	}
}
