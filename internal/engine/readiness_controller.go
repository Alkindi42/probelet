package engine

import (
	"context"
	"sync"
	"time"
)

const initialReadinessDelayReason = "initial readiness delay"

type Readiness interface {
	Get() (bool, string)
	Set(bool, string)
}

type ReadinessController struct {
	store       *ReadinessStore
	mu          sync.Mutex
	delayCancel context.CancelFunc
	delayGen    uint64
}

// NewReadinessController returns a new ReadinessController using the given store.
func NewReadinessController(store *ReadinessStore) *ReadinessController {
	return &ReadinessController{store: store}
}

// StartInitialDelay marks the application as not ready and schedules it to
// become ready after the given delay.
//
// If delay is less than or equal to zero, the method returns immediately.
//
// If another delay is already running, it is canceled and replaced.
// The delay is also canceled if the provided context is done.
func (c *ReadinessController) StartInitialDelay(ctx context.Context, delay time.Duration) {
	if delay <= 0 {
		return
	}

	c.mu.Lock()
	if c.delayCancel != nil {
		c.delayCancel()
	}
	delayCtx, cancel := context.WithCancel(ctx)
	c.delayCancel = cancel
	c.delayGen++
	myGen := c.delayGen
	c.mu.Unlock()

	c.store.Set(false, initialReadinessDelayReason)

	go func() {
		t := time.NewTimer(delay)
		defer t.Stop()

		select {
		case <-t.C:
			c.mu.Lock()
			stillActive := (c.delayGen == myGen)
			if stillActive {
				c.delayCancel = nil
			}
			c.mu.Unlock()

			if !stillActive {
				return
			}
			c.store.Set(true, "")

		case <-delayCtx.Done():
			c.mu.Lock()
			if c.delayGen == myGen {
				c.delayCancel = nil
			}
			c.mu.Unlock()
			return
		}
	}()
}

// CancelInitialDelay cancels any running initial readiness delay.
//
// If no delay is active, it is a no-op.
func (c *ReadinessController) CancelInitialDelay() {
	c.mu.Lock()
	cancel := c.delayCancel
	c.delayCancel = nil
	c.delayGen++ // invalidate any in-flight timer
	c.mu.Unlock()

	if cancel != nil {
		cancel()
	}
}

func (c *ReadinessController) Set(ready bool, reason string) {
	c.CancelInitialDelay()
	c.store.Set(ready, reason)
}

// Get returns the current readiness state and reason.
func (c *ReadinessController) Get() (bool, string) {
	return c.store.Get()
}
