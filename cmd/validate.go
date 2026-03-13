// validate.go - validate phone numbers via the kwtSMS API.
// Endpoint: POST /API/validate/
// Arguments: one or more phone numbers (space, comma, or mixed-separated)
// Related files: internal/api/client.go, internal/sanitize/sanitize.go, internal/output/output.go
package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/boxlinknet/kwtsms-cli/internal/api"
	"github.com/boxlinknet/kwtsms-cli/internal/output"
	"github.com/boxlinknet/kwtsms-cli/internal/sanitize"
)

var validateCmd = &cobra.Command{
	Use:   "validate <NUMBERS>...",
	Short: "Validate phone numbers",
	Long: `Validate one or more phone numbers via kwtSMS.

Numbers can be space-separated, comma-separated, or mixed:
  kwtsms-cli validate 96598765432 96512345678
  kwtsms-cli validate 96598765432,96512345678
  kwtsms-cli validate 96598765432,96512345678 96599999`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runValidate(args)
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(args []string) error {
	// Join all positional args then re-split on commas, spaces, newlines
	joined := strings.Join(args, ",")
	phones, err := sanitize.SanitizePhones(joined)
	if err != nil {
		return fmt.Errorf("invalid phone number(s): %w", err)
	}

	numbers := strings.Join(phones, ",")
	resp, err := api.ValidateNumbers(Username, Password, numbers)
	if err != nil {
		output.PrintError(err)
		return err
	}
	if jsonFlag {
		return output.PrintJSON(resp)
	}
	output.PrintValidate(resp)
	return nil
}
