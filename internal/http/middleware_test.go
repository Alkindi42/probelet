package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

func TestRequireToken(t *testing.T) {
	t.Run("missing token returns unauthorized", func(t *testing.T) {
		called := false

		handler := apphttp.RequireToken("secret-token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
		}
		if called {
			t.Fatal("wrapped handler should not have been called")
		}
	})

	t.Run("invalid token returns unauthorized", func(t *testing.T) {
		called := false

		handler := apphttp.RequireToken("secret-token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("X-Probelet-Token", "wrong-token")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusUnauthorized {
			t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
		}
		if called {
			t.Fatal("wrapped handler should not have been called")
		}
	})

	t.Run("valid token calls wrapped handler", func(t *testing.T) {
		called := false

		handler := apphttp.RequireToken("secret-token", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			w.WriteHeader(http.StatusOK)
		}))

		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.Header.Set("X-Probelet-Token", "secret-token")
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
		if !called {
			t.Fatal("wrapped handler should have been called")
		}
	})
}
