package calculator_test

import (
	"errors"
	"math"
	"testing"

	"sezzle-calculator/internal/calculator"
)

func TestAddOp(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive", 2, 3, 5, nil},
		{"negative plus positive", -2, 3, 1, nil},
		{"zeros", 0, 0, 0, nil},
		{"decimals", 1.5, 2.25, 3.75, nil},
		{"float precision trade-off", 0.1, 0.2, 0.30000000000000004, nil},
		{"overflow", math.MaxFloat64, math.MaxFloat64, 0, calculator.ErrResultOutOfRange},
	}
	r := calculator.NewDefaultRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Calculate("add", []float64{tt.a, tt.b})
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
