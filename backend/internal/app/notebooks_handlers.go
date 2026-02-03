package app

import (
	"errors"
	"net/http"
	"strings"
)

const (
	notebookNameMaxLen        = 255
	notebookDescriptionMaxLen = 10000
)

// CreateNotebookHandler handles POST /v1/notebooks
func (app *Application) CreateNotebookHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var input struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Emoji       *string `json:"emoji"`
		Color       *string `json:"color"`
		Position    *int32  `json:"position"`
	}

	if err := app.ReadJSON(w, r, &input); err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate input
	fieldErrors := make(map[string]string)
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		fieldErrors["name"] = "Name is required"
	} else if len(input.Name) > notebookNameMaxLen {
		fieldErrors["name"] = "Name must not exceed 255 characters"
	}

	if input.Description != nil && len(*input.Description) > notebookDescriptionMaxLen {
		fieldErrors["description"] = "Description must not exceed 10,000 characters"
	}

	if len(fieldErrors) > 0 {
		app.RespondFieldErrors(w, r, fieldErrors)
		return
	}

	notebook, err := app.CreateNotebook(r.Context(), user.ID, CreateNotebookParams{
		Name:        input.Name,
		Description: input.Description,
		Emoji:       input.Emoji,
		Color:       input.Color,
		Position:    input.Position,
	})
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("create notebook", err))
		return
	}

	app.RespondJSON(w, r, http.StatusCreated, toNotebookResponse(notebook))
}

// ListNotebooksHandler handles GET /v1/notebooks
func (app *Application) ListNotebooksHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebooks, err := app.ListNotebooks(r.Context(), user.ID)
	if err != nil {
		app.RespondServerError(w, r, ErrDBQuery("list notebooks", err))
		return
	}

	// Convert to response format
	response := make([]NotebookResponse, len(notebooks))
	for i, n := range notebooks {
		response[i] = toNotebookResponse(n)
	}

	app.RespondJSON(w, r, http.StatusOK, response)
}

// GetNotebookHandler handles GET /v1/notebooks/{id}
func (app *Application) GetNotebookHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebookID, ok := app.PathUUID(w, r, "id")
	if !ok {
		return
	}

	notebook, err := app.GetNotebook(r.Context(), user.ID, notebookID)
	if err != nil {
		if errors.Is(err, errNotebookNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Notebook not found")
			return
		}
		app.RespondServerError(w, r, ErrDBQuery("get notebook", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, toNotebookResponse(notebook))
}

// UpdateNotebookHandler handles PATCH /v1/notebooks/{id}
func (app *Application) UpdateNotebookHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebookID, ok := app.PathUUID(w, r, "id")
	if !ok {
		return
	}

	var input struct {
		Name             *string        `json:"name"`
		Description      OptionalString `json:"description"`
		Emoji            OptionalString `json:"emoji"`
		Color            OptionalString `json:"color"`
		Position         *int32         `json:"position"`
		DesiredRetention *float64       `json:"desired_retention"`
	}

	if err := app.ReadJSON(w, r, &input); err != nil {
		app.RespondError(w, r, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate fields
	fieldErrors := make(map[string]string)

	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		input.Name = &name
		if name == "" {
			fieldErrors["name"] = "Name cannot be empty"
		} else if len(name) > notebookNameMaxLen {
			fieldErrors["name"] = "Name must not exceed 255 characters"
		}
	}

	if input.DesiredRetention != nil {
		if *input.DesiredRetention <= 0 || *input.DesiredRetention >= 1 {
			fieldErrors["desired_retention"] = "Desired retention must be between 0 and 1 (exclusive)"
		}
	}

	if input.Description.Set && input.Description.Value != nil && len(*input.Description.Value) > notebookDescriptionMaxLen {
		fieldErrors["description"] = "Description must not exceed 10,000 characters"
	}

	if len(fieldErrors) > 0 {
		app.RespondFieldErrors(w, r, fieldErrors)
		return
	}

	notebook, err := app.UpdateNotebook(r.Context(), user.ID, notebookID, UpdateNotebookParams{
		Name:             input.Name,
		Description:      input.Description,
		Emoji:            input.Emoji,
		Color:            input.Color,
		Position:         input.Position,
		DesiredRetention: input.DesiredRetention,
	})
	if err != nil {
		if errors.Is(err, errNotebookNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Notebook not found")
			return
		}
		app.RespondServerError(w, r, ErrDBQuery("update notebook", err))
		return
	}

	app.RespondJSON(w, r, http.StatusOK, toNotebookResponse(notebook))
}

// DeleteNotebookHandler handles DELETE /v1/notebooks/{id}
// This performs a soft delete (archive) rather than permanent deletion.
func (app *Application) DeleteNotebookHandler(w http.ResponseWriter, r *http.Request) {
	user := app.GetUser(r.Context())
	if user.IsAnonymous() {
		app.RespondError(w, r, http.StatusUnauthorized, "Not authenticated")
		return
	}

	notebookID, ok := app.PathUUID(w, r, "id")
	if !ok {
		return
	}

	err := app.ArchiveNotebook(r.Context(), user.ID, notebookID)
	if err != nil {
		if errors.Is(err, errNotebookNotFound) {
			app.RespondError(w, r, http.StatusNotFound, "Notebook not found")
			return
		}
		app.RespondServerError(w, r, ErrDBQuery("archive notebook", err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
