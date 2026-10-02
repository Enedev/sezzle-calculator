package calculator

type DivideOp struct{}

func (DivideOp) Name() string { return "divide" }

func (DivideOp) Arity() Arity { return Arity{Min: 2, Max: 2} }

func (DivideOp) Apply(operands []float64) (float64, error) {
	if operands[1] == 0 {
		return 0, ErrDivisionByZero
	}
	return operands[0] / operands[1], nil
}
