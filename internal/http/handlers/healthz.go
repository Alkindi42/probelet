// Package handlers provides HTTP handlers for the application's API endpoints
package handlers

import (
	"net/http"

	"github.com/Alkindi42/probelet/internal/http/response"
)

func NewHealthzHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, "healthy", nil)
	})
}
