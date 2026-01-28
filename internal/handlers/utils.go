package handlers

import (
	"net/http"

	apierrors "github.com/m-szczepanski/gocalc-api/internal/errors"
	"github.com/m-szczepanski/gocalc-api/internal/models"
	"github.com/m-szczepanski/gocalc-api/internal/validation"
	"github.com/m-szczepanski/gocalc-api/pkg/calculations"
)

func BMIHandler(w http.ResponseWriter, r *http.Request) {
	if err := validation.ValidateMethod(r.Method, http.MethodPost); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	var req models.BMIRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeErrorWithDetails(w, r, apierrors.InvalidInput("invalid request body").WithError(err))
		return
	}

	if err := validation.ValidateBMIRequest(&req); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	result, err := calculations.CalculateBMI(req.Weight, req.WeightUnit, req.Height, req.HeightUnit)
	if err != nil {
		writeErrorWithDetails(w, r, apierrors.ValidationError("calculation error", err.Error()))
		return
	}

	response := models.BMIResponse{
		BMI:      result.BMI,
		Category: string(result.Category),
	}

	if err := writeSuccessResponse(w, r, response); err != nil {
		// Error already logged, headers likely already sent
		return
	}
}

func UnitConversionHandler(w http.ResponseWriter, r *http.Request) {
	if err := validation.ValidateMethod(r.Method, http.MethodPost); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	var req models.UnitConversionRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeErrorWithDetails(w, r, apierrors.InvalidInput("invalid request body").WithError(err))
		return
	}

	if err := validation.ValidateUnitConversionRequest(&req); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	result, err := calculations.ConvertUnit(req.Value, req.FromUnit, req.ToUnit, req.UnitType)
	if err != nil {
		writeErrorWithDetails(w, r, apierrors.ValidationError("conversion error", err.Error()))
		return
	}

	response := models.UnitConversionResponse{
		Result:   result,
		FromUnit: req.FromUnit,
		ToUnit:   req.ToUnit,
		UnitType: req.UnitType,
	}

	if err := writeSuccessResponse(w, r, response); err != nil {
		// Error already logged, headers likely already sent
		return
	}
}
