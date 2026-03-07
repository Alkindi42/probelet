package http_test

import "github.com/Alkindi42/probelet/internal/engine"

type FakeReadiness struct {
	ready  bool
	reason string
}

var _ engine.Readiness = (*FakeReadiness)(nil)

func (f *FakeReadiness) Get() (bool, string) {
	return f.ready, f.reason
}

func (f *FakeReadiness) Set(ready bool, reason string) {
	f.ready = ready
	f.reason = reason
}
