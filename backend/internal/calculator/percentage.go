package calculator

// PercentageOp computes "a percent of b": a*b/100. Multiplying before
// dividing (rather than (a/100)*b) avoids introducing a rounding error from
// a/100 for the common case where a*b/100 lands on an exact float64, e.g.
// percentage(7, 100) == 7, not 7.000000000000001.
type PercentageOp struct{}

func (PercentageOp) Name() string { return "percentage" }

func (PercentageOp) Arity() Arity { return Arity{Min: 2, Max: 2} }

func (PercentageOp) Apply(operands []float64) (float64, error) {
	return operands[0] * operands[1] / 100, nil
}
