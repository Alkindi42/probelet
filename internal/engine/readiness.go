package engine

import "sync"

// ReadinessStore holds the in-memory readiness state of the application.
// It is safe for concurrent use.
type ReadinessStore struct {
	ready  bool
	reason string
	mu     sync.RWMutex
}

// Set updates the readiness state.
func (r *ReadinessStore) Set(ready bool, reason string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.ready = ready
	r.reason = reason
}

// Get returns the current readiness state and reason.
func (r *ReadinessStore) Get() (ready bool, reason string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.ready, r.reason
}

// NewReadinessStore creates a ReadinessStore initialized as ready.
func NewReadinessStore() *ReadinessStore {
	return &ReadinessStore{ready: true, reason: ""}
}
