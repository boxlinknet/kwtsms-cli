// update.go - check for a newer kwtsms-cli release on GitHub.
// Endpoint: https://api.github.com/repos/boxlinknet/kwtsms-cli/releases/latest
// No credentials required.
// Related files: internal/update/update.go, cmd/root.go
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/boxlinknet/kwtsms-cli/internal/output"
	"github.com/boxlinknet/kwtsms-cli/internal/update"
)

var updateCmd = &cobra.Command{
	Use:   "version-check",
	Short: "Check for a newer version of kwtsms-cli",
	Long:  "Compares the installed version against the latest GitHub release and prints a notice if an update is available.",
	// Override PersistentPreRunE — no credentials needed for this command.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := update.CheckNow(version)
		if err != nil {
			return err
		}
		if jsonFlag {
			return output.PrintJSON(result)
		}
		if result.UpToDate {
			fmt.Printf("kwtsms-cli %s is up to date.\n", result.Current)
		} else {
			fmt.Printf("A new version is available: %s\nDownload: %s\n", result.Latest, result.DownloadURL)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
