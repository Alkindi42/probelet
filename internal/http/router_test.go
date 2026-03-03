package http_test

import (
	"encoding/json"
	"net/http"
	"runtime"
	"strconv"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

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

// #############
// /stress/cpu #
// #############

func TestStressCPU(t *testing.T) {
	type stressData struct {
		Cores    int    `json:"cores"`
		Duration string `json:"duration"`
	}

	maxCores := runtime.GOMAXPROCS(0)

	tests := []struct {
		name        string
		path        string
		wantStatus  int
		wantOK      bool
		wantMessage string
		wantData    *stressData
	}{
		{
			name:        "missing_duration_400",
			path:        "/stress/cpu",
			wantStatus:  http.StatusBadRequest,
			wantOK:      false,
			wantMessage: "required",
			wantData:    nil,
		},
		{
			name:        "invalid_duration_format_400",
			path:        "/stress/cpu?duration=abc",
			wantStatus:  http.StatusBadRequest,
			wantOK:      false,
			wantMessage: "invalid",
			wantData:    nil,
		},
		{
			name:        "zero_duration_400",
			path:        "/stress/cpu?duration=0s",
			wantStatus:  http.StatusBadRequest,
			wantOK:      false,
			wantMessage: "duration must be greater than 0",
			wantData:    nil,
		},
		{
			name:        "negative_duration_400",
			path:        "/stress/cpu?duration=-1s",
			wantStatus:  http.StatusBadRequest,
			wantOK:      false,
			wantMessage: "duration must be greater than 0",
			wantData:    nil,
		},
		{
			name:        "duration_above_max_400",
			path:        "/stress/cpu?duration=2m1s",
			wantStatus:  http.StatusBadRequest,
			wantOK:      false,
			wantMessage: "duration must be <= 2m0s",
			wantData:    nil,
		},
		{
			name:        "invalid_cores_nonint_400",
			path:        "/stress/cpu?duration=10ms&cores=abc",
			wantStatus:  http.StatusBadRequest,
			wantOK:      false,
			wantMessage: "cores must be a number between 1 and " + strconv.Itoa(maxCores) + " or equal to 'max'",
			wantData:    nil,
		},
		{
			name:        "invalid_cores_zero_400",
			path:        "/stress/cpu?duration=10ms&cores=0",
			wantStatus:  http.StatusBadRequest,
			wantOK:      false,
			wantMessage: "cores must be a number between 1 and " + strconv.Itoa(maxCores) + " or equal to 'max'",
			wantData:    nil,
		},
		{
			name:        "invalid_cores_too_large_400",
			path:        "/stress/cpu?duration=10ms&cores=999999",
			wantStatus:  http.StatusBadRequest,
			wantOK:      false,
			wantMessage: "cores must be a number between 1 and " + strconv.Itoa(maxCores) + " or equal to 'max'",
			wantData:    nil,
		},
		// ---- success cases ----
		{
			name:        "duration_only_defaults_cores_200",
			path:        "/stress/cpu?duration=10ms",
			wantStatus:  http.StatusOK,
			wantOK:      true,
			wantMessage: "done",
			wantData: &stressData{
				Cores:    1,
				Duration: "10ms",
			},
		},
		{
			name:        "cores_max_duration_200",
			path:        "/stress/cpu?duration=10ms&cores=max",
			wantStatus:  http.StatusOK,
			wantOK:      true,
			wantMessage: "done",
			wantData: &stressData{
				Cores:    maxCores,
				Duration: "10ms",
			},
		},
		{
			name:        "cores_explicit_duration_200",
			path:        "/stress/cpu?duration=10ms&cores=2",
			wantStatus:  http.StatusOK,
			wantOK:      true,
			wantMessage: "done",
			wantData: &stressData{
				Cores:    2,
				Duration: "10ms",
			},
		},
	}

	server := apphttp.NewRouter()

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := doJSON(t, server, http.MethodGet, tc.path, nil)
			env := assertJSONResponse(t, rr, tc.wantStatus, tc.wantOK, tc.wantMessage)

			if tc.wantData == nil {
				assertNoData(t, env)
				return
			}

			var got stressData
			if err := json.Unmarshal(env.Data, &got); err != nil {
				t.Fatalf("failed to decode data: %v; data=%s", err, string(env.Data))
			}

			if got.Cores != tc.wantData.Cores {
				t.Fatalf("expected data.cores=%d, got %d", tc.wantData.Cores, got.Cores)
			}
			if got.Duration != tc.wantData.Duration {
				t.Fatalf("expected data.duration=%q, got %q", tc.wantData.Duration, got.Duration)
			}
		})
	}
}
