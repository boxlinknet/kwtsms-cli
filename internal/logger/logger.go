// Package logger provides append-only JSON-line logging for send operations.
// Each send writes one line to the log file: timestamp, numbers, credits charged,
// balance after, and all message IDs. Errors are also logged.
// Logging is opt-in: all functions silently do nothing if path is empty.
// Credentials and message text are never written to the log.
// Related files: cmd/send.go, cmd/root.go, internal/config/config.go
package logger

import (
	"encoding/json"
	"os"
	"time"

	"github.com/boxlinknet/kwtsms-cli/internal/api"
)

type entry struct {
	Time    string   `json:"time"`
	Numbers int64    `json:"numbers,omitempty"`
	Charged int64    `json:"charged,omitempty"`
	Balance int64    `json:"balance,omitempty"`
	MsgIDs  []string `json:"msg_ids,omitempty"`
	Error   string   `json:"error,omitempty"`
}

// LogSend appends a JSON line recording a successful send to the log file.
// Aggregates totals across all batches. Silently does nothing if path is empty.
func LogSend(path string, responses []*api.SendResponse) {
	if path == "" {
		return
	}
	e := entry{Time: time.Now().UTC().Format(time.RFC3339)}
	for _, r := range responses {
		e.Numbers += r.Numbers
		e.Charged += r.PointsCharged
		e.MsgIDs = append(e.MsgIDs, r.MsgID)
	}
	if len(responses) > 0 {
		e.Balance = responses[len(responses)-1].BalanceAfter
	}
	appendEntry(path, e)
}

// LogError appends a JSON line recording a send failure to the log file.
// Silently does nothing if path is empty.
func LogError(path string, err error) {
	if path == "" {
		return
	}
	appendEntry(path, entry{
		Time:  time.Now().UTC().Format(time.RFC3339),
		Error: err.Error(),
	})
}

func appendEntry(path string, e entry) {
	data, err := json.Marshal(e)
	if err != nil {
		return
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	_, _ = f.Write(append(data, '\n'))
}
