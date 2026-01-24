package handlers

import (
	"encoding/base64"
	"io"
	"net"
	"net/http"
	"unicode/utf8"

	"github.com/Alkindi42/probelet/internal/http/response"
)

const maxEchoBody = 64 << 10

// NewEchoAnyHandler returns an HTTP handler that echoes back request
// information for debugging purposes.
func NewEchoAnyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			bodyStr       string
			bodyBytes     int
			bodyTruncated bool
		)

		clientIP := r.RemoteAddr
		if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			clientIP = host
		}

		body := map[string]any{
			"content":      "",
			"bytes":        0,
			"is_truncated": false,
		}

		if r.Body != nil {
			b, err := io.ReadAll(io.LimitReader(r.Body, maxEchoBody+1))
			if err != nil {
				response.JSONError(w, http.StatusBadRequest, "failed to read request body")
				return
			}

			if len(b) > maxEchoBody {
				bodyTruncated = true
				b = b[:maxEchoBody]
			}

			bodyBytes = len(b)

			if utf8.Valid(b) {
				bodyStr = string(b)
			} else {
				bodyStr = base64.StdEncoding.EncodeToString(b)
			}

			body = map[string]any{
				"content":      bodyStr,
				"bytes":        bodyBytes,
				"is_truncated": bodyTruncated,
			}
		}

		response.JSON(w, http.StatusOK, "echo", map[string]any{
			"method":         r.Method,
			"path":           r.URL.Path,
			"query":          r.URL.Query(),
			"raw_query":      r.URL.RawQuery,
			"body":           body,
			"content_length": r.ContentLength,
			"headers":        r.Header,
			"host":           r.Host,
			"proto":          r.Proto,
			"user_agent":     r.UserAgent(),
			"client_ip":      clientIP,
			"remote_addr":    r.RemoteAddr,
			"content_type":   r.Header.Get("Content-Type"),
			"x_real_ip":      r.Header.Get("X-Real-IP"),
			"forwarded_for":  r.Header.Get("X-Forwarded-For"),
		})
	})
}
