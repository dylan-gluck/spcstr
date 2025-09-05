package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "1.0.0"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:     "spcstr",
	Short:   "spcstr - a CLI/TUI tool for Claude Code session observability",
	Long:    `spcstr provides real-time observability for Claude Code sessions through hook integration and an interactive terminal interface.`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("spcstr v%s\n", version)
		fmt.Println("Use 'spcstr --help' for more information.")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of spcstr",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("spcstr v%s\n", version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
