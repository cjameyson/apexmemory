package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"apexmemory.ai/internal/app"
)

// resolveDSN returns the database DSN from flags or environment.
func resolveDSN(flagDSN string) (string, error) {
	if flagDSN != "" {
		return flagDSN, nil
	}
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn, nil
	}
	if dsn := os.Getenv("PG_APP_DSN"); dsn != "" {
		return dsn, nil
	}
	return "", fmt.Errorf("missing database DSN: use -dsn flag or DATABASE_URL/PG_APP_DSN env var")
}

// connectApp creates and connects an Application instance.
func connectApp(dsn string) (*app.Application, context.Context, context.CancelFunc, error) {
	cfg := app.Config{Env: "development"}
	cfg.DB.DSN = dsn
	cfg.DB.MaxOpenConns = 5

	application := app.New(cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	if err := application.ConnectDB(ctx); err != nil {
		cancel()
		return nil, nil, nil, fmt.Errorf("connect db: %w", err)
	}

	return application, ctx, cancel, nil
}
