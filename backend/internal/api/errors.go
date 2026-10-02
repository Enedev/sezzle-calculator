package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"sezzle-calculator/internal/calculator"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("failed to encode response", "error", err)
	}
}

func writeError(w http.ResponseWriter, status int, code ErrorCode, message string) {
	writeJSON(w, status, ErrorResponse{Error: ErrorBody{Code: code, Message: message}})
}

// domainErrorToAPI maps a calculator domain error to its API error code,
// HTTP status, and a human-readable message.
func domainErrorToAPI(err error) (code ErrorCode, status int, message string) {
	switch {
	case errors.Is(err, calculator.ErrUnknownOperation):
		return ErrCodeUnknownOperation, http.StatusBadRequest, "unknown operation"
	case errors.Is(err, calculator.ErrInvalidOperandCount):
		return ErrCodeInvalidOperandCount, http.StatusBadRequest, "invalid operand count for this operation"
	case errors.Is(err, calculator.ErrDivisionByZero):
		return ErrCodeDivisionByZero, http.StatusUnprocessableEntity, "division by zero"
	case errors.Is(err, calculator.ErrMathDomain):
		return ErrCodeMathDomainError, http.StatusUnprocessableEntity, "operation is undefined for the given operands"
	case errors.Is(err, calculator.ErrResultOutOfRange):
		return ErrCodeResultOutOfRange, http.StatusUnprocessableEntity, "result is not a finite number"
	default:
		return ErrCodeInternal, http.StatusInternalServerError, "internal error"
	}
}
