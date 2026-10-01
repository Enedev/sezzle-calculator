package calculator_test

import (
	"errors"
	"math"
	"testing"

	"sezzle-calculator/internal/calculator"
)

func TestMultiplyOp(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive", 4, 5, 20, nil},
		{"negative times positive", -4, 5, -20, nil},
		{"zero", 0, 123, 0, nil},
		{"overflow", math.MaxFloat64, 2, 0, calculator.ErrResultOutOfRange},
	}
	r := calculator.NewDefaultRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Calculate("multiply", []float64{tt.a, tt.b})
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiplyOp_NegativeZeroNormalizedToPositiveZero(t *testing.T) {
	r := calculator.NewDefaultRegistry()
	got, err := r.Calculate("multiply", []float64{-4, 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0 {
		t.Fatalf("got %v, want 0", got)
	}
	if math.Signbit(got) {
		t.Fatalf("got negative zero, want positive zero so the API never serializes \"-0\"")
	}
}
