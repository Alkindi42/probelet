package handlers

import (
	"encoding/json"
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
func NewReadyzGetHandler(readinessStore *engine.ReadinessStore) http.Handler {
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

// NewReadyzPostHandler returns an HTTP handler that updates the readiness state.
func NewReadyzPostHandler(readinessStore *engine.ReadinessStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

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
		readinessStore.Set(readyRequest.Ready, readyRequest.Reason)

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
