package app

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"apexmemory.ai/internal/db"
	"apexmemory.ai/internal/testutil"
)

// testApp creates a test application connected to the test database.
// It truncates all tables before returning to ensure test isolation.
func testApp(t *testing.T) *Application {
	t.Helper()

	pool, err := testutil.Pool()
	if err != nil {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}

	// Truncate tables for test isolation
	testutil.TruncateTablesT(t, context.Background(), pool)

	config := Config{
		Port: 4000,
		Env:  "test",
	}
	config.DB.DSN = testutil.TestDSN()
	config.DB.MaxOpenConns = 5
	config.DB.MinIdleConns = 2

	app := New(config)
	app.RateLimiters = NewTestRateLimiters()
	app.DB = pool
	app.Queries = db.New(pool)

	t.Cleanup(func() {
		app.RateLimiters.Stop()
	})

	return app
}

// jsonRequest creates a test HTTP request with JSON body.
func jsonRequest(t *testing.T, method, path string, body any) *http.Request {
	t.Helper()

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		bodyReader = strings.NewReader(string(jsonBody))
	}

	req := httptest.NewRequest(method, path, bodyReader)
	req.Header.Set("Content-Type", "application/json")
	return req
}

// decodeResponse decodes a JSON response body.
func decodeResponse(t *testing.T, rr *httptest.ResponseRecorder, dst any) {
	t.Helper()

	if err := json.Unmarshal(rr.Body.Bytes(), dst); err != nil {
		t.Fatalf("failed to decode response: %v\nBody: %s", err, rr.Body.String())
	}
}

func TestHealthcheckHandler(t *testing.T) {
	app := testApp(t)

	req := httptest.NewRequest("GET", "/v1/healthcheck", nil)
	rr := httptest.NewRecorder()

	app.HealthcheckHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var resp map[string]any
	decodeResponse(t, rr, &resp)

	if resp["status"] != "available" {
		t.Errorf("expected status 'available', got %v", resp["status"])
	}
}
