package validation

import (
	"math"
	"testing"

	"github.com/m-szczepanski/gocalc-api/internal/errors"
	"github.com/m-szczepanski/gocalc-api/internal/models"
)

func TestValidateMathRequest(t *testing.T) {
	tests := []struct {
		name         string
		req          *models.MathRequest
		expectError  bool
		expectedCode string
	}{
		{
			name:        "valid request",
			req:         &models.MathRequest{A: 5, B: 3},
			expectError: false,
		},
		{
			name:        "zero values",
			req:         &models.MathRequest{A: 0, B: 0},
			expectError: false,
		},
		{
			name:        "negative values",
			req:         &models.MathRequest{A: -5, B: -3},
			expectError: false,
		},
		{
			name:        "decimal values",
			req:         &models.MathRequest{A: 1.5, B: 2.5},
			expectError: false,
		},
		{
			name:         "nil request",
			req:          nil,
			expectError:  true,
			expectedCode: errors.ErrCodeInvalidInput,
		},
		{
			name:         "NaN in A",
			req:          &models.MathRequest{A: math.NaN(), B: 3},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name:         "NaN in B",
			req:          &models.MathRequest{A: 5, B: math.NaN()},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name:         "Inf in A",
			req:          &models.MathRequest{A: math.Inf(1), B: 3},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name:         "Inf in B",
			req:          &models.MathRequest{A: 5, B: math.Inf(-1)},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMathRequest(tt.req)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				} else if err.Code != tt.expectedCode {
					t.Errorf("error code = %s, want %s", err.Code, tt.expectedCode)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateDivisionRequest(t *testing.T) {
	tests := []struct {
		name         string
		req          *models.MathRequest
		expectError  bool
		expectedCode string
	}{
		{
			name:        "valid division",
			req:         &models.MathRequest{A: 6, B: 2},
			expectError: false,
		},
		{
			name:         "division by zero",
			req:          &models.MathRequest{A: 5, B: 0},
			expectError:  true,
			expectedCode: errors.ErrCodeDivisionByZero,
		},
		{
			name:        "non-zero divisor",
			req:         &models.MathRequest{A: 10, B: 0.5},
			expectError: false,
		},
		{
			name:         "NaN divisor",
			req:          &models.MathRequest{A: 5, B: math.NaN()},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDivisionRequest(tt.req)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				} else if err.Code != tt.expectedCode {
					t.Errorf("error code = %s, want %s", err.Code, tt.expectedCode)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

func TestValidateMethod(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		allowedMethod string
		expectError   bool
	}{
		{
			name:          "matching method",
			method:        "POST",
			allowedMethod: "POST",
			expectError:   false,
		},
		{
			name:          "non-matching method",
			method:        "GET",
			allowedMethod: "POST",
			expectError:   true,
		},
		{
			name:          "PUT not allowed",
			method:        "PUT",
			allowedMethod: "POST",
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMethod(tt.method, tt.allowedMethod)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				} else if err.Code != errors.ErrCodeMethodNotAllowed {
					t.Errorf("error code = %s, want %s", err.Code, errors.ErrCodeMethodNotAllowed)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}
