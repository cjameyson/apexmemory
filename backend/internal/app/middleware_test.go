//go:build integration

package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestLogRequests_ValidRequestID_Accepted(t *testing.T) {
	app := testApp(t)

	validID := uuid.Must(uuid.NewV7()).String()

	handler := app.LogRequests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/v1/healthcheck", nil)
	req.Header.Set("X-Request-ID", validID)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	got := rr.Header().Get("X-Request-ID")
	if got != validID {
		t.Errorf("expected X-Request-ID %q, got %q", validID, got)
	}
}

func TestLogRequests_InvalidRequestID_Ignored(t *testing.T) {
	app := testApp(t)

	handler := app.LogRequests(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/v1/healthcheck", nil)
	req.Header.Set("X-Request-ID", "not-a-uuid")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	got := rr.Header().Get("X-Request-ID")
	if got == "not-a-uuid" {
		t.Error("expected invalid X-Request-ID to be replaced, but it was echoed back")
	}

	// Verify a valid UUID was generated instead
	if _, err := uuid.Parse(got); err != nil {
		t.Errorf("expected a valid UUID in X-Request-ID, got %q", got)
	}
}
