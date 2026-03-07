package http_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

// ################
// /status/{code} #
// ################
func TestStatus(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantStatus  int
		wantOK      bool
		wantMessage string // exact or substring (your assert supports both)
		wantCode    int
	}{
		{"valid_202", "/status/202", 202, true, strings.ToLower(http.StatusText(http.StatusAccepted)), 202},
		{"valid_504", "/status/504", 504, true, strings.ToLower(http.StatusText(http.StatusGatewayTimeout)), 504},
		{"invalid_99", "/status/99", 400, false, "invalid status code", 0},
		{"invalid_600", "/status/600", 400, false, "invalid status code", 0},
		{"invalid_nonint", "/status/abc", 400, false, "invalid status code", 0},
		// duation
		{"valid_503_with_duration", "/status/503?duration=5ms", 503, true, strings.ToLower(http.StatusText(http.StatusServiceUnavailable)), 503},
		{"invalid_duration_format", "/status/200?duration=abc", 400, false, "invalid", 0},
		{"duration_too_large", "/status/200?duration=20m", 400, false, "<=", 0},
	}

	fakeReadiness := FakeReadiness{ready: true}
	server := apphttp.NewRouter(&fakeReadiness)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := doJSON(t, server, http.MethodGet, tc.path, nil)
			env := assertJSONResponse(t, rr, tc.wantStatus, tc.wantOK, tc.wantMessage)

			if tc.wantCode == 0 {
				assertNoData(t, env)
				return
			}

			var data struct {
				Code int `json:"code"`
			}
			if err := json.Unmarshal(env.Data, &data); err != nil {
				t.Fatalf("invalid data json: %v; data=%s", err, string(env.Data))
			}
			if data.Code != tc.wantCode {
				t.Fatalf("expected data.code=%d, got %d", tc.wantCode, data.Code)
			}
		})
	}
}
