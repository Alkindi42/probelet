package http_test

import (
	"net/http"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

func TestDelay(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantStatus  int
		wantOK      bool
		wantMessage string
	}{
		{"missing_duration", "/delay", http.StatusBadRequest, false, "Missing"},
		{"invalid_duration", "/delay?duration=foo", http.StatusBadRequest, false, "invalid duration"},
		{"duration_too_long", "/delay?duration=10m", http.StatusBadRequest, false, "must be <="},
		{"ok", "/delay?duration=5ms", http.StatusOK, true, "done"},
	}

	fakeReadiness := FakeReadiness{ready: true}
	server := apphttp.NewRouter(&fakeReadiness, apphttp.RouterConfig{})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := doJSON(t, server, http.MethodGet, tc.path, nil)
			assertJSONResponse(t, rr, tc.wantStatus, tc.wantOK, tc.wantMessage)
		})
	}
}
