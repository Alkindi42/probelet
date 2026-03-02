package handlers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/Alkindi42/probelet/internal/engine"
	"github.com/Alkindi42/probelet/internal/http/response"
	"k8s.io/apimachinery/pkg/api/resource"
)

var maxCores = runtime.GOMAXPROCS(0)

const (
	maxMemorySizeBytes      int64         = 1 << 30
	maxMemoryStressDuration time.Duration = 5 * time.Minute
	maxCPUStressDuration                  = 2 * time.Minute
)

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

		duration, err := parseDurationParam(durationStr, maxCPUStressDuration)
		if err != nil {
			response.JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if err := engine.StressCPU(r.Context(), cores, duration); err != nil {
			if errors.Is(err, context.Canceled) {
				slog.Warn("client disconnected during cpu stress", "duration", duration.String())
				return
			}

			// Any other error is unexpected.
			response.JSONError(w, http.StatusInternalServerError, "cpu stress failed")
			return
		}

		response.JSON(w, http.StatusOK, "done", map[string]any{
			"cores":    cores,
			"duration": duration.String(),
		})
	})
}

// NewStressMemoryGetHandler returns an HTTP handler that triggers memory stress
// for a given duration and sizes.
func NewStressMemoryGetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sizeStr := r.URL.Query().Get("size")
		durationStr := r.URL.Query().Get("duration")

		duration, err := parseDurationParam(durationStr, maxMemoryStressDuration)
		if err != nil {
			response.JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		if sizeStr == "" {
			response.JSONError(
				w,
				http.StatusBadRequest,
				"size is required (e.g. 64Mi, 128Mi, 1Gi)",
			)
			return
		}
		size, err := resource.ParseQuantity(sizeStr)
		if err != nil {
			response.JSONError(
				w,
				http.StatusBadRequest,
				"invalid size (examples: 64Mi, 128Mi, 1Gi)",
			)
			return
		}

		sizeBytes := size.Value()

		if sizeBytes <= 0 {
			response.JSONError(w, http.StatusBadRequest, "size must be greater than 0")
			return
		}
		if sizeBytes > maxMemorySizeBytes {
			maxQ := resource.NewQuantity(maxMemorySizeBytes, resource.BinarySI)
			response.JSONError(
				w,
				http.StatusBadRequest,
				fmt.Sprintf("size must be <= %s", maxQ.String()),
			)
			return
		}

		if err := engine.StressMemory(r.Context(), sizeBytes, duration); err != nil {
			if errors.Is(err, context.Canceled) {
				slog.Warn("client disconnected during memory stress", "duration", duration.String())
				return
			}
			response.JSONError(
				w,
				http.StatusInternalServerError,
				"memory stress failed",
			)
			return
		}

		response.JSON(w, http.StatusOK, "done", map[string]any{
			"size":     sizeStr,
			"bytes":    sizeBytes,
			"duration": duration.String(),
		})
	})
}
