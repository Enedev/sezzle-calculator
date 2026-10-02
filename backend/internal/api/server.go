package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"sezzle-calculator/internal/calculator"
)

// newAPIHandler builds just the /healthz, /api/v1/* routes plus the JSON
// 404/405 envelope for unmatched /api/ paths — no CORS/logging/recovery yet,
// and no awareness of static files. Shared by NewRouter and
// NewRouterWithStatic so both stay in sync automatically.
func newAPIHandler(registry *calculator.Registry) http.Handler {
	h := NewHandlers(registry)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /api/v1/operations", h.Operations)
	mux.HandleFunc("POST /api/v1/calculate", h.Calculate)

	return apiRouteErrorsMiddleware(mux)
}

func withCommonMiddleware(handler http.Handler) http.Handler {
	handler = CORSMiddleware(handler)
	handler = LoggingMiddleware(handler)
	handler = RecoverMiddleware(handler)
	return handler
}

// NewRouter wires every route and middleware for the calculator API alone —
// used standalone in local/dev runs and by every existing test.
func NewRouter(registry *calculator.Registry) http.Handler {
	return withCommonMiddleware(newAPIHandler(registry))
}

// NewRouterWithStatic additionally serves a built frontend (e.g. Vite's
// dist/) from staticDir for any path that isn't under /api/ or /healthz,
// falling back to staticDir/index.html for unmatched paths so client-side
// routing works. /api/ and /healthz always take precedence and are never
// shadowed by a static file or the SPA fallback.
func NewRouterWithStatic(registry *calculator.Registry, staticDir string) http.Handler {
	return withCommonMiddleware(withStaticFallback(newAPIHandler(registry), staticDir))
}

// Serve starts the HTTP server on addr and blocks until either it fails to
// serve or ctx is canceled, in which case it shuts down gracefully. Signal
// handling is the caller's responsibility (see cmd/server), which keeps this
// function testable with an ordinary context instead of real OS signals. If
// staticDir is empty, the server is API-only (no frontend to serve).
func Serve(ctx context.Context, addr string, registry *calculator.Registry, staticDir string) error {
	var handler http.Handler
	if staticDir == "" {
		handler = NewRouter(registry)
	} else {
		handler = NewRouterWithStatic(registry, staticDir)
	}

	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	serveErr := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serveErr <- err
			return
		}
		serveErr <- nil
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
