package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sezzle-calculator/internal/api"
	"sezzle-calculator/internal/calculator"
)

// contractNewHandler builds a fresh router backed by the default calculator
// registry, exactly as an external consumer of the package would.
func contractNewHandler() http.Handler {
	return api.NewRouter(calculator.NewDefaultRegistry())
}

// contractDo issues a request against the handler using httptest, without
// opening a real network listener.
func contractDo(t *testing.T, handler http.Handler, method, path string, body []byte, contentType *string) *httptest.ResponseRecorder {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	} else {
		reader = bytes.NewReader(nil)
	}
	req := httptest.NewRequest(method, path, reader)
	if contentType != nil {
		req.Header.Set("Content-Type", *contentType)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func contractJSONContentType() *string {
	v := "application/json"
	return &v
}

// contractDecodeErrorResponse decodes the response body into an
// api.ErrorResponse, failing the test if that's not possible.
func contractDecodeErrorResponse(t *testing.T, rec *httptest.ResponseRecorder) api.ErrorResponse {
	t.Helper()
	var out api.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("failed to decode error response: %v; body=%s", err, rec.Body.String())
	}
	return out
}

func contractDecodeCalculateResponse(t *testing.T, rec *httptest.ResponseRecorder) api.CalculateResponse {
	t.Helper()
	var out api.CalculateResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("failed to decode calculate response: %v; body=%s", err, rec.Body.String())
	}
	return out
}

// contractCalcCase describes one black-box calculate request/response case.
type contractCalcCase struct {
	name        string
	rawBody     string // raw JSON body sent verbatim; if empty, built from operation/operands
	operation   string
	operands    []float64
	contentType *string // nil = use application/json; pass explicit string (incl. empty) to override

	wantStatus int
	wantResult float64
	wantCode   api.ErrorCode
	checkExact bool // if true, assert exact float equality on result (no delta)
}

func contractBuildBody(operation string, operands []float64) []byte {
	type reqBody struct {
		Operation string    `json:"operation"`
		Operands  []float64 `json:"operands"`
	}
	b, _ := json.Marshal(reqBody{Operation: operation, Operands: operands})
	return b
}

