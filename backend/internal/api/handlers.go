package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strings"

	"sezzle-calculator/internal/calculator"
)

// maxRequestBodyBytes caps the request body for /calculate. A handful of
// floats and an operation name comfortably fit well under 1 KiB.
const maxRequestBodyBytes = 1024

type Handlers struct {
	registry *calculator.Registry
}

func NewHandlers(registry *calculator.Registry) *Handlers {
	return &Handlers{registry: registry}
}

func (h *Handlers) Healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handlers) Operations(w http.ResponseWriter, r *http.Request) {
	ops := h.registry.List()
	infos := make([]OperationInfo, len(ops))
	for i, op := range ops {
		a := op.Arity()
		infos[i] = OperationInfo{Name: op.Name(), Arity: ArityInfo{Min: a.Min, Max: a.Max}}
	}
	writeJSON(w, http.StatusOK, OperationsResponse{Operations: infos})
}

func (h *Handlers) Calculate(w http.ResponseWriter, r *http.Request) {
	if !hasJSONContentType(r) {
		writeError(w, http.StatusUnsupportedMediaType, ErrCodeUnsupportedMediaType,
			"Content-Type must be application/json")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	var req CalculateRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			writeError(w, http.StatusRequestEntityTooLarge, ErrCodeRequestTooLarge,
				"request body too large")
			return
		}
		code, message := decodeJSONError(err)
		writeError(w, http.StatusBadRequest, code, message)
		return
	}
	if dec.More() {
		writeError(w, http.StatusBadRequest, ErrCodeInvalidJSON,
			"request body must contain a single JSON object")
		return
	}

	// MISSING_FIELD is reserved for an absent/empty `operation`. An
	// absent or empty `operands` array is instead left to flow into
	// Registry.Calculate, which already rejects it as INVALID_OPERAND_COUNT
	// (every operation's Arity.Min is >= 1) — one less special case here,
	// and a single consistent code for "wrong number of operands."
	if req.Operation == "" {
		writeError(w, http.StatusBadRequest, ErrCodeMissingField, "operation is required")
		return
	}

	operands := make([]float64, len(req.Operands))
	for i, p := range req.Operands {
		if p == nil {
			writeError(w, http.StatusBadRequest, ErrCodeInvalidJSON,
				fmt.Sprintf("operand at index %d must be a number, not null", i))
			return
		}
		operands[i] = *p
	}

	result, err := h.registry.Calculate(req.Operation, operands)
	if err != nil {
		code, status, message := domainErrorToAPI(err)
		writeError(w, status, code, message)
		return
	}

	writeJSON(w, http.StatusOK, CalculateResponse{Result: result})
}

func hasJSONContentType(r *http.Request) bool {
	ct := r.Header.Get("Content-Type")
	if ct == "" {
		return false
	}
	mediaType, _, err := mime.ParseMediaType(ct)
	if err != nil {
		return false
	}
	return mediaType == "application/json"
}

// decodeJSONError classifies a json.Decoder.Decode error into an API error
// code with a friendly message. encoding/json has no exported error type for
// "unknown field" (it's a plain error built with fmt.Errorf), so that one
// case is matched on its message text; everything else collapses to a
// generic malformed-JSON response rather than leaking Go's internal wording.
func decodeJSONError(err error) (ErrorCode, string) {
	if strings.Contains(err.Error(), "unknown field") {
		return ErrCodeUnknownField, "request body contains an unknown field"
	}
	return ErrCodeInvalidJSON, "request body is not valid JSON"
}
