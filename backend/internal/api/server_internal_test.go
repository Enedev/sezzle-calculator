package api

import (
	"net/http"
	"testing"
)

// Guards against Slowloris-style connection exhaustion: without these, a
// client that sends headers/body slowly (or never finishes) can hold a
// connection open indefinitely.
func TestNewHTTPServer_SetsTimeouts(t *testing.T) {
	srv := newHTTPServer(":0", http.NewServeMux())

	if srv.ReadHeaderTimeout <= 0 {
		t.Error("expected a positive ReadHeaderTimeout")
	}
	if srv.ReadTimeout <= 0 {
		t.Error("expected a positive ReadTimeout")
	}
	if srv.WriteTimeout <= 0 {
		t.Error("expected a positive WriteTimeout")
	}
	if srv.IdleTimeout <= 0 {
		t.Error("expected a positive IdleTimeout")
	}
}
