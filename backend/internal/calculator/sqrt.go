package calculator

import "math"

type SqrtOp struct{}

func (SqrtOp) Name() string { return "sqrt" }

func (SqrtOp) Arity() Arity { return Arity{Min: 1, Max: 1} }

func (SqrtOp) Apply(operands []float64) (float64, error) {
	if operands[0] < 0 {
		return 0, ErrMathDomain
	}
	return math.Sqrt(operands[0]), nil
}
