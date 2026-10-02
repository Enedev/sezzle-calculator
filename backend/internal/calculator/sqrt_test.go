package calculator_test

import (
	"errors"
	"testing"

	"sezzle-calculator/internal/calculator"
)

func TestSqrtOp(t *testing.T) {
	tests := []struct {
		name    string
		x       float64
		want    float64
		wantErr error
	}{
		{"perfect square", 9, 3, nil},
		{"zero", 0, 0, nil},
		{"non-perfect square", 2, 1.4142135623730951, nil},
		{"negative", -4, 0, calculator.ErrMathDomain},
	}
	r := calculator.NewDefaultRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Calculate("sqrt", []float64{tt.x})
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
