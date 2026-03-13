package engine

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestFlaky(t *testing.T) {
	origFloat64 := randomFloat64
	origInt63n := randomInt63n
	defer func() {
		randomFloat64 = origFloat64
		randomInt63n = origInt63n
	}()

	t.Run("error_path", func(t *testing.T) {
		randomFloat64 = func() float64 { return 0.1 }

		result, err := Flaky(context.Background(), FlakyConfig{
			ErrorRate: 0.2,
			DelayRate: 0.3,
			Status:    http.StatusServiceUnavailable,
			MaxDelay:  5 * time.Millisecond,
		})
		if err != nil {
			t.Fatalf("Flaky() error = %v", err)
		}
		if result.Status != http.StatusServiceUnavailable {
			t.Fatalf("expected status=%d, got %d", http.StatusServiceUnavailable, result.Status)
		}
		if result.Delayed {
			t.Fatal("expected delayed=false")
		}
		if result.Delay != 0 {
			t.Fatalf("expected delay=0, got %v", result.Delay)
		}
	})

	t.Run("delay_path", func(t *testing.T) {
		randomFloat64 = func() float64 { return 0.4 }
		randomInt63n = func(n int64) int64 { return int64(2 * time.Millisecond) }

		result, err := Flaky(context.Background(), FlakyConfig{
			ErrorRate: 0.2,
			DelayRate: 0.3,
			Status:    http.StatusServiceUnavailable,
			MaxDelay:  5 * time.Millisecond,
		})
		if err != nil {
			t.Fatalf("Flaky() error = %v", err)
		}
		if result.Status != http.StatusOK {
			t.Fatalf("expected status=%d, got %d", http.StatusOK, result.Status)
		}
		if !result.Delayed {
			t.Fatal("expected delayed=true")
		}
		if result.Delay != 2*time.Millisecond {
			t.Fatalf("expected delay=%v, got %v", 2*time.Millisecond, result.Delay)
		}
	})

	t.Run("success_path", func(t *testing.T) {
		randomFloat64 = func() float64 { return 0.9 }

		result, err := Flaky(context.Background(), FlakyConfig{
			ErrorRate: 0.2,
			DelayRate: 0.3,
			Status:    http.StatusServiceUnavailable,
			MaxDelay:  5 * time.Millisecond,
		})
		if err != nil {
			t.Fatalf("Flaky() error = %v", err)
		}
		if result.Status != http.StatusOK {
			t.Fatalf("expected status=%d, got %d", http.StatusOK, result.Status)
		}
		if result.Delayed {
			t.Fatal("expected delayed=false")
		}
		if result.Delay != 0 {
			t.Fatalf("expected delay=0, got %v", result.Delay)
		}
	})

	t.Run("delay_path_context_canceled", func(t *testing.T) {
		randomFloat64 = func() float64 { return 0.4 }
		randomInt63n = func(n int64) int64 { return int64(50 * time.Millisecond) }

		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := Flaky(ctx, FlakyConfig{
			ErrorRate: 0.2,
			DelayRate: 0.3,
			Status:    http.StatusServiceUnavailable,
			MaxDelay:  100 * time.Millisecond,
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	})
}
