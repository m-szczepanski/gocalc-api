package errors

import (
	"fmt"
	"net/http"
)

// Error codes used in API responses
const (
	ErrCodeInvalidInput     = "INVALID_INPUT"
	ErrCodeValidationError  = "VALIDATION_ERROR"
	ErrCodeDivisionByZero   = "DIVISION_BY_ZERO"
	ErrCodeMethodNotAllowed = "METHOD_NOT_ALLOWED"
	ErrCodeInternalError    = "INTERNAL_ERROR"
)

// APIError represents a structured error returned by the API.
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	err     error  // underlying error for logging/debugging
}

// NewAPIError creates a new APIError with the given code and message.
func NewAPIError(code, message string) *APIError {
	return &APIError{
		Code:    code,
		Message: message,
	}
}

// WithDetails adds details to the error and returns itself for chaining.
func (e *APIError) WithDetails(details string) *APIError {
	e.Details = details
	return e
}

// WithError wraps an underlying error for debugging purposes.
func (e *APIError) WithError(err error) *APIError {
	e.err = err
	return e
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// HTTPStatus returns the appropriate HTTP status code for this error.
func (e *APIError) HTTPStatus() int {
	switch e.Code {
	case ErrCodeValidationError, ErrCodeInvalidInput:
		return http.StatusBadRequest
	case ErrCodeDivisionByZero:
		return http.StatusBadRequest
	case ErrCodeMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case ErrCodeInternalError:
		return http.StatusInternalServerError
	default:
		return http.StatusBadRequest
	}
}

// Common error constructors for convenience

// InvalidInput returns a validation error for invalid input.
func InvalidInput(message string) *APIError {
	return NewAPIError(ErrCodeInvalidInput, message)
}

// ValidationError returns a validation error with details.
func ValidationError(message, details string) *APIError {
	return NewAPIError(ErrCodeValidationError, message).WithDetails(details)
}

// DivisionByZero returns a division by zero error.
func DivisionByZero() *APIError {
	return NewAPIError(ErrCodeDivisionByZero, "division by zero is not allowed")
}

// MethodNotAllowed returns a method not allowed error.
func MethodNotAllowed(method string) *APIError {
	return NewAPIError(ErrCodeMethodNotAllowed, fmt.Sprintf("method %s not allowed", method))
}

// InternalError returns an internal server error.
func InternalError(message string) *APIError {
	return NewAPIError(ErrCodeInternalError, message)
}
