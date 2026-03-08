package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Alkindi42/probelet/internal/app"
	"github.com/Alkindi42/probelet/internal/http/response"
)

// NewStressCPUGetHandler returns an HTTP handler that triggers CPU stress
// for a given duration and number of cores.
func NewStressCPUGetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		coresStr := r.URL.Query().Get("cores")
		durationStr := r.URL.Query().Get("duration")

		req := app.CPUStressRequest{
			Cores:    coresStr,
			Duration: durationStr,
		}
		result, err := app.RunCPUStress(r.Context(), req)
		if err != nil {
			var validationErr *app.ValidationError

			switch {
			case errors.As(err, &validationErr):
				response.JSONError(w, http.StatusBadRequest, err.Error())
				return

			case errors.Is(err, context.Canceled):
				slog.Warn("client disconnected during cpu stress", "duration", durationStr)
				return

			default:
				response.JSONError(w, http.StatusInternalServerError, "cpu stress failed")
				return
			}
		}

		response.JSON(w, http.StatusOK, "done", map[string]any{
			"cores":    result.Cores,
			"duration": result.Duration.String(),
		})
	})
}

// NewStressMemoryGetHandler returns an HTTP handler that triggers memory stress
// for a given duration and sizes.
func NewStressMemoryGetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sizeStr := r.URL.Query().Get("size")
		durationStr := r.URL.Query().Get("duration")

		req := app.MemoryStressRequest{
			Size:     sizeStr,
			Duration: durationStr,
		}
		result, err := app.RunMemoryStress(r.Context(), req)
		if err != nil {
			var validationErr *app.ValidationError

			switch {
			case errors.As(err, &validationErr):
				response.JSONError(w, http.StatusBadRequest, err.Error())
				return

			case errors.Is(err, context.Canceled):
				slog.Warn("client disconnected during memory stress", "duration", durationStr)
				return

			default:
				response.JSONError(w, http.StatusInternalServerError, "memory stress failed")
				return
			}
		}

		response.JSON(w, http.StatusOK, "done", map[string]any{
			"size":     result.Size,
			"bytes":    result.Bytes,
			"duration": result.Duration.String(),
		})
	})
}
