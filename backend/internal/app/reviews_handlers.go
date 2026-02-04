package app

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// GetStudyCardsHandler handles GET /v1/reviews/study
func (app *Application) GetStudyCardsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	limit := int32(20)
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 100 {
			limit = int32(n)
		}
	}

	var notebookID *uuid.UUID
	if v := r.URL.Query().Get("notebook_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			app.RespondError(w, r, http.StatusBadRequest, "invalid notebook_id")
			return
		}
		notebookID = &id
	}

	cards, err := app.getStudyCards(r.Context(), user.ID, notebookID, limit)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("get study cards", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, cards)
}

// GetPracticeCardsHandler handles GET /v1/reviews/practice
func (app *Application) GetPracticeCardsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	limit, offset := parsePagination(r)

	var notebookID *uuid.UUID
	if v := r.URL.Query().Get("notebook_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			app.RespondError(w, r, http.StatusBadRequest, "invalid notebook_id")
			return
		}
		notebookID = &id
	}

	cards, total, err := app.getPracticeCards(r.Context(), user.ID, notebookID, limit, offset)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("get practice cards", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, NewPageResponse(cards, total, limit, offset))
}

// SubmitReviewHandler handles POST /v1/reviews
func (app *Application) SubmitReviewHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req ReviewRequest
	if err := app.ReadJSON(w, r, &req); err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.ID == uuid.Nil {
		app.RespondError(w, r, http.StatusBadRequest, "id is required")
		return
	}
	if req.CardID == uuid.Nil {
		app.RespondError(w, r, http.StatusBadRequest, "card_id is required")
		return
	}
	if req.Rating == "" {
		app.RespondError(w, r, http.StatusBadRequest, "rating is required")
		return
	}

	resp, err := app.submitReview(r.Context(), user.ID, req)
	if err != nil {
		if errors.Is(err, errCardNotReviewable) {
			app.RespondError(w, r, http.StatusNotFound, "Card not found")
			return
		}
		if errors.Is(err, errInvalidRating) || errors.Is(err, errInvalidMode) {
			app.RespondError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		app.RespondServerError(w, r, ErrDBTransaction("submit review", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, resp)
}

// GetStudyCountsHandler handles GET /v1/reviews/study-counts
func (app *Application) GetStudyCountsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	resp, err := app.getStudyCounts(r.Context(), user.ID)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("get study counts", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, resp)
}

// UndoReviewHandler handles DELETE /v1/reviews/{id}
func (app *Application) UndoReviewHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	reviewID, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "invalid review id")
		return
	}

	resp, err := app.undoReview(r.Context(), user.ID, reviewID)
	if err != nil {
		if errors.Is(err, errReviewNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Review not found")
			return
		}
		if errors.Is(err, errReviewNotLatest) {
			app.RespondError(w, r, http.StatusConflict, "Cannot undo: this is not the most recent review for this card")
			return
		}
		if errors.Is(err, errCardAlreadyDeleted) {
			app.RespondError(w, r, http.StatusGone, "Cannot undo: card has been deleted")
			return
		}
		app.RespondServerError(w, r, ErrDBTransaction("undo review", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, resp)
}

// parseDate parses a YYYY-MM-DD date from a query parameter.
func parseDate(r *http.Request, param string) (*time.Time, error) {
	v := r.URL.Query().Get(param)
	if v == "" {
		return nil, nil
	}
	t, err := time.Parse("2006-01-02", v)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetReviewSummaryHandler handles GET /v1/reviews/summary
func (app *Application) GetReviewSummaryHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var notebookID *uuid.UUID
	if v := r.URL.Query().Get("notebook_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			app.RespondError(w, r, http.StatusBadRequest, "invalid notebook_id")
			return
		}
		notebookID = &id
	}

	date, err := parseDate(r, "date")
	if err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
		return
	}

	resp, err := app.getReviewSummary(r.Context(), user.ID, notebookID, date)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("get review summary", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, resp)
}

// GetReviewHistoryHandler handles GET /v1/notebooks/{notebook_id}/reviews
func (app *Application) GetReviewHistoryHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebookID, err := uuid.Parse(r.PathValue("notebook_id"))
	if err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "invalid notebook_id")
		return
	}

	date, err := parseDate(r, "date")
	if err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "invalid date format, expected YYYY-MM-DD")
		return
	}

	limit, offset := parsePagination(r)

	items, total, err := app.getReviewHistory(r.Context(), user.ID, notebookID, date, limit, offset)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("get review history", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, NewPageResponse(items, total, limit, offset))
}
