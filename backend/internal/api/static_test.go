package api_test

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"sezzle-calculator/internal/api"
	"sezzle-calculator/internal/calculator"
)

func newTestStaticDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>spa shell</html>"), 0o644); err != nil {
		t.Fatalf("failed to write index.html: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "app.js"), []byte("console.log('hi')"), 0o644); err != nil {
		t.Fatalf("failed to write app.js: %v", err)
	}
	return dir
}

func newTestStaticRouter(t *testing.T) http.Handler {
	return api.NewRouterWithStatic(calculator.NewDefaultRegistry(), newTestStaticDir(t))
}

func TestStatic_ServesRealFilesDirectly(t *testing.T) {
	handler := newTestStaticRouter(t)
	rec := doRequest(t, handler, http.MethodGet, "/app.js", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.String() != "console.log('hi')" {
		t.Fatalf("body = %q, want the real file contents", rec.Body.String())
	}
}

func TestStatic_FallsBackToIndexHTMLForUnknownPaths(t *testing.T) {
	handler := newTestStaticRouter(t)
	rec := doRequest(t, handler, http.MethodGet, "/some/client-side/route", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.String() != "<html>spa shell</html>" {
		t.Fatalf("body = %q, want the SPA shell", rec.Body.String())
	}
}

func TestStatic_RootServesIndexHTML(t *testing.T) {
	handler := newTestStaticRouter(t)
	rec := doRequest(t, handler, http.MethodGet, "/", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if rec.Body.String() != "<html>spa shell</html>" {
		t.Fatalf("body = %q, want the SPA shell", rec.Body.String())
	}
}

// /api/ must always win over the static fallback, even when it's an
// unmatched /api/ path that would otherwise just be "not a real file" to the
// static handler. It must get the JSON 404 envelope, never the SPA shell.
func TestStatic_APIPrefixTakesPrecedenceOverStaticFallback(t *testing.T) {
	handler := newTestStaticRouter(t)
	rec := doRequest(t, handler, http.MethodGet, "/api/v1/nonexistent", "", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
	got := decodeError(t, rec)
	if got.Error.Code != api.ErrCodeNotFound {
		t.Fatalf("error code = %q, want %q (should be the JSON envelope, not the SPA shell)", got.Error.Code, api.ErrCodeNotFound)
	}
}

func TestStatic_RealAPIRouteStillWorks(t *testing.T) {
	handler := newTestStaticRouter(t)
	rec := doRequest(t, handler, http.MethodPost, "/api/v1/calculate", "application/json", `{"operation":"add","operands":[2,3]}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestStatic_HealthzTakesPrecedenceOverStaticFallback(t *testing.T) {
	handler := newTestStaticRouter(t)
	rec := doRequest(t, handler, http.MethodGet, "/healthz", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
}
