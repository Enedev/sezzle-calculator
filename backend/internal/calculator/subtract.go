package calculator

type SubtractOp struct{}

func (SubtractOp) Name() string { return "subtract" }

func (SubtractOp) Arity() Arity { return Arity{Min: 2, Max: 2} }

func (SubtractOp) Apply(operands []float64) (float64, error) {
	return operands[0] - operands[1], nil
}
