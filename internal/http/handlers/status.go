package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Alkindi42/probelet/internal/http/response"
)

const (
	maxStatusDelay = 2 * time.Minute
)

func parseStatusCode(codeStr string) (int, bool) {
	code, err := strconv.Atoi(codeStr)

	if err != nil || code < 100 || code > 599 {
		return 0, false
	}

	return code, true
}

// NewStatusHandler returns an HTTP handler that responds with the given
// HTTP status code.
func NewStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		codeStr := r.PathValue("code")

		httpStatusCode, ok := parseStatusCode(codeStr)
		if !ok {
			response.JSONError(w, http.StatusBadRequest, "invalid status code")
			return
		}

		durationStr := r.URL.Query().Get("duration")
		duration, err := parseOptionalDurationParam(durationStr, maxStatusDelay)
		if err != nil {
			response.JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if duration > 0 {
			select {
			case <-time.After(duration):
				slog.Info("status delay finished", "duration", duration.String())

			case <-r.Context().Done():
				slog.Warn("client disconnected before delay finished", "duration", duration.String())
				return
			}
		}

		response.JSON(
			w,
			httpStatusCode,
			strings.ToLower(http.StatusText(httpStatusCode)),
			map[string]int{"code": httpStatusCode},
		)
	})
}
