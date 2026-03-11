package http

import (
	"net/http"

	"github.com/Alkindi42/probelet/internal/engine"
	"github.com/Alkindi42/probelet/internal/http/handlers"
)

// RouterConfig defines HTTP router configuration.
type RouterConfig struct {
	// ProbeletToken enables token protection for selected endpoints when non-empty.
	ProbeletToken string
}

// NewRouter returns the application's root HTTP handler.
func NewRouter(readiness engine.Readiness, cfg RouterConfig) http.Handler {
	mux := http.NewServeMux()

	rootHandler := handlers.NewRootGetHandler()

	echoHandler := handlers.NewEchoAnyHandler()
	statusHandler := handlers.NewStatusHandler()
	delayHandler := handlers.NewDelayGetHandler()
	// Health probes
	healthzHandler := handlers.NewHealthzHandler()
	readyzGetHandler := handlers.NewReadyzGetHandler(readiness)
	readyzPostHandler := handlers.NewReadyzPostHandler(readiness)
	// Stress
	stressCPUHandler := handlers.NewStressCPUGetHandler()
	stressMemoryHandler := handlers.NewStressMemoryGetHandler()
	// Docs
	docsHandler := handlers.NewDocsHandler()
	openAPIHandler := handlers.NewOpenAPIHandler()

	probeletToken := cfg.ProbeletToken
	protected := probeletToken != ""

	if protected {
		echoHandler = RequireToken(probeletToken, echoHandler)
		statusHandler = RequireToken(probeletToken, statusHandler)
		delayHandler = RequireToken(probeletToken, delayHandler)
		readyzPostHandler = RequireToken(probeletToken, readyzPostHandler)
		stressCPUHandler = RequireToken(probeletToken, stressCPUHandler)
		stressMemoryHandler = RequireToken(probeletToken, stressMemoryHandler)
	}

	mux.Handle("GET /{$}", rootHandler)
	mux.Handle("GET /healthz", healthzHandler)
	mux.Handle("GET /readyz", readyzGetHandler)
	mux.Handle("POST /readyz", readyzPostHandler)
	mux.Handle("/echo", echoHandler)
	mux.Handle("GET /delay", delayHandler)
	mux.Handle("GET /status/{code}", statusHandler)
	mux.Handle("GET /stress/cpu", stressCPUHandler)
	mux.Handle("GET /stress/memory", stressMemoryHandler)
	mux.Handle("GET /docs", docsHandler)
	mux.Handle("GET /openapi.yaml", openAPIHandler)

	var handler http.Handler = mux
	handler = Logger(handler)
	handler = RequestID(handler)

	return handler
}
