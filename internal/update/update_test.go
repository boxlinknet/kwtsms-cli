package update

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func serve(t *testing.T, tag string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"tag_name": tag})
	}))
	t.Cleanup(srv.Close)
	return srv
}

func checkWith(t *testing.T, srv *httptest.Server, current string) (*Result, error) {
	t.Helper()
	orig := releaseURL
	releaseURL = srv.URL
	t.Cleanup(func() { releaseURL = orig })
	return check(current)
}

func TestCheck_NewerAvailable(t *testing.T) {
	srv := serve(t, "v1.1.0")
	result, err := checkWith(t, srv, "v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UpToDate {
		t.Error("expected UpToDate=false")
	}
	if result.Latest != "v1.1.0" {
		t.Errorf("Latest = %q, want v1.1.0", result.Latest)
	}
	if result.Current != "v1.0.0" {
		t.Errorf("Current = %q, want v1.0.0", result.Current)
	}
	if !strings.Contains(result.DownloadURL, "github.com") {
		t.Errorf("DownloadURL missing github.com: %q", result.DownloadURL)
	}
}

func TestCheck_AlreadyLatest(t *testing.T) {
	srv := serve(t, "v1.0.1")
	result, err := checkWith(t, srv, "v1.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.UpToDate {
		t.Error("expected UpToDate=true")
	}
	if result.DownloadURL != "" {
		t.Errorf("DownloadURL should be empty when up to date, got %q", result.DownloadURL)
	}
}

func TestCheck_DevBuild(t *testing.T) {
	srv := serve(t, "v1.0.1")
	_, err := checkWith(t, srv, "dev")
	if err == nil {
		t.Error("expected error for dev build")
	}
}

func TestCheck_EmptyVersion(t *testing.T) {
	srv := serve(t, "v1.0.1")
	_, err := checkWith(t, srv, "")
	if err == nil {
		t.Error("expected error for empty version")
	}
}

func TestCheck_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)
	_, err := checkWith(t, srv, "v1.0.0")
	if err == nil {
		t.Error("expected error on server 500")
	}
}

func TestCheck_MalformedTagFromAPI(t *testing.T) {
	srv := serve(t, "<script>alert(1)</script>")
	_, err := checkWith(t, srv, "v1.0.0")
	if err == nil {
		t.Error("expected error for malformed tag from API")
	}
}

func TestCheck_MalformedCurrentVersion(t *testing.T) {
	srv := serve(t, "v1.1.0")
	_, err := checkWith(t, srv, "not-a-version")
	if err == nil {
		t.Error("expected error for malformed current version")
	}
}

func TestCheckNow_NewerAvailable(t *testing.T) {
	srv := serve(t, "v9.9.9")
	orig := releaseURL
	releaseURL = srv.URL
	t.Cleanup(func() { releaseURL = orig })

	result, err := CheckNow("v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.UpToDate {
		t.Error("expected UpToDate=false")
	}
	if result.Latest != "v9.9.9" {
		t.Errorf("Latest = %q, want v9.9.9", result.Latest)
	}
}
