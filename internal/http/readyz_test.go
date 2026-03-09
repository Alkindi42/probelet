package http_test

import (
	"encoding/json"
	"net/http"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

// #########
// /readyz #
// #########

func TestReadyz_StateTransitions_POST(t *testing.T) {
	fakeReadiness := FakeReadiness{ready: true}
	server := apphttp.NewRouter(&fakeReadiness)

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

func TestReadyzPost_InvalidPayload_POST(t *testing.T) {
	fakeReadiness := FakeReadiness{ready: true}
	server := apphttp.NewRouter(&fakeReadiness)

	// Unknown field should fail due to DisallowUnknownFields()
	rr := doJSON(t, server, http.MethodPost, "/readyz", map[string]any{
		"ready":   false,
		"reason":  "x",
		"unknown": "field",
	})

	_ = assertJSONResponse(t, rr, http.StatusBadRequest, false, "invalid payload")
}
