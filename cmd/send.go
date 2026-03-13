// send.go - send an SMS message via the kwtSMS API.
// Endpoint: POST /API/send/
// Required flags: --to, --message
// Optional flags: --sender, --test
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

var (
	sendTo      string
	sendMessage string
	sendSender  string
	sendTest    bool
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Send an SMS message",
	Long:  "Send an SMS to one or more phone numbers via kwtSMS.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSend()
	},
}

func init() {
	sendCmd.Flags().StringVarP(&sendTo, "to", "t", "", "recipient phone number(s), comma-separated (required)")
	sendCmd.Flags().StringVarP(&sendMessage, "message", "m", "", "message text (required)")
	sendCmd.Flags().StringVarP(&sendSender, "sender", "s", "", "sender ID (overrides config default)")
	sendCmd.Flags().BoolVar(&sendTest, "test", false, "test mode: queue message without delivery")
	_ = sendCmd.MarkFlagRequired("to")
	_ = sendCmd.MarkFlagRequired("message")
	rootCmd.AddCommand(sendCmd)
}

func runSend() error {
	// Sanitize phone numbers
	phones, err := sanitize.SanitizePhones(sendTo)
	if err != nil {
		return fmt.Errorf("invalid phone number(s): %w", err)
	}

	// Sanitize message
	message, err := sanitize.SanitizeMessage(sendMessage)
	if err != nil {
		return fmt.Errorf("invalid message: %w", err)
	}

	// Resolve sender: --sender flag > config default > error
	sender := sendSender
	if sender == "" {
		sender = DefaultSender
	}
	if sender == "" {
		return fmt.Errorf("no sender ID provided. Use --sender or set a default via 'kwtsms-cli setup'")
	}
	sender, err = sanitize.SanitizeSenderID(sender)
	if err != nil {
		return fmt.Errorf("invalid sender ID: %w", err)
	}

	numbers := strings.Join(phones, ",")
	resp, err := api.SendSMS(Username, Password, sender, numbers, message, sendTest)
	if err != nil {
		return err
	}
	if jsonFlag {
		return output.PrintJSON(resp)
	}
	output.PrintSend(resp)
	return nil
}
