//go:build integration

package testutil

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// execCommand is a thin wrapper so tests can substitute it if needed.
var execCommand = exec.CommandContext

// repoRoot returns the absolute path to the repository root (apexmemory/).
func repoRoot() string {
	// This file is at backend/internal/testutil/container.go
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "..", "..", "..")
}

// StartContainer spins up a disposable PostgreSQL container with migrations
// and app_code applied. Call the returned cleanup function in TestMain after
// tests complete.
func StartContainer(ctx context.Context) (cleanup func(), err error) {
	root := repoRoot()

	const (
		superUser = "postgres"
		superPass = "postgres"
		dbName    = "apexmemory_test"
		migrator  = "migrator"
		migratorP = "migrator"
		appUser   = "appuser"
		appUserP  = "appuser"
		appSchema = "app"
		codeSchema = "app_code"
	)

	container, err := postgres.Run(ctx,
		"postgres:18beta2",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(superUser),
		postgres.WithPassword(superPass),
		postgres.WithInitScripts(filepath.Join(root, "backend", "db", "init", "01-local-setup.sh")),
		testcontainers.WithEnv(map[string]string{
			"PG_MIGRATOR_USER":     migrator,
			"PG_MIGRATOR_PASSWORD": migratorP,
			"PG_APP_USER":          appUser,
			"PG_APP_PASSWORD":      appUserP,
			"PG_APP_SCHEMA":        appSchema,
			"PG_APP_CODE_SCHEMA":   codeSchema,
		}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("start postgres container: %w", err)
	}

	cleanupFn := func() {
		_ = container.Terminate(context.Background())
	}

	// Get connection string for migrator to run migrations
	host, err := container.Host(ctx)
	if err != nil {
		cleanupFn()
		return nil, fmt.Errorf("get container host: %w", err)
	}
	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		cleanupFn()
		return nil, fmt.Errorf("get container port: %w", err)
	}

	migratorDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		migrator, migratorP, host, port.Port(), dbName)

	// Run tern migrations first (migrations are self-contained and create any
	// functions they need inline), then install app_code to bring functions
	// to their current version. This matches the production deploy order.
	if err := runTernMigrate(ctx, root, migratorDSN, appSchema); err != nil {
		cleanupFn()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	// Install app_code (latest version of functions/triggers)
	if err := installAppCode(ctx, container, root, superUser, dbName, migrator, appSchema, codeSchema); err != nil {
		cleanupFn()
		return nil, fmt.Errorf("install app_code: %w", err)
	}

	// Grant TRUNCATE to app user (needed for test isolation via TruncateTables).
	// The init script's DEFAULT PRIVILEGES only cover SELECT/INSERT/UPDATE/DELETE.
	_, _, err = container.Exec(ctx, []string{
		"psql", "-U", superUser, "-d", dbName, "-c",
		fmt.Sprintf("GRANT TRUNCATE ON ALL TABLES IN SCHEMA %s TO %s;", appSchema, appUser),
	})
	if err != nil {
		cleanupFn()
		return nil, fmt.Errorf("grant truncate: %w", err)
	}

	// Set the DSN for the app user (used by tests)
	appDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		appUser, appUserP, host, port.Port(), dbName)
	SetDSN(appDSN)

	return cleanupFn, nil
}

// runTernMigrate executes tern migrate using the host-installed tern binary.
func runTernMigrate(ctx context.Context, root, dsn, appSchema string) error {
	migrDir := filepath.Join(root, "backend", "db", "migrations")

	// Parse DSN to extract components for tern env vars
	// tern uses its own config format; we set env vars and use the config file
	// But simpler: just call tern with a connection string via env override
	// tern.conf uses template env vars, so we set them directly.

	// Extract host:port from DSN
	// DSN format: postgres://user:pass@host:port/dbname?sslmode=disable
	parts := strings.SplitAfter(dsn, "@")
	hostPortDB := parts[1] // host:port/dbname?sslmode=disable
	hostPort := strings.Split(hostPortDB, "/")[0]
	hostParts := strings.Split(hostPort, ":")
	pgHost := hostParts[0]
	pgPort := hostParts[1]

	dbParts := strings.SplitN(strings.Split(hostPortDB, "/")[1], "?", 2)
	pgDB := dbParts[0]

	// Extract user:pass
	userPass := strings.TrimPrefix(strings.Split(dsn, "@")[0], "postgres://")
	userParts := strings.SplitN(userPass, ":", 2)
	pgUser := userParts[0]
	pgPass := userParts[1]

	env := map[string]string{
		"PG_HOST":              pgHost,
		"PG_PORT":              pgPort,
		"PG_DATABASE":          pgDB,
		"PG_MIGRATOR_USER":     pgUser,
		"PG_MIGRATOR_PASSWORD": pgPass,
		"PG_APP_SCHEMA":        appSchema,
	}

	// Build env slice
	envSlice := os.Environ()
	for k, v := range env {
		envSlice = append(envSlice, k+"="+v)
	}

	confPath := filepath.Join(root, "backend", "tern.conf")

	// Use exec.Command to run tern
	cmd := execCommand(ctx, "tern", "migrate",
		"--config", confPath,
		"--migrations", migrDir,
	)
	cmd.Env = envSlice
	cmd.Dir = filepath.Join(root, "backend")

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("tern migrate failed: %w\noutput: %s", err, string(output))
	}
	return nil
}

// installAppCode runs the app_code SQL functions directly via psql in the container.
func installAppCode(ctx context.Context, container *postgres.PostgresContainer, root, superUser, dbName, migrator, appSchema, codeSchema string) error {
	// Read each SQL file and execute it
	codeDir := filepath.Join(root, "backend", "db", "code", "funcs")
	entries, err := os.ReadDir(codeDir)
	if err != nil {
		return fmt.Errorf("read code dir: %w", err)
	}

	// Ensure app_code schema exists, owned by migrator
	_, _, err = container.Exec(ctx, []string{
		"psql", "-U", superUser, "-d", dbName, "-c",
		fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s AUTHORIZATION %s", codeSchema, migrator),
	})
	if err != nil {
		return fmt.Errorf("create code schema: %w", err)
	}

	// Execute each function file as MIGRATOR (so migrator owns the functions,
	// matching production behavior where tern runs as migrator)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		sqlBytes, err := os.ReadFile(filepath.Join(codeDir, entry.Name()))
		if err != nil {
			return fmt.Errorf("read %s: %w", entry.Name(), err)
		}

		_, _, execErr := container.Exec(ctx, []string{
			"psql", "-U", migrator, "-d", dbName, "-c", string(sqlBytes),
		})
		if execErr != nil {
			return fmt.Errorf("exec %s: %w", entry.Name(), execErr)
		}
	}

	// Grant execute on all functions in app_code to app user
	_, _, err = container.Exec(ctx, []string{
		"psql", "-U", superUser, "-d", dbName, "-c",
		fmt.Sprintf("GRANT USAGE ON SCHEMA %s TO appuser; GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA %s TO appuser;", codeSchema, codeSchema),
	})
	if err != nil {
		return fmt.Errorf("grant app_code privileges: %w", err)
	}

	return nil
}