func TestContract_Calculate(t *testing.T) {
	handler := contractNewHandler()

	cases := []contractCalcCase{
		// --- happy paths, one per operation ---
		{name: "add happy path", operation: "add", operands: []float64{2, 3}, wantStatus: 200, wantResult: 5, checkExact: true},
		{name: "subtract happy path", operation: "subtract", operands: []float64{10, 4}, wantStatus: 200, wantResult: 6, checkExact: true},
		{name: "multiply happy path", operation: "multiply", operands: []float64{6, 7}, wantStatus: 200, wantResult: 42, checkExact: true},
		{name: "divide happy path", operation: "divide", operands: []float64{10, 4}, wantStatus: 200, wantResult: 2.5, checkExact: true},
		{name: "power happy path", operation: "power", operands: []float64{2, 10}, wantStatus: 200, wantResult: 1024, checkExact: true},
		{name: "sqrt happy path", operation: "sqrt", operands: []float64{16}, wantStatus: 200, wantResult: 4, checkExact: true},

		// --- percentage semantics: "a percent of b" ---
		{name: "percentage 20 of 50", operation: "percentage", operands: []float64{20, 50}, wantStatus: 200, wantResult: 10, checkExact: true},
		{name: "percentage 100 of 7", operation: "percentage", operands: []float64{100, 7}, wantStatus: 200, wantResult: 7, checkExact: true},
		{name: "percentage 7 of 100 exact", operation: "percentage", operands: []float64{7, 100}, wantStatus: 200, wantResult: 7, checkExact: true},

		// --- MISSING_FIELD: operation absent or empty ---
		{name: "missing operation field entirely", rawBody: `{"operands":[1,2]}`, wantStatus: 400, wantCode: "MISSING_FIELD"},
		{name: "empty operation string", rawBody: `{"operation":"","operands":[1,2]}`, wantStatus: 400, wantCode: "MISSING_FIELD"},

		// --- UNKNOWN_OPERATION ---
		{name: "unknown operation name", operation: "modulo", operands: []float64{1, 2}, wantStatus: 400, wantCode: "UNKNOWN_OPERATION"},

		// --- INVALID_OPERAND_COUNT, including absent/empty operands ---
		{name: "operands absent entirely", rawBody: `{"operation":"add"}`, wantStatus: 400, wantCode: "INVALID_OPERAND_COUNT"},
		{name: "operands empty array", rawBody: `{"operation":"add","operands":[]}`, wantStatus: 400, wantCode: "INVALID_OPERAND_COUNT"},
		{name: "add with 1 operand", operation: "add", operands: []float64{1}, wantStatus: 400, wantCode: "INVALID_OPERAND_COUNT"},
		{name: "add with 3 operands", operation: "add", operands: []float64{1, 2, 3}, wantStatus: 400, wantCode: "INVALID_OPERAND_COUNT"},
		{name: "sqrt with 2 operands", operation: "sqrt", operands: []float64{1, 2}, wantStatus: 400, wantCode: "INVALID_OPERAND_COUNT"},
		{name: "sqrt with 0 operands", rawBody: `{"operation":"sqrt","operands":[]}`, wantStatus: 400, wantCode: "INVALID_OPERAND_COUNT"},
		{name: "percentage with 1 operand", operation: "percentage", operands: []float64{1}, wantStatus: 400, wantCode: "INVALID_OPERAND_COUNT"},

		// --- INVALID_JSON: malformed syntax, null in operands, string number, trailing data ---
		{name: "malformed json syntax", rawBody: `{"operation":"add","operands":[1,2]`, wantStatus: 400, wantCode: "INVALID_JSON"},
		{name: "null inside operands rejected", rawBody: `{"operation":"divide","operands":[1,null]}`, wantStatus: 400, wantCode: "INVALID_JSON"},
		{name: "number sent as json string", rawBody: `{"operation":"add","operands":["2",3]}`, wantStatus: 400, wantCode: "INVALID_JSON"},
		{name: "trailing data after json object", rawBody: `{"operation":"add","operands":[1,2]}{"operation":"add","operands":[3,4]}`, wantStatus: 400, wantCode: "INVALID_JSON"},

		// --- UNKNOWN_FIELD ---
		{name: "unknown extra field", rawBody: `{"operation":"add","operands":[1,2],"extra":"nope"}`, wantStatus: 400, wantCode: "UNKNOWN_FIELD"},

		// --- DIVISION_BY_ZERO: positive/0, negative/0, 0/0 ---
		{name: "divide positive by zero", operation: "divide", operands: []float64{5, 0}, wantStatus: 422, wantCode: "DIVISION_BY_ZERO"},
		{name: "divide negative by zero", operation: "divide", operands: []float64{-5, 0}, wantStatus: 422, wantCode: "DIVISION_BY_ZERO"},
		{name: "divide zero by zero", operation: "divide", operands: []float64{0, 0}, wantStatus: 422, wantCode: "DIVISION_BY_ZERO"},

		// --- MATH_DOMAIN_ERROR ---
		{name: "sqrt of negative number", operation: "sqrt", operands: []float64{-4}, wantStatus: 422, wantCode: "MATH_DOMAIN_ERROR"},
		{name: "power negative base fractional exponent", operation: "power", operands: []float64{-8, 0.5}, wantStatus: 422, wantCode: "MATH_DOMAIN_ERROR"},

		// --- RESULT_OUT_OF_RANGE ---
		{name: "power overflow to infinity", operation: "power", operands: []float64{10, 400}, wantStatus: 422, wantCode: "RESULT_OUT_OF_RANGE"},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var body []byte
			if tc.rawBody != "" {
				body = []byte(tc.rawBody)
			} else {
				body = contractBuildBody(tc.operation, tc.operands)
			}

			ct := tc.contentType
			if ct == nil {
				ct = contractJSONContentType()
			}

			rec := contractDo(t, handler, http.MethodPost, "/api/v1/calculate", body, ct)

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}

			if tc.wantStatus == http.StatusOK {
				got := contractDecodeCalculateResponse(t, rec)
				if tc.checkExact {
					if got.Result != tc.wantResult {
						t.Fatalf("result = %v, want exact %v", got.Result, tc.wantResult)
					}
				} else if got.Result != tc.wantResult {
					t.Fatalf("result = %v, want %v", got.Result, tc.wantResult)
				}
				return
			}

			errResp := contractDecodeErrorResponse(t, rec)
			if errResp.Error.Code != tc.wantCode {
				t.Fatalf("error.code = %q, want %q (message=%q)", errResp.Error.Code, tc.wantCode, errResp.Error.Message)
			}
			if errResp.Error.Message == "" {
				t.Fatalf("expected non-empty error message")
			}
		})
	}
}

