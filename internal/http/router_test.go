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
		wantMessage string
		wantCode    int
	}{
		{"valid_202", "/status/202", 202, true, strings.ToLower(http.StatusText(202)), 202},
		{"valid_504", "/status/504", 504, true, strings.ToLower(http.StatusText(504)), 504},
		{"invalid_99", "/status/99", 400, false, "invalid status code", 0},
		{"invalid_600", "/status/600", 400, false, "invalid status code", 0},
		{"invalid_nonint", "/status/abc", 400, false, "invalid status code", 0},
	}

	server := apphttp.NewRouter()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			rr := doJSON(t, server, http.MethodGet, tc.path, nil)

			env := assertJSONResponse(t, rr, tc.wantStatus, tc.wantOK, tc.wantMessage)

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
				assertNoData(t, env)
			}
		})
	}
}

// #########
// /readyz #
// #########

func TestReadyz_StateTransitions(t *testing.T) {
	server := apphttp.NewRouter()

	// 1) Check the default state (ready).
	rr := doJSON(t, server, http.MethodGet, "/readyz", nil)
	env := assertJSONResponse(t, rr, http.StatusOK, true, "ready")
	assertNoData(t, env)

	// 2) Set not ready with a reason
	rr = doJSON(t, server, http.MethodPost, "/readyz", map[string]any{
		"ready":  false,
		"reason": "db down",
	})
	env = assertJSONResponse(t, rr, http.StatusOK, true, "updated")

	var postData struct {
		Ready  bool   `json:"ready"`
		Reason string `json:"reason"`
	}
	if err := json.Unmarshal(env.Data, &postData); err != nil {
		t.Fatalf("failed to decode data: %v; data=%s", err, string(env.Data))
	}
	if postData.Ready != false || postData.Reason != "db down" {
		t.Fatalf("unexpected post data: %+v", postData)
	}

	// 3) GET should be 503 with reason as message
	rr = doJSON(t, server, http.MethodGet, "/readyz", nil)
	env = assertJSONResponse(t, rr, http.StatusServiceUnavailable, false, "db down")
	assertNoData(t, env)

	// 4) Set ready=true should clear reason
	rr = doJSON(t, server, http.MethodPost, "/readyz", map[string]any{
		"ready":  true,
		"reason": "should be cleared",
	})
	env = assertJSONResponse(t, rr, http.StatusOK, true, "updated")

	if err := json.Unmarshal(env.Data, &postData); err != nil {
		t.Fatalf("failed to decode data: %v; data=%s", err, string(env.Data))
	}
	if postData.Ready != true || postData.Reason != "" {
		t.Fatalf("expected ready=true and empty reason, got %+v", postData)
	}

	// 5) GET again should be ready
	rr = doJSON(t, server, http.MethodGet, "/readyz", nil)
	env = assertJSONResponse(t, rr, http.StatusOK, true, "ready")
	assertNoData(t, env)
}

func TestReadyzPost_InvalidPayload(t *testing.T) {
	server := apphttp.NewRouter()

	// Unknown field should fail due to DisallowUnknownFields()
	rr := doJSON(t, server, http.MethodPost, "/readyz", map[string]any{
		"ready":   false,
		"reason":  "x",
		"unknown": "field",
	})

	_ = assertJSONResponse(t, rr, http.StatusBadRequest, false, "invalid payload")
}
