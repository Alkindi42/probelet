package response

import (
	"encoding/json"
	"net/http"
)

// Envelope is the standard wrapper used for JSON API responses.
type Envelope struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// writeJSON writes a JSON response with the given HTTP status code and payload.
func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(payload)
}

// JSON writes a standardized JSON response using the Envelope format.
func JSON(w http.ResponseWriter, status int, message string, data any) {
	writeJSON(
		w,
		status,
		Envelope{
			OK:      true,
			Message: message,
			Data:    data,
		},
	)
}

// JSONError writes a standardized JSON error response using the Envelope format.
func JSONError(w http.ResponseWriter, status int, message string) {
	writeJSON(
		w,
		status,
		Envelope{
			OK:      false,
			Message: message,
		},
	)
}
