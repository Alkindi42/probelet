package http

import (
	"net/http"

	"github.com/Alkindi42/probelet/internal/http/handlers"
)

// NewRouter returns the application's root HTTP handler.
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /http/status/{code}", handlers.NewStatusHandler())
	mux.Handle("GET /http/delay", handlers.NewDelayHandler())

	return Logger(mux)
}