func TestContract_Calculate_ContentTypeHandling(t *testing.T) {
	handler := contractNewHandler()
	body := contractBuildBody("add", []float64{1, 2})

	t.Run("missing content type rejected", func(t *testing.T) {
		rec := contractDo(t, handler, http.MethodPost, "/api/v1/calculate", body, nil)
		if rec.Code != http.StatusUnsupportedMediaType {
			t.Fatalf("status = %d, want 415; body=%s", rec.Code, rec.Body.String())
		}
		errResp := contractDecodeErrorResponse(t, rec)
		if errResp.Error.Code != "UNSUPPORTED_MEDIA_TYPE" {
			t.Fatalf("error.code = %q, want UNSUPPORTED_MEDIA_TYPE", errResp.Error.Code)
		}
	})

	t.Run("text/plain rejected", func(t *testing.T) {
		ct := "text/plain"
		rec := contractDo(t, handler, http.MethodPost, "/api/v1/calculate", body, &ct)
		if rec.Code != http.StatusUnsupportedMediaType {
			t.Fatalf("status = %d, want 415; body=%s", rec.Code, rec.Body.String())
		}
		errResp := contractDecodeErrorResponse(t, rec)
		if errResp.Error.Code != "UNSUPPORTED_MEDIA_TYPE" {
			t.Fatalf("error.code = %q, want UNSUPPORTED_MEDIA_TYPE", errResp.Error.Code)
		}
	})

	t.Run("malformed content type rejected", func(t *testing.T) {
		ct := ";;;not-a-real-type"
		rec := contractDo(t, handler, http.MethodPost, "/api/v1/calculate", body, &ct)
		if rec.Code != http.StatusUnsupportedMediaType {
			t.Fatalf("status = %d, want 415; body=%s", rec.Code, rec.Body.String())
		}
		errResp := contractDecodeErrorResponse(t, rec)
		if errResp.Error.Code != "UNSUPPORTED_MEDIA_TYPE" {
			t.Fatalf("error.code = %q, want UNSUPPORTED_MEDIA_TYPE", errResp.Error.Code)
		}
	})

	t.Run("application/json accepted", func(t *testing.T) {
		ct := "application/json"
		rec := contractDo(t, handler, http.MethodPost, "/api/v1/calculate", body, &ct)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
	})

	t.Run("application/json with charset accepted", func(t *testing.T) {
		ct := "application/json; charset=utf-8"
		rec := contractDo(t, handler, http.MethodPost, "/api/v1/calculate", body, &ct)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
		got := contractDecodeCalculateResponse(t, rec)
		if got.Result != 3 {
			t.Fatalf("result = %v, want 3", got.Result)
		}
	})

	t.Run("application/json with uppercase and spacing variations accepted", func(t *testing.T) {
		ct := "Application/JSON ; charset=UTF-8"
		rec := contractDo(t, handler, http.MethodPost, "/api/v1/calculate", body, &ct)
		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
		}
	})
}

