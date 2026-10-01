package calculator_test

import (
	"errors"
	"math"
	"testing"

	"sezzle-calculator/internal/calculator"
)

func TestPercentageOp(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"20 percent of 50", 20, 50, 10, nil},
		{"100 percent", 100, 7, 7, nil},
		{"zero percent", 0, 50, 0, nil},
		{"negative percent", -50, 10, -5, nil},
		{"overflow", math.MaxFloat64, 300, 0, calculator.ErrResultOutOfRange},
	}
	r := calculator.NewDefaultRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Calculate("percentage", []float64{tt.a, tt.b})
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
