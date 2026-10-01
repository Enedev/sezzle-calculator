package calculator

// Arity declares how many operands an Operation accepts.
type Arity struct {
	Min, Max int
}

// Operation is a single calculator strategy: pure math, no HTTP awareness.
type Operation interface {
	Name() string
	Arity() Arity
	Apply(operands []float64) (float64, error)
}
