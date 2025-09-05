package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Display the current version of Spec⭐️ along with build information.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Fprintf(cmd.OutOrStdout(), "Spec⭐️ version %s\n", Version())
		fmt.Fprintln(cmd.OutOrStdout(), "Build: development")
		fmt.Fprintln(cmd.OutOrStdout(), "Go version: 1.21.0")
		return nil
	},
}
