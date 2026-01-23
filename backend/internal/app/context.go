package app

import (
	"context"
	"net/http"
)

// contextKey is a custom type for context keys to avoid collisions.
type contextKey string

// Context keys for request-scoped values.
const (
	userContextKey = contextKey("user")
	loggerKey      = contextKey("logger")
	logContextKey  = contextKey("log_context")
)

func (app *Application) WithUser(r *http.Request, user *AppUser) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *Application) GetUser(ctx context.Context) *AppUser {
	user, ok := ctx.Value(userContextKey).(*AppUser)
	if !ok {
		return AnonymousUser
	}
	return user
}
