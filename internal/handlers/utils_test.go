package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/m-szczepanski/gocalc-api/internal/middleware"
	"github.com/m-szczepanski/gocalc-api/internal/models"
)

func TestBMIHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           *models.BMIRequest
		expectedStatus int
		expectedBMI    float64
		expectedCat    string
		expectError    bool
	}{
		{
			name:   "valid BMI calculation - kg and m",
			method: http.MethodPost,
			body: &models.BMIRequest{
				Weight:     70,
				WeightUnit: "kg",
				Height:     1.75,
				HeightUnit: "m",
			},
			expectedStatus: http.StatusOK,
			expectedBMI:    22.86,
			expectedCat:    "normal",
			expectError:    false,
		},
		{
			name:   "valid BMI calculation - lb and ft",
			method: http.MethodPost,
			body: &models.BMIRequest{
				Weight:     154,
				WeightUnit: "lb",
				Height:     5.74,
				HeightUnit: "ft",
			},
			expectedStatus: http.StatusOK,
			expectedBMI:    22.82,
			expectedCat:    "normal",
			expectError:    false,
		},
		{
			name:   "underweight BMI",
			method: http.MethodPost,
			body: &models.BMIRequest{
				Weight:     50,
				WeightUnit: "kg",
				Height:     175,
				HeightUnit: "cm",
			},
			expectedStatus: http.StatusOK,
			expectedBMI:    16.33,
			expectedCat:    "underweight",
			expectError:    false,
		},
		{
			name:   "overweight BMI",
			method: http.MethodPost,
			body: &models.BMIRequest{
				Weight:     85,
				WeightUnit: "kg",
				Height:     1.75,
				HeightUnit: "m",
			},
			expectedStatus: http.StatusOK,
			expectedBMI:    27.76,
			expectedCat:    "overweight",
			expectError:    false,
		},
		{
			name:           "method not allowed",
			method:         http.MethodGet,
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
		},
		{
			name:   "zero weight",
			method: http.MethodPost,
			body: &models.BMIRequest{
				Weight:     0,
				WeightUnit: "kg",
				Height:     1.75,
				HeightUnit: "m",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "negative weight",
			method: http.MethodPost,
			body: &models.BMIRequest{
				Weight:     -70,
				WeightUnit: "kg",
				Height:     1.75,
				HeightUnit: "m",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "invalid weight unit",
			method: http.MethodPost,
			body: &models.BMIRequest{
				Weight:     70,
				WeightUnit: "invalid",
				Height:     1.75,
				HeightUnit: "m",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if tt.body != nil {
				body, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(tt.method, "/api/utils/bmi", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Add request ID to context
			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-request-id")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			BMIHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectError {
				var errResp models.APIErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
					t.Errorf("failed to decode error response: %v", err)
				}
				return
			}

			var successResp models.SuccessResponse
			if err := json.NewDecoder(w.Body).Decode(&successResp); err != nil {
				t.Errorf("failed to decode success response: %v", err)
				return
			}

			var bmiResp models.BMIResponse
			dataBytes, err := json.Marshal(successResp.Data)
			if err != nil {
				t.Errorf("failed to marshal data: %v", err)
				return
			}
			if err := json.Unmarshal(dataBytes, &bmiResp); err != nil {
				t.Errorf("failed to unmarshal BMI response: %v", err)
				return
			}

			if !floatEquals(bmiResp.BMI, tt.expectedBMI, 0.01) {
				t.Errorf("expected BMI %.2f, got %.2f", tt.expectedBMI, bmiResp.BMI)
			}

			if bmiResp.Category != tt.expectedCat {
				t.Errorf("expected category %s, got %s", tt.expectedCat, bmiResp.Category)
			}
		})
	}
}

func TestUnitConversionHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		body           *models.UnitConversionRequest
		expectedStatus int
		expectedResult float64
		expectError    bool
	}{
		{
			name:   "valid weight conversion - kg to lb",
			method: http.MethodPost,
			body: &models.UnitConversionRequest{
				Value:    10,
				FromUnit: "kg",
				ToUnit:   "lb",
				UnitType: "weight",
			},
			expectedStatus: http.StatusOK,
			expectedResult: 22.046226,
			expectError:    false,
		},
		{
			name:   "valid temperature conversion - C to F",
			method: http.MethodPost,
			body: &models.UnitConversionRequest{
				Value:    0,
				FromUnit: "C",
				ToUnit:   "F",
				UnitType: "temperature",
			},
			expectedStatus: http.StatusOK,
			expectedResult: 32,
			expectError:    false,
		},
		{
			name:   "valid height conversion - m to cm",
			method: http.MethodPost,
			body: &models.UnitConversionRequest{
				Value:    1.75,
				FromUnit: "m",
				ToUnit:   "cm",
				UnitType: "height",
			},
			expectedStatus: http.StatusOK,
			expectedResult: 175,
			expectError:    false,
		},
		{
			name:   "valid distance conversion - km to mi",
			method: http.MethodPost,
			body: &models.UnitConversionRequest{
				Value:    1,
				FromUnit: "km",
				ToUnit:   "mi",
				UnitType: "distance",
			},
			expectedStatus: http.StatusOK,
			expectedResult: 0.621371,
			expectError:    false,
		},
		{
			name:   "valid volume conversion - L to ml",
			method: http.MethodPost,
			body: &models.UnitConversionRequest{
				Value:    1,
				FromUnit: "L",
				ToUnit:   "ml",
				UnitType: "volume",
			},
			expectedStatus: http.StatusOK,
			expectedResult: 1000,
			expectError:    false,
		},
		{
			name:           "method not allowed",
			method:         http.MethodGet,
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
			expectError:    true,
		},
		{
			name:   "invalid unit type",
			method: http.MethodPost,
			body: &models.UnitConversionRequest{
				Value:    10,
				FromUnit: "kg",
				ToUnit:   "lb",
				UnitType: "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "invalid from unit",
			method: http.MethodPost,
			body: &models.UnitConversionRequest{
				Value:    10,
				FromUnit: "invalid",
				ToUnit:   "lb",
				UnitType: "weight",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "invalid to unit",
			method: http.MethodPost,
			body: &models.UnitConversionRequest{
				Value:    10,
				FromUnit: "kg",
				ToUnit:   "invalid",
				UnitType: "weight",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:   "empty from_unit",
			method: http.MethodPost,
			body: &models.UnitConversionRequest{
				Value:    10,
				FromUnit: "",
				ToUnit:   "lb",
				UnitType: "weight",
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body []byte
			var err error
			if tt.body != nil {
				body, err = json.Marshal(tt.body)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req := httptest.NewRequest(tt.method, "/api/utils/unit-conversion", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Add request ID to context
			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-request-id")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			UnitConversionHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectError {
				var errResp models.APIErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&errResp); err != nil {
					t.Errorf("failed to decode error response: %v", err)
				}
				return
			}

			var successResp models.SuccessResponse
			if err := json.NewDecoder(w.Body).Decode(&successResp); err != nil {
				t.Errorf("failed to decode success response: %v", err)
				return
			}

			var convResp models.UnitConversionResponse
			dataBytes, err := json.Marshal(successResp.Data)
			if err != nil {
				t.Errorf("failed to marshal data: %v", err)
				return
			}
			if err := json.Unmarshal(dataBytes, &convResp); err != nil {
				t.Errorf("failed to unmarshal conversion response: %v", err)
				return
			}

			if !floatEquals(convResp.Result, tt.expectedResult, 0.01) {
				t.Errorf("expected result %.6f, got %.6f", tt.expectedResult, convResp.Result)
			}

			if convResp.FromUnit != tt.body.FromUnit {
				t.Errorf("expected from_unit %s, got %s", tt.body.FromUnit, convResp.FromUnit)
			}

			if convResp.ToUnit != tt.body.ToUnit {
				t.Errorf("expected to_unit %s, got %s", tt.body.ToUnit, convResp.ToUnit)
			}

			if convResp.UnitType != tt.body.UnitType {
				t.Errorf("expected unit_type %s, got %s", tt.body.UnitType, convResp.UnitType)
			}
		})
	}
}

// floatEquals checks if two float64 values are approximately equal within a tolerance
func floatEquals(a, b, tolerance float64) bool {
	return math.Abs(a-b) <= tolerance
}
