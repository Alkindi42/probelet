package response

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents the JSON structure used for API error responses.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Envelope is the standard wrapper used for successful JSON responses.
type Envelope struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// JSONError writes a JSON-formatted error response with the given HTTP status code.
func JSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// JSON writes a JSON response with the given HTTP status code and payload.
func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

// JSONResponse writes a standardized JSON response using the Envelope format
func JSONResponse(w http.ResponseWriter, status int, message string, data any) {
	JSON(w, status, Envelope{Message: message, Data: data})
}
