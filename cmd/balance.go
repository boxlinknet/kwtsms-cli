// balance.go - show account balance.
// Endpoint: POST /API/balance/
// Related files: internal/api/client.go, internal/output/output.go
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/boxlinknet/kwtsms-cli/internal/api"
	"github.com/boxlinknet/kwtsms-cli/internal/output"
)

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Show account balance",
	Long:  "Display available and purchased SMS credits on your kwtSMS account.",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.GetBalance(Username, Password)
		if err != nil {
			return err
		}
		if jsonFlag {
			return output.PrintJSON(resp)
		}
		output.PrintBalance(resp)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(balanceCmd)
}
