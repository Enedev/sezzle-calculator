package calculator

type AddOp struct{}

func (AddOp) Name() string { return "add" }

func (AddOp) Arity() Arity { return Arity{Min: 2, Max: 2} }

func (AddOp) Apply(operands []float64) (float64, error) {
	return operands[0] + operands[1], nil
}
