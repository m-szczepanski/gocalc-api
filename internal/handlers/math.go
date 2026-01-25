package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	apierrors "github.com/m-szczepanski/gocalc-api/internal/errors"
	"github.com/m-szczepanski/gocalc-api/internal/models"
	"github.com/m-szczepanski/gocalc-api/internal/validation"
	"github.com/m-szczepanski/gocalc-api/pkg/calculations"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if err := validation.ValidateMethod(r.Method, http.MethodPost); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	var req models.MathRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeErrorWithDetails(w, r, apierrors.InvalidInput("invalid request body").WithError(err))
		return
	}

	if err := validation.ValidateMathRequest(&req); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	result := calculations.Add(req.A, req.B)
	if err := writeSuccessResponse(w, r, models.MathResponse{Result: result}); err != nil {
		// Error already logged, headers likely already sent
		return
	}
}

func SubtractHandler(w http.ResponseWriter, r *http.Request) {
	if err := validation.ValidateMethod(r.Method, http.MethodPost); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	var req models.MathRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeErrorWithDetails(w, r, apierrors.InvalidInput("invalid request body").WithError(err))
		return
	}

	if err := validation.ValidateMathRequest(&req); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	result := calculations.Subtract(req.A, req.B)
	if err := writeSuccessResponse(w, r, models.MathResponse{Result: result}); err != nil {
		// Error already logged, headers likely already sent
		return
	}
}

func MultiplyHandler(w http.ResponseWriter, r *http.Request) {
	if err := validation.ValidateMethod(r.Method, http.MethodPost); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	var req models.MathRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeErrorWithDetails(w, r, apierrors.InvalidInput("invalid request body").WithError(err))
		return
	}

	if err := validation.ValidateMathRequest(&req); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	result := calculations.Multiply(req.A, req.B)
	if err := writeSuccessResponse(w, r, models.MathResponse{Result: result}); err != nil {
		// Error already logged, headers likely already sent
		return
	}
}

func DivideHandler(w http.ResponseWriter, r *http.Request) {
	if err := validation.ValidateMethod(r.Method, http.MethodPost); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	var req models.MathRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeErrorWithDetails(w, r, apierrors.InvalidInput("invalid request body").WithError(err))
		return
	}

	if err := validation.ValidateDivisionRequest(&req); err != nil {
		writeErrorWithDetails(w, r, err)
		return
	}

	result, _ := calculations.Divide(req.A, req.B)

	if err := writeSuccessResponse(w, r, models.MathResponse{Result: result}); err != nil {
		// Error already logged, headers likely already sent
		return
	}
}

func decodeJSONBody(body io.ReadCloser, target interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(target)
}
