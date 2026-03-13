package logger

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/boxlinknet/kwtsms-cli/internal/api"
)

func TestLogSend_EmptyPath(t *testing.T) {
	// Must not panic or create any file
	LogSend("", []*api.SendResponse{{MsgID: "abc", Numbers: 1, PointsCharged: 1, BalanceAfter: 100}})
}

func TestLogError_EmptyPath(t *testing.T) {
	// Must not panic or create any file
	LogError("", errors.New("some error"))
}

func TestLogSend_WritesJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")

	LogSend(path, []*api.SendResponse{
		{MsgID: "aaa", Numbers: 5, PointsCharged: 5, BalanceAfter: 995},
	})

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	var e entry
	if err := json.Unmarshal(data[:len(data)-1], &e); err != nil {
		t.Fatalf("json.Unmarshal error: %v, raw: %s", err, data)
	}
	if e.Numbers != 5 {
		t.Errorf("Numbers = %d, want 5", e.Numbers)
	}
	if e.Charged != 5 {
		t.Errorf("Charged = %d, want 5", e.Charged)
	}
	if e.Balance != 995 {
		t.Errorf("Balance = %d, want 995", e.Balance)
	}
	if len(e.MsgIDs) != 1 || e.MsgIDs[0] != "aaa" {
		t.Errorf("MsgIDs = %v, want [aaa]", e.MsgIDs)
	}
	if e.Time == "" {
		t.Error("Time should not be empty")
	}
}

func TestLogSend_MultiBatch(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")

	LogSend(path, []*api.SendResponse{
		{MsgID: "aaa", Numbers: 200, PointsCharged: 200, BalanceAfter: 800},
		{MsgID: "bbb", Numbers: 50, PointsCharged: 50, BalanceAfter: 750},
	})

	data, _ := os.ReadFile(path)
	var e entry
	json.Unmarshal(data[:len(data)-1], &e)

	if e.Numbers != 250 {
		t.Errorf("Numbers = %d, want 250", e.Numbers)
	}
	if e.Balance != 750 {
		t.Errorf("Balance = %d, want 750 (last batch)", e.Balance)
	}
	if len(e.MsgIDs) != 2 {
		t.Errorf("MsgIDs len = %d, want 2", len(e.MsgIDs))
	}
}

func TestLogError_WritesJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")

	LogError(path, errors.New("[ERR010] Zero balance."))

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}

	var e entry
	json.Unmarshal(data[:len(data)-1], &e)
	if e.Error != "[ERR010] Zero balance." {
		t.Errorf("Error = %q, want [ERR010] Zero balance.", e.Error)
	}
	if e.Time == "" {
		t.Error("Time should not be empty")
	}
}

func TestLogSend_AppendMultipleLines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.log")

	LogSend(path, []*api.SendResponse{{MsgID: "aaa", Numbers: 1, PointsCharged: 1, BalanceAfter: 99}})
	LogSend(path, []*api.SendResponse{{MsgID: "bbb", Numbers: 2, PointsCharged: 2, BalanceAfter: 97}})

	f, _ := os.Open(path)
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) != 2 {
		t.Fatalf("got %d lines, want 2", len(lines))
	}

	var e1, e2 entry
	json.Unmarshal([]byte(lines[0]), &e1)
	json.Unmarshal([]byte(lines[1]), &e2)

	if e1.MsgIDs[0] != "aaa" {
		t.Errorf("line 1 MsgID = %q, want aaa", e1.MsgIDs[0])
	}
	if e2.MsgIDs[0] != "bbb" {
		t.Errorf("line 2 MsgID = %q, want bbb", e2.MsgIDs[0])
	}
}
