package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Alkindi42/probelet/internal/engine"
	"k8s.io/apimachinery/pkg/api/resource"
)

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

const (
	maxMemoryStressDuration time.Duration = 5 * time.Minute
	maxMemorySizeBytes      int64         = 1 * GiB
)

// RunMemoryStress validates the request and runs the memory stress workload.
//
// It returns a ValidationError if the input parameters are invalid.
// If the context is canceled, context.Canceled is returned.
func RunMemoryStress(ctx context.Context, req MemoryStressRequest) (*MemoryStressResult, error) {
	result := MemoryStressResult{
		Size: req.Size,
	}

	duration, err := ParseDurationParam(req.Duration, maxMemoryStressDuration)
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
