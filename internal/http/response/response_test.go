package response_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alkindi42/probelet/internal/http/response"
)

type envelope struct {
	OK      bool            `json:"ok"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func decodeEnvelope(t *testing.T, body string) envelope {
	t.Helper()
	var env envelope
	if err := json.Unmarshal([]byte(body), &env); err != nil {
		t.Fatalf("failed to decode json: %v\nbody=%s", err, body)
	}
	return env
}

func TestJSON_SetsHeadersStatusAndBody(t *testing.T) {
	rr := httptest.NewRecorder()

	payload := map[string]int{"code": 201}
	response.JSON(rr, http.StatusCreated, "Created", payload)

	if rr.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("expected Content-Type application/json; charset=utf-8, got %q", ct)
	}

	if cc := rr.Header().Get("Cache-Control"); cc != "no-store" {
		t.Fatalf("expected Cache-Control no-store, got %q", cc)
	}

	if cc := rr.Header().Get("X-Content-Type-Options"); cc != "nosniff" {
		t.Fatalf("expected X-Content-Type-Options nosniff, got %q", cc)
	}

	env := decodeEnvelope(t, rr.Body.String())
	if env.OK != true {
		t.Fatalf("expected ok=true, got %v", env.OK)
	}
	if env.Message != "Created" {
		t.Fatalf("expected message %q, got %q", "Created", env.Message)
	}

	var got map[string]int
	if err := json.Unmarshal(env.Data, &got); err != nil {
		t.Fatalf("failed to decode data: %v", err)
	}
	if got["code"] != 201 {
		t.Fatalf("expected data.code=201, got %v", got["code"])
	}
}

func TestJSON_OmitsDataWhenNil(t *testing.T) {
	rr := httptest.NewRecorder()

	response.JSON(rr, http.StatusOK, "healthy", nil)

	env := decodeEnvelope(t, rr.Body.String())
	if env.OK != true {
		t.Fatalf("expected ok=true, got %v", env.OK)
	}
	if env.Message != "healthy" {
		t.Fatalf("expected message %q, got %q", "healthy", env.Message)
	}

	if len(env.Data) != 0 {
		t.Fatalf("expected data to be omitted, got %s", string(env.Data))
	}
}

func TestJSONError_SetsHeadersStatusAndBody(t *testing.T) {
	rr := httptest.NewRecorder()

	response.JSONError(rr, http.StatusBadRequest, "invalid status code")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("expected Content-Type application/json; charset=utf-8, got %q", ct)
	}

	if cc := rr.Header().Get("Cache-Control"); cc != "no-store" {
		t.Fatalf("expected Cache-Control no-store, got %q", cc)
	}

	if cc := rr.Header().Get("X-Content-Type-Options"); cc != "nosniff" {
		t.Fatalf("expected X-Content-Type-Options nosniff, got %q", cc)
	}

	env := decodeEnvelope(t, rr.Body.String())
	if env.OK != false {
		t.Fatalf("expected ok=false, got %v", env.OK)
	}
	if env.Message != "invalid status code" {
		t.Fatalf("expected message %q, got %q", "invalid status code", env.Message)
	}

	if len(env.Data) != 0 {
		t.Fatalf("expected data to be omitted, got %s", string(env.Data))
	}
}
