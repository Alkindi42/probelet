package http_test

import (
	"encoding/json"
	"net/http"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

func TestStressMemory(t *testing.T) {
	type stressData struct {
		Size     string `json:"size"`
		Bytes    int64  `json:"bytes"`
		Duration string `json:"duration"`
	}

	tests := []struct {
		name        string
		path        string
		wantStatus  int
		wantOK      bool
		wantMessage string
		wantData    *stressData
	}{
		{"duration_required", "/stress/memory?size=1Mi", http.StatusBadRequest, false, "required", nil},
		{"duration_invalid", "/stress/memory?duration=abc&size=1Mi", http.StatusBadRequest, false, "invalid", nil},
		{"duration_zero", "/stress/memory?duration=0s&size=1Mi", http.StatusBadRequest, false, "greater than 0", nil},
		{"duration_too_large", "/stress/memory?duration=10m&size=1Mi", http.StatusBadRequest, false, "<=", nil},
		{"size_required", "/stress/memory?duration=5ms", http.StatusBadRequest, false, "size is required", nil},
		{"size_invalid", "/stress/memory?duration=5ms&size=bad", http.StatusBadRequest, false, "invalid size", nil},
		{"size_zero", "/stress/memory?duration=5ms&size=0", http.StatusBadRequest, false, "greater than 0", nil},
		{"size_too_large", "/stress/memory?duration=5ms&size=2Gi", http.StatusBadRequest, false, "size must be", nil},
		{
			"ok",
			"/stress/memory?duration=5ms&size=1Mi",
			http.StatusOK,
			true,
			"done",
			&stressData{Size: "1Mi", Bytes: 1 << 20, Duration: "5ms"},
		},
	}

	fakeReadiness := FakeReadiness{ready: true}
	server := apphttp.NewRouter(&fakeReadiness)

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

			if got.Size != tc.wantData.Size {
				t.Fatalf("expected data.size=%q, got %q", tc.wantData.Size, got.Size)
			}
			if got.Bytes != tc.wantData.Bytes {
				t.Fatalf("expected data.bytes=%d, got %d", tc.wantData.Bytes, got.Bytes)
			}
			if got.Duration != tc.wantData.Duration {
				t.Fatalf("expected data.duration=%q, got %q", tc.wantData.Duration, got.Duration)
			}
		})
	}
}
