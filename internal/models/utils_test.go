package models

import (
	"encoding/json"
	"testing"
)

func TestBMIRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected BMIRequest
		wantErr  bool
	}{
		{
			name: "valid request with kg and m",
			json: `{"weight": 70, "weight_unit": "kg", "height": 1.75, "height_unit": "m"}`,
			expected: BMIRequest{
				Weight:     70,
				WeightUnit: "kg",
				Height:     1.75,
				HeightUnit: "m",
			},
			wantErr: false,
		},
		{
			name: "valid request with lb and ft",
			json: `{"weight": 154, "weight_unit": "lb", "height": 5.74, "height_unit": "ft"}`,
			expected: BMIRequest{
				Weight:     154,
				WeightUnit: "lb",
				Height:     5.74,
				HeightUnit: "ft",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			json:    `{"weight": "invalid"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req BMIRequest
			err := json.Unmarshal([]byte(tt.json), &req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if req.Weight != tt.expected.Weight {
				t.Errorf("expected weight %v, got %v", tt.expected.Weight, req.Weight)
			}
			if req.WeightUnit != tt.expected.WeightUnit {
				t.Errorf("expected weight_unit %v, got %v", tt.expected.WeightUnit, req.WeightUnit)
			}
			if req.Height != tt.expected.Height {
				t.Errorf("expected height %v, got %v", tt.expected.Height, req.Height)
			}
			if req.HeightUnit != tt.expected.HeightUnit {
				t.Errorf("expected height_unit %v, got %v", tt.expected.HeightUnit, req.HeightUnit)
			}
		})
	}
}

func TestBMIResponseJSON(t *testing.T) {
	response := BMIResponse{
		BMI:      22.86,
		Category: "normal",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `{"bmi":22.86,"category":"normal"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}

func TestUnitConversionRequestJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected UnitConversionRequest
		wantErr  bool
	}{
		{
			name: "valid weight conversion",
			json: `{"value": 10, "from_unit": "kg", "to_unit": "lb", "unit_type": "weight"}`,
			expected: UnitConversionRequest{
				Value:    10,
				FromUnit: "kg",
				ToUnit:   "lb",
				UnitType: "weight",
			},
			wantErr: false,
		},
		{
			name: "valid temperature conversion",
			json: `{"value": 0, "from_unit": "C", "to_unit": "F", "unit_type": "temperature"}`,
			expected: UnitConversionRequest{
				Value:    0,
				FromUnit: "C",
				ToUnit:   "F",
				UnitType: "temperature",
			},
			wantErr: false,
		},
		{
			name:    "invalid json",
			json:    `{"value": "invalid"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req UnitConversionRequest
			err := json.Unmarshal([]byte(tt.json), &req)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if req.Value != tt.expected.Value {
				t.Errorf("expected value %v, got %v", tt.expected.Value, req.Value)
			}
			if req.FromUnit != tt.expected.FromUnit {
				t.Errorf("expected from_unit %v, got %v", tt.expected.FromUnit, req.FromUnit)
			}
			if req.ToUnit != tt.expected.ToUnit {
				t.Errorf("expected to_unit %v, got %v", tt.expected.ToUnit, req.ToUnit)
			}
			if req.UnitType != tt.expected.UnitType {
				t.Errorf("expected unit_type %v, got %v", tt.expected.UnitType, req.UnitType)
			}
		})
	}
}

func TestUnitConversionResponseJSON(t *testing.T) {
	response := UnitConversionResponse{
		Result:   22.046226,
		FromUnit: "kg",
		ToUnit:   "lb",
		UnitType: "weight",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := `{"result":22.046226,"from_unit":"kg","to_unit":"lb","unit_type":"weight"}`
	if string(data) != expected {
		t.Errorf("expected %s, got %s", expected, string(data))
	}
}
