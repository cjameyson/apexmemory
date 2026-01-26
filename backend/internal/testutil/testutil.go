// Package testutil provides shared test infrastructure for integration tests.
// It manages a connection pool to the test database and provides helpers for
// test isolation through table truncation.
package testutil

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool     *pgxpool.Pool
	poolOnce sync.Once
	poolErr  error
)

// TestDSN returns the test database connection string.
// Returns empty string if TEST_DATABASE_URL is not set.
func TestDSN() string {
	return os.Getenv("TEST_DATABASE_URL")
}

// Pool returns a shared connection pool for the test database.
// The pool is created once and reused across all tests.
// Returns nil if TEST_DATABASE_URL is not set.
func Pool() (*pgxpool.Pool, error) {
	poolOnce.Do(func() {
		dsn := TestDSN()
		if dsn == "" {
			poolErr = fmt.Errorf("TEST_DATABASE_URL not set")
			return
		}

		ctx := context.Background()
		config, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			poolErr = fmt.Errorf("parse test DSN: %w", err)
			return
		}

		// Conservative pool settings for tests
		config.MaxConns = 10
		config.MinConns = 2

		pool, poolErr = pgxpool.NewWithConfig(ctx, config)
		if poolErr != nil {
			poolErr = fmt.Errorf("connect to test database: %w", poolErr)
			return
		}

		// Verify connection
		if err := pool.Ping(ctx); err != nil {
			pool.Close()
			pool = nil
			poolErr = fmt.Errorf("ping test database: %w", err)
		}
	})

	return pool, poolErr
}

// Close closes the shared connection pool.
// Should be called in TestMain after all tests complete.
func Close() {
	if pool != nil {
		pool.Close()
	}
}

// TruncateTables truncates all data tables in the app schema.
// Uses TRUNCATE ... CASCADE to handle foreign key constraints.
// This is called between tests to ensure isolation.
func TruncateTables(ctx context.Context, p *pgxpool.Pool) error {
	// Order matters for foreign keys, or use CASCADE
	// Using a single TRUNCATE with CASCADE is more efficient
	_, err := p.Exec(ctx, `
		TRUNCATE TABLE
			app.notebooks,
			app.user_sessions,
			app.auth_identities,
			app.users
		CASCADE
	`)
	return err
}

// TruncateTablesT is a test helper that truncates tables and fails the test on error.
func TruncateTablesT(t interface {
	Helper()
	Fatalf(format string, args ...any)
}, ctx context.Context, p *pgxpool.Pool) {
	t.Helper()
	if err := TruncateTables(ctx, p); err != nil {
		t.Fatalf("failed to truncate tables: %v", err)
	}
}

// MustPool returns the shared pool or panics.
// Useful in TestMain where we want to fail fast.
func MustPool() *pgxpool.Pool {
	p, err := Pool()
	if err != nil {
		panic(fmt.Sprintf("testutil: %v", err))
	}
	return p
}
