package calculator_test

import (
	"errors"
	"testing"

	"sezzle-calculator/internal/calculator"
)

func TestPowerOp(t *testing.T) {
	tests := []struct {
		name    string
		base    float64
		exp     float64
		want    float64
		wantErr error
	}{
		{"positive integer exponent", 2, 10, 1024, nil},
		{"negative exponent", 2, -1, 0.5, nil},
		{"zero to the zero", 0, 0, 1, nil},
		{"zero to negative exponent", 0, -1, 0, calculator.ErrResultOutOfRange},
		{"fractional root of negative base", -8, 0.5, 0, calculator.ErrMathDomain},
		{"overflow", 10, 400, 0, calculator.ErrResultOutOfRange},
	}
	r := calculator.NewDefaultRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Calculate("power", []float64{tt.base, tt.exp})
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
