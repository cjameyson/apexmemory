package app

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (app *Application) LogRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate request ID (UUIDv7 for time-ordering)
		requestID := uuid.Must(uuid.NewV7()).String()
		if incomingID := r.Header.Get("X-Request-ID"); incomingID != "" {
			requestID = incomingID
		}

		traceID := extractTraceID(r)
		if traceID == "" {
			traceID = requestID
		}

		w.Header().Set("X-Request-ID", requestID)
		w.Header().Set("X-Trace-ID", traceID)

		logCtx := &RequestLogContext{
			RequestID: requestID,
			TraceID:   traceID,
			UserID:    anonymousUserID,
			Method:    r.Method,
			Path:      r.URL.Path,
		}

		logger := app.buildRequestLogger(logCtx)
		ctx := WithLogContext(r.Context(), logCtx)
		ctx = WithLogger(ctx, logger)
		r = r.WithContext(ctx)

		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)

		// Rebuild logger with final user_id (may have been set by Authenticate)
		finalLogCtx := GetLogContext(r.Context())
		if finalLogCtx == nil {
			finalLogCtx = logCtx
		}
		app.buildRequestLogger(finalLogCtx).Info("request_completed",
			"status", sw.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// statusWriter captures the HTTP status code for logging.
// Implements http.Flusher and http.Hijacker by delegation to support SSE and WebSockets.
type statusWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (sw *statusWriter) WriteHeader(code int) {
	if !sw.wroteHeader {
		sw.status = code
		sw.wroteHeader = true
	}
	sw.ResponseWriter.WriteHeader(code)
}

func (sw *statusWriter) Write(b []byte) (int, error) {
	if !sw.wroteHeader {
		sw.wroteHeader = true
	}
	return sw.ResponseWriter.Write(b)
}

// Flush implements http.Flusher for SSE support.
func (sw *statusWriter) Flush() {
	if f, ok := sw.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Hijack implements http.Hijacker for WebSocket support.
func (sw *statusWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if h, ok := sw.ResponseWriter.(http.Hijacker); ok {
		return h.Hijack()
	}
	return nil, nil, fmt.Errorf("hijacking not supported")
}

func (app *Application) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")

				// Log panic with stable shape for queryability
				logger := GetLogger(r.Context())
				logger.Error("panic",
					"event", "panic",
					"panic_value", fmt.Sprintf("%v", err),
					"stack", string(debug.Stack()),
				)

				app.RespondServerError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// Authenticate populates the request context with the authenticated user.
// If no token or invalid token is provided, the user is set to AnonymousUser.
// This allows public endpoints to work while still populating user context when available.
// Use RequireAuth middleware on routes that must reject anonymous users.
func (app *Application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		token := app.extractTokenFromHeader(r)
		if token == "" {
			r = app.WithUser(r, AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		user, err := app.ValidateSession(r.Context(), token)
		if err != nil {
			if errors.Is(err, ErrInvalidToken) {
				r = app.WithUser(r, AnonymousUser)
				next.ServeHTTP(w, r)
				return
			}
			app.RespondServerError(w, r, err)
			return
		}

		// Add user_id to logger for authenticated users
		if logCtx := GetLogContext(r.Context()); logCtx != nil {
			logCtx.UserID = user.ID.String()
		}

		r = app.WithUser(r, user)
		next.ServeHTTP(w, r)
	})
}

// RequireAuth is a middleware that rejects requests from anonymous users.
// Use this to wrap handlers that require authentication, providing defense-in-depth
// beyond the per-handler IsAnonymous() checks.
func (app *Application) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.GetUser(r.Context())
		if user.IsAnonymous() {
			app.RespondError(w, r, http.StatusUnauthorized, "Authentication required")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Application) extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}

	return ""
}

func (app *Application) buildRequestLogger(logCtx *RequestLogContext) *slog.Logger {
	return app.Logger.With(
		"request_id", logCtx.RequestID,
		"trace_id", logCtx.TraceID,
		"user_id", logCtx.UserID,
		"method", logCtx.Method,
		"path", logCtx.Path,
	)
}

func extractTraceID(r *http.Request) string {
	if traceID := traceIDFromTraceparent(r.Header.Get("traceparent")); traceID != "" {
		return traceID
	}
	if traceID := r.Header.Get("X-Trace-ID"); traceID != "" {
		return traceID
	}
	return ""
}

func traceIDFromTraceparent(header string) string {
	if header == "" {
		return ""
	}
	parts := strings.Split(header, "-")
	if len(parts) < 4 {
		return ""
	}
	traceID := parts[1]
	if len(traceID) != 32 || !isHexString(traceID) || traceID == strings.Repeat("0", 32) {
		return ""
	}
	return traceID
}

func isHexString(value string) bool {
	for i := 0; i < len(value); i++ {
		ch := value[i]
		if ch >= '0' && ch <= '9' {
			continue
		}
		if ch >= 'a' && ch <= 'f' {
			continue
		}
		if ch >= 'A' && ch <= 'F' {
			continue
		}
		return false
	}
	return true
}
