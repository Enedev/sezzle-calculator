package calculator

import (
	"fmt"
	"math"
	"sort"
)

// Registry holds Operations by name (the strategy pattern). Adding a new
// operation is one new Operation type plus one Register call here.
type Registry struct {
	ops map[string]Operation
}

func NewRegistry() *Registry {
	return &Registry{ops: make(map[string]Operation)}
}

// Register adds an operation under its own name. Registering the same name
// twice is a programmer error and panics rather than silently overwriting.
func (r *Registry) Register(op Operation) {
	name := op.Name()
	if _, exists := r.ops[name]; exists {
		panic(fmt.Sprintf("calculator: operation %q already registered", name))
	}
	r.ops[name] = op
}

func (r *Registry) Get(name string) (Operation, bool) {
	op, ok := r.ops[name]
	return op, ok
}

// List returns every registered operation sorted by name, so callers like
// GET /api/v1/operations get a stable order.
func (r *Registry) List() []Operation {
	names := make([]string, 0, len(r.ops))
	for name := range r.ops {
		names = append(names, name)
	}
	sort.Strings(names)

	ops := make([]Operation, 0, len(names))
	for _, name := range names {
		ops = append(ops, r.ops[name])
	}
	return ops
}

// Calculate is the single entry point for evaluating an operation by name.
// It owns every math/contract rule shared across operations: operand-count
// validation against the operation's declared Arity, and classification of
// non-finite results (NaN -> ErrMathDomain, +/-Inf -> ErrResultOutOfRange).
// It also normalizes a -0 result to 0 so the API never serializes "-0".
func (r *Registry) Calculate(name string, operands []float64) (float64, error) {
	op, ok := r.Get(name)
	if !ok {
		return 0, ErrUnknownOperation
	}

	arity := op.Arity()
	if len(operands) < arity.Min || len(operands) > arity.Max {
		return 0, ErrInvalidOperandCount
	}

	result, err := op.Apply(operands)
	if err != nil {
		return 0, err
	}

	switch {
	case math.IsNaN(result):
		return 0, ErrMathDomain
	case math.IsInf(result, 0):
		return 0, ErrResultOutOfRange
	case result == 0:
		return 0, nil
	default:
		return result, nil
	}
}
