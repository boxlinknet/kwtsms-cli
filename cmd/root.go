// Package cmd implements all kwtsms-cli subcommands using cobra.
// root.go defines the root command, global flags, and viper config binding.
// All subcommands access credentials via the package-level Username and Password vars,
// which are populated by PersistentPreRunE before any subcommand runs.
// Related files: main.go, internal/config/config.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/boxlinknet/kwtsms-cli/internal/config"
	"github.com/boxlinknet/kwtsms-cli/internal/output"
	"github.com/boxlinknet/kwtsms-cli/internal/sanitize"
)

// Credentials resolved by PersistentPreRunE and used by all subcommands.
var (
	Username      string
	Password      string
	DefaultSender string
)

// Global flag values
var (
	jsonFlag     bool
	configFlag   string
	usernameFlag string
	passwordFlag string
)

var rootCmd = &cobra.Command{
	Use:           "kwtsms-cli",
	Short:         "kwtSMS command-line interface",
	Long:          "Send SMS, check balance, and manage your kwtSMS account from the terminal.",
	Version:       "0.1.0",
	SilenceUsage:  true, // don't print usage on API or runtime errors
	SilenceErrors: true, // let Execute() handle error printing once
	// PersistentPreRunE runs before every subcommand except setup (which overrides it).
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return loadCredentials()
	},
}

// loadCredentials resolves credentials from config file, env vars, and flags.
// Priority (highest wins): inline flags > env vars > config file.
func loadCredentials() error {
	// Configure viper: read config file
	cfgPath, err := config.ConfigPath(configFlag)
	if err != nil {
		return err
	}
	viper.SetConfigFile(cfgPath)
	viper.SetConfigType("toml")

	// Silent read - config file may not exist yet (user hasn't run setup)
	_ = viper.ReadInConfig()

	// Bind env vars (KWTSMS_USERNAME, KWTSMS_PASSWORD, KWTSMS_SENDER)
	viper.SetEnvPrefix("KWTSMS")
	viper.AutomaticEnv()

	// Resolve username: flag > env > config
	u := viper.GetString("username")
	if usernameFlag != "" {
		u = usernameFlag
	}
	p := viper.GetString("password")
	if passwordFlag != "" {
		p = passwordFlag
	}
	s := viper.GetString("sender")

	// Sanitize resolved values
	if u != "" {
		u, err = sanitize.SanitizeConfigValue(u)
		if err != nil {
			return fmt.Errorf("invalid username: %w", err)
		}
	}
	if p != "" {
		p, err = sanitize.SanitizeConfigValue(p)
		if err != nil {
			return fmt.Errorf("invalid password: %w", err)
		}
	}
	if s != "" {
		s, err = sanitize.SanitizeConfigValue(s)
		if err != nil {
			s = "" // non-fatal: sender may be provided per-command
		}
	}

	if u == "" || p == "" {
		return fmt.Errorf("no credentials found. Run 'kwtsms-cli setup' to configure")
	}

	Username = u
	Password = p
	DefaultSender = s
	return nil
}

// Execute runs the root command. Called from main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.PrintError(err)
		os.Exit(1)
	}
}

func init() {
	// Persistent flags available to all subcommands
	rootCmd.PersistentFlags().BoolVar(&jsonFlag, "json", false, "output as JSON")
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "", "override config file path")
	rootCmd.PersistentFlags().StringVar(&usernameFlag, "username", "", "API username (overrides config and env)")
	rootCmd.PersistentFlags().StringVar(&passwordFlag, "password", "", "API password (overrides config and env)")
}
