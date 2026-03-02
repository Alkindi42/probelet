package engine

import (
	"context"
	"errors"
	"os"
	"runtime"
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

// freeMemory
func freeMemory(blocks [][]byte) {
	for i := range blocks {
		blocks[i] = nil
	}

	runtime.GC()
}

// StressMemory allocate memory on a specified size for the given duration.
// The allocation is canceled when the context is done.
func StressMemory(ctx context.Context, size int64, duration time.Duration) error {
	var blocks [][]byte

	defer freeMemory(blocks)

	pageSize := os.Getpagesize()
	const blockSize = 32 * 1024 * 1024

	remaining := size

	for remaining > 0 {
		if err := ctx.Err(); err != nil {
			return err
		}

		allocSize := min(blockSize, remaining)
		block := make([]byte, allocSize)

		for i := 0; i < len(block); i += pageSize {
			block[i] = 42
			if i%(pageSize*256) == 0 {
				if err := ctx.Err(); err != nil {
					return err
				}
			}
		}

		// Touch the last byte in case len(b) is not a factor of pagSize.
		block[len(block)-1] = 42

		blocks = append(blocks, block)
		remaining -= allocSize
	}

	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	<-ctx.Done()

	err := ctx.Err()
	if errors.Is(err, context.DeadlineExceeded) {
		return nil
	}
	return err
}
