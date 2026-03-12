package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Alkindi42/probelet/internal/engine"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	defaultBaseDir                      = "probelet"
	maxDiskSizeBytes      int64         = 5 * GiB
	maxDiskStressDuration time.Duration = 5 * time.Minute
)

// DiskStressRequest contains raw input parameters for a disk stress request.
type DiskStressRequest struct {
	Size     string
	Duration string
}

// DiskStressResult contains validated parameters used to run the disk stress.
type DiskStressResult struct {
	Size     string
	Bytes    int64
	Duration time.Duration
}

// RunDiskStress validates the request and runs the disk stress workload.
//
// It returns a ValidationError if the input parameters are invalid.
// If the context is canceled, context.Canceled is returned.
func RunDiskStress(ctx context.Context, req DiskStressRequest) (*DiskStressResult, error) {
	result := DiskStressResult{
		Size: req.Size,
	}

	duration, err := ParseDurationParam(req.Duration, maxDiskStressDuration)
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
	if sizeBytes > maxDiskSizeBytes {
		maxQ := resource.NewQuantity(maxDiskSizeBytes, resource.BinarySI)
		return nil, &ValidationError{Message: fmt.Sprintf("size must be <= %s", maxQ.String())}
	}
	result.Bytes = sizeBytes

	baseDir := filepath.Join(os.TempDir(), defaultBaseDir)
	if err := os.MkdirAll(baseDir, 0o700); err != nil {
		return nil, fmt.Errorf("create base temp dir: %w", err)
	}

	workDir, err := os.MkdirTemp(baseDir, "stress-disk-*")
	if err != nil {
		return nil, fmt.Errorf("create stress temp dir: %w", err)
	}
	defer func() {
		_ = os.RemoveAll(workDir)
	}()

	if err := engine.StressDisk(ctx, result.Bytes, result.Duration, workDir); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		return nil, fmt.Errorf("stress disk failed: %w", err)
	}

	return &result, nil
}
