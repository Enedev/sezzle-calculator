package calculator_test

import (
	"errors"
	"testing"

	"sezzle-calculator/internal/calculator"
)

func TestNewDefaultRegistry_HasSevenOperations(t *testing.T) {
	r := calculator.NewDefaultRegistry()
	want := []string{"add", "subtract", "multiply", "divide", "power", "sqrt", "percentage"}
	ops := r.List()
	if len(ops) != len(want) {
		t.Fatalf("expected %d operations, got %d", len(want), len(ops))
	}
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
		{"sqrt wrong arity", "sqrt", []float64{1, 2}},
		{"percentage wrong arity", "percentage", []float64{1}},
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
