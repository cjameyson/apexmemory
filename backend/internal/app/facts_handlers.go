package app

import (
	"encoding/json"
	"errors"
	"net/http"

	"apexmemory.ai/internal/db"
)

// CreateFactHandler handles POST /v1/notebooks/{notebook_id}/facts
func (app *Application) CreateFactHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebookID, ok := app.PathUUID(w, r, "notebook_id")
	if !ok {
		return
	}

	// Verify notebook exists and belongs to user
	if _, err := app.GetNotebook(r.Context(), user.ID, notebookID); err != nil {
		if errors.Is(err, errNotebookNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Notebook not found")
			return
		}
		app.RespondServerError(w, r, ErrDBQuery("get notebook", err))
		return
	}

	var input struct {
		FactType string          `json:"fact_type"`
		Content  json.RawMessage `json:"content"`
	}

	if err := app.ReadJSON(w, r, &input); err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if input.FactType == "" {
		input.FactType = "basic"
	}

	if input.Content == nil {
		app.RespondError(w, r, http.StatusBadRequest, "content is required")
		return
	}

	fact, cards, err := app.CreateFact(r.Context(), user.ID, notebookID, input.FactType, input.Content)
	if err != nil {
		var validationErr *FactValidationError
		if errors.As(err, &validationErr) {
			app.RespondError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		app.RespondServerError(w, r, ErrDBTransaction("create fact", err))
		return
	}

	resp := FactDetailResponse{
		FactResponse: toFactResponse(fact),
	}
	resp.Cards = make([]CardResponse, len(cards))
	for i, c := range cards {
		resp.Cards[i] = toCardResponse(c)
	}

	app.RespondJSON(w, r, http.StatusCreated, resp)
}

// ListFactsHandler handles GET /v1/notebooks/{notebook_id}/facts
// Query params: limit, offset, type (fact_type filter), q (search), stats=true (include stats)
func (app *Application) ListFactsHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebookID, ok := app.PathUUID(w, r, "notebook_id")
	if !ok {
		return
	}

	limit, offset := parsePagination(r)
	q := r.URL.Query()

	sort, err := parseSort(r, "created", "updated")
	if err != nil {
		app.RespondError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	filters := FactListFilters{
		FactType: q.Get("type"),
		Search:   q.Get("q"),
		Sort:     sort,
	}

	facts, total, err := app.ListFactsFiltered(r.Context(), user.ID, notebookID, filters, limit, offset)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("list facts", err))
		return
	}

	data := make([]FactFilteredResponse, len(facts))
	for i, n := range facts {
		data[i] = toFactFilteredResponse(n)
	}

	resp := NewPageResponse(data, total, limit, offset)

	// If stats requested, include them
	if q.Get("stats") == "true" {
		stats, err := app.GetFactStats(r.Context(), user.ID, notebookID)
		if err != nil {
			app.RespondServerError(w, r, ErrDBQuery("get fact stats", err))
			return
		}
		app.RespondJSON(w, r, http.StatusOK, struct {
			PageResponse[FactFilteredResponse]
			Stats FactStatsResponse `json:"stats"`
		}{
			PageResponse: resp,
			Stats: FactStatsResponse{
				TotalFacts: stats.TotalFacts,
				TotalCards: stats.TotalCards,
				TotalDue:   stats.TotalDue,
				ByType: FactStatsTypeBreakdown{
					Basic:          stats.BasicCount,
					Cloze:          stats.ClozeCount,
					ImageOcclusion: stats.ImageOcclusionCount,
				},
			},
		})
		return
	}

	app.RespondJSON(w, r, http.StatusOK, resp)
}

// GetFactHandler handles GET /v1/notebooks/{notebook_id}/facts/{id}
func (app *Application) GetFactHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebookID, ok := app.PathUUID(w, r, "notebook_id")
	if !ok {
		return
	}

	factID, ok := app.PathUUID(w, r, "id")
	if !ok {
		return
	}

	fact, err := app.GetFact(r.Context(), user.ID, notebookID, factID)
	if err != nil {
		if errors.Is(err, errFactNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Fact not found")
			return
		}
		app.RespondServerError(w, r, ErrDBQuery("get fact", err))
		return
	}

	// Fetch cards for this fact
	cards, err := app.Queries.ListCardsByFact(r.Context(), db.ListCardsByFactParams{
		UserID: user.ID,
		FactID: factID,
	})
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("list cards for fact", err))
		return
	}

	resp := FactDetailResponse{
		FactResponse: toFactResponse(fact),
	}
	resp.Cards = make([]CardResponse, len(cards))
	for i, c := range cards {
		resp.Cards[i] = toCardResponse(c)
	}

	app.RespondJSON(w, r, http.StatusOK, resp)
}

// UpdateFactHandler handles PATCH /v1/notebooks/{notebook_id}/facts/{id}
func (app *Application) UpdateFactHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebookID, ok := app.PathUUID(w, r, "notebook_id")
	if !ok {
		return
	}

	factID, ok := app.PathUUID(w, r, "id")
	if !ok {
		return
	}

	var input struct {
		Content json.RawMessage `json:"content"`
	}

	if err := app.ReadJSON(w, r, &input); err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if input.Content == nil {
		app.RespondError(w, r, http.StatusBadRequest, "content is required")
		return
	}

	result, err := app.UpdateFact(r.Context(), user.ID, notebookID, factID, input.Content)
	if err != nil {
		if errors.Is(err, errFactNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Fact not found")
			return
		}
		var validationErr *FactValidationError
		if errors.As(err, &validationErr) {
			app.RespondError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		app.RespondServerError(w, r, ErrDBTransaction("update fact", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, map[string]any{
		"fact_id":   result.Fact.ID,
		"created":   result.Created,
		"deleted":   result.Deleted,
		"unchanged": result.Unchanged,
	})
}

// DeleteFactHandler handles DELETE /v1/notebooks/{notebook_id}/facts/{id}
func (app *Application) DeleteFactHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebookID, ok := app.PathUUID(w, r, "notebook_id")
	if !ok {
		return
	}

	factID, ok := app.PathUUID(w, r, "id")
	if !ok {
		return
	}

	err := app.DeleteFact(r.Context(), user.ID, notebookID, factID)
	if err != nil {
		if errors.Is(err, errFactNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Fact not found")
			return
		}
		app.RespondServerError(w, r, ErrDBQuery("delete fact", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
