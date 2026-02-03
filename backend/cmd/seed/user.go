package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"apexmemory.ai/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// Test agent credentials - fixed values for consistent testing
var (
	testAgentEmail    = "test-agent@apexmemory.ai"
	testAgentPassword = "--DozensOfUs123!"
	testAgentUsername = "test-agent"
	testAgentDisplay  = "Test Agent"
	testAgentUUID     = uuid.MustParse("019501a0-0000-7000-8000-000000000001")
)

func runUserCmd(args []string) {
	fs := flag.NewFlagSet("user", flag.ExitOnError)
	var dsn string
	fs.StringVar(&dsn, "dsn", "", "PostgreSQL DSN (or set DATABASE_URL env var)")
	fs.Usage = func() {
		fmt.Println(`seed user - Create/verify the test agent account

Usage:
  seed user [flags]

Flags:`)
		fs.PrintDefaults()
		fmt.Println(`
Test agent account credentials:
  Email:    test-agent@apexmemory.ai
  Password: --DozensOfUs123!
  Username: test-agent
  User ID:  019501a0-0000-7000-8000-000000000001`)
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if err := createTestUser(dsn); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func createTestUser(flagDSN string) error {
	dsn, err := resolveDSN(flagDSN)
	if err != nil {
		return err
	}

	application, ctx, cancel, err := connectApp(dsn)
	if err != nil {
		return err
	}
	defer cancel()
	defer application.CloseDB()

	// Hash the password
	hashedPassword, err := application.HashPassword(testAgentPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	// Create or update the user
	user, err := application.Queries.UpsertTestUser(ctx, db.UpsertTestUserParams{
		ID:          testAgentUUID,
		Email:       testAgentEmail,
		Username:    testAgentUsername,
		DisplayName: pgtype.Text{String: testAgentDisplay, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}

	// Create or update the auth identity
	err = application.Queries.UpsertTestAuthIdentity(ctx, db.UpsertTestAuthIdentityParams{
		UserID:         user.ID,
		ProviderUserID: user.ID.String(),
		Email:          pgtype.Text{String: testAgentEmail, Valid: true},
		PasswordHash:   pgtype.Text{String: hashedPassword, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("upsert auth identity: %w", err)
	}

	slog.Info("test agent account ready",
		"user_id", user.ID,
		"email", user.Email,
		"username", user.Username,
	)

	return nil
}
