// senderid.go - list approved sender IDs on the account.
// Endpoint: POST /API/senderid/
// Related files: internal/api/client.go, internal/output/output.go
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/boxlinknet/kwtsms-cli/internal/api"
	"github.com/boxlinknet/kwtsms-cli/internal/output"
)

var senderidCmd = &cobra.Command{
	Use:   "senderid",
	Short: "List approved sender IDs",
	Long:  "Display all sender IDs approved on your kwtSMS account.",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := api.GetSenderID(Username, Password)
		if err != nil {
			output.PrintError(err)
			return err
		}
		if jsonFlag {
			return output.PrintJSON(resp)
		}
		output.PrintSenderID(resp)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(senderidCmd)
}
