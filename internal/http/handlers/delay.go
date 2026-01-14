package handlers

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Alkindi42/probelet/internal/http/response"
)

// NewDelayHandler returns an HTTP handler that waits for a given duration
// provided via the "duration" query parameter before responding.
func NewDelayHandler() http.Handler {
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

		slog.Info("delay requested", "duration", durationStr)

		timer := time.NewTimer(duration)

		defer timer.Stop()

		select {
		case <-timer.C:
			response.JSONResponse(w, http.StatusOK, "Delay finished", map[string]string{
				"requested_duration": durationStr,
			})
		case <-r.Context().Done():
			slog.Warn("client disconnected during delay", "duration", durationStr)
			return
		}
	})
}
