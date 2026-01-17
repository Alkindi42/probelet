package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

type envelope struct {
	OK      bool            `json:"ok"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func decodeEnvelope(t *testing.T, body []byte) envelope {
	t.Helper()
	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		t.Fatalf("failed to decode json: %v\nbody=%s", err, body)
	}
	return env
}

func TestStatus(t *testing.T) {
	server := apphttp.NewRouter()

	tests := []struct {
		name        string
		path        string
		wantStatus  int
		wantOK      bool
		wantMessage string
		wantCode    int
	}{
		{"valid_202", "/status/202", 202, true, strings.ToLower(http.StatusText(202)), 202},
		{"valid_504", "/status/504", 504, true, strings.ToLower(http.StatusText(504)), 504},
		{"invalid_99", "/status/99", 400, false, "invalid status code", 0},
		{"invalid_600", "/status/600", 400, false, "invalid status code", 0},
		{"invalid_nonint", "/status/abc", 400, false, "invalid status code", 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rr := httptest.NewRecorder()

			server.ServeHTTP(rr, req)

			if rr.Code != tc.wantStatus {
				t.Fatalf("expected %d, got %d", tc.wantStatus, rr.Code)
			}
			if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
				t.Fatalf("expected content-type json, got %q", ct)
			}

			env := decodeEnvelope(t, rr.Body.Bytes())

			if env.OK != tc.wantOK {
				t.Fatalf("expected ok=%v, got %v", tc.wantOK, env.OK)
			}
			if env.Message != tc.wantMessage {
				t.Fatalf("expected message=%q, got %q", tc.wantMessage, env.Message)
			}

			if tc.wantCode != 0 {
				var data struct {
					Code int `json:"code"`
				}
				if err := json.Unmarshal(env.Data, &data); err != nil {
					t.Fatalf("invalid data json: %v", err)
				}
				if data.Code != tc.wantCode {
					t.Fatalf("expected data.code=%d, got %d", tc.wantCode, data.Code)
				}
			} else {
				if len(env.Data) != 0 {
					t.Fatalf("expected no data, got %s", string(env.Data))
				}
			}
		})
	}
}
