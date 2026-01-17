package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Alkindi42/probelet/internal/http/handlers"
)

func TestHealthz(t *testing.T) {
	h := handlers.NewHealthzHandler()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	ct := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected json content-type, got %q", ct)
	}

	cache := rr.Header().Get("Cache-Control")
	if cache != "no-store" {
		t.Fatalf("expected Cache-Control no-store, got %q", cache)
	}

	if cc := rr.Header().Get("X-Content-Type-Options"); cc != "nosniff" {
		t.Fatalf("expected X-Content-Type-Options nosniff, got %q", cc)
	}

	if !strings.Contains(rr.Body.String(), `"message":"healthy"`) {
		t.Fatalf("unexpected body: %s", rr.Body.String())
	}
}
