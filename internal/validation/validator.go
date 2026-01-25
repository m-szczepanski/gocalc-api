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
