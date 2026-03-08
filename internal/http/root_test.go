package http_test

import (
	"encoding/json"
	"net/http"
	"testing"

	apphttp "github.com/Alkindi42/probelet/internal/http"
)

func TestRoot_GET(t *testing.T) {
	fakeReadiness := FakeReadiness{ready: true}
	server := apphttp.NewRouter(&fakeReadiness)

	rr := doJSON(t, server, http.MethodGet, "/", nil)
	env := assertJSONResponse(t, rr, http.StatusOK, true, "probelet")

	if rid := rr.Header().Get("X-Request-Id"); rid == "" {
		t.Fatalf("expected X-Request-Id header to be present")
	}

	var data struct {
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		BuildDate string `json:"build_date"`
		Docs      string `json:"docs"`
		OpenAPI   string `json:"openapi"`
	}

	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("invalid data json: %v", err)
	}

	if data.Docs != "/docs" {
		t.Fatalf("expected docs=/docs, got %q", data.Docs)
	}
	if data.OpenAPI != "/openapi.yaml" {
		t.Fatalf("expected openapi=/openapi.yaml, got %q", data.OpenAPI)
	}
	if data.Version == "" {
		t.Fatal("expected version to be set")
	}
	if data.Commit == "" {
		t.Fatal("expected commit to be set")
	}
	if data.BuildDate == "" {
		t.Fatal("expected build_date to be set")
	}
}
