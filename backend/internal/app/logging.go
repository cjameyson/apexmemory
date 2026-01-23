package app

import (
	"context"
	"log/slog"
)

const anonymousUserID = "anonymous"

// RequestLogContext stores request-scoped fields that may be updated mid-request.
type RequestLogContext struct {
	RequestID string
	TraceID   string
	UserID    string
	Method    string
	Path      string
}

// GetLogContext retrieves the request log context.
func GetLogContext(ctx context.Context) *RequestLogContext {
	if logCtx, ok := ctx.Value(logContextKey).(*RequestLogContext); ok {
		return logCtx
	}
	return nil
}

// WithLogContext stores the request log context.
func WithLogContext(ctx context.Context, logCtx *RequestLogContext) context.Context {
	return context.WithValue(ctx, logContextKey, logCtx)
}

// GetLogger retrieves the request-scoped logger from context.
// Falls back to slog.Default() if no logger is stored.
func GetLogger(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// WithLogger stores a logger in context.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
