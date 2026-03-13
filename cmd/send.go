// send.go - send an SMS message via the kwtSMS API.
// Endpoint: POST /API/send/
// Required flags: --to, --message
// Optional flags: --sender, --test
// Related files: internal/api/client.go, internal/sanitize/sanitize.go, internal/output/output.go
package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/boxlinknet/kwtsms-cli/internal/api"
	"github.com/boxlinknet/kwtsms-cli/internal/logger"
	"github.com/boxlinknet/kwtsms-cli/internal/output"
	"github.com/boxlinknet/kwtsms-cli/internal/sanitize"
)

const maxBatchSize = 200

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

	responses, err := sendBulk(Username, Password, sender, message, phones, sendTest)
	if err != nil {
		logger.LogError(LogFile, err)
		return err
	}
	logger.LogSend(LogFile, responses)
	if jsonFlag {
		if len(responses) == 1 {
			return output.PrintJSON(responses[0])
		}
		return output.PrintJSON(responses)
	}
	output.PrintSendResults(responses)
	return nil
}

// sendBulk splits phones into batches of maxBatchSize and sends each batch,
// sleeping 500ms between batches to stay within rate limits.
// Stops and returns the error on the first failed batch.
func sendBulk(username, password, sender, message string, phones []string, test bool) ([]*api.SendResponse, error) {
	var responses []*api.SendResponse
	for i := 0; i < len(phones); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(phones) {
			end = len(phones)
		}
		numbers := strings.Join(phones[i:end], ",")
		resp, err := api.SendSMS(username, password, sender, numbers, message, test)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
		if end < len(phones) {
			time.Sleep(500 * time.Millisecond)
		}
	}
	return responses, nil
}
