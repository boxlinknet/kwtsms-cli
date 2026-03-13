// Package cmd implements all kwtsms-cli subcommands using cobra.
// root.go defines the root command, global flags, and viper config binding.
// All subcommands access credentials via rootCmd.PersistentPreRunE.
// Related files: main.go, internal/config/config.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kwtsms-cli",
	Short: "kwtSMS command-line interface",
	Long:  "Send SMS, check balance, and manage your kwtSMS account from the terminal.",
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
