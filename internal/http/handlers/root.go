package handlers

import (
	"net/http"

	"github.com/Alkindi42/probelet/internal/http/response"
	"github.com/Alkindi42/probelet/internal/version"
)

// NewRootGetHandler returns an HTTP handler that exposes basic
// Probelet build and discovery information at the API root.
func NewRootGetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, "probelet", map[string]string{
			"version":    version.Version,
			"commit":     version.Commit,
			"build_date": version.BuildDate,
			"docs":       "/docs",
			"openapi":    "/openapi.yaml",
		})
	})
}
