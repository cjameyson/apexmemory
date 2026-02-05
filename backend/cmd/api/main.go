package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"apexmemory.ai/internal/app"
	"apexmemory.ai/internal/storage"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	var cfg app.Config
	var trustedProxies string
	var deprecatedMaxIdleConns int

	flag.IntVar(&cfg.Port, "port", 4000, "API server port")
	flag.StringVar(&cfg.Env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.DB.DSN, "dsn", "", "PostgreSQL DSN (or set DATABASE_URL env var)")
	flag.IntVar(&cfg.DB.MaxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.DB.MinIdleConns, "db-min-idle-conns", 0, "PostgreSQL min idle connections (pgxpool MinConns)")
	flag.IntVar(&deprecatedMaxIdleConns, "db-max-idle-conns", 0, "DEPRECATED: use -db-min-idle-conns")
	flag.DurationVar(&cfg.DB.MaxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")
	flag.StringVar(&trustedProxies, "trusted-proxies", "", "Comma-separated list of trusted proxy IPs/CIDRs (or set TRUSTED_PROXIES env var)")
	flag.Parse()

	// DSN from env var takes precedence if flag not provided
	if cfg.DB.DSN == "" {
		cfg.DB.DSN = os.Getenv("DATABASE_URL")
	}
	if cfg.DB.DSN == "" {
		cfg.DB.DSN = os.Getenv("PG_APP_DSN")
	}
	if cfg.DB.DSN == "" {
		fmt.Fprintf(os.Stderr, "Error: -dsn flag or DATABASE_URL/PG_APP_DSN environment variable is required\n")
		return fmt.Errorf("missing database DSN")
	}

	// Trusted proxies from env var takes precedence if flag not provided
	if trustedProxies == "" {
		trustedProxies = os.Getenv("TRUSTED_PROXIES")
	}
	if trustedProxies != "" {
		for _, proxy := range strings.Split(trustedProxies, ",") {
			proxy = strings.TrimSpace(proxy)
			if proxy != "" {
				cfg.TrustedProxies = append(cfg.TrustedProxies, proxy)
			}
		}
	}

	if cfg.DB.MinIdleConns == 0 && deprecatedMaxIdleConns > 0 {
		cfg.DB.MinIdleConns = deprecatedMaxIdleConns
	}
	if deprecatedMaxIdleConns > 0 {
		fmt.Fprintln(os.Stderr, "warning: -db-max-idle-conns is deprecated; use -db-min-idle-conns")
	}

	// Storage path from env var (default: ./data/assets)
	cfg.StoragePath = os.Getenv("STORAGE_PATH")
	if cfg.StoragePath == "" {
		cfg.StoragePath = "./data/assets"
	}

	application := app.New(cfg)

	store, err := storage.NewLocalStorage(cfg.StoragePath)
	if err != nil {
		application.Logger.Error("failed to initialize storage",
			"error", err,
			"error_code", "STORAGE_INIT_FAILED",
			"path", cfg.StoragePath,
			"retryable", false,
			"remediation_hint", "Check that the storage path is writable",
		)
		return err
	}
	application.Storage = store

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.ConnectDB(ctx); err != nil {
		application.Logger.Error("failed to connect to database", "error", err)
		return err
	}
	defer application.CloseDB()

	// Start background jobs (session cleanup, etc.)
	application.StartBackgroundJobs()
	defer application.Shutdown()

	mux := application.Routes()

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(application.Logger.Handler(), slog.LevelError),
	}

	// Graceful shutdown: listen for SIGINT/SIGTERM
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	errCh := make(chan error, 1)

	// Start server in goroutine
	go func() {
		application.Logger.Info("starting server", "addr", srv.Addr, "env", cfg.Env)
		errCh <- srv.ListenAndServe()
	}()

	// Wait for shutdown signal or server error
	select {
	case sig := <-shutdownChan:
		application.Logger.Info("received shutdown signal", "signal", sig.String())
	case err := <-errCh:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			application.Logger.Error("server error", "error", err)
			return err
		}
		// Server closed without signal (unusual but possible)
		application.Logger.Info("server stopped unexpectedly")
		return nil
	}

	// Allow force-quit on second signal
	forceQuitChan := make(chan os.Signal, 1)
	signal.Notify(forceQuitChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-forceQuitChan
		application.Logger.Warn("received second signal, forcing exit", "signal", sig.String())
		os.Exit(1)
	}()

	// Give in-flight requests time to complete
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	application.Logger.Info("shutting down server", "timeout", "30s")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		application.Logger.Error("graceful shutdown failed", "error", err)
		// Don't return yet - still need to drain errCh and run cleanup
	}

	// Wait for ListenAndServe to return (it will after Shutdown completes)
	if err := <-errCh; err != nil && !errors.Is(err, http.ErrServerClosed) {
		application.Logger.Error("server error during shutdown", "error", err)
		return err
	}

	application.Logger.Info("server stopped gracefully")
	return nil
}
