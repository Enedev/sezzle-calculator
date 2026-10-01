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

func newTestRouter() http.Handler {
	return api.NewRouter(calculator.NewDefaultRegistry())
}

func doRequest(t *testing.T, handler http.Handler, method, path, contentType, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) api.ErrorResponse {
	t.Helper()
	var got api.ErrorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("response is not a valid ErrorResponse: %v (body=%s)", err, rec.Body.String())
	}
	return got
}

func TestCalculate_AllOperations(t *testing.T) {
	handler := newTestRouter()
	tests := []struct {
		name string
		body string
		want float64
	}{
		{"add", `{"operation":"add","operands":[2,3]}`, 5},
		{"subtract", `{"operation":"subtract","operands":[5,3]}`, 2},
		{"multiply", `{"operation":"multiply","operands":[4,5]}`, 20},
		{"divide", `{"operation":"divide","operands":[10,2]}`, 5},
		{"power", `{"operation":"power","operands":[2,10]}`, 1024},
		{"sqrt", `{"operation":"sqrt","operands":[9]}`, 3},
		{"percentage", `{"operation":"percentage","operands":[20,50]}`, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(t, handler, http.MethodPost, "/api/v1/calculate", "application/json", tt.body)
			if rec.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
			}
			var got api.CalculateResponse
			if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
				t.Fatalf("invalid response body: %v", err)
			}
			if got.Result != tt.want {
				t.Fatalf("result = %v, want %v", got.Result, tt.want)
			}
		})
	}
}

