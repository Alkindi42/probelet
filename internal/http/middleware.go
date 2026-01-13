package http

import (
	"log/slog"
	"net/http"
	"time"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the endpoint
		next.ServeHTTP(w, r)

		durationNs := time.Since(start)

		slog.Info("http request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", durationNs),
			slog.String("remote_ip", r.RemoteAddr),
			slog.Int64("duration_ms", durationNs.Milliseconds()),
		)
	})
}
