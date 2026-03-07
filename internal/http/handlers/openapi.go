package handlers

import (
	"net/http"

	"github.com/Alkindi42/probelet/internal/http/assets"
)

// NewOpenAPIHandler returns an HTTP handler that serves the embedded
// OpenAPI specification for the Probelet API in YAML format.
func NewOpenAPIHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := assets.Files.ReadFile("openapi.yaml")
		if err != nil {
			http.Error(w, "openapi.yaml not available", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})
}
