package calculator_test

import (
	"errors"
	"testing"

	"sezzle-calculator/internal/calculator"
)

func TestDivideOp(t *testing.T) {
	tests := []struct {
		name    string
		a, b    float64
		want    float64
		wantErr error
	}{
		{"positive", 10, 2, 5, nil},
		{"negative divisor", 10, -2, -5, nil},
		{"negative numerator by zero", -10, 0, 0, calculator.ErrDivisionByZero},
		{"zero by zero", 0, 0, 0, calculator.ErrDivisionByZero},
		{"positive by zero", 7, 0, 0, calculator.ErrDivisionByZero},
	}
	r := calculator.NewDefaultRegistry()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.Calculate("divide", []float64{tt.a, tt.b})
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
