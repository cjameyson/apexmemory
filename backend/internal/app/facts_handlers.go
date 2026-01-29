package app

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"apexmemory.ai/internal/db"
)

const (
	defaultPageLimit = 50
	maxPageLimit     = 100
)

// parsePagination extracts limit and offset from query params with defaults.
func parsePagination(r *http.Request) (limit, offset int32) {
	limit = defaultPageLimit
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			limit = int32(n)
		}
	}
	if limit > maxPageLimit {
		limit = maxPageLimit
	}

	if v := r.URL.Query().Get("offset"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n >= 0 {
			offset = int32(n)
		}
	}
	return
}

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
		// Validation errors from content parsing
		if isValidationError(err) {
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

	facts, total, err := app.ListFacts(r.Context(), user.ID, notebookID, limit, offset)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("list facts", err))
		return
	}

	data := make([]FactResponse, len(facts))
	for i, n := range facts {
		data[i] = toFactListResponse(n)
	}

	app.RespondJSON(w, r, http.StatusOK, map[string]any{
		"data":  data,
		"total": total,
	})
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
		if isValidationError(err) {
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

// isValidationError checks if an error is a content validation error (not a DB error).
func isValidationError(err error) bool {
	msg := err.Error()
	// These prefixes cover all validation errors from validateFactContent and extractElementIDs
	for _, prefix := range []string{
		"invalid content JSON",
		"content must have",
		"fact must generate",
		"fact exceeds maximum",
		"unsupported fact type",
		"cloze fact must contain",
		"image occlusion fact must contain",
		"image occlusion region missing",
		"duplicate region id",
		"invalid cloze element_id",
		"invalid image occlusion element_id",
		"basic fact element_id",
		"failed to parse",
	} {
		if len(msg) >= len(prefix) && msg[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}
