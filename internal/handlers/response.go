package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	apierrors "github.com/m-szczepanski/gocalc-api/internal/errors"
	"github.com/m-szczepanski/gocalc-api/internal/middleware"
	"github.com/m-szczepanski/gocalc-api/internal/models"
)

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, statusCode int, message string) error {
	return writeJSON(w, statusCode, models.ErrorResponse{Error: message})
}

func writeErrorWithDetails(w http.ResponseWriter, r *http.Request, apiErr *apierrors.APIError) {
	requestID := middleware.ExtractRequestID(r.Context())
	statusCode := apiErr.HTTPStatus()

	resp := models.NewAPIErrorResponse(apiErr.Code, apiErr.Message, apiErr.Details, requestID)
	if err := writeJSON(w, statusCode, resp); err != nil {
		slog.Error("failed to encode error response",
			"error", err,
			"request_id", requestID,
			"status_code", statusCode,
		)
	}
}

func writeSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}) error {
	requestID := middleware.ExtractRequestID(r.Context())

	resp := models.NewSuccessResponse(data, requestID)
	if err := writeJSON(w, http.StatusOK, resp); err != nil {
		slog.Error("failed to encode success response",
			"error", err,
			"request_id", requestID,
		)
		return err
	}
	return nil
}
