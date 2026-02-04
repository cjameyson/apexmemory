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

	// Facts
	mux.Handle("POST /v1/notebooks/{notebook_id}/facts", protected.ThenFunc(app.CreateFactHandler))
	mux.Handle("GET /v1/notebooks/{notebook_id}/facts", protected.ThenFunc(app.ListFactsHandler))
	mux.Handle("GET /v1/notebooks/{notebook_id}/facts/{id}", protected.ThenFunc(app.GetFactHandler))
	mux.Handle("PATCH /v1/notebooks/{notebook_id}/facts/{id}", protected.ThenFunc(app.UpdateFactHandler))
	mux.Handle("DELETE /v1/notebooks/{notebook_id}/facts/{id}", protected.ThenFunc(app.DeleteFactHandler))

	// Cards
	mux.Handle("GET /v1/notebooks/{notebook_id}/cards", protected.ThenFunc(app.ListCardsHandler))
	mux.Handle("GET /v1/notebooks/{notebook_id}/cards/{id}", protected.ThenFunc(app.GetCardHandler))

	// Reviews
	mux.Handle("GET /v1/reviews/study", protected.ThenFunc(app.GetStudyCardsHandler))
	mux.Handle("GET /v1/reviews/practice", protected.ThenFunc(app.GetPracticeCardsHandler))
	mux.Handle("GET /v1/reviews/study-counts", protected.ThenFunc(app.GetStudyCountsHandler))
	mux.Handle("GET /v1/reviews/summary", protected.ThenFunc(app.GetReviewSummaryHandler))
	mux.Handle("POST /v1/reviews", protected.ThenFunc(app.SubmitReviewHandler))
	mux.Handle("DELETE /v1/reviews/{id}", protected.ThenFunc(app.UndoReviewHandler))
	mux.Handle("GET /v1/notebooks/{notebook_id}/reviews", protected.ThenFunc(app.GetReviewHistoryHandler))

	return global.Then(mux)
}
