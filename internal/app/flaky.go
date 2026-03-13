package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/Alkindi42/probelet/internal/engine"
)

type FlakyRequest struct {
	DelayRate string
	ErrorRate string
	MaxDelay  string
	Status    string
}

type FlakyResult struct {
	Status  int
	Delayed bool
	Delay   time.Duration
}

const maxFlakyDelay = 1 * time.Minute

// RunFlaky validates the request parameters and executes a flaky response
// simulation.
//
// A ValidationError is returned when one or more input parameters are
// invalid. If the context is canceled while executing the simulation,
// context.Canceled is returned.
func RunFlaky(ctx context.Context, req FlakyRequest) (*FlakyResult, error) {
	errorRate, err := ParseRateParam(req.ErrorRate)
	if err != nil {
		return nil, &ValidationError{Message: fmt.Sprintf("error_rate %s", err.Error())}
	}
	delayRate, err := ParseRateParam(req.DelayRate)
	if err != nil {
		return nil, &ValidationError{Message: fmt.Sprintf("delay_rate %s", err.Error())}
	}
	if delayRate == 0 && errorRate == 0 {
		return nil, &ValidationError{Message: "at least one of 'error_rate' or 'delay_rate' must be greater than 0"}
	}
	if errorRate+delayRate > 1 {
		return nil, &ValidationError{
			Message: "the sum of 'error_rate' and 'delay_rate' must be <= 1",
		}
	}

	var md time.Duration

	if delayRate > 0 {
		md, err = ParseDurationParam(req.MaxDelay, maxFlakyDelay)
		if err != nil {
			return nil, &ValidationError{Message: err.Error()}
		}
	}

	status := http.StatusServiceUnavailable
	if req.Status != "" {
		s, err := strconv.Atoi(req.Status)
		if err != nil {
			return nil, &ValidationError{Message: "invalid status code"}
		}
		if s < 100 || s > 599 {
			return nil, &ValidationError{Message: "status code must be between 100 and 599"}
		}
		status = s
	}

	result, err := engine.Flaky(ctx, engine.FlakyConfig{
		ErrorRate: errorRate,
		DelayRate: delayRate,
		Status:    status,
		MaxDelay:  md,
	})
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, err
		}
		return nil, fmt.Errorf("flaky failed: %w", err)
	}

	return &FlakyResult{
		Status:  result.Status,
		Delayed: result.Delayed,
		Delay:   result.Delay,
	}, nil
}
