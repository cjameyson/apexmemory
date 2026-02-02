package app

import (
	"errors"
	"net/http"

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
		if n, err := parseInt32(v); err == nil && n > 0 && n <= 100 {
			limit = n
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
		if err.Error() == "invalid rating: "+`"`+req.Rating+`"` || err.Error() == "invalid mode: "+`"`+req.Mode+`"` {
			app.RespondError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		app.RespondServerError(w, r, ErrDBTransaction("submit review", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, resp)
}

// parseInt32 parses a string to int32.
func parseInt32(s string) (int32, error) {
	var n int32
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0, errors.New("not a number")
		}
		n = n*10 + int32(c-'0')
	}
	return n, nil
}
