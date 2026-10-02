package api_test

import (
	"context"
	"testing"
	"time"

	"sezzle-calculator/internal/api"
	"sezzle-calculator/internal/calculator"
)

func TestServe_ShutsDownGracefullyOnContextCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := api.Serve(ctx, "127.0.0.1:0", calculator.NewDefaultRegistry(), "")
	if err != nil {
		t.Fatalf("Serve returned error: %v", err)
	}
}
