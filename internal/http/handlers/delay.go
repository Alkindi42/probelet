// Package handlers provides HTTP handlers for the application's API endpoints.
package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Alkindi42/probelet/internal/http/response"
)

const maxDelay = 5 * time.Minute

// NewDelayGetHandler returns an HTTP handler that waits for a given duration
// provided via the "duration" query parameter before responding.
func NewDelayGetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		durationStr := r.URL.Query().Get("duration")

		if durationStr == "" {
			response.JSONError(w, http.StatusBadRequest, "Missing 'duration' parameter (e.g., ?duration=2s)")
			return
		}

		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			response.JSONError(w, http.StatusBadRequest, "Invalid duration format. Use '5s', '100ms', etc.")
			return
		}

		if duration > maxDelay {
			response.JSONError(
				w,
				http.StatusBadRequest,
				fmt.Sprintf("duration exceeds maximum allowed (%.1fm)", maxDelay.Minutes()),
			)
			return
		}

		slog.Info("delay requested", "duration", durationStr)

		timer := time.NewTimer(duration)

		defer timer.Stop()

		select {
		case <-timer.C:
			response.JSON(
				w, http.StatusOK,
				"done",
				map[string]string{"duration": durationStr},
			)
		case <-r.Context().Done():
			slog.Warn("client disconnected during delay", "duration", durationStr)
			return
		}
	})
}
