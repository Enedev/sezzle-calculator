package calculator

// PercentageOp computes "a percent of b": (a/100) * b.
type PercentageOp struct{}

func (PercentageOp) Name() string { return "percentage" }

func (PercentageOp) Arity() Arity { return Arity{Min: 2, Max: 2} }

func (PercentageOp) Apply(operands []float64) (float64, error) {
	return (operands[0] / 100) * operands[1], nil
}
