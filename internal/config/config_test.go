package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestConfigPath_Override(t *testing.T) {
	got, err := ConfigPath("/custom/path/config.toml")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "/custom/path/config.toml" {
		t.Errorf("got %q, want /custom/path/config.toml", got)
	}
}

func TestConfigPath_Default(t *testing.T) {
	got, err := ConfigPath("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(got, filepath.Join("kwtsms-cli", "kwtsms-cli.toml")) {
		t.Errorf("got %q, want path ending in kwtsms-cli/kwtsms-cli.toml", got)
	}
}

func TestDefaultLogFile(t *testing.T) {
	got, err := DefaultLogFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "kwtsms-cli.log" {
		t.Errorf("got %q, want kwtsms-cli.log", got)
	}
}

func TestWrite_CreatesFileWithContent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kwtsms-cli.toml")

	err := Write(path, "myuser", "mypass", "MY-SENDER", "kwtsms-cli.log")
	if err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, `username = "myuser"`) {
		t.Errorf("missing username in config:\n%s", content)
	}
	if !strings.Contains(content, `password = "mypass"`) {
		t.Errorf("missing password in config:\n%s", content)
	}
	if !strings.Contains(content, `sender   = "MY-SENDER"`) {
		t.Errorf("missing sender in config:\n%s", content)
	}
	if !strings.Contains(content, `log_file = "kwtsms-cli.log"`) {
		t.Errorf("missing log_file in config:\n%s", content)
	}
}

func TestWrite_NoLogFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "kwtsms-cli.toml")

	err := Write(path, "myuser", "mypass", "MY-SENDER", "")
	if err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "log_file") {
		t.Errorf("log_file should not appear when logFilePath is empty:\n%s", string(data))
	}
}

func TestWrite_Permissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("permission bits not enforced on Windows")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "kwtsms-cli.toml")

	err := Write(path, "u", "p", "s", "")
	if err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("file permissions = %o, want 0600", info.Mode().Perm())
	}
}

func TestWrite_CreatesDirectory(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "subdir", "kwtsms-cli.toml")

	err := Write(path, "u", "p", "s", "")
	if err != nil {
		t.Fatalf("Write() should create parent dirs, got error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("config file not created: %v", err)
	}
}
