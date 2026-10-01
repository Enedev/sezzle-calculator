package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"sezzle-calculator/internal/api"
	"sezzle-calculator/internal/calculator"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	registry := calculator.NewDefaultRegistry()

	slog.Info("starting server", "addr", addr)
	if err := api.Serve(ctx, addr, registry); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
