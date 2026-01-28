package calculations

import (
	"fmt"
	"math"
	"strings"
)

var globalRegistry = NewUnitRegistry()

type BMICategory string

const (
	BMIUnderweight   BMICategory = "underweight"
	BMINormalWeight  BMICategory = "normal"
	BMIOverweight    BMICategory = "overweight"
	BMIObesityClass1 BMICategory = "obesity_class_1"
	BMIObesityClass2 BMICategory = "obesity_class_2"
	BMIObesityClass3 BMICategory = "obesity_class_3"
)

type BMIResult struct {
	BMI      float64
	Category BMICategory
}

func CalculateBMI(weight float64, weightUnit string, height float64, heightUnit string) (*BMIResult, error) {
	// Validate inputs - check for NaN/Inf first
	if math.IsNaN(weight) || math.IsInf(weight, 0) {
		return nil, fmt.Errorf("weight must be a valid number")
	}
	if math.IsNaN(height) || math.IsInf(height, 0) {
		return nil, fmt.Errorf("height must be a valid number")
	}
	if weight <= 0 {
		return nil, fmt.Errorf("weight must be greater than zero")
	}
	if height <= 0 {
		return nil, fmt.Errorf("height must be greater than zero")
	}

	// Convert weight to kg
	weightKg, err := globalRegistry.ConvertToBaseUnit(weight, UnitTypeWeight, weightUnit)
	if err != nil {
		return nil, fmt.Errorf("invalid weight unit: %w", err)
	}

	// Convert height to m
	heightM, err := globalRegistry.ConvertToBaseUnit(height, UnitTypeHeight, heightUnit)
	if err != nil {
		return nil, fmt.Errorf("invalid height unit: %w", err)
	}

	// Calculate BMI
	bmi := weightKg / (heightM * heightM)

	// Determine category using unrounded BMI to avoid boundary misclassification
	category := categorizeBMI(bmi)

	// Round BMI for display
	roundedBMI := math.Round(bmi*100) / 100

	return &BMIResult{
		BMI:      roundedBMI,
		Category: category,
	}, nil
}

func categorizeBMI(bmi float64) BMICategory {
	switch {
	case bmi < 18.5:
		return BMIUnderweight
	case bmi < 25:
		return BMINormalWeight
	case bmi < 30:
		return BMIOverweight
	case bmi < 35:
		return BMIObesityClass1
	case bmi < 40:
		return BMIObesityClass2
	default:
		return BMIObesityClass3
	}
}

// ConvertUnit performs unit conversion for various unit types
// Supports: weight, height, temperature, distance, volume
func ConvertUnit(value float64, fromUnit, toUnit string, unitType string) (float64, error) {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return 0, fmt.Errorf("value must be a valid number")
	}

	// Normalize and validate unit type
	if !IsValidUnitType(unitType) {
		return 0, fmt.Errorf("invalid unit type: %s (valid types: weight, height, temperature, distance, volume)", unitType)
	}

	// Use the normalized version for UnitType conversion
	normalizedUnitType := strings.ToLower(strings.TrimSpace(unitType))
	ut := UnitType(normalizedUnitType)

	if !globalRegistry.IsValidUnit(ut, fromUnit) {
		validUnits := globalRegistry.GetValidUnits(ut)
		return 0, fmt.Errorf("invalid source unit '%s' for type '%s' (valid units: %v)", fromUnit, unitType, validUnits)
	}
	if !globalRegistry.IsValidUnit(ut, toUnit) {
		validUnits := globalRegistry.GetValidUnits(ut)
		return 0, fmt.Errorf("invalid target unit '%s' for type '%s' (valid units: %v)", toUnit, unitType, validUnits)
	}

	result, err := globalRegistry.Convert(value, ut, fromUnit, toUnit)
	if err != nil {
		return 0, fmt.Errorf("conversion failed: %w", err)
	}

	result = math.Round(result*1000000) / 1000000

	return result, nil
}
