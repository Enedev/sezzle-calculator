package api_test

import (
	"net/http"
	"testing"

	"sezzle-calculator/internal/api"
)

// Requests under /api/ that don't match any route, or match the wrong
// method, must get the same JSON error envelope as every other error —
// not Go's default plain-text 404/405 bodies.

func TestUnknownAPIRoute_ReturnsJSONNotFound(t *testing.T) {
	handler := newTestRouter()
	rec := doRequest(t, handler, http.MethodGet, "/api/v1/nonexistent", "", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	got := decodeError(t, rec)
	if got.Error.Code != api.ErrCodeNotFound {
		t.Fatalf("error code = %q, want %q", got.Error.Code, api.ErrCodeNotFound)
	}
}

func TestWrongMethodOnAPIRoute_ReturnsJSONMethodNotAllowedWithAllowHeader(t *testing.T) {
	handler := newTestRouter()
	rec := doRequest(t, handler, http.MethodGet, "/api/v1/calculate", "", "")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
	if allow := rec.Header().Get("Allow"); allow == "" {
		t.Fatal("expected a non-empty Allow header")
	}
	got := decodeError(t, rec)
	if got.Error.Code != api.ErrCodeMethodNotAllowed {
		t.Fatalf("error code = %q, want %q", got.Error.Code, api.ErrCodeMethodNotAllowed)
	}
}

func TestRootPath_IsNotWrappedInJSON(t *testing.T) {
	// "/" is reserved for serving the frontend starting in Phase 4 and must
	// keep Go's default behavior, not the API's JSON error envelope.
	handler := newTestRouter()
	rec := doRequest(t, handler, http.MethodGet, "/", "", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	ct := rec.Header().Get("Content-Type")
	if ct == "application/json" {
		t.Fatalf("expected root path to NOT use the JSON error envelope, got Content-Type %q", ct)
	}
}
