package app

import (
	"os"
	"testing"

	"apexmemory.ai/internal/testutil"
)

// TestMain runs before all tests in the app package.
// It sets up the shared test database connection and ensures cleanup.
func TestMain(m *testing.M) {
	// Check if test database is available
	if testutil.TestDSN() == "" {
		// Run tests anyway - individual tests will skip if they need DB
		os.Exit(m.Run())
	}

	// Verify we can connect to test database
	_, err := testutil.Pool()
	if err != nil {
		// Print error but don't fail - let individual tests skip
		os.Stderr.WriteString("warning: test database not available: " + err.Error() + "\n")
		os.Exit(m.Run())
	}

	// Run all tests
	code := m.Run()

	// Cleanup
	testutil.Close()

	os.Exit(code)
}
