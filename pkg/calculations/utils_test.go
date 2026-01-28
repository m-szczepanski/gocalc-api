package calculations

import (
	"math"
	"testing"
)

func TestCalculateBMI(t *testing.T) {
	tests := []struct {
		name         string
		weight       float64
		weightUnit   string
		height       float64
		heightUnit   string
		expectedBMI  float64
		expectedCat  BMICategory
		expectError  bool
		errorMessage string
	}{
		{
			name:        "normal BMI - kg and m",
			weight:      70,
			weightUnit:  "kg",
			height:      1.75,
			heightUnit:  "m",
			expectedBMI: 22.86,
			expectedCat: BMINormalWeight,
			expectError: false,
		},
		{
			name:        "normal BMI - lb and ft",
			weight:      154,
			weightUnit:  "lb",
			height:      5.74,
			heightUnit:  "ft",
			expectedBMI: 22.82,
			expectedCat: BMINormalWeight,
			expectError: false,
		},
		{
			name:        "underweight - kg and cm",
			weight:      50,
			weightUnit:  "kg",
			height:      175,
			heightUnit:  "cm",
			expectedBMI: 16.33,
			expectedCat: BMIUnderweight,
			expectError: false,
		},
		{
			name:        "overweight - kg and m",
			weight:      85,
			weightUnit:  "kg",
			height:      1.75,
			heightUnit:  "m",
			expectedBMI: 27.76,
			expectedCat: BMIOverweight,
			expectError: false,
		},
		{
			name:        "obesity class 1 - kg and m",
			weight:      95,
			weightUnit:  "kg",
			height:      1.75,
			heightUnit:  "m",
			expectedBMI: 31.02,
			expectedCat: BMIObesityClass1,
			expectError: false,
		},
		{
			name:        "obesity class 2 - kg and m",
			weight:      110,
			weightUnit:  "kg",
			height:      1.75,
			heightUnit:  "m",
			expectedBMI: 35.92,
			expectedCat: BMIObesityClass2,
			expectError: false,
		},
		{
			name:        "obesity class 3 - kg and m",
			weight:      130,
			weightUnit:  "kg",
			height:      1.75,
			heightUnit:  "m",
			expectedBMI: 42.45,
			expectedCat: BMIObesityClass3,
			expectError: false,
		},
		{
			name:         "zero weight",
			weight:       0,
			weightUnit:   "kg",
			height:       1.75,
			heightUnit:   "m",
			expectError:  true,
			errorMessage: "weight must be greater than zero",
		},
		{
			name:         "negative weight",
			weight:       -70,
			weightUnit:   "kg",
			height:       1.75,
			heightUnit:   "m",
			expectError:  true,
			errorMessage: "weight must be greater than zero",
		},
		{
			name:         "zero height",
			weight:       70,
			weightUnit:   "kg",
			height:       0,
			heightUnit:   "m",
			expectError:  true,
			errorMessage: "height must be greater than zero",
		},
		{
			name:         "negative height",
			weight:       70,
			weightUnit:   "kg",
			height:       -1.75,
			heightUnit:   "m",
			expectError:  true,
			errorMessage: "height must be greater than zero",
		},
		{
			name:         "NaN weight",
			weight:       math.NaN(),
			weightUnit:   "kg",
			height:       1.75,
			heightUnit:   "m",
			expectError:  true,
			errorMessage: "weight must be a valid number",
		},
		{
			name:         "Inf height",
			weight:       70,
			weightUnit:   "kg",
			height:       math.Inf(1),
			heightUnit:   "m",
			expectError:  true,
			errorMessage: "height must be a valid number",
		},
		{
			name:         "invalid weight unit",
			weight:       70,
			weightUnit:   "invalid",
			height:       1.75,
			heightUnit:   "m",
			expectError:  true,
			errorMessage: "invalid weight unit",
		},
		{
			name:         "invalid height unit",
			weight:       70,
			weightUnit:   "kg",
			height:       175,
			heightUnit:   "invalid",
			expectError:  true,
			errorMessage: "invalid height unit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateBMI(tt.weight, tt.weightUnit, tt.height, tt.heightUnit)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Errorf("expected result but got nil")
				return
			}

			if math.Abs(result.BMI-tt.expectedBMI) > 0.01 {
				t.Errorf("expected BMI %.2f, got %.2f", tt.expectedBMI, result.BMI)
			}

			if result.Category != tt.expectedCat {
				t.Errorf("expected category %s, got %s", tt.expectedCat, result.Category)
			}
		})
	}
}

func TestCategorizeBMI(t *testing.T) {
	tests := []struct {
		name     string
		bmi      float64
		expected BMICategory
	}{
		{"extreme underweight", 15.0, BMIUnderweight},
		{"borderline underweight", 18.4, BMIUnderweight},
		{"borderline normal", 18.5, BMINormalWeight},
		{"normal weight", 22.0, BMINormalWeight},
		{"borderline overweight", 24.9, BMINormalWeight},
		{"overweight", 25.0, BMIOverweight},
		{"high overweight", 29.9, BMIOverweight},
		{"obesity class 1 start", 30.0, BMIObesityClass1},
		{"obesity class 1", 32.0, BMIObesityClass1},
		{"obesity class 2 start", 35.0, BMIObesityClass2},
		{"obesity class 2", 37.0, BMIObesityClass2},
		{"obesity class 3 start", 40.0, BMIObesityClass3},
		{"extreme obesity", 50.0, BMIObesityClass3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			category := categorizeBMI(tt.bmi)
			if category != tt.expected {
				t.Errorf("expected category %s for BMI %.1f, got %s", tt.expected, tt.bmi, category)
			}
		})
	}
}

