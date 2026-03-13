// Package update checks for a newer kwtsms-cli release on GitHub.
// All errors are silently discarded — a failed check is invisible to the user.
package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// semverRE matches valid semver version strings with optional v prefix: v1.2.3 or 1.2.3
var semverRE = regexp.MustCompile(`^v?\d+\.\d+\.\d+$`)

// releaseURL is a var so tests can override it with an httptest.Server URL.
var releaseURL = "https://api.github.com/repos/boxlinknet/kwtsms-cli/releases/latest"

const downloadBase = "https://github.com/boxlinknet/kwtsms-cli/releases/latest/download/"

// Result holds the outcome of a version check.
type Result struct {
	Current     string `json:"current"`
	Latest      string `json:"latest"`
	UpToDate    bool   `json:"up_to_date"`
	DownloadURL string `json:"download_url"`
}

// CheckNow fetches the latest GitHub release and returns a Result.
// On any error the returned Result has Latest empty and UpToDate false.
func CheckNow(current string) (*Result, error) {
	return check(current)
}

// platformBinary returns the binary filename for the current OS and architecture.
func platformBinary() string {
	switch runtime.GOOS {
	case "linux":
		switch runtime.GOARCH {
		case "amd64":
			return "kwtsms-cli-linux-x64"
		case "arm64":
			return "kwtsms-cli-linux-arm64"
		case "arm":
			return "kwtsms-cli-linux-armv7"
		}
	case "darwin":
		switch runtime.GOARCH {
		case "amd64":
			return "kwtsms-cli-macos-x64"
		case "arm64":
			return "kwtsms-cli-macos-arm64"
		}
	case "windows":
		if runtime.GOARCH == "amd64" {
			return "kwtsms-cli-windows-x64.exe"
		}
	}
	return ""
}

func check(current string) (*Result, error) {
	if current == "dev" || current == "" || !semverRE.MatchString(current) {
		return &Result{Current: current}, fmt.Errorf("version-check not available for this build")
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(releaseURL)
	if err != nil {
		return nil, fmt.Errorf("cannot reach GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response from GitHub: HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !semverRE.MatchString(release.TagName) {
		return nil, fmt.Errorf("unexpected tag format from GitHub: %q", release.TagName)
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	cur := strings.TrimPrefix(current, "v")
	upToDate := latest == cur

	downloadURL := ""
	if !upToDate {
		if binary := platformBinary(); binary != "" {
			downloadURL = downloadBase + binary
		} else {
			downloadURL = "https://github.com/boxlinknet/kwtsms-cli/releases/latest"
		}
	}

	return &Result{
		Current:     current,
		Latest:      release.TagName,
		UpToDate:    upToDate,
		DownloadURL: downloadURL,
	}, nil
}
