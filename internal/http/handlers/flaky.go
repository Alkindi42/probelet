package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Alkindi42/probelet/internal/app"
	"github.com/Alkindi42/probelet/internal/http/response"
)

func NewFlakyGetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		statusStr := r.URL.Query().Get("status")
		maxDelayStr := r.URL.Query().Get("max_delay")
		delayRateStr := r.URL.Query().Get("delay_rate")
		errorRateStr := r.URL.Query().Get("error_rate")

		req := app.FlakyRequest{
			DelayRate: delayRateStr,
			ErrorRate: errorRateStr,
			MaxDelay:  maxDelayStr,
			Status:    statusStr,
		}

		result, err := app.RunFlaky(r.Context(), req)
		if err != nil {

			var validationErr *app.ValidationError

			switch {
			case errors.As(err, &validationErr):
				response.JSONError(w, http.StatusBadRequest, err.Error())
				return

			case errors.Is(err, context.Canceled):
				slog.Warn(
					"client disconnected during flaky simulation",
					"delay_rate", delayRateStr,
					"error_rate", errorRateStr,
					"max_delay", maxDelayStr,
				)
				return

			default:
				response.JSONError(w, http.StatusInternalServerError, "flaky failed")
				return
			}
		}

		message := "done"
		if result.Status != http.StatusOK {
			message = strings.ToLower(http.StatusText(result.Status))
		}
		response.JSON(w, result.Status, message, map[string]any{
			"delay":   result.Delay.String(),
			"delayed": result.Delayed,
			"status":  result.Status,
		})
	})
}
