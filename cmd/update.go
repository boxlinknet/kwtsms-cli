// update.go - check for a newer kwtsms-cli release on GitHub.
// Endpoint: https://api.github.com/repos/boxlinknet/kwtsms-cli/releases/latest
// No credentials required.
// Related files: internal/update/update.go, cmd/root.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/boxlinknet/kwtsms-cli/internal/update"
)

var updateCmd = &cobra.Command{
	Use:   "version-check",
	Short: "Check for a newer version of kwtsms-cli",
	Long:  "Compares the installed version against the latest GitHub release and prints a notice if an update is available.",
	// Override PersistentPreRunE — no credentials needed for this command.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		msg := update.CheckNow(version)
		if msg == "" {
			fmt.Printf("kwtsms-cli %s is up to date.\n", version)
		} else {
			fmt.Println(msg)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
