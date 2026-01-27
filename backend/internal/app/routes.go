package app

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	// Middleware chains
	global := alice.New(app.RecoverPanic, app.LogRequests, app.Authenticate)
	protected := alice.New(app.RequireAuth, app.RateLimitByUser(app.RateLimiters.API))
	authChain := alice.New(app.RateLimitByIP(app.RateLimiters.Auth))
	registerChain := alice.New(app.RateLimitByIP(app.RateLimiters.Register))

	// Public endpoints
	mux.HandleFunc("GET /v1/healthcheck", app.HealthcheckHandler)
	mux.Handle("POST /v1/auth/register", registerChain.ThenFunc(app.RegisterHandler))
	mux.Handle("POST /v1/auth/login", authChain.ThenFunc(app.LoginHandler))

	// Protected endpoints
	mux.Handle("POST /v1/auth/logout", protected.ThenFunc(app.LogoutHandler))
	mux.Handle("POST /v1/auth/logout-all", protected.ThenFunc(app.LogoutAllHandler))
	mux.Handle("GET /v1/auth/me", protected.ThenFunc(app.GetCurrentUserHandler))

	// Notebooks
	mux.Handle("POST /v1/notebooks", protected.ThenFunc(app.CreateNotebookHandler))
	mux.Handle("GET /v1/notebooks", protected.ThenFunc(app.ListNotebooksHandler))
	mux.Handle("GET /v1/notebooks/{id}", protected.ThenFunc(app.GetNotebookHandler))
	mux.Handle("PATCH /v1/notebooks/{id}", protected.ThenFunc(app.UpdateNotebookHandler))
	mux.Handle("DELETE /v1/notebooks/{id}", protected.ThenFunc(app.DeleteNotebookHandler))

	return global.Then(mux)
}
