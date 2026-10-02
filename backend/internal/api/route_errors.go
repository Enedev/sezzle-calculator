package api

import (
	"bytes"
	"net/http"
	"strings"
)

// apiRouteErrorsMiddleware rewrites the default plain-text 404/405 bodies
// ServeMux produces into the same JSON error envelope as every other
// response — but only for paths under /api/. "/" (and anything else outside
// /api/) is left alone since it's reserved for serving the frontend from
// Phase 4 onward.
func apiRouteErrorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r)
			return
		}

		buf := newBufferingResponseWriter()
		next.ServeHTTP(buf, r)

		switch buf.status {
		case http.StatusNotFound:
			writeError(w, http.StatusNotFound, ErrCodeNotFound, "resource not found")
		case http.StatusMethodNotAllowed:
			if allow := buf.header.Get("Allow"); allow != "" {
				w.Header().Set("Allow", allow)
			}
			writeError(w, http.StatusMethodNotAllowed, ErrCodeMethodNotAllowed, "method not allowed")
		default:
			for key, values := range buf.header {
				w.Header()[key] = values
			}
			w.WriteHeader(buf.status)
			w.Write(buf.body.Bytes())
		}
	})
}

// bufferingResponseWriter lets us inspect the status ServeMux would have
// sent before committing anything to the real ResponseWriter.
type bufferingResponseWriter struct {
	header http.Header
	status int
	body   bytes.Buffer
}

func newBufferingResponseWriter() *bufferingResponseWriter {
	return &bufferingResponseWriter{header: make(http.Header), status: http.StatusOK}
}

func (w *bufferingResponseWriter) Header() http.Header { return w.header }

func (w *bufferingResponseWriter) WriteHeader(status int) { w.status = status }

func (w *bufferingResponseWriter) Write(b []byte) (int, error) { return w.body.Write(b) }
