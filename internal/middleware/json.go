package middleware

import (
	"encoding/json"
	"net/http"
)

func encodeJSON(w http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}
