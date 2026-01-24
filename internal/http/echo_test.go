package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

type echoBody struct {
	Content     string `json:"content"`
	Bytes       int    `json:"bytes"`
	IsTruncated bool   `json:"is_truncated"`
}

type echoData struct {
	Method        string      `json:"method"`
	Path          string      `json:"path"`
	Query         url.Values  `json:"query"`
	RawQuery      string      `json:"raw_query"`
	Body          echoBody    `json:"body"`
	ContentLength int64       `json:"content_length"`
	Headers       http.Header `json:"headers"`
	Host          string      `json:"host"`
	Proto         string      `json:"proto"`
	UserAgent     string      `json:"user_agent"`
	ClientIP      string      `json:"client_ip"`
	RemoteAddr    string      `json:"remote_addr"`
	XRealIP       string      `json:"x_real_ip"`
	ForwardedFor  string      `json:"forwarded_for"`
}

// #######
// /echo #
// #######

func TestEcho_GET_Default(t *testing.T) {
	server := apphttp.NewRouter()

	rr := doJSON(t, server, http.MethodGet, "/echo?q=foo", nil)
	env := assertJSONResponse(t, rr, http.StatusOK, true, "echo")

	var data echoData
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("invalid data json: %v", err)
	}

	if data.Method != http.MethodGet {
		t.Fatalf("expected method GET, got %q", data.Method)
	}
	if data.Path != "/echo" {
		t.Fatalf("expected path /echo, got %q", data.Path)
	}
	if data.Query.Get("q") != "foo" {
		t.Fatalf("expected query.q=foo, got %q", data.Query.Get("q"))
	}

	// body defaults = CONTRACT
	if data.Body.Bytes != 0 || data.Body.Content != "" || data.Body.IsTruncated {
		t.Fatalf("unexpected body defaults: %+v", data.Body)
	}
}

func TestEcho_POST_BodyAndHeaders(t *testing.T) {
	server := apphttp.NewRouter()

	reqBody := bytes.NewBufferString(`{"foo":"bar"}`)
	req := httptest.NewRequest(http.MethodPost, "/echo", reqBody)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Foo", "bar")
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	req.Header.Set("X-Real-IP", "5.6.7.8")

	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	env := assertJSONResponse(t, rr, http.StatusOK, true, "echo")

	var data echoData
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("invalid data json: %v; data=%s", err, string(env.Data))
	}

	if data.Method != http.MethodPost {
		t.Fatalf("expected method=%q, got %q", http.MethodPost, data.Method)
	}
	if data.Body.Content != `{"foo":"bar"}` {
		t.Fatalf("expected body.content=%q, got %q", `{"foo":"bar"}`, data.Body.Content)
	}
	if data.Body.IsTruncated {
		t.Fatalf("expected body.is_truncated=false, got true")
	}
	if data.Headers.Get("X-Foo") != "bar" {
		t.Fatalf("expected header X-Foo=%q, got %q", "bar", data.Headers.Get("X-Foo"))
	}
	if data.ForwardedFor != "1.2.3.4" {
		t.Fatalf("expected forwarded_for=%q, got %q", "1.2.3.4", data.ForwardedFor)
	}
	if data.XRealIP != "5.6.7.8" {
		t.Fatalf("expected x_real_ip=%q, got %q", "5.6.7.8", data.XRealIP)
	}
}

func TestEcho_POST_Truncation(t *testing.T) {
	server := apphttp.NewRouter()

	const max = 64 << 10
	payload := bytes.Repeat([]byte("a"), max+1)

	req := httptest.NewRequest(http.MethodPost, "/echo", bytes.NewReader(payload))
	rr := httptest.NewRecorder()
	server.ServeHTTP(rr, req)

	env := assertJSONResponse(t, rr, http.StatusOK, true, "echo")

	var data echoData
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("invalid data json: %v", err)
	}

	if !data.Body.IsTruncated {
		t.Fatalf("expected body to be truncated")
	}
	if data.Body.Bytes != max {
		t.Fatalf("expected body.bytes=%d, got %d", max, data.Body.Bytes)
	}
}
