package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/m-szczepanski/gocalc-api/internal/models"
	"github.com/m-szczepanski/gocalc-api/pkg/calculations"
)

func AddHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.MathRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result := calculations.Add(req.A, req.B)
	writeJSON(w, http.StatusOK, models.MathResponse{Result: result})
}

func SubtractHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.MathRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result := calculations.Subtract(req.A, req.B)
	writeJSON(w, http.StatusOK, models.MathResponse{Result: result})
}

func MultiplyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.MathRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result := calculations.Multiply(req.A, req.B)
	writeJSON(w, http.StatusOK, models.MathResponse{Result: result})
}

func DivideHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req models.MathRequest
	if err := decodeJSONBody(r.Body, &req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := calculations.Divide(req.A, req.B)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, models.MathResponse{Result: result})
}

func decodeJSONBody(body io.ReadCloser, target interface{}) error {
	defer body.Close()
	return json.NewDecoder(body).Decode(target)
}
