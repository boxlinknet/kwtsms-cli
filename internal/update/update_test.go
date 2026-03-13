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

func checkWith(t *testing.T, srv *httptest.Server, current string) string {
	t.Helper()
	orig := releaseURL
	releaseURL = srv.URL
	t.Cleanup(func() { releaseURL = orig })
	return check(current)
}

func TestCheck_NewerAvailable(t *testing.T) {
	srv := serve(t, "v1.1.0")
	msg := checkWith(t, srv, "v1.0.0")
	if msg == "" {
		t.Fatal("expected notice, got empty string")
	}
	if !strings.Contains(msg, "v1.1.0") {
		t.Errorf("notice missing new version: %q", msg)
	}
	if !strings.Contains(msg, "github.com") {
		t.Errorf("notice missing release URL: %q", msg)
	}
}

func TestCheck_AlreadyLatest(t *testing.T) {
	srv := serve(t, "v1.0.1")
	msg := checkWith(t, srv, "v1.0.1")
	if msg != "" {
		t.Errorf("expected empty notice when up to date, got %q", msg)
	}
}

func TestCheck_DevBuild(t *testing.T) {
	srv := serve(t, "v1.0.1")
	msg := checkWith(t, srv, "dev")
	if msg != "" {
		t.Errorf("expected empty notice for dev build, got %q", msg)
	}
}

func TestCheck_EmptyVersion(t *testing.T) {
	srv := serve(t, "v1.0.1")
	msg := checkWith(t, srv, "")
	if msg != "" {
		t.Errorf("expected empty notice for empty version, got %q", msg)
	}
}

func TestCheck_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	t.Cleanup(srv.Close)
	msg := checkWith(t, srv, "v1.0.0")
	if msg != "" {
		t.Errorf("expected silent failure on server error, got %q", msg)
	}
}

func TestCheck_MalformedTagFromAPI(t *testing.T) {
	srv := serve(t, "<script>alert(1)</script>")
	msg := checkWith(t, srv, "v1.0.0")
	if msg != "" {
		t.Errorf("expected empty notice for malformed tag, got %q", msg)
	}
}

func TestCheck_MalformedCurrentVersion(t *testing.T) {
	srv := serve(t, "v1.1.0")
	msg := checkWith(t, srv, "not-a-version")
	if msg != "" {
		t.Errorf("expected empty notice for malformed current version, got %q", msg)
	}
}

func TestCheckNow_NewerAvailable(t *testing.T) {
	srv := serve(t, "v9.9.9")
	msg := checkWith(t, srv, "v1.0.0")
	if msg == "" {
		t.Error("expected notice from CheckNow, got empty string")
	}
}
