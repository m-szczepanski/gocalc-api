package errors

import (
	"net/http"
	"testing"
)

func TestNewAPIError(t *testing.T) {
	err := NewAPIError("TEST_CODE", "test message")

	if err.Code != "TEST_CODE" {
		t.Errorf("Code = %s, want TEST_CODE", err.Code)
	}
	if err.Message != "test message" {
		t.Errorf("Message = %s, want test message", err.Message)
	}
}

func TestAPIErrorWithDetails(t *testing.T) {
	err := NewAPIError("TEST_CODE", "message").WithDetails("error details")

	if err.Details != "error details" {
		t.Errorf("Details = %s, want error details", err.Details)
	}
}

func TestAPIErrorWithError(t *testing.T) {
	originalErr := NewAPIError("ORIGINAL", "original error")
	err := NewAPIError("TEST_CODE", "message").WithError(originalErr)

	if err.err != originalErr {
		t.Errorf("underlying error not set correctly")
	}
}

func TestAPIErrorError(t *testing.T) {
	tests := []struct {
		name string
		err  *APIError
		want string
	}{
		{
			name: "without underlying error",
			err:  NewAPIError("CODE", "message"),
			want: "[CODE] message",
		},
		{
			name: "with underlying error",
			err:  NewAPIError("CODE", "message").WithError(NewAPIError("INNER", "inner")),
			want: "[CODE] message:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if tt.want != "" && got != tt.want {
				if !contains(got, tt.want) {
					t.Errorf("Error() = %s, want to contain %s", got, tt.want)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	for i := range s {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		code       string
		wantStatus int
	}{
		{ErrCodeValidationError, http.StatusBadRequest},
		{ErrCodeInvalidInput, http.StatusBadRequest},
		{ErrCodeDivisionByZero, http.StatusBadRequest},
		{ErrCodeMethodNotAllowed, http.StatusMethodNotAllowed},
		{ErrCodeInternalError, http.StatusInternalServerError},
		{"UNKNOWN_CODE", http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.code, func(t *testing.T) {
			err := NewAPIError(tt.code, "message")
			if status := err.HTTPStatus(); status != tt.wantStatus {
				t.Errorf("HTTPStatus() = %d, want %d", status, tt.wantStatus)
			}
		})
	}
}

func TestErrorConstructors(t *testing.T) {
	tests := []struct {
		name     string
		errFunc  func() *APIError
		wantCode string
	}{
		{
			name:     "InvalidInput",
			errFunc:  func() *APIError { return InvalidInput("test") },
			wantCode: ErrCodeInvalidInput,
		},
		{
			name:     "DivisionByZero",
			errFunc:  func() *APIError { return DivisionByZero() },
			wantCode: ErrCodeDivisionByZero,
		},
		{
			name:     "MethodNotAllowed",
			errFunc:  func() *APIError { return MethodNotAllowed("POST") },
			wantCode: ErrCodeMethodNotAllowed,
		},
		{
			name:     "InternalError",
			errFunc:  func() *APIError { return InternalError("test") },
			wantCode: ErrCodeInternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc()
			if err.Code != tt.wantCode {
				t.Errorf("Code = %s, want %s", err.Code, tt.wantCode)
			}
		})
	}
}
