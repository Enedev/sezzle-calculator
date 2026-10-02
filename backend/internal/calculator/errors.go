package calculator

import "errors"

var (
	ErrUnknownOperation    = errors.New("unknown operation")
	ErrInvalidOperandCount = errors.New("invalid operand count")
	ErrDivisionByZero      = errors.New("division by zero")
	ErrMathDomain          = errors.New("math domain error")
	ErrResultOutOfRange    = errors.New("result out of range")
)
