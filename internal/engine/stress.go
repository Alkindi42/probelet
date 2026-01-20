package engine

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

var cpuSink uint64

// StressCPU burns CPU on the specified number of cores for the given duration.
// The stress is canceled when the context is done.
func StressCPU(ctx context.Context, cores int, duration time.Duration) error {

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(cores)

	for core := range cores {
		go func(workerID int) {
			defer wg.Done()

			const step = 4096

			x := uint64(workerID + 1)
			for i := 0; ; i++ {
				// cpu work
				x ^= x << 13
				x ^= x >> 7
				x ^= x << 17

				if i%step == 0 && ctx.Err() != nil {
					atomic.AddUint64(&cpuSink, x)
					return
				}
			}
		}(core)
	}

	wg.Wait()

	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return nil
	}

	return ctx.Err()
}
