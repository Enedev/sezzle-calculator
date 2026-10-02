package api

// CalculateRequest is the wire shape for POST /api/v1/calculate.
//
// Operands is []*float64 (not []float64) so a literal JSON null inside the
// array can be detected explicitly. encoding/json silently leaves a
// non-pointer numeric destination untouched when the source value is null,
// which would turn `{"operands":[1,null]}` into operands=[1,0] with no error
// — exactly the kind of silent zero that could trigger a division by zero
// the caller never intended. Decoding into pointers surfaces the null as a
// nil element we can reject explicitly.
type CalculateRequest struct {
	Operation string     `json:"operation"`
	Operands  []*float64 `json:"operands"`
}

type CalculateResponse struct {
	Result float64 `json:"result"`
}

// ArityInfo mirrors calculator.Arity for the wire format. Kept separate from
// calculator.Arity so the domain package stays free of JSON tags/HTTP concerns.
type ArityInfo struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type OperationInfo struct {
	Name  string    `json:"name"`
	Arity ArityInfo `json:"arity"`
}

// OperationsResponse wraps the operation list in an object (rather than a
// bare array) so the contract can grow additional fields later without
// breaking clients that deserialize this response.
type OperationsResponse struct {
	Operations []OperationInfo `json:"operations"`
}

type ErrorCode string

const (
	ErrCodeInvalidJSON          ErrorCode = "INVALID_JSON"
	ErrCodeUnknownField         ErrorCode = "UNKNOWN_FIELD"
	ErrCodeMissingField         ErrorCode = "MISSING_FIELD"
	ErrCodeUnknownOperation     ErrorCode = "UNKNOWN_OPERATION"
	ErrCodeInvalidOperandCount  ErrorCode = "INVALID_OPERAND_COUNT"
	ErrCodeUnsupportedMediaType ErrorCode = "UNSUPPORTED_MEDIA_TYPE"
	ErrCodeRequestTooLarge      ErrorCode = "REQUEST_TOO_LARGE"
	ErrCodeNotFound             ErrorCode = "NOT_FOUND"
	ErrCodeMethodNotAllowed     ErrorCode = "METHOD_NOT_ALLOWED"
	ErrCodeDivisionByZero       ErrorCode = "DIVISION_BY_ZERO"
	ErrCodeMathDomainError      ErrorCode = "MATH_DOMAIN_ERROR"
	ErrCodeResultOutOfRange     ErrorCode = "RESULT_OUT_OF_RANGE"
	ErrCodeInternal             ErrorCode = "INTERNAL_ERROR"
)

type ErrorBody struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}
