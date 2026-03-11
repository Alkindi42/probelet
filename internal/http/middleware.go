package http

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Alkindi42/probelet/internal/http/response"
)

type requestIDKey struct{}

const (
	requestIDHeader     = "X-Request-Id"
	probeletTokenHeader = "X-Probelet-Token"
)

// GetRequestID returns the request ID from ctx, if present.
func GetRequestID(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey{}).(string)
	return id, ok
}

// newRequestID generates a cryptographically random request ID.
func newRequestID() string {
	var b [16]byte

	if _, err := rand.Read(b[:]); err != nil {
		// fallback
		return strconv.FormatInt(time.Now().UnixNano(), 10)
	}
	return hex.EncodeToString(b[:])
}

// RequestID is middleware that ensures every request has a stable request ID.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do not inject on echo endpoints to keep request inspection faithful.
		if r.URL.Path == "/echo" || r.URL.Path == "/echo/" {
			id := strings.TrimSpace(r.Header.Get(requestIDHeader))
			if id != "" {
				ctx := context.WithValue(r.Context(), requestIDKey{}, id)
				w.Header().Set(requestIDHeader, id)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		id := strings.TrimSpace(r.Header.Get(requestIDHeader))

		if id == "" {
			id = newRequestID()
		}

		ctx := context.WithValue(r.Context(), requestIDKey{}, id)
		w.Header().Set(requestIDHeader, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireToken is middleware that protects an HTTP handler with a static token.
//
// The client must provide the expected token in the X-Probelet-Token header.
// If the header is missing or invalid, the middleware responds with 401 Unauthorized.
func RequireToken(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(probeletTokenHeader)
		if token == "" {
			response.JSONError(w, http.StatusUnauthorized, "missing "+probeletTokenHeader+" header")
			return
		}

		if subtle.ConstantTimeCompare([]byte(token), []byte(secret)) != 1 {
			response.JSONError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Logger logs basic HTTP request information.
// If a request ID is present in the context, it is included for correlation.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Call the endpoint
		next.ServeHTTP(w, r)

		duration := time.Since(start)

		attrs := []slog.Attr{
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Duration("duration", duration),
			slog.String("remote_ip", r.RemoteAddr),
			slog.Int64("duration_ms", duration.Milliseconds()),
		}

		if rid, ok := GetRequestID(r.Context()); ok {
			attrs = append(attrs, slog.String("request_id", rid))
		}

		slog.LogAttrs(
			r.Context(),
			slog.LevelInfo,
			"http request",
			attrs...,
		)
	})
}
