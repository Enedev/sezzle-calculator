package calculator_test

import (
	"errors"
	"math"
	"testing"

	"sezzle-calculator/internal/calculator"
)

// stubOp is a controllable Operation used to test Registry behavior in
// isolation from any real operation's math.
type stubOp struct {
	name  string
	arity calculator.Arity
	fn    func([]float64) (float64, error)
}

func (s stubOp) Name() string            { return s.name }
func (s stubOp) Arity() calculator.Arity { return s.arity }
func (s stubOp) Apply(operands []float64) (float64, error) {
	return s.fn(operands)
}

func constOp(name string, arity calculator.Arity, result float64) stubOp {
	return stubOp{name: name, arity: arity, fn: func([]float64) (float64, error) { return result, nil }}
}

func TestRegistry_UnknownOperation(t *testing.T) {
	r := calculator.NewRegistry()
	_, err := r.Calculate("foo", []float64{1, 2})
	if !errors.Is(err, calculator.ErrUnknownOperation) {
		t.Fatalf("expected ErrUnknownOperation, got %v", err)
	}
}

func TestRegistry_InvalidOperandCount(t *testing.T) {
	r := calculator.NewRegistry()
	r.Register(constOp("binary", calculator.Arity{Min: 2, Max: 2}, 0))

	tests := []struct {
		name     string
		operands []float64
	}{
		{"too few", []float64{1}},
		{"too many", []float64{1, 2, 3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := r.Calculate("binary", tt.operands)
			if !errors.Is(err, calculator.ErrInvalidOperandCount) {
				t.Fatalf("expected ErrInvalidOperandCount, got %v", err)
			}
		})
	}
}

func TestRegistry_Calculate_Dispatch(t *testing.T) {
	r := calculator.NewRegistry()
	r.Register(constOp("five", calculator.Arity{Min: 2, Max: 2}, 5))

	result, err := r.Calculate("five", []float64{2, 3})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != 5 {
		t.Fatalf("expected 5, got %v", result)
	}
}

func TestRegistry_Calculate_NaNBecomesMathDomainError(t *testing.T) {
	r := calculator.NewRegistry()
	r.Register(stubOp{
		name:  "nan",
		arity: calculator.Arity{Min: 1, Max: 1},
		fn:    func([]float64) (float64, error) { return math.NaN(), nil },
	})

	_, err := r.Calculate("nan", []float64{1})
	if !errors.Is(err, calculator.ErrMathDomain) {
		t.Fatalf("expected ErrMathDomain, got %v", err)
	}
}

func TestRegistry_Calculate_InfBecomesResultOutOfRange(t *testing.T) {
	r := calculator.NewRegistry()
	r.Register(stubOp{
		name:  "inf",
		arity: calculator.Arity{Min: 1, Max: 1},
		fn:    func([]float64) (float64, error) { return math.Inf(1), nil },
	})

	_, err := r.Calculate("inf", []float64{1})
	if !errors.Is(err, calculator.ErrResultOutOfRange) {
		t.Fatalf("expected ErrResultOutOfRange, got %v", err)
	}
}

func TestRegistry_Calculate_NegativeZeroNormalizedToPositiveZero(t *testing.T) {
	r := calculator.NewRegistry()
	r.Register(stubOp{
		name:  "negzero",
		arity: calculator.Arity{Min: 1, Max: 1},
		fn:    func([]float64) (float64, error) { return math.Copysign(0, -1), nil },
	})

	got, err := r.Calculate("negzero", []float64{1})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != 0 || math.Signbit(got) {
		t.Fatalf("got %v (signbit=%v), want positive zero", got, math.Signbit(got))
	}
}

func TestRegistry_Register_DuplicatePanics(t *testing.T) {
	r := calculator.NewRegistry()
	op := constOp("dup", calculator.Arity{Min: 1, Max: 1}, 0)
	r.Register(op)

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic on duplicate registration")
		}
	}()
	r.Register(op)
}

func TestRegistry_List_SortedByName(t *testing.T) {
	r := calculator.NewRegistry()
	r.Register(constOp("zebra", calculator.Arity{Min: 1, Max: 1}, 0))
	r.Register(constOp("alpha", calculator.Arity{Min: 1, Max: 1}, 0))
	r.Register(constOp("mike", calculator.Arity{Min: 1, Max: 1}, 0))

	ops := r.List()
	got := make([]string, len(ops))
	for i, op := range ops {
		got[i] = op.Name()
	}
	want := []string{"alpha", "mike", "zebra"}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("List() = %v, want %v", got, want)
		}
	}
}
