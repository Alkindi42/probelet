package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"context"

	"github.com/Alkindi42/probelet/internal/engine"
	"github.com/Alkindi42/probelet/internal/http/response"
)

var maxCores = runtime.GOMAXPROCS(0)

const maxStressDuration = 2 * time.Minute

// NewStressCPUGetHandler returns an HTTP handler that triggers CPU stress
// for a given duration and number of cores.
func NewStressCPUGetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		coresStr := r.URL.Query().Get("cores")
		durationStr := r.URL.Query().Get("duration")

		cores := 1
		if coresStr == "max" {
			cores = maxCores
		} else if coresStr != "" {
			c, err := strconv.Atoi(coresStr)

			if err != nil || c <= 0 || c > maxCores {
				response.JSONError(
					w,
					http.StatusBadRequest,
					fmt.Sprintf("cores must be a number between 1 and %d or equal to 'max'", maxCores),
				)
				return
			}
			cores = c
		}

		// Parse duration (required)
		if durationStr == "" {
			response.JSONError(
				w,
				http.StatusBadRequest,
				"duration query parameter is required (e.g. 100ms, 5s, 2m)",
			)
			return
		}

		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			response.JSONError(
				w,
				http.StatusBadRequest,
				"invalid duration query parameter (examples: 100ms, 5s, 2m)",
			)
			return
		}
		if duration <= 0 {
			response.JSONError(w, http.StatusBadRequest, "duration must be greater than 0")
			return
		}
		if duration < 0 {
			response.JSONError(
				w,
				http.StatusBadRequest,
				"duration must be greater than 0",
			)
			return
		}

		if duration > maxStressDuration {
			response.JSONError(
				w,
				http.StatusBadRequest,
				fmt.Sprintf("duration must be <= %s", maxStressDuration),
			)
			return
		}

		if err := engine.StressCPU(r.Context(), cores, duration); err != nil {
			if errors.Is(err, context.Canceled) {
				slog.Warn("client disconnected during cpu stress", "duration", duration.String())
				return
			}
		}

		response.JSON(w, http.StatusOK, "done", map[string]any{
			"cores":    cores,
			"duration": duration.String(),
		})
	})
}
