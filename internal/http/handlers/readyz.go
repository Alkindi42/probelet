package handlers

import (
	"net/http"

	"github.com/Alkindi42/probelet/internal/engine"
	"github.com/Alkindi42/probelet/internal/http/response"
)

// NewReadyzHandler returns an HTTP handler that reports the current readiness
// state of the application, intended for use as a Kubernetes readiness probe.
func NewReadyzHandler(readinessStore *engine.ReadinessStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ready, reason := readinessStore.Get()

		if !ready {
			response.JSONError(
				w,
				http.StatusServiceUnavailable,
				reason,
			)
		} else {
			response.JSON(w, http.StatusOK, "ready", nil)
		}
	})
}
