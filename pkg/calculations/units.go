package calculations

import (
	"fmt"
	"strings"
)

type UnitType string

const (
	UnitTypeWeight      UnitType = "weight"
	UnitTypeHeight      UnitType = "height"
	UnitTypeTemperature UnitType = "temperature"
	UnitTypeDistance    UnitType = "distance"
	UnitTypeVolume      UnitType = "volume"
)

func ValidUnitTypes() []UnitType {
	return []UnitType{
		UnitTypeWeight,
		UnitTypeHeight,
		UnitTypeTemperature,
		UnitTypeDistance,
		UnitTypeVolume,
	}
}

func IsValidUnitType(unitType string) bool {
	normalized := strings.ToLower(strings.TrimSpace(unitType))
	for _, valid := range ValidUnitTypes() {
		if string(valid) == normalized {
			return true
		}
	}
	return false
}

type WeightUnit string

const (
	WeightKilogram WeightUnit = "kg"
	WeightGram     WeightUnit = "g"
	WeightPound    WeightUnit = "lb"
	WeightOunce    WeightUnit = "oz"
)

type HeightUnit string

const (
	HeightMeter      HeightUnit = "m"
	HeightCentimeter HeightUnit = "cm"
	HeightFoot       HeightUnit = "ft"
	HeightInch       HeightUnit = "in"
)

type TemperatureUnit string

const (
	TemperatureCelsius    TemperatureUnit = "C"
	TemperatureFahrenheit TemperatureUnit = "F"
	TemperatureKelvin     TemperatureUnit = "K"
)

type DistanceUnit string

const (
	DistanceMeter     DistanceUnit = "m"
	DistanceKilometer DistanceUnit = "km"
	DistanceMile      DistanceUnit = "mi"
	DistanceFoot      DistanceUnit = "ft"
	DistanceYard      DistanceUnit = "yd"
)

type VolumeUnit string

const (
	VolumeLiter      VolumeUnit = "L"
	VolumeMilliliter VolumeUnit = "ml"
	VolumeGallon     VolumeUnit = "gal"
	VolumeFluidOunce VolumeUnit = "fl_oz"
)

type UnitRegistry struct {
	conversions map[UnitType]map[string]float64
}

func NewUnitRegistry() *UnitRegistry {
	registry := &UnitRegistry{
		conversions: make(map[UnitType]map[string]float64),
	}

	registry.conversions[UnitTypeWeight] = map[string]float64{
		string(WeightKilogram): 1.0,
		string(WeightGram):     0.001,
		string(WeightPound):    0.45359237,
		string(WeightOunce):    0.0283495231,
	}

	registry.conversions[UnitTypeHeight] = map[string]float64{
		string(HeightMeter):      1.0,
		string(HeightCentimeter): 0.01,
		string(HeightFoot):       0.3048,
		string(HeightInch):       0.0254,
	}

	registry.conversions[UnitTypeDistance] = map[string]float64{
		string(DistanceMeter):     1.0,
		string(DistanceKilometer): 1000.0,
		string(DistanceMile):      1609.344,
		string(DistanceFoot):      0.3048,
		string(DistanceYard):      0.9144,
	}

	registry.conversions[UnitTypeVolume] = map[string]float64{
		string(VolumeLiter):      1.0,
		string(VolumeMilliliter): 0.001,
		string(VolumeGallon):     3.78541,
		string(VolumeFluidOunce): 0.0295735,
		"l":                      1.0,
	}

	registry.conversions[UnitTypeTemperature] = map[string]float64{
		string(TemperatureCelsius):    1.0,
		string(TemperatureFahrenheit): 1.0,
		string(TemperatureKelvin):     1.0,
		"c":                           1.0,
		"f":                           1.0,
		"k":                           1.0,
	}

	return registry
}

func (r *UnitRegistry) GetConversionFactor(unitType UnitType, unit string) (float64, error) {
	normalized := strings.ToLower(strings.TrimSpace(unit))

	typeConversions, exists := r.conversions[unitType]
	if !exists {
		return 0, fmt.Errorf("unsupported unit type: %s", unitType)
	}

	factor, exists := typeConversions[normalized]
	if !exists {
		return 0, fmt.Errorf("unsupported unit '%s' for type '%s'", unit, unitType)
	}

	return factor, nil
}

func (r *UnitRegistry) IsValidUnit(unitType UnitType, unit string) bool {
	_, err := r.GetConversionFactor(unitType, unit)
	return err == nil
}

func (r *UnitRegistry) GetValidUnits(unitType UnitType) []string {
	typeConversions, exists := r.conversions[unitType]
	if !exists {
		return []string{}
	}

	units := make([]string, 0, len(typeConversions))
	for unit := range typeConversions {
		units = append(units, unit)
	}
	return units
}

func (r *UnitRegistry) ConvertToBaseUnit(value float64, unitType UnitType, fromUnit string) (float64, error) {
	factor, err := r.GetConversionFactor(unitType, fromUnit)
	if err != nil {
		return 0, err
	}
	return value * factor, nil
}

func (r *UnitRegistry) ConvertFromBaseUnit(value float64, unitType UnitType, toUnit string) (float64, error) {
	factor, err := r.GetConversionFactor(unitType, toUnit)
	if err != nil {
		return 0, err
	}
	return value / factor, nil
}

func (r *UnitRegistry) Convert(value float64, unitType UnitType, fromUnit, toUnit string) (float64, error) {
	// Special handling for temperature
	if unitType == UnitTypeTemperature {
		return convertTemperature(value, fromUnit, toUnit)
	}

	// Convert to base unit
	baseValue, err := r.ConvertToBaseUnit(value, unitType, fromUnit)
	if err != nil {
		return 0, err
	}

	// Convert from base unit to target unit
	result, err := r.ConvertFromBaseUnit(baseValue, unitType, toUnit)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func convertTemperature(value float64, fromUnit, toUnit string) (float64, error) {
	from := strings.ToUpper(strings.TrimSpace(fromUnit))
	to := strings.ToUpper(strings.TrimSpace(toUnit))

	if from == to {
		return value, nil
	}

	var celsius float64
	switch from {
	case "C":
		celsius = value
	case "F":
		celsius = (value - 32) * 5 / 9
	case "K":
		celsius = value - 273.15
	default:
		return 0, fmt.Errorf("unsupported temperature unit: %s", fromUnit)
	}

	var result float64
	switch to {
	case "C":
		result = celsius
	case "F":
		result = celsius*9/5 + 32
	case "K":
		result = celsius + 273.15
	default:
		return 0, fmt.Errorf("unsupported temperature unit: %s", toUnit)
	}

	return result, nil
}
