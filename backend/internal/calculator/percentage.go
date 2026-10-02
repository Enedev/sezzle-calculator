package calculator

import "math"

// PercentageOp computes "a percent of b". a*b/100 (multiplying before
// dividing) avoids a rounding error from a/100 for the common case where the
// true result is an exact float64, e.g. percentage(7, 100) == 7, not
// 7.000000000000001. But multiplying first can overflow to +-Inf for huge
// operands even when the true a*b/100 is finite and in range, so that case
// falls back to (a/100)*b, which shrinks a before multiplying.
type PercentageOp struct{}

func (PercentageOp) Name() string { return "percentage" }

func (PercentageOp) Arity() Arity { return Arity{Min: 2, Max: 2} }

func (PercentageOp) Apply(operands []float64) (float64, error) {
	a, b := operands[0], operands[1]
	result := a * b / 100
	if math.IsInf(result, 0) {
		return (a / 100) * b, nil
	}
	return result, nil
}
