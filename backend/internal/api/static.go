package api

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// withStaticFallback serves files out of staticDir for any request that
// isn't under /api/ or /healthz, and falls back to staticDir/index.html for
// paths that don't match a real file — the standard single-page-app
// routing pattern. /api/ and /healthz requests are always handed straight
// to next, so an unmatched /api/ path still gets the JSON 404 from
// apiRouteErrorsMiddleware instead of being swallowed by the SPA fallback.
func withStaticFallback(next http.Handler, staticDir string) http.Handler {
	fileServer := http.FileServer(http.Dir(staticDir))
	indexPath := filepath.Join(staticDir, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/healthz" {
			next.ServeHTTP(w, r)
			return
		}

		cleanPath := path.Clean(r.URL.Path)
		fullPath := filepath.Join(staticDir, cleanPath)
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		http.ServeFile(w, r, indexPath)
	})
}
