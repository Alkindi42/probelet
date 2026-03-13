package engine

import (
	"context"
	"math/rand"
	"net/http"
	"time"
)

type FlakyConfig struct {
	DelayRate float64
	ErrorRate float64
	Status    int
	MaxDelay  time.Duration
}

type FlakyResult struct {
	Status  int
	Delayed bool
	Delay   time.Duration
}

var (
	randomFloat64 = rand.Float64
	randomInt63n  = rand.Int63n
)

// Flaky simulates an intermittent HTTP response behavior.
//
// Based on the configured error and delay rates, the function randomly
// decides whether the request should:
//   - return an error with the configured HTTP status
//   - introduce a delay up to MaxDelay
//   - return immediately with a successful response
//
// The function respects the provided context and returns context.Canceled
// if the operation is interrupted.
func Flaky(ctx context.Context, cfg FlakyConfig) (*FlakyResult, error) {
	r := randomFloat64()

	// Error path
	if r < cfg.ErrorRate {
		return &FlakyResult{Status: cfg.Status, Delay: 0, Delayed: false}, nil
	}

	// Delay path
	if r < cfg.ErrorRate+cfg.DelayRate {
		delay := time.Duration(randomInt63n(cfg.MaxDelay.Nanoseconds() + 1))

		timer := time.NewTimer(delay)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-timer.C:
			return &FlakyResult{
				Status:  http.StatusOK,
				Delayed: true,
				Delay:   delay,
			}, nil
		}
	}

	// Success path
	return &FlakyResult{
		Status:  http.StatusOK,
		Delayed: false,
		Delay:   0,
	}, nil
}
