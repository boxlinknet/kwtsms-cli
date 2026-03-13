// Package update checks for a newer kwtsms-cli release on GitHub.
// The check runs in a goroutine so it never blocks command execution.
// All errors are silently discarded — a failed check is invisible to the user.
package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// semverRE matches valid semver version strings with optional v prefix: v1.2.3 or 1.2.3
var semverRE = regexp.MustCompile(`^v?\d+\.\d+\.\d+$`)

// releaseURL is a var so tests can override it with an httptest.Server URL.
var releaseURL = "https://api.github.com/repos/boxlinknet/kwtsms-cli/releases/latest"

// CheckAsync starts a background goroutine that fetches the latest GitHub
// release and sends a human-readable notice to the returned channel.
// The channel receives exactly one value: a non-empty notice string if a
// newer version is available, or an empty string otherwise.
// Call this before running the command; read the channel after.
func CheckAsync(current string) <-chan string {
	ch := make(chan string, 1)
	go func() {
		ch <- check(current)
	}()
	return ch
}

func check(current string) string {
	// Skip update check for dev builds and unversioned binaries.
	if current == "dev" || current == "" {
		return ""
	}

	// Validate current version looks like semver before making any request.
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

	// Cap response body to 4KB — we only need the tag_name field.
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

	// Validate the tag from the API is a proper semver string before using it.
	if !semverRE.MatchString(release.TagName) {
		return ""
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	cur := strings.TrimPrefix(current, "v")

	if latest == cur {
		return ""
	}

	return fmt.Sprintf("\nA new version is available: %s → %s",
		release.TagName,
		"https://github.com/boxlinknet/kwtsms-cli/releases/latest",
	)
}
