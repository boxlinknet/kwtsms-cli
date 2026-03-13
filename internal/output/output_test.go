package output

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/boxlinknet/kwtsms-cli/internal/api"
)

// captureStdout redirects os.Stdout during f() and returns what was printed.
func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	r, w, _ := os.Pipe()
	orig := os.Stdout
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = orig
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		n    int64
		want string
	}{
		{0, "0"},
		{1, "1"},
		{999, "999"},
		{1000, "1,000"},
		{1234, "1,234"},
		{12345, "12,345"},
		{123456, "123,456"},
		{1234567, "1,234,567"},
		{-1234, "-1,234"},
		{1000000, "1,000,000"},
	}
	for _, tt := range tests {
		got := FormatNumber(tt.n)
		if got != tt.want {
			t.Errorf("FormatNumber(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

func TestPrintBalance(t *testing.T) {
	out := captureStdout(t, func() {
		PrintBalance(&api.BalanceResponse{Available: 1234, Purchased: 5000})
	})
	if !strings.Contains(out, "1,234") {
		t.Errorf("output missing available balance: %q", out)
	}
	if !strings.Contains(out, "5,000") {
		t.Errorf("output missing purchased balance: %q", out)
	}
}

func TestPrintSenderID(t *testing.T) {
	out := captureStdout(t, func() {
		PrintSenderID(&api.SenderIDResponse{SenderID: []string{"MY-BRAND", "KWT-SMS"}})
	})
	if !strings.Contains(out, "MY-BRAND") {
		t.Errorf("output missing MY-BRAND: %q", out)
	}
	if !strings.Contains(out, "KWT-SMS") {
		t.Errorf("output missing KWT-SMS: %q", out)
	}
}

func TestPrintSend_SingleBatch(t *testing.T) {
	out := captureStdout(t, func() {
		PrintSendResults([]*api.SendResponse{
			{MsgID: "abc123", Numbers: 1, PointsCharged: 1, BalanceAfter: 999},
		})
	})
	if !strings.Contains(out, "Sent") {
		t.Errorf("missing 'Sent': %q", out)
	}
	if !strings.Contains(out, "abc123") {
		t.Errorf("missing MsgID: %q", out)
	}
	if !strings.Contains(out, "999") {
		t.Errorf("missing balance: %q", out)
	}
}

func TestPrintSendResults_MultiBatch(t *testing.T) {
	out := captureStdout(t, func() {
		PrintSendResults([]*api.SendResponse{
			{MsgID: "aaa111", Numbers: 200, PointsCharged: 200, BalanceAfter: 800},
			{MsgID: "bbb222", Numbers: 50, PointsCharged: 50, BalanceAfter: 750},
		})
	})
	if !strings.Contains(out, "250") {
		t.Errorf("missing aggregated numbers (250): %q", out)
	}
	if !strings.Contains(out, "750") {
		t.Errorf("missing final balance (750): %q", out)
	}
	if !strings.Contains(out, "aaa111") {
		t.Errorf("missing first MsgID: %q", out)
	}
	if !strings.Contains(out, "bbb222") {
		t.Errorf("missing second MsgID: %q", out)
	}
}

func TestPrintJSON(t *testing.T) {
	out := captureStdout(t, func() {
		PrintJSON(map[string]interface{}{"result": "OK", "available": 100})
	})
	if !strings.Contains(out, `"result"`) {
		t.Errorf("missing result key: %q", out)
	}
	if !strings.Contains(out, `"available"`) {
		t.Errorf("missing available key: %q", out)
	}
}

func TestPrintError(t *testing.T) {
	r, w, _ := os.Pipe()
	orig := os.Stderr
	os.Stderr = w
	PrintError(fmt.Errorf("something went wrong"))
	w.Close()
	os.Stderr = orig
	var buf bytes.Buffer
	io.Copy(&buf, r)
	if !strings.Contains(buf.String(), "something went wrong") {
		t.Errorf("stderr missing error message: %q", buf.String())
	}
}
