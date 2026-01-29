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

// CreateNoteHandler handles POST /v1/notebooks/{notebook_id}/notes
func (app *Application) CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
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
		NoteType string          `json:"note_type"`
		Content  json.RawMessage `json:"content"`
	}

	if err := app.ReadJSON(w, r, &input); err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if input.NoteType == "" {
		input.NoteType = "basic"
	}

	if input.Content == nil {
		app.RespondError(w, r, http.StatusBadRequest, "content is required")
		return
	}

	note, cards, err := app.CreateNote(r.Context(), user.ID, notebookID, input.NoteType, input.Content)
	if err != nil {
		// Validation errors from content parsing
		if isValidationError(err) {
			app.RespondError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		app.RespondServerError(w, r, ErrDBTransaction("create note", err))
		return
	}

	resp := NoteDetailResponse{
		NoteResponse: toNoteResponse(note),
	}
	resp.Cards = make([]CardResponse, len(cards))
	for i, c := range cards {
		resp.Cards[i] = toCardResponse(c)
	}

	app.RespondJSON(w, r, http.StatusCreated, resp)
}

// ListNotesHandler handles GET /v1/notebooks/{notebook_id}/notes
func (app *Application) ListNotesHandler(w http.ResponseWriter, r *http.Request) {
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

	notes, total, err := app.ListNotes(r.Context(), user.ID, notebookID, limit, offset)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("list notes", err))
		return
	}

	data := make([]NoteResponse, len(notes))
	for i, n := range notes {
		data[i] = toNoteListResponse(n)
	}

	app.RespondJSON(w, r, http.StatusOK, map[string]any{
		"data":  data,
		"total": total,
	})
}

// GetNoteHandler handles GET /v1/notebooks/{notebook_id}/notes/{id}
func (app *Application) GetNoteHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebookID, ok := app.PathUUID(w, r, "notebook_id")
	if !ok {
		return
	}

	noteID, ok := app.PathUUID(w, r, "id")
	if !ok {
		return
	}

	note, err := app.GetNote(r.Context(), user.ID, notebookID, noteID)
	if err != nil {
		if errors.Is(err, errNoteNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Note not found")
			return
		}
		app.RespondServerError(w, r, ErrDBQuery("get note", err))
		return
	}

	// Fetch cards for this note
	cards, err := app.Queries.ListCardsByNote(r.Context(), db.ListCardsByNoteParams{
		UserID: user.ID,
		NoteID: noteID,
	})
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("list cards for note", err))
		return
	}

	resp := NoteDetailResponse{
		NoteResponse: toNoteResponse(note),
	}
	resp.Cards = make([]CardResponse, len(cards))
	for i, c := range cards {
		resp.Cards[i] = toCardResponse(c)
	}

	app.RespondJSON(w, r, http.StatusOK, resp)
}

// UpdateNoteHandler handles PATCH /v1/notebooks/{notebook_id}/notes/{id}
func (app *Application) UpdateNoteHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebookID, ok := app.PathUUID(w, r, "notebook_id")
	if !ok {
		return
	}

	noteID, ok := app.PathUUID(w, r, "id")
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

	result, err := app.UpdateNote(r.Context(), user.ID, notebookID, noteID, input.Content)
	if err != nil {
		if errors.Is(err, errNoteNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Note not found")
			return
		}
		if isValidationError(err) {
			app.RespondError(w, r, http.StatusBadRequest, err.Error())
			return
		}
		app.RespondServerError(w, r, ErrDBTransaction("update note", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, map[string]any{
		"note_id":   result.Note.ID,
		"created":   result.Created,
		"deleted":   result.Deleted,
		"unchanged": result.Unchanged,
	})
}

// DeleteNoteHandler handles DELETE /v1/notebooks/{notebook_id}/notes/{id}
func (app *Application) DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	noteID, ok := app.PathUUID(w, r, "id")
	if !ok {
		return
	}

	err := app.DeleteNote(r.Context(), user.ID, noteID)
	if err != nil {
		if errors.Is(err, errNoteNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Note not found")
			return
		}
		app.RespondServerError(w, r, ErrDBQuery("delete note", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// isValidationError checks if an error is a content validation error (not a DB error).
func isValidationError(err error) bool {
	msg := err.Error()
	// These prefixes cover all validation errors from validateNoteContent and extractElementIDs
	for _, prefix := range []string{
		"invalid content JSON",
		"content must have",
		"note must generate",
		"note exceeds maximum",
		"unsupported note type",
		"cloze note must contain",
		"image occlusion note must contain",
		"image occlusion region missing",
		"duplicate region id",
		"invalid cloze element_id",
		"invalid image occlusion element_id",
		"basic note element_id",
		"failed to parse",
	} {
		if len(msg) >= len(prefix) && msg[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}
