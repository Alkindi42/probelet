package http_test

import (
	"net/http"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

// func TestDelay(t *testing.T) {
// 	tests := []struct {
// 		name        string
// 		path        string
// 		wantStatus  int
// 		wantOK      bool
// 		wantMessage string
// 	}{
// 		{"without_duration_400", "/delay", http.StatusBadRequest, false, "Missing 'duration' parameter (e.g., ?duration=2s)"},
// 		{"invalid_duration_400", "/delay?duration=foo", http.StatusBadRequest, false, "Invalid duration format. Use '5s', '100ms', etc."},
// 		{"duration_too_long_400", "/delay?duration=10m", http.StatusBadRequest, false, "duration exceeds maximum allowed (5.0m)"},
// 		{"ok_200", "/delay?duration=5ms", http.StatusOK, true, "done"},
// 	}
//
// 	server := apphttp.NewRouter()
//
// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			rr := doJSON(t, server, http.MethodGet, tc.path, nil)
// 			assertJSONResponse(t, rr, tc.wantStatus, tc.wantOK, tc.wantMessage)
// 		})
// 	}
// }

func TestDelay(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		wantStatus  int
		wantOK      bool
		wantMessage string
	}{
		{"missing_duration", "/delay", http.StatusBadRequest, false, "Missing"},
		{"invalid_duration", "/delay?duration=foo", http.StatusBadRequest, false, "format"},
		{"duration_too_long", "/delay?duration=10m", http.StatusBadRequest, false, "maximum"},
		{"ok", "/delay?duration=5ms", http.StatusOK, true, "done"},
	}

	server := apphttp.NewRouter()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := doJSON(t, server, http.MethodGet, tc.path, nil)
			assertJSONResponse(t, rr, tc.wantStatus, tc.wantOK, tc.wantMessage)
		})
	}
}
