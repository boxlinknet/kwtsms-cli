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

// CheckNow fetches the latest GitHub release synchronously and returns a
// human-readable notice if a newer version is available, or an empty string.
func CheckNow(current string) string {
	return check(current)
}

// platformBinary returns the binary filename for the current OS and architecture.
// Returns an empty string for unsupported platforms.
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

func check(current string) string {
	if current == "dev" || current == "" {
		return ""
	}
	if !semverRE.MatchString(current) {
		return ""
	}

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(releaseURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if err != nil {
		return ""
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(body, &release); err != nil {
		return ""
	}

	if !semverRE.MatchString(release.TagName) {
		return ""
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	cur := strings.TrimPrefix(current, "v")
	if latest == cur {
		return ""
	}

	binary := platformBinary()
	if binary == "" {
		return fmt.Sprintf("A new version is available: %s\nDownload: %s",
			release.TagName,
			"https://github.com/boxlinknet/kwtsms-cli/releases/latest",
		)
	}
	return fmt.Sprintf("A new version is available: %s\nDownload: %s%s",
		release.TagName,
		downloadBase,
		binary,
	)
}
