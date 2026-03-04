// Package handlers provides HTTP handlers for the application's API endpoints.
package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Alkindi42/probelet/internal/engine"
	"github.com/Alkindi42/probelet/internal/http/response"
)

// ReadyRequest represents the JSON payload used to update the readiness state.
type ReadyRequest struct {
	Ready  bool   `json:"ready"`
	Reason string `json:"reason,omitempty"`
}

// NewReadyzGetHandler returns an HTTP handler that reports the current readiness
// state of the application, intended for use as a Kubernetes readiness probe.
func NewReadyzGetHandler(readiness engine.Readiness) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ready, reason := readiness.Get()

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

// NewReadyzPostHandler returns an HTTP handler that updates the readiness state.
func NewReadyzPostHandler(readiness engine.Readiness) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := r.Body.Close(); err != nil {
				slog.Warn("close request body", "err", err, "method", r.Method, "path", r.URL.Path)
			}
		}()

		var readyRequest ReadyRequest

		dec := json.NewDecoder(r.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&readyRequest); err != nil {
			response.JSONError(w, http.StatusBadRequest, "invalid payload")
			return
		}

		if readyRequest.Ready {
			readyRequest.Reason = ""
		}
		readiness.Set(readyRequest.Ready, readyRequest.Reason)

		response.JSON(
			w,
			http.StatusOK,
			"updated",
			map[string]any{
				"ready":  readyRequest.Ready,
				"reason": readyRequest.Reason,
			},
		)
	})
}
