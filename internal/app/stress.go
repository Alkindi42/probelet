package app

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/Alkindi42/probelet/internal/engine"
	"k8s.io/apimachinery/pkg/api/resource"
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

// MemoryStressResult contains validated parameters used to run the memory stress.
type MemoryStressResult struct {
	Size     string
	Bytes    int64
	Duration time.Duration
}

var maxCores = runtime.GOMAXPROCS(0)

const (
	maxCPUStressDuration                  = 2 * time.Minute
	maxMemorySizeBytes      int64         = 1 << 30
	maxMemoryStressDuration time.Duration = 5 * time.Minute
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

// RunMemoryStress validates the request and runs the memory stress workload.
//
// It returns a ValidationError if the input parameters are invalid.
// If the context is canceled, context.Canceled is returned.
func RunMemoryStress(ctx context.Context, req MemoryStressRequest) (*MemoryStressResult, error) {
	result := MemoryStressResult{
		Size: req.Size,
	}

	duration, err := parseDurationParam(req.Duration, maxMemoryStressDuration)
	if err != nil {
		return nil, &ValidationError{
			Message: err.Error(),
		}
	}
	result.Duration = duration

	if req.Size == "" {
		return nil, &ValidationError{Message: "size is required (e.g. 64Mi, 128Mi, 1Gi)"}
	}
	size, err := resource.ParseQuantity(req.Size)
	if err != nil {
		return nil, &ValidationError{Message: "invalid size (examples: 64Mi, 128Mi, 1Gi)"}
	}

	sizeBytes := size.Value()
	if sizeBytes <= 0 {
		return nil, &ValidationError{Message: "size must be greater than 0"}
	}
	if sizeBytes > maxMemorySizeBytes {
		maxQ := resource.NewQuantity(maxMemorySizeBytes, resource.BinarySI)
		return nil, &ValidationError{Message: fmt.Sprintf("size must be <= %s", maxQ.String())}
	}
	result.Bytes = sizeBytes

	if err := engine.StressMemory(ctx, result.Bytes, result.Duration); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		return nil, fmt.Errorf("stress memory failed: %w", err)
	}

	return &result, nil
}
