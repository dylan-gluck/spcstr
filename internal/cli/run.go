package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the spcstr TUI",
	Long: `Launch the Spec⭐️ terminal user interface to view and manage
Claude Code sessions, browse documentation, and monitor agent activities.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Placeholder implementation
		fmt.Fprintln(cmd.OutOrStdout(), "Launching spcstr TUI...")
		fmt.Fprintln(cmd.OutOrStdout(), "This command will be implemented in a future story.")
		return nil
	},
}

func init() {
	// Command-specific flags
	runCmd.Flags().StringP("session", "s", "", "load specific session by ID")
	runCmd.Flags().BoolP("plan", "p", false, "start in Plan view")
	runCmd.Flags().BoolP("observe", "o", false, "start in Observe view")
}
