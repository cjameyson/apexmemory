//go:build integration

package app

import (
	"context"
	"os"
	"testing"

	"apexmemory.ai/internal/testutil"
)

// TestMain runs before all tests in the app package.
// It starts a disposable PostgreSQL container, applies migrations and app_code,
// then runs all tests.
func TestMain(m *testing.M) {
	ctx := context.Background()

	cleanup, err := testutil.StartContainer(ctx)
	if err != nil {
		os.Stderr.WriteString("failed to start test container: " + err.Error() + "\n")
		os.Exit(1)
	}

	// Verify we can connect
	if _, err := testutil.Pool(); err != nil {
		os.Stderr.WriteString("failed to connect to test database: " + err.Error() + "\n")
		cleanup()
		os.Exit(1)
	}

	code := m.Run()

	testutil.Close()
	cleanup()

	os.Exit(code)
}
