package http_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

func doRequest(
	t *testing.T,
	server http.Handler,
	method, path string,
	body []byte,
	token string,
) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("X-Probelet-Token", token)
	}

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	return rr
}

func TestRouter_AuthenticationPolicy(t *testing.T) {
	t.Run("without_token_configured_protected_routes_remain_accessible", func(t *testing.T) {
		fakeReadiness := FakeReadiness{ready: true}
		server := apphttp.NewRouter(&fakeReadiness, apphttp.RouterConfig{})

		tests := []struct {
			name       string
			method     string
			path       string
			body       []byte
			wantStatus int
		}{
			{"delay_is_accessible", http.MethodGet, "/delay?duration=1ms", nil, http.StatusOK},
			{"echo_is_accessible", http.MethodGet, "/echo", nil, http.StatusOK},
			{"readyz_post_is_accessible", http.MethodPost, "/readyz", []byte(`{"ready":false,"reason":"maintenance"}`), http.StatusOK},
			{"stress_cpu_is_accessible", http.MethodGet, "/stress/cpu?duration=1ms&cores=1", nil, http.StatusOK},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				rr := doRequest(t, server, tc.method, tc.path, tc.body, "")

				if rr.Code != tc.wantStatus {
					t.Fatalf("expected status %d, got %d, body=%s", tc.wantStatus, rr.Code, rr.Body.String())
				}
			})
		}
	})

	t.Run("with_token_configured_public_routes_remain_accessible", func(t *testing.T) {
		fakeReadiness := FakeReadiness{ready: true}
		server := apphttp.NewRouter(&fakeReadiness, apphttp.RouterConfig{
			ProbeletToken: "secret-token",
		})

		tests := []struct {
			name       string
			method     string
			path       string
			wantStatus int
		}{
			{"root_is_public", http.MethodGet, "/", http.StatusOK},
			{"healthz_is_public", http.MethodGet, "/healthz", http.StatusOK},
			{"readyz_get_is_public", http.MethodGet, "/readyz", http.StatusOK},
			{"docs_is_public", http.MethodGet, "/docs", http.StatusOK},
			{"openapi_is_public", http.MethodGet, "/openapi.yaml", http.StatusOK},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				rr := doRequest(t, server, tc.method, tc.path, nil, "")

				if rr.Code != tc.wantStatus {
					t.Fatalf("expected status %d, got %d, body=%s", tc.wantStatus, rr.Code, rr.Body.String())
				}
			})
		}
	})

	t.Run("with_token_configured_protected_routes_require_valid_token", func(t *testing.T) {
		fakeReadiness := FakeReadiness{ready: true}
		server := apphttp.NewRouter(&fakeReadiness, apphttp.RouterConfig{
			ProbeletToken: "secret-token",
		})

		tests := []struct {
			name       string
			method     string
			path       string
			body       []byte
			token      string
			wantStatus int
		}{
			{"delay_without_token_is_unauthorized", http.MethodGet, "/delay?duration=1ms", nil, "", http.StatusUnauthorized},
			{"delay_with_invalid_token_is_unauthorized", http.MethodGet, "/delay?duration=1ms", nil, "wrong-token", http.StatusUnauthorized},
			{"delay_with_valid_token_is_allowed", http.MethodGet, "/delay?duration=1ms", nil, "secret-token", http.StatusOK},
			{"echo_without_token_is_unauthorized", http.MethodGet, "/echo", nil, "", http.StatusUnauthorized},
			{"status_with_valid_token_is_allowed", http.MethodGet, "/status/503", nil, "secret-token", http.StatusServiceUnavailable},
			{"readyz_post_without_token_is_unauthorized", http.MethodPost, "/readyz", []byte(`{"ready":false,"reason":"maintenance"}`), "", http.StatusUnauthorized},
			{"readyz_post_with_valid_token_is_allowed", http.MethodPost, "/readyz", []byte(`{"ready":false,"reason":"maintenance"}`), "secret-token", http.StatusOK},
			{"stress_cpu_without_token_is_unauthorized", http.MethodGet, "/stress/cpu?duration=1ms&cores=1", nil, "", http.StatusUnauthorized},
			{"stress_cpu_with_valid_token_is_allowed", http.MethodGet, "/stress/cpu?duration=1ms&cores=1", nil, "secret-token", http.StatusOK},
			{"stress_memory_without_token_is_unauthorized", http.MethodGet, "/stress/memory?duration=1ms&size=1Mi", nil, "", http.StatusUnauthorized},
			{"stress_memory_with_valid_token_is_allowed", http.MethodGet, "/stress/memory?duration=1ms&size=1Mi", nil, "secret-token", http.StatusOK},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				rr := doRequest(t, server, tc.method, tc.path, tc.body, tc.token)

				if rr.Code != tc.wantStatus {
					t.Fatalf("expected status %d, got %d, body=%s", tc.wantStatus, rr.Code, rr.Body.String())
				}
			})
		}
	})
}
