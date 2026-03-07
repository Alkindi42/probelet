package handlers

import (
	"net/http"

	"github.com/Alkindi42/probelet/internal/http/assets"
)

// NewDocsHandler returns an HTTP handler that serves the embedded
// API documentation UI for the Probelet OpenAPI specification.
func NewDocsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := assets.Files.ReadFile("docs.html")
		if err != nil {
			http.Error(w, "docs.yaml not available", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	})
}
