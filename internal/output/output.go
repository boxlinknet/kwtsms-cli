// Package output handles all CLI output formatting.
// Human-readable output goes to stdout. Errors go to stderr.
// Use --json flag for machine-readable JSON output.
// Related files: cmd/balance.go, cmd/send.go, cmd/senderid.go, cmd/coverage.go
package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/boxlinknet/kwtsms-cli/internal/api"
)

// PrintJSON pretty-prints any value as indented JSON to stdout.
func PrintJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// FormatNumber formats an integer with thousands separators (comma-separated).
// Example: 1234567 -> "1,234,567"
func FormatNumber(n int64) string {
	s := fmt.Sprintf("%d", n)
	if n < 0 {
		s = s[1:]
	}
	// Insert commas every 3 digits from the right
	var result []byte
	for i, ch := range []byte(s) {
		if i > 0 && (len(s)-i)%3 == 0 {
			result = append(result, ',')
		}
		result = append(result, ch)
	}
	if n < 0 {
		return "-" + string(result)
	}
	return string(result)
}

// PrintBalance prints available and purchased credits in human-readable format.
func PrintBalance(resp *api.BalanceResponse) {
	fmt.Printf("Available:  %s\n", FormatNumber(resp.Available))
	fmt.Printf("Purchased:  %s\n", FormatNumber(resp.Purchased))
}

// PrintSenderID prints one sender ID per line.
func PrintSenderID(resp *api.SenderIDResponse) {
	for _, id := range resp.SenderID {
		fmt.Println(id)
	}
}

// PrintCoverage prints coverage data from raw JSON.
// Iterates over any map or array in the response and prints entries.
func PrintCoverage(raw json.RawMessage) {
	// Try to parse as a generic map to find coverage data
	var envelope map[string]json.RawMessage
	if err := json.Unmarshal(raw, &envelope); err != nil {
		fmt.Println(string(raw))
		return
	}

	// Look for coverage field (dynamic key name from API)
	for key, val := range envelope {
		if key == "result" || key == "code" || key == "description" {
			continue
		}
		// Try as array of objects
		var entries []map[string]interface{}
		if err := json.Unmarshal(val, &entries); err == nil {
			for _, entry := range entries {
				var parts []string
				for k, v := range entry {
					parts = append(parts, fmt.Sprintf("%s: %v", k, v))
				}
				fmt.Println(strings.Join(parts, "  "))
			}
			return
		}
		// Try as array of strings
		var strs []string
		if err := json.Unmarshal(val, &strs); err == nil {
			for _, s := range strs {
				fmt.Println(s)
			}
			return
		}
		// Fallback: print raw value
		fmt.Println(string(val))
	}
}

// PrintSend prints the send result in the standard 5-line format.
// Lines: Sent / Numbers / Charged / Balance / MsgID
func PrintSend(resp *api.SendResponse) {
	fmt.Println("Sent")
	fmt.Printf("Numbers:    %s\n", FormatNumber(resp.Numbers))
	fmt.Printf("Charged:    %s\n", FormatNumber(resp.PointsCharged))
	fmt.Printf("Balance:    %s\n", FormatNumber(resp.BalanceAfter))
	fmt.Printf("MsgID:      %s\n", resp.MsgID)
}

// PrintSendResults prints aggregated results for one or more batches.
// For a single batch it delegates to PrintSend. For multiple batches it sums
// Numbers and Charged, shows the final balance, and lists each MsgID.
func PrintSendResults(responses []*api.SendResponse) {
	if len(responses) == 1 {
		PrintSend(responses[0])
		return
	}
	var totalNumbers, totalCharged int64
	for _, r := range responses {
		totalNumbers += r.Numbers
		totalCharged += r.PointsCharged
	}
	fmt.Println("Sent")
	fmt.Printf("Numbers:    %s\n", FormatNumber(totalNumbers))
	fmt.Printf("Charged:    %s\n", FormatNumber(totalCharged))
	fmt.Printf("Balance:    %s\n", FormatNumber(responses[len(responses)-1].BalanceAfter))
	for i, r := range responses {
		if i == 0 {
			fmt.Printf("MsgID:      %s\n", r.MsgID)
		} else {
			fmt.Printf("            %s\n", r.MsgID)
		}
	}
}

// PrintError prints an error message to stderr.
func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
}
