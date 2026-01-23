package app

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	// Middleware shortcuts
	authRL := app.RateLimitByIP(app.RateLimiters.Auth)
	registerRL := app.RateLimitByIP(app.RateLimiters.Register)
	apiRL := app.RateLimitByUser(app.RateLimiters.API)

	// Public endpoints
	mux.HandleFunc("GET /v1/healthcheck", app.HealthcheckHandler)
	mux.Handle("POST /v1/auth/register", registerRL(http.HandlerFunc(app.RegisterHandler)))
	mux.Handle("POST /v1/auth/login", authRL(http.HandlerFunc(app.LoginHandler)))

	// Protected endpoints (wrapped with RequireAuth for defense-in-depth)
	mux.Handle("POST /v1/auth/logout", app.RequireAuth(apiRL(http.HandlerFunc(app.LogoutHandler))))
	mux.Handle("POST /v1/auth/logout-all", app.RequireAuth(apiRL(http.HandlerFunc(app.LogoutAllHandler))))
	mux.Handle("GET /v1/auth/me", app.RequireAuth(apiRL(http.HandlerFunc(app.GetCurrentUserHandler))))

	return app.RecoverPanic(app.LogRequests(app.Authenticate(mux)))
}
