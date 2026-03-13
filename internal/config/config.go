// Package config handles the kwtsms-cli configuration file.
// Config file location: ~/.config/kwtsms-cli/config.toml (Linux/macOS)
//
//	%APPDATA%\kwtsms-cli\config.toml (Windows)
//
// Credential priority (highest wins): CLI flags > env vars > config file.
// Viper handles the priority merging in cmd/root.go.
// Related files: cmd/root.go, cmd/setup.go
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

const (
	appName    = "kwtsms-cli"
	configFile = "config.toml"
)

// ConfigDir returns the platform-appropriate config directory for kwtsms-cli.
// Uses os.UserConfigDir() which returns:
//
//	Linux/macOS: $HOME/.config
//	Windows:     %APPDATA%
func ConfigDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("cannot determine config directory: %w", err)
	}
	return filepath.Join(base, appName), nil
}

// ConfigPath returns the full path to the config file.
// If override is non-empty, it is returned as-is.
func ConfigPath(override string) (string, error) {
	if override != "" {
		return override, nil
	}
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, configFile), nil
}

// Write creates the config directory if needed and writes the TOML config file.
// On Unix systems it sets file permissions to 0600 (user read/write only)
// to protect credentials from other users.
func Write(path, username, password, sender string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create config directory %s: %w", dir, err)
	}

	content := fmt.Sprintf("username = %q\npassword = %q\nsender   = %q\n",
		username, password, sender)

	// Write with restricted permissions
	perm := os.FileMode(0600)
	if err := os.WriteFile(path, []byte(content), perm); err != nil {
		return fmt.Errorf("failed to write config file %s: %w", path, err)
	}

	// Enforce 0600 on Unix (WriteFile respects umask so we chmod explicitly)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(path, 0600); err != nil {
			return fmt.Errorf("failed to set config file permissions: %w", err)
		}
	}

	return nil
}
