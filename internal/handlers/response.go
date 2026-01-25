package handlers

import (
	"encoding/json"
	"net/http"

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
