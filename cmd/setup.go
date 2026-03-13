// setup.go - interactive configuration wizard.
// The only command that reads from stdin interactively.
// Verifies credentials against the API before writing the config file.
// Writes the config file to the platform-appropriate path.
// Skips PersistentPreRunE credential check (no credentials needed to run setup).
// Related files: internal/config/config.go, internal/api/client.go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/boxlinknet/kwtsms-cli/internal/api"
	"github.com/boxlinknet/kwtsms-cli/internal/config"
	"github.com/boxlinknet/kwtsms-cli/internal/sanitize"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure API credentials interactively",
	Long:  "Interactive wizard to set your kwtSMS API username, password, and default sender ID.",
	// Override PersistentPreRunE so setup does not require existing credentials.
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSetup()
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("kwtSMS CLI Setup")
	fmt.Println("----------------")

	// Read username
	fmt.Print("API Username: ")
	rawUsername, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read username: %w", err)
	}
	username, err := sanitize.SanitizeConfigValue(strings.TrimSpace(rawUsername))
	if err != nil {
		return fmt.Errorf("invalid username: %w", err)
	}

	// Read password (masked)
	fmt.Print("API Password: ")
	rawPassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // newline after masked input
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	password, err := sanitize.SanitizeConfigValue(strings.TrimSpace(string(rawPassword)))
	if err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	// Verify credentials and fetch sender IDs
	fmt.Println("Verifying credentials...")
	senderResp, err := api.GetSenderID(username, password)
	if err != nil {
		return fmt.Errorf("credential verification failed: %w", err)
	}

	if len(senderResp.SenderID) == 0 {
		return fmt.Errorf("no sender IDs found on this account. Please contact kwtSMS support")
	}

	// Display sender ID list
	fmt.Println("\nAvailable sender IDs:")
	for i, id := range senderResp.SenderID {
		fmt.Printf("  [%d] %s\n", i+1, id)
	}

	// Prompt for selection
	fmt.Printf("\nSelect default sender ID [1]: ")
	rawSelection, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read selection: %w", err)
	}
	rawSelection = strings.TrimSpace(rawSelection)

	// Default to first if empty
	idx := 1
	if rawSelection != "" {
		idx, err = strconv.Atoi(rawSelection)
		if err != nil || idx < 1 || idx > len(senderResp.SenderID) {
			return fmt.Errorf("invalid selection %q, must be between 1 and %d", rawSelection, len(senderResp.SenderID))
		}
	}
	chosenSender := senderResp.SenderID[idx-1]

	sender, err := sanitize.SanitizeSenderID(chosenSender)
	if err != nil {
		return fmt.Errorf("invalid sender ID: %w", err)
	}

	// Write config file
	cfgPath, err := config.ConfigPath(configFlag)
	if err != nil {
		return err
	}
	if err := config.Write(cfgPath, username, password, sender); err != nil {
		return err
	}

	fmt.Printf("\nConfig saved to: %s\n", cfgPath)
	return nil
}
