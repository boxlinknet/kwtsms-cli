// coverage.go - list active country coverage and prefixes.
// Endpoint: POST /API/coverage/
// Related files: internal/api/client.go, internal/output/output.go
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/boxlinknet/kwtsms-cli/internal/api"
	"github.com/boxlinknet/kwtsms-cli/internal/output"
)

var coverageCmd = &cobra.Command{
	Use:   "coverage",
	Short: "List active country coverage",
	Long:  "Display all active country prefixes available on your kwtSMS account.",
	RunE: func(cmd *cobra.Command, args []string) error {
		raw, err := api.GetCoverage(Username, Password)
		if err != nil {
			output.PrintError(err)
			return err
		}
		if jsonFlag {
			return output.PrintJSON(raw)
		}
		output.PrintCoverage(raw)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(coverageCmd)
}
