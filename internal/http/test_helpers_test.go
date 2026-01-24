package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// envelope mirrors the JSON response envelope returned by the API.
type envelope struct {
	OK      bool            `json:"ok"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
}

// decodeEnvelope unmarshals an API response body into an envelope.
func decodeEnvelope(t *testing.T, body []byte) envelope {
	t.Helper()

	var env envelope
	if err := json.Unmarshal(body, &env); err != nil {
		t.Fatalf("failed to decode json: %v\nbody=%s", err, string(body))
	}
	return env
}

// doJSON sends an HTTP request to server with an optional JSON payload and returns the recorded response.
func doJSON(t *testing.T, server http.Handler, method, path string, payload any) *httptest.ResponseRecorder {
	t.Helper()

	var body *bytes.Reader
	if payload == nil {
		body = bytes.NewReader(nil)
	} else {
		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatalf("failed to marshal payload: %v", err)
		}
		body = bytes.NewReader(b)
	}

	req := httptest.NewRequest(method, path, body)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rr := httptest.NewRecorder()

	server.ServeHTTP(rr, req)
	return rr
}

// assertJSONResponse asserts status and envelope fields, and returns the decoded envelope.
func assertJSONResponse(t *testing.T, rr *httptest.ResponseRecorder, wantStatus int, wantOK bool, wantMsg string) envelope {
	t.Helper()

	if rr.Code != wantStatus {
		t.Fatalf("expected status %d, got %d; body=%s", wantStatus, rr.Code, rr.Body.String())
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json; charset=utf-8" {
		t.Fatalf("expected Content-Type application/json; charset=utf-8, got %q", ct)
	}

	env := decodeEnvelope(t, rr.Body.Bytes())
	if env.OK != wantOK {
		t.Fatalf("expected ok=%v, got %v; body=%s", wantOK, env.OK, rr.Body.String())
	}
	if env.Message != wantMsg && !strings.Contains(env.Message, wantMsg) {
		t.Fatalf("expected message=%q, got %q; body=%s", wantMsg, env.Message, rr.Body.String())
	}
	return env
}

// assertNoData asserts that the envelope has no data field.
func assertNoData(t *testing.T, env envelope) {
	t.Helper()

	if len(env.Data) != 0 {
		t.Fatalf("expected no data, got %s", string(env.Data))
	}
}
