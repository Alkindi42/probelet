package app

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/Alkindi42/probelet/internal/engine"
)

// CPUStressRequest contains raw input parameters for a CPU stress request.
type CPUStressRequest struct {
	Cores    string
	Duration string
}

// CPUStressResult contains validated parameters used to run the CPU stress.
type CPUStressResult struct {
	Cores    int
	Duration time.Duration
}

// MemoryStressRequest contains raw input parameters for a memory stress request.
type MemoryStressRequest struct {
	Size     string
	Duration string
}

var maxCores = runtime.GOMAXPROCS(0)

const (
	maxCPUStressDuration = 2 * time.Minute
)

// RunCPUStress validates the request and runs the CPU stress workload.
//
// It returns a ValidationError if the input parameters are invalid.
// If the context is canceled, context.Canceled is returned.
func RunCPUStress(ctx context.Context, req CPUStressRequest) (*CPUStressResult, error) {
	result := CPUStressResult{
		Cores: 1,
	}

	if req.Cores == "max" {
		result.Cores = maxCores
	} else if req.Cores != "" {
		c, err := strconv.Atoi(req.Cores)

		if err != nil || c <= 0 || c > maxCores {
			return nil, &ValidationError{
				Message: fmt.Sprintf("cores must be a number between 1 and %d or equal to 'max'", maxCores),
			}
		}
		result.Cores = c
	}

	duration, err := parseDurationParam(req.Duration, maxCPUStressDuration)
	if err != nil {
		return nil, &ValidationError{
			Message: err.Error(),
		}
	}
	result.Duration = duration

	if err := engine.StressCPU(ctx, result.Cores, result.Duration); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		return nil, fmt.Errorf("stress cpu failed: %w", err)
	}

	return &result, nil
}
