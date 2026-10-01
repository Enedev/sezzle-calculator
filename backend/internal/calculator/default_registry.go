package calculator

// NewDefaultRegistry returns a Registry pre-populated with every built-in operation.
func NewDefaultRegistry() *Registry {
	r := NewRegistry()
	r.Register(AddOp{})
	r.Register(SubtractOp{})
	r.Register(MultiplyOp{})
	r.Register(DivideOp{})
	return r
}
