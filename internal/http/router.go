// Package http provides the main HTTP routing and middleware setup for the application.package http
package http

import (
	"net/http"

	"github.com/Alkindi42/probelet/internal/engine"
	"github.com/Alkindi42/probelet/internal/http/handlers"
)

// NewRouter returns the application's root HTTP handler.
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/echo", handlers.NewEchoAnyHandler())
	mux.Handle("GET /delay", handlers.NewDelayHandler())
	mux.Handle("GET /status/{code}", handlers.NewStatusHandler())
	// Liveness
	mux.Handle("GET /healthz", handlers.NewHealthzHandler())
	// Readiness
	readinessStore := engine.NewReadinessStore()
	mux.Handle("GET /readyz", handlers.NewReadyzGetHandler(readinessStore))
	mux.Handle("POST /readyz", handlers.NewReadyzPostHandler(readinessStore))
	// Stress
	mux.Handle("GET /stress/cpu", handlers.NewStressCPUGetHandler())

	return Logger(mux)
}