func TestContract_Calculate_RequestTooLarge(t *testing.T) {
	handler := contractNewHandler()

	// Build an oversized body (> ~1 KiB) using a huge operation name padded
	// with whitespace inside a valid-looking JSON structure, so the only
	// reason for rejection is pure size, not JSON validity. We use operands
	// padding via a long unknown-operation string is not ideal since that
	// could trip other validation first; instead we bloat via many operand
	// entries that are themselves syntactically valid numbers followed by
	// an early reject check for size being prioritized.
	var sb strings.Builder
	sb.WriteString(`{"operation":"add","operands":[1`)
	for sb.Len() < 2048 {
		sb.WriteString(",1")
	}
	sb.WriteString("]}")
	body := []byte(sb.String())

	rec := contractDo(t, handler, http.MethodPost, "/api/v1/calculate", body, contractJSONContentType())
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413; body length=%d", rec.Code, len(body))
	}
	errResp := contractDecodeErrorResponse(t, rec)
	if errResp.Error.Code != "REQUEST_TOO_LARGE" {
		t.Fatalf("error.code = %q, want REQUEST_TOO_LARGE", errResp.Error.Code)
	}
}

func TestContract_Operations(t *testing.T) {
	handler := contractNewHandler()

	rec := contractDo(t, handler, http.MethodGet, "/api/v1/operations", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}

	var out api.OperationsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("failed to decode operations response: %v; body=%s", err, rec.Body.String())
	}

	wantNames := []string{"add", "divide", "multiply", "percentage", "power", "sqrt", "subtract"}
	wantArity := map[string]api.ArityInfo{
		"add":        {Min: 2, Max: 2},
		"subtract":   {Min: 2, Max: 2},
		"multiply":   {Min: 2, Max: 2},
		"divide":     {Min: 2, Max: 2},
		"power":      {Min: 2, Max: 2},
		"sqrt":       {Min: 1, Max: 1},
		"percentage": {Min: 2, Max: 2},
	}

	if len(out.Operations) != len(wantNames) {
		t.Fatalf("operations count = %d, want %d (got %+v)", len(out.Operations), len(wantNames), out.Operations)
	}

	var gotNames []string
	for _, op := range out.Operations {
		gotNames = append(gotNames, op.Name)
		wantA, ok := wantArity[op.Name]
		if !ok {
			t.Fatalf("unexpected operation in response: %q", op.Name)
		}
		if op.Arity != wantA {
			t.Fatalf("operation %q arity = %+v, want %+v", op.Name, op.Arity, wantA)
		}
	}

	if !sortedAlphabetically(gotNames) {
		t.Fatalf("operations not sorted alphabetically by name: %+v", gotNames)
	}

	for _, want := range wantNames {
		found := false
		for _, got := range gotNames {
			if got == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("missing expected operation %q in response: %+v", want, gotNames)
		}
	}

	// Confirm the response is wrapped in an object with an "operations" key,
	// not a bare JSON array.
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatalf("operations response is not a JSON object (expected {\"operations\":[...]}, got array or other): %v; body=%s", err, rec.Body.String())
	}
	if _, ok := raw["operations"]; !ok {
		t.Fatalf(`expected top-level "operations" key, got keys: %+v`, raw)
	}
}

func sortedAlphabetically(names []string) bool {
	for i := 1; i < len(names); i++ {
		if names[i-1] > names[i] {
			return false
		}
	}
	return true
}

func TestContract_Operations_WrongMethod(t *testing.T) {
	handler := contractNewHandler()
	rec := contractDo(t, handler, http.MethodPost, "/api/v1/operations", nil, nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405; body=%s", rec.Code, rec.Body.String())
	}
}

func TestContract_Healthz(t *testing.T) {
	handler := contractNewHandler()
	rec := contractDo(t, handler, http.MethodGet, "/healthz", nil, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
}

func TestContract_NotFound(t *testing.T) {
	handler := contractNewHandler()
	rec := contractDo(t, handler, http.MethodGet, "/this/path/does/not/exist", nil, nil)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404; body=%s", rec.Code, rec.Body.String())
	}
}

func TestContract_MethodNotAllowed(t *testing.T) {
	handler := contractNewHandler()
	rec := contractDo(t, handler, http.MethodGet, "/api/v1/calculate", nil, nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405; body=%s", rec.Code, rec.Body.String())
	}
}
