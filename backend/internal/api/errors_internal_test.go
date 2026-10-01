package api

import (
	"errors"
	"net/http"
	"testing"
)

// Exercises the default branch of domainErrorToAPI, which real callers
// never hit today (Registry.Calculate only ever returns the five known
// calculator sentinel errors) but which exists as a safety net against a
// future operation returning something unmapped.
func TestDomainErrorToAPI_UnmappedErrorIsInternal(t *testing.T) {
	code, status, _ := domainErrorToAPI(errors.New("some future unmapped error"))
	if code != ErrCodeInternal {
		t.Fatalf("code = %q, want %q", code, ErrCodeInternal)
	}
	if status != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", status, http.StatusInternalServerError)
	}
}
