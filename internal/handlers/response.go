package handlers

import (
	"context"
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
	requestID := extractRequestID(r.Context())
	statusCode := apiErr.HTTPStatus()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := models.NewAPIErrorResponse(apiErr.Code, apiErr.Message, apiErr.Details, requestID)
	json.NewEncoder(w).Encode(resp)
}

// writeSuccessResponse writes a successful response with request metadata.
func writeSuccessResponse(w http.ResponseWriter, r *http.Request, data interface{}) error {
	requestID := extractRequestID(r.Context())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := models.NewSuccessResponse(data, requestID)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		slog.Error("failed to encode success response",
			"error", err,
			"request_id", requestID,
		)
		return err
	}
	return nil
}

func extractRequestID(ctx context.Context) string {
	requestID, ok := ctx.Value(middleware.RequestIDKey).(string)
	if !ok {
		return "unknown"
	}
	return requestID
}
