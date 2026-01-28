package validation

import (
	"fmt"
	"math"

	"github.com/m-szczepanski/gocalc-api/internal/errors"
	"github.com/m-szczepanski/gocalc-api/internal/models"
)

func ValidateMathRequest(req *models.MathRequest) *errors.APIError {
	if req == nil {
		return errors.InvalidInput("request body is required")
	}

	if math.IsNaN(req.A) || math.IsInf(req.A, 0) {
		return errors.ValidationError(
			"invalid a",
			fmt.Sprintf("a must be a valid number, got %v", req.A),
		)
	}

	if math.IsNaN(req.B) || math.IsInf(req.B, 0) {
		return errors.ValidationError(
			"invalid b",
			fmt.Sprintf("b must be a valid number, got %v", req.B),
		)
	}

	return nil
}

func ValidateDivisionRequest(req *models.MathRequest) *errors.APIError {
	if apiErr := ValidateMathRequest(req); apiErr != nil {
		return apiErr
	}

	// Check for division by zero
	// Simple equality check is sufficient because:
	// 1. ValidateMathRequest already checked for Inf and NaN
	// 2. Go's == operator treats -0.0 and +0.0 as equal
	// 3. Very small values (e.g., 1e-308) are valid divisors in floating-point arithmetic
	if req.B == 0 {
		return errors.DivisionByZero()
	}

	return nil
}

func ValidateMethod(method, allowedMethod string) *errors.APIError {
	if method != allowedMethod {
		return errors.MethodNotAllowed(method)
	}
	return nil
}

func ValidateVATRequest(req *models.VATRequest) *errors.APIError {
	if req == nil {
		return errors.InvalidInput("request body is required")
	}

	if math.IsNaN(req.Amount) || math.IsInf(req.Amount, 0) {
		return errors.ValidationError(
			"invalid amount",
			fmt.Sprintf("amount must be a valid number, got %v", req.Amount),
		)
	}

	if req.Amount < 0 {
		return errors.ValidationError(
			"invalid amount",
			"amount cannot be negative",
		)
	}

	if math.IsNaN(req.Rate) || math.IsInf(req.Rate, 0) {
		return errors.ValidationError(
			"invalid rate",
			fmt.Sprintf("rate must be a valid number, got %v", req.Rate),
		)
	}

	if req.Rate < 0 {
		return errors.ValidationError(
			"invalid rate",
			"rate cannot be negative",
		)
	}

	return nil
}

func ValidateCompoundInterestRequest(req *models.CompoundInterestRequest) *errors.APIError {
	if req == nil {
		return errors.InvalidInput("request body is required")
	}

	if math.IsNaN(req.Principal) || math.IsInf(req.Principal, 0) {
		return errors.ValidationError(
			"invalid principal",
			fmt.Sprintf("principal must be a valid number, got %v", req.Principal),
		)
	}

	if req.Principal < 0 {
		return errors.ValidationError(
			"invalid principal",
			"principal cannot be negative",
		)
	}

	if math.IsNaN(req.Rate) || math.IsInf(req.Rate, 0) {
		return errors.ValidationError(
			"invalid rate",
			fmt.Sprintf("rate must be a valid number, got %v", req.Rate),
		)
	}

	if req.Rate < 0 {
		return errors.ValidationError(
			"invalid rate",
			"rate cannot be negative",
		)
	}

	if math.IsNaN(req.Time) || math.IsInf(req.Time, 0) {
		return errors.ValidationError(
			"invalid time",
			fmt.Sprintf("time must be a valid number, got %v", req.Time),
		)
	}

	if req.Time < 0 {
		return errors.ValidationError(
			"invalid time",
			"time cannot be negative",
		)
	}

	if req.CompoundFrequency <= 0 {
		return errors.ValidationError(
			"invalid compound_frequency",
			"compound_frequency must be positive",
		)
	}

	return nil
}

func ValidateLoanPaymentRequest(req *models.LoanPaymentRequest) *errors.APIError {
	if req == nil {
		return errors.InvalidInput("request body is required")
	}

	if math.IsNaN(req.Principal) || math.IsInf(req.Principal, 0) {
		return errors.ValidationError(
			"invalid principal",
			fmt.Sprintf("principal must be a valid number, got %v", req.Principal),
		)
	}

	if req.Principal < 0 {
		return errors.ValidationError(
			"invalid principal",
			"principal cannot be negative",
		)
	}

	if math.IsNaN(req.AnnualRate) || math.IsInf(req.AnnualRate, 0) {
		return errors.ValidationError(
			"invalid annual_rate",
			fmt.Sprintf("annual_rate must be a valid number, got %v", req.AnnualRate),
		)
	}

	if req.AnnualRate < 0 {
		return errors.ValidationError(
			"invalid annual_rate",
			"annual_rate cannot be negative",
		)
	}

	if math.IsNaN(req.Years) || math.IsInf(req.Years, 0) {
		return errors.ValidationError(
			"invalid years",
			fmt.Sprintf("years must be a valid number, got %v", req.Years),
		)
	}

	if req.Years <= 0 {
		return errors.ValidationError(
			"invalid years",
			"years must be positive",
		)
	}

	if req.PaymentsPerYear <= 0 {
		return errors.ValidationError(
			"invalid payments_per_year",
			"payments_per_year must be positive",
		)
	}

	return nil
}
