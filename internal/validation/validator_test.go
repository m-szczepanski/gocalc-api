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
func TestValidateVATRequest(t *testing.T) {
	tests := []struct {
		name         string
		req          *models.VATRequest
		expectError  bool
		expectedCode string
	}{
		{
			name:        "valid VAT request exclusive",
			req:         &models.VATRequest{Amount: 100, Rate: 23, Inclusive: false},
			expectError: false,
		},
		{
			name:        "valid VAT request inclusive",
			req:         &models.VATRequest{Amount: 123, Rate: 23, Inclusive: true},
			expectError: false,
		},
		{
			name:        "zero amount and rate",
			req:         &models.VATRequest{Amount: 0, Rate: 0, Inclusive: false},
			expectError: false,
		},
		{
			name:         "nil request",
			req:          nil,
			expectError:  true,
			expectedCode: errors.ErrCodeInvalidInput,
		},
		{
			name:         "negative amount",
			req:          &models.VATRequest{Amount: -100, Rate: 23, Inclusive: false},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name:         "negative rate",
			req:          &models.VATRequest{Amount: 100, Rate: -5, Inclusive: false},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name:         "NaN amount",
			req:          &models.VATRequest{Amount: math.NaN(), Rate: 23, Inclusive: false},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name:         "Inf amount",
			req:          &models.VATRequest{Amount: math.Inf(1), Rate: 23, Inclusive: false},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name:         "NaN rate",
			req:          &models.VATRequest{Amount: 100, Rate: math.NaN(), Inclusive: false},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name:         "Inf rate",
			req:          &models.VATRequest{Amount: 100, Rate: math.Inf(1), Inclusive: false},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVATRequest(tt.req)

			if (err != nil) != tt.expectError {
				t.Errorf("ValidateVATRequest() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if tt.expectError && err.Code != tt.expectedCode {
				t.Errorf("ValidateVATRequest() error code = %v, expected %v", err.Code, tt.expectedCode)
			}
		})
	}
}

func TestValidateCompoundInterestRequest(t *testing.T) {
	tests := []struct {
		name         string
		req          *models.CompoundInterestRequest
		expectError  bool
		expectedCode string
	}{
		{
			name: "valid request monthly",
			req: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              5,
				Time:              10,
				CompoundFrequency: 12,
			},
			expectError: false,
		},
		{
			name: "valid request annual",
			req: &models.CompoundInterestRequest{
				Principal:         5000,
				Rate:              3,
				Time:              5,
				CompoundFrequency: 1,
			},
			expectError: false,
		},
		{
			name: "zero rate and time",
			req: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              0,
				Time:              0,
				CompoundFrequency: 12,
			},
			expectError: false,
		},
		{
			name:         "nil request",
			req:          nil,
			expectError:  true,
			expectedCode: errors.ErrCodeInvalidInput,
		},
		{
			name: "negative principal",
			req: &models.CompoundInterestRequest{
				Principal:         -1000,
				Rate:              5,
				Time:              10,
				CompoundFrequency: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "negative rate",
			req: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              -5,
				Time:              10,
				CompoundFrequency: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "negative time",
			req: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              5,
				Time:              -10,
				CompoundFrequency: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "zero compound frequency",
			req: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              5,
				Time:              10,
				CompoundFrequency: 0,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "negative compound frequency",
			req: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              5,
				Time:              10,
				CompoundFrequency: -12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "NaN principal",
			req: &models.CompoundInterestRequest{
				Principal:         math.NaN(),
				Rate:              5,
				Time:              10,
				CompoundFrequency: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "Inf rate",
			req: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              math.Inf(1),
				Time:              10,
				CompoundFrequency: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "NaN time",
			req: &models.CompoundInterestRequest{
				Principal:         1000,
				Rate:              5,
				Time:              math.NaN(),
				CompoundFrequency: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCompoundInterestRequest(tt.req)

			if (err != nil) != tt.expectError {
				t.Errorf("ValidateCompoundInterestRequest() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if tt.expectError && err.Code != tt.expectedCode {
				t.Errorf("ValidateCompoundInterestRequest() error code = %v, expected %v", err.Code, tt.expectedCode)
			}
		})
	}
}

func TestValidateLoanPaymentRequest(t *testing.T) {
	tests := []struct {
		name         string
		req          *models.LoanPaymentRequest
		expectError  bool
		expectedCode string
	}{
		{
			name: "valid monthly payments",
			req: &models.LoanPaymentRequest{
				Principal:       300000,
				AnnualRate:      4.5,
				Years:           30,
				PaymentsPerYear: 12,
			},
			expectError: false,
		},
		{
			name: "valid quarterly payments",
			req: &models.LoanPaymentRequest{
				Principal:       50000,
				AnnualRate:      5,
				Years:           10,
				PaymentsPerYear: 4,
			},
			expectError: false,
		},
		{
			name: "zero interest rate",
			req: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      0,
				Years:           5,
				PaymentsPerYear: 12,
			},
			expectError: false,
		},
		{
			name:         "nil request",
			req:          nil,
			expectError:  true,
			expectedCode: errors.ErrCodeInvalidInput,
		},
		{
			name: "negative principal",
			req: &models.LoanPaymentRequest{
				Principal:       -10000,
				AnnualRate:      5,
				Years:           5,
				PaymentsPerYear: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "negative annual rate",
			req: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      -5,
				Years:           5,
				PaymentsPerYear: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "zero years",
			req: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      5,
				Years:           0,
				PaymentsPerYear: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "negative years",
			req: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      5,
				Years:           -5,
				PaymentsPerYear: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "zero payments per year",
			req: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      5,
				Years:           5,
				PaymentsPerYear: 0,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "negative payments per year",
			req: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      5,
				Years:           5,
				PaymentsPerYear: -12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "NaN principal",
			req: &models.LoanPaymentRequest{
				Principal:       math.NaN(),
				AnnualRate:      5,
				Years:           5,
				PaymentsPerYear: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "Inf annual rate",
			req: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      math.Inf(1),
				Years:           5,
				PaymentsPerYear: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
		{
			name: "NaN years",
			req: &models.LoanPaymentRequest{
				Principal:       10000,
				AnnualRate:      5,
				Years:           math.NaN(),
				PaymentsPerYear: 12,
			},
			expectError:  true,
			expectedCode: errors.ErrCodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateLoanPaymentRequest(tt.req)

			if (err != nil) != tt.expectError {
				t.Errorf("ValidateLoanPaymentRequest() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if tt.expectError && err.Code != tt.expectedCode {
				t.Errorf("ValidateLoanPaymentRequest() error code = %v, expected %v", err.Code, tt.expectedCode)
			}
		})
	}
}
