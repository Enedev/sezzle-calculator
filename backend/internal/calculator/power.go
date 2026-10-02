package calculator

import "math"

type PowerOp struct{}

func (PowerOp) Name() string { return "power" }

func (PowerOp) Arity() Arity { return Arity{Min: 2, Max: 2} }

func (PowerOp) Apply(operands []float64) (float64, error) {
	return math.Pow(operands[0], operands[1]), nil
}