func TestConvertUnit(t *testing.T) {
	tests := []struct {
		name         string
		value        float64
		fromUnit     string
		toUnit       string
		unitType     string
		expected     float64
		expectError  bool
		errorMessage string
	}{
		{
			name:     "kg to lb",
			value:    10,
			fromUnit: "kg",
			toUnit:   "lb",
			unitType: "weight",
			expected: 22.046226,
		},
		{
			name:     "lb to kg",
			value:    22.046226,
			fromUnit: "lb",
			toUnit:   "kg",
			unitType: "weight",
			expected: 10,
		},
		{
			name:     "g to kg",
			value:    1000,
			fromUnit: "g",
			toUnit:   "kg",
			unitType: "weight",
			expected: 1,
		},
		{
			name:     "oz to g",
			value:    1,
			fromUnit: "oz",
			toUnit:   "g",
			unitType: "weight",
			expected: 28.349523,
		},
		{
			name:     "m to cm",
			value:    1.75,
			fromUnit: "m",
			toUnit:   "cm",
			unitType: "height",
			expected: 175,
		},
		{
			name:     "ft to m",
			value:    5.74,
			fromUnit: "ft",
			toUnit:   "m",
			unitType: "height",
			expected: 1.749552,
		},
		{
			name:     "in to cm",
			value:    12,
			fromUnit: "in",
			toUnit:   "cm",
			unitType: "height",
			expected: 30.48,
		},
		{
			name:     "km to m",
			value:    1,
			fromUnit: "km",
			toUnit:   "m",
			unitType: "distance",
			expected: 1000,
		},
		{
			name:     "mi to km",
			value:    1,
			fromUnit: "mi",
			toUnit:   "km",
			unitType: "distance",
			expected: 1.609344,
		},
		{
			name:     "C to F",
			value:    0,
			fromUnit: "C",
			toUnit:   "F",
			unitType: "temperature",
			expected: 32,
		},
		{
			name:     "F to C",
			value:    32,
			fromUnit: "F",
			toUnit:   "C",
			unitType: "temperature",
			expected: 0,
		},
		{
			name:     "C to K",
			value:    0,
			fromUnit: "C",
			toUnit:   "K",
			unitType: "temperature",
			expected: 273.15,
		},
		{
			name:     "F to K",
			value:    32,
			fromUnit: "F",
			toUnit:   "K",
			unitType: "temperature",
			expected: 273.15,
		},
		{
			name:     "K to C",
			value:    273.15,
			fromUnit: "K",
			toUnit:   "C",
			unitType: "temperature",
			expected: 0,
		},
		{
			name:     "L to ml",
			value:    1,
			fromUnit: "L",
			toUnit:   "ml",
			unitType: "volume",
			expected: 1000,
		},
		{
			name:     "gal to L",
			value:    1,
			fromUnit: "gal",
			toUnit:   "L",
			unitType: "volume",
			expected: 3.78541,
		},
		{
			name:     "fl_oz to ml",
			value:    1,
			fromUnit: "fl_oz",
			toUnit:   "ml",
			unitType: "volume",
			expected: 29.5735,
		},
		{
			name:     "kg to kg",
			value:    10,
			fromUnit: "kg",
			toUnit:   "kg",
			unitType: "weight",
			expected: 10,
		},
		{
			name:         "invalid unit type",
			value:        10,
			fromUnit:     "kg",
			toUnit:       "lb",
			unitType:     "invalid",
			expectError:  true,
			errorMessage: "invalid unit type",
		},
		{
			name:         "invalid from unit",
			value:        10,
			fromUnit:     "invalid",
			toUnit:       "kg",
			unitType:     "weight",
			expectError:  true,
			errorMessage: "invalid source unit",
		},
		{
			name:         "invalid to unit",
			value:        10,
			fromUnit:     "kg",
			toUnit:       "invalid",
			unitType:     "weight",
			expectError:  true,
			errorMessage: "invalid target unit",
		},
		{
			name:         "NaN value",
			value:        math.NaN(),
			fromUnit:     "kg",
			toUnit:       "lb",
			unitType:     "weight",
			expectError:  true,
			errorMessage: "value must be a valid number",
		},
		{
			name:         "Inf value",
			value:        math.Inf(1),
			fromUnit:     "kg",
			toUnit:       "lb",
			unitType:     "weight",
			expectError:  true,
			errorMessage: "value must be a valid number",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertUnit(tt.value, tt.fromUnit, tt.toUnit, tt.unitType)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("expected %.6f, got %.6f", tt.expected, result)
			}
		})
	}
}
