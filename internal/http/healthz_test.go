package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

func TestHealthz_GeneratesRequestID(t *testing.T) {
	fakeReadiness := FakeReadiness{ready: true}
	server := apphttp.NewRouter(&fakeReadiness)

	rr := doJSON(t, server, http.MethodGet, "/healthz", nil)
	_ = assertJSONResponse(t, rr, http.StatusOK, true, "healthy")

	if rid := rr.Header().Get("X-Request-Id"); rid == "" {
		t.Fatalf("expected X-Request-Id header to be set")
	}
}

func TestHealthz_PreservesIncomingRequestID(t *testing.T) {
	fakeReadiness := FakeReadiness{ready: true}
	server := apphttp.NewRouter(&fakeReadiness)

	requestID := "alkindi"

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set("X-Request-Id", requestID)

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	_ = assertJSONResponse(t, rr, http.StatusOK, true, "healthy")

	if rid := rr.Header().Get("X-Request-Id"); rid != requestID {
		t.Fatalf("expected X-Request-Id to be preserved, got %q", rid)
	}
}
