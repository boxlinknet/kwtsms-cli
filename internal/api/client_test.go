package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// serve returns an httptest.Server that responds with the given status code and JSON body.
// It also sets baseURL to point at the test server and restores it on cleanup.
func serve(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	data, _ := json.Marshal(body)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		w.Write(data)
	}))
	orig := baseURL
	baseURL = srv.URL + "/"
	t.Cleanup(func() {
		srv.Close()
		baseURL = orig
	})
	return srv
}

func TestGetBalance_OK(t *testing.T) {
	serve(t, 200, map[string]interface{}{
		"result":    "OK",
		"available": 1234,
		"purchased": 5000,
	})
	resp, err := GetBalance("user", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Available != 1234 {
		t.Errorf("Available = %d, want 1234", resp.Available)
	}
	if resp.Purchased != 5000 {
		t.Errorf("Purchased = %d, want 5000", resp.Purchased)
	}
}

func TestGetBalance_APIError(t *testing.T) {
	serve(t, 200, map[string]interface{}{
		"result":      "ERROR",
		"code":        "ERR003",
		"description": "Authentication error.",
	})
	_, err := GetBalance("bad", "creds")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "[ERR003] Authentication error." {
		t.Errorf("error = %q, want [ERR003] Authentication error.", err.Error())
	}
}

func TestGetBalance_HTTPError(t *testing.T) {
	serve(t, 500, map[string]interface{}{})
	_, err := GetBalance("user", "pass")
	if err == nil {
		t.Fatal("expected error on HTTP 500, got nil")
	}
}

func TestGetSenderID_OK(t *testing.T) {
	serve(t, 200, map[string]interface{}{
		"result":   "OK",
		"senderid": []string{"MY-BRAND", "KWT-SMS"},
	})
	resp, err := GetSenderID("user", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.SenderID) != 2 {
		t.Errorf("got %d sender IDs, want 2", len(resp.SenderID))
	}
	if resp.SenderID[0] != "MY-BRAND" {
		t.Errorf("SenderID[0] = %q, want MY-BRAND", resp.SenderID[0])
	}
}

func TestGetSenderID_APIError(t *testing.T) {
	serve(t, 200, map[string]interface{}{
		"result":      "ERROR",
		"code":        "ERR003",
		"description": "Authentication error.",
	})
	_, err := GetSenderID("bad", "creds")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSendSMS_OK(t *testing.T) {
	serve(t, 200, map[string]interface{}{
		"result":          "OK",
		"msg-id":          "abc123",
		"numbers":         1,
		"points-charged":  1,
		"balance-after":   999,
	})
	resp, err := SendSMS("user", "pass", "MY-BRAND", "96598765432", "Hello", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.MsgID != "abc123" {
		t.Errorf("MsgID = %q, want abc123", resp.MsgID)
	}
	if resp.Numbers != 1 {
		t.Errorf("Numbers = %d, want 1", resp.Numbers)
	}
	if resp.BalanceAfter != 999 {
		t.Errorf("BalanceAfter = %d, want 999", resp.BalanceAfter)
	}
}

func TestSendSMS_TestMode(t *testing.T) {
	var gotBody map[string]interface{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": "OK", "msg-id": "x", "numbers": 1,
			"points-charged": 1, "balance-after": 100,
		})
	}))
	orig := baseURL
	baseURL = srv.URL + "/"
	t.Cleanup(func() { srv.Close(); baseURL = orig })

	SendSMS("user", "pass", "MY-BRAND", "96598765432", "Hello", true)
	if gotBody["test"] != "1" {
		t.Errorf("test field = %v, want \"1\"", gotBody["test"])
	}
}

func TestSendSMS_APIError(t *testing.T) {
	serve(t, 200, map[string]interface{}{
		"result":      "ERROR",
		"code":        "ERR010",
		"description": "Zero balance.",
	})
	_, err := SendSMS("user", "pass", "MY-BRAND", "96598765432", "Hello", false)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "[ERR010] Zero balance." {
		t.Errorf("error = %q, want [ERR010] Zero balance.", err.Error())
	}
}

func TestGetCoverage_OK(t *testing.T) {
	serve(t, 200, map[string]interface{}{
		"result":   "OK",
		"coverage": []string{"965", "966", "971"},
	})
	raw, err := GetCoverage("user", "pass")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if raw == nil {
		t.Error("expected non-nil raw message")
	}
}

func TestGetCoverage_APIError(t *testing.T) {
	serve(t, 200, map[string]interface{}{
		"result":      "ERROR",
		"code":        "ERR033",
		"description": "No active coverage.",
	})
	_, err := GetCoverage("user", "pass")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUserAgentHeader(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": "OK", "available": 100, "purchased": 1000,
		})
	}))
	orig := baseURL
	baseURL = srv.URL + "/"
	origUA := UserAgent
	UserAgent = "kwtsms-cli/test"
	t.Cleanup(func() { srv.Close(); baseURL = orig; UserAgent = origUA })

	GetBalance("user", "pass")
	if gotUA != "kwtsms-cli/test" {
		t.Errorf("User-Agent = %q, want kwtsms-cli/test", gotUA)
	}
}
