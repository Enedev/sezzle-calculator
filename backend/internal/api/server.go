package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"sezzle-calculator/internal/calculator"
)

// NewRouter wires every route and middleware for the calculator API.
func NewRouter(registry *calculator.Registry) http.Handler {
	h := NewHandlers(registry)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.Healthz)
	mux.HandleFunc("GET /api/v1/operations", h.Operations)
	mux.HandleFunc("POST /api/v1/calculate", h.Calculate)

	var handler http.Handler = mux
	handler = CORSMiddleware(handler)
	handler = LoggingMiddleware(handler)
	handler = RecoverMiddleware(handler)
	return handler
}

// Serve starts the HTTP server on addr and blocks until either it fails to
// serve or ctx is canceled, in which case it shuts down gracefully. Signal
// handling is the caller's responsibility (see cmd/server), which keeps this
// function testable with an ordinary context instead of real OS signals.
func Serve(ctx context.Context, addr string, registry *calculator.Registry) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: NewRouter(registry),
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
