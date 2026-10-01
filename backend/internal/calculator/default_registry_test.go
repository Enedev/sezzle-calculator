package calculator_test

import (
	"errors"
	"testing"

	"sezzle-calculator/internal/calculator"
)

func TestNewDefaultRegistry_HasArithmeticOperations(t *testing.T) {
	r := calculator.NewDefaultRegistry()
	want := []string{"add", "subtract", "multiply", "divide"}
	for _, name := range want {
		if _, ok := r.Get(name); !ok {
			t.Fatalf("expected operation %q to be registered", name)
		}
	}
}

func TestNewDefaultRegistry_InvalidOperandCount(t *testing.T) {
	r := calculator.NewDefaultRegistry()
	tests := []struct {
		name     string
		op       string
		operands []float64
	}{
		{"add too few", "add", []float64{1}},
		{"add too many", "add", []float64{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := r.Calculate(tt.op, tt.operands)
			if !errors.Is(err, calculator.ErrInvalidOperandCount) {
				t.Fatalf("expected ErrInvalidOperandCount, got %v", err)
			}
		})
	}
}
