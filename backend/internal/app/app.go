package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"

	"apexmemory.ai/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

const Version = "1.0.0"

type Application struct {
	Config         Config
	Logger         *slog.Logger
	DB             *pgxpool.Pool
	Queries        *db.Queries
	RateLimiters   *RateLimiters
	BackgroundJobs *BackgroundJobs
	trustedProxies *TrustedProxyChecker
}

func New(config Config) *Application {
	// Configure logger based on environment
	var handler slog.Handler
	if config.Env == "development" || config.Env == "" {
		// Dev: text format with source locations for easier debugging
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	} else {
		// Prod/staging: JSON format for machine parsing, no source (too noisy)
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	logger := slog.New(handler).With(
		"service", "apexmemory-api",
		"version", Version,
		"env", config.Env,
	)

	proxyChecker, err := NewTrustedProxyChecker(config.TrustedProxies)
	if err != nil {
		// In production/staging, misconfigured proxies is a critical error
		// that could silently break rate limiting. Fail fast.
		if config.Env == "production" || config.Env == "staging" {
			logger.Error("invalid trusted proxy configuration", "error", err)
			os.Exit(1)
		}
		// In development, log warning and fall back to empty (trust nothing)
		logger.Warn("invalid trusted proxy configuration, falling back to RemoteAddr only", "error", err)
		proxyChecker, _ = NewTrustedProxyChecker(nil)
	}

	if len(config.TrustedProxies) > 0 {
		logger.Info("trusted proxies configured", "proxies", config.TrustedProxies)
	}

	return &Application{
		Config:         config,
		Logger:         logger,
		RateLimiters:   NewRateLimiters(),
		trustedProxies: proxyChecker,
	}
}

// Shutdown performs graceful cleanup of application resources.
func (app *Application) Shutdown() {
	if app.BackgroundJobs != nil {
		app.BackgroundJobs.Stop()
	}
	if app.RateLimiters != nil {
		app.RateLimiters.Stop()
	}
}

// StartBackgroundJobs initializes and starts background job processing.
func (app *Application) StartBackgroundJobs() {
	app.BackgroundJobs = NewBackgroundJobs(app)
	app.BackgroundJobs.Start()
	app.Logger.Info("background jobs started")
}

func (app *Application) ConnectDB(ctx context.Context) error {
	poolConfig, err := pgxpool.ParseConfig(app.Config.DB.DSN)
	if err != nil {
		return fmt.Errorf("failed to parse database DSN: %w", err)
	}

	if app.Config.DB.MaxOpenConns > 0 {
		poolConfig.MaxConns = int32(app.Config.DB.MaxOpenConns)
	}
	if app.Config.DB.MaxIdleTime > 0 {
		poolConfig.MaxConnIdleTime = app.Config.DB.MaxIdleTime
	}
	if app.Config.DB.MinIdleConns > 0 {
		minConns := int32(app.Config.DB.MinIdleConns)
		if poolConfig.MaxConns > 0 && minConns > poolConfig.MaxConns {
			minConns = poolConfig.MaxConns
		}
		poolConfig.MinConns = minConns
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	app.DB = pool
	app.Queries = db.New(pool)

	app.Logger.Info("database connection established", "dsn", maskDSN(app.Config.DB.DSN))
	return nil
}

func (app *Application) CloseDB() {
	if app.DB != nil {
		app.DB.Close()
		app.Logger.Info("database connection closed")
	}
}

func maskDSN(dsn string) string {
	u, err := url.Parse(dsn)
	if err != nil || u.User == nil {
		return "***"
	}
	if _, hasPass := u.User.Password(); hasPass {
		u.User = url.UserPassword(u.User.Username(), "***")
	}
	return u.String()
}

// HealthcheckHandler returns the health status of the API.
func (app *Application) HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	dbStatus := "available"
	if app.DB != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
		defer cancel()

		if err := app.DB.Ping(ctx); err != nil {
			dbStatus = "unavailable"
		}
	} else {
		dbStatus = "not configured"
	}

	status := "available"
	if dbStatus != "available" {
		status = "degraded"
	}

	data := map[string]any{
		"status": status,
		"system_info": map[string]any{
			"environment": app.Config.Env,
			"version":     Version,
		},
		"database": dbStatus,
	}

	app.RespondJSON(w, r, http.StatusOK, data)
}
