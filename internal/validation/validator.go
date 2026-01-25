package validation

import (
	"fmt"
	"math"

	"github.com/m-szczepanski/gocalc-api/internal/errors"
	"github.com/m-szczepanski/gocalc-api/internal/models"
)

// ValidateMathRequest validates a MathRequest for arithmetic operations.
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

// ValidateDivisionRequest validates a MathRequest for division specifically.
// It includes all standard validations plus division-specific checks.
func ValidateDivisionRequest(req *models.MathRequest) *errors.APIError {
	// Run standard validations first
	if apiErr := ValidateMathRequest(req); apiErr != nil {
		return apiErr
	}

	// Check for division by zero
	if req.B == 0 {
		return errors.DivisionByZero()
	}

	return nil
}

// ValidateMethod validates that the HTTP method is allowed.
func ValidateMethod(method, allowedMethod string) *errors.APIError {
	if method != allowedMethod {
		return errors.MethodNotAllowed(method)
	}
	return nil
}
