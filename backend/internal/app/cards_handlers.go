package app

import (
	"errors"
	"net/http"

	"apexmemory.ai/internal/db"
)

// ListCardsHandler handles GET /v1/notebooks/{notebook_id}/cards
func (app *Application) ListCardsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.MustUser(r)

	notebookID, ok := app.PathUUID(w, r, "notebook_id")
	if !ok {
		return
	}

	limit, offset := parsePagination(r)

	// Parse optional state filter
	var state db.NullAppCardState
	if s := r.URL.Query().Get("state"); s != "" {
		switch db.AppCardState(s) {
		case db.AppCardStateNew, db.AppCardStateLearning, db.AppCardStateReview, db.AppCardStateRelearning:
			state = db.NullAppCardState{AppCardState: db.AppCardState(s), Valid: true}
		default:
			app.RespondError(w, r, http.StatusBadRequest, "invalid state filter")
			return
		}
	}

	cards, total, err := app.ListCards(r.Context(), user.ID, notebookID, state, limit, offset)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("list cards", err))
		return
	}

	data := make([]CardResponse, len(cards))
	for i, c := range cards {
		data[i] = toCardResponse(c)
	}

	app.RespondJSON(w, r, http.StatusOK, NewPageResponse(data, total, limit, offset))
}

// GetCardHandler handles GET /v1/notebooks/{notebook_id}/cards/{id}
func (app *Application) GetCardHandler(w http.ResponseWriter, r *http.Request) {
	user := app.MustUser(r)

	notebookID, ok := app.PathUUID(w, r, "notebook_id")
	if !ok {
		return
	}

	cardID, ok := app.PathUUID(w, r, "id")
	if !ok {
		return
	}

	card, err := app.GetCard(r.Context(), user.ID, notebookID, cardID)
	if err != nil {
		if errors.Is(err, errCardNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Card not found")
			return
		}
		app.RespondServerError(w, r, ErrDBQuery("get card", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, toCardResponse(card))
}