func TestCalculate_ContentTypeWithCharsetAccepted(t *testing.T) {
	handler := newTestRouter()
	rec := doRequest(t, handler, http.MethodPost, "/api/v1/calculate",
		"application/json; charset=utf-8", `{"operation":"add","operands":[1,2]}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestCalculate_ErrorCases(t *testing.T) {
	handler := newTestRouter()
	tests := []struct {
		name        string
		contentType string
		body        string
		wantStatus  int
		wantCode    api.ErrorCode
	}{
		{
			"division by zero",
			"application/json", `{"operation":"divide","operands":[1,0]}`,
			http.StatusUnprocessableEntity, api.ErrCodeDivisionByZero,
		},
		{
			"sqrt of negative",
			"application/json", `{"operation":"sqrt","operands":[-4]}`,
			http.StatusUnprocessableEntity, api.ErrCodeMathDomainError,
		},
		{
			"fractional power of negative base is a math domain error",
			"application/json", `{"operation":"power","operands":[-8,0.5]}`,
			http.StatusUnprocessableEntity, api.ErrCodeMathDomainError,
		},
		{
			"overflow is result out of range",
			"application/json", `{"operation":"power","operands":[10,400]}`,
			http.StatusUnprocessableEntity, api.ErrCodeResultOutOfRange,
		},
		{
			"wrong operand count",
			"application/json", `{"operation":"add","operands":[1]}`,
			http.StatusBadRequest, api.ErrCodeInvalidOperandCount,
		},
		{
			"operands missing maps to invalid operand count, not missing field",
			"application/json", `{"operation":"add"}`,
			http.StatusBadRequest, api.ErrCodeInvalidOperandCount,
		},
		{
			"operands empty array maps to invalid operand count, not missing field",
			"application/json", `{"operation":"add","operands":[]}`,
			http.StatusBadRequest, api.ErrCodeInvalidOperandCount,
		},
		{
			"unknown operation",
			"application/json", `{"operation":"modulo","operands":[1,2]}`,
			http.StatusBadRequest, api.ErrCodeUnknownOperation,
		},
		{
			"malformed json",
			"application/json", `{"operation":"add","operands":[1,2]`,
			http.StatusBadRequest, api.ErrCodeInvalidJSON,
		},
		{
			"unknown field",
			"application/json", `{"operation":"add","operands":[1,2],"extra":true}`,
			http.StatusBadRequest, api.ErrCodeUnknownField,
		},
		{
			"missing operation field",
			"application/json", `{"operands":[1,2]}`,
			http.StatusBadRequest, api.ErrCodeMissingField,
		},
		{
			"empty operation field",
			"application/json", `{"operation":"","operands":[1,2]}`,
			http.StatusBadRequest, api.ErrCodeMissingField,
		},
		{
			"null operand is rejected, not silently zeroed",
			"application/json", `{"operation":"divide","operands":[1,null]}`,
			http.StatusBadRequest, api.ErrCodeInvalidJSON,
		},
		{
			"operand sent as a json string",
			"application/json", `{"operation":"add","operands":["2",3]}`,
			http.StatusBadRequest, api.ErrCodeInvalidJSON,
		},
		{
			"trailing data after the json object",
			"application/json", `{"operation":"add","operands":[1,2]}{"operation":"add","operands":[1,2]}`,
			http.StatusBadRequest, api.ErrCodeInvalidJSON,
		},
		{
			"wrong content type",
			"text/plain", `{"operation":"add","operands":[1,2]}`,
			http.StatusUnsupportedMediaType, api.ErrCodeUnsupportedMediaType,
		},
		{
			"missing content type",
			"", `{"operation":"add","operands":[1,2]}`,
			http.StatusUnsupportedMediaType, api.ErrCodeUnsupportedMediaType,
		},
		{
			"malformed content type",
			"application/json; charset", `{"operation":"add","operands":[1,2]}`,
			http.StatusUnsupportedMediaType, api.ErrCodeUnsupportedMediaType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := doRequest(t, handler, http.MethodPost, "/api/v1/calculate", tt.contentType, tt.body)
			if rec.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d (body = %s)", rec.Code, tt.wantStatus, rec.Body.String())
			}
			got := decodeError(t, rec)
			if got.Error.Code != tt.wantCode {
				t.Fatalf("error code = %q, want %q", got.Error.Code, tt.wantCode)
			}
		})
	}
}

func TestCalculate_OversizedBody(t *testing.T) {
	handler := newTestRouter()
	// 1 KiB cap: pad operands with a huge array of zeros to exceed it.
	var buf bytes.Buffer
	buf.WriteString(`{"operation":"add","operands":[`)
	for i := 0; i < 500; i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString("0")
	}
	buf.WriteString(`]}`)

	rec := doRequest(t, handler, http.MethodPost, "/api/v1/calculate", "application/json", buf.String())
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want %d (body = %s)", rec.Code, http.StatusRequestEntityTooLarge, rec.Body.String())
	}
	got := decodeError(t, rec)
	if got.Error.Code != api.ErrCodeRequestTooLarge {
		t.Fatalf("error code = %q, want %q", got.Error.Code, api.ErrCodeRequestTooLarge)
	}
}

func TestOperations_ListsAllSevenSorted(t *testing.T) {
	handler := newTestRouter()
	rec := doRequest(t, handler, http.MethodGet, "/api/v1/operations", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	var got api.OperationsResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid response body: %v", err)
	}
	want := []string{"add", "divide", "multiply", "percentage", "power", "sqrt", "subtract"}
	if len(got.Operations) != len(want) {
		t.Fatalf("got %d operations, want %d", len(got.Operations), len(want))
	}
	for i, name := range want {
		if got.Operations[i].Name != name {
			t.Fatalf("operations[%d].Name = %q, want %q", i, got.Operations[i].Name, name)
		}
	}
}

func TestHealthz(t *testing.T) {
	handler := newTestRouter()
	rec := doRequest(t, handler, http.MethodGet, "/healthz", "", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}

func TestUnknownRoute_Returns404(t *testing.T) {
	handler := newTestRouter()
	rec := doRequest(t, handler, http.MethodGet, "/no-such-route", "", "")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestWrongMethod_Returns405(t *testing.T) {
	handler := newTestRouter()
	rec := doRequest(t, handler, http.MethodGet, "/api/v1/calculate", "", "")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}
