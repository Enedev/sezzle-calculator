package calculator

type MultiplyOp struct{}

func (MultiplyOp) Name() string { return "multiply" }

func (MultiplyOp) Arity() Arity { return Arity{Min: 2, Max: 2} }

func (MultiplyOp) Apply(operands []float64) (float64, error) {
	return operands[0] * operands[1], nil
}
