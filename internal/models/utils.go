package models

type BMIRequest struct {
	Weight     float64 `json:"weight"`
	WeightUnit string  `json:"weight_unit"`
	Height     float64 `json:"height"`
	HeightUnit string  `json:"height_unit"`
}

type BMIResponse struct {
	BMI      float64 `json:"bmi"`
	Category string  `json:"category"`
}

type UnitConversionRequest struct {
	Value    float64 `json:"value"`
	FromUnit string  `json:"from_unit"`
	ToUnit   string  `json:"to_unit"`
	UnitType string  `json:"unit_type"`
}

type UnitConversionResponse struct {
	Result   float64 `json:"result"`
	FromUnit string  `json:"from_unit"`
	ToUnit   string  `json:"to_unit"`
	UnitType string  `json:"unit_type"`
}
