package app

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"apexmemory.ai/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// FSRSSettings represents the FSRS algorithm configuration for a notebook.
// Go code is the source of truth for default values.
type FSRSSettings struct {
	DesiredRetention float64   `json:"desired_retention"`
	Version          string    `json:"version"`
	Params           []float64 `json:"params"`
	LearningSteps    []int     `json:"learning_steps"`
	RelearningSteps  []int     `json:"relearning_steps"`
	MaximumInterval  int       `json:"maximum_interval"`
	EnableFuzzing    bool      `json:"enable_fuzzing"`
}

// DefaultFSRSSettings returns the default FSRS v6 settings.
// This is the canonical source of truth for FSRS defaults.
func DefaultFSRSSettings() FSRSSettings {
	return FSRSSettings{
		DesiredRetention: 0.9,
		Version:          "6",
		Params: []float64{
			0.212, 1.2931, 2.3065, 8.2956, 6.4133, 0.8334, 3.0194, 0.001,
			1.8722, 0.1666, 0.796, 1.4835, 0.0614, 0.2629, 1.6483, 0.6014,
			1.8729, 0.5425, 0.0912, 0.0658, 0.1542,
		},
		LearningSteps:   []int{60, 600},
		RelearningSteps: []int{600},
		MaximumInterval: 36500,
		EnableFuzzing:   true,
	}
}

// NotebookResponse is the API response representation of a notebook.
type NotebookResponse struct {
	ID           uuid.UUID    `json:"id"`
	Name         string       `json:"name"`
	Description  *string      `json:"description"`
	Emoji        *string      `json:"emoji"`
	Color        *string      `json:"color"`
	FSRSSettings FSRSSettings `json:"fsrs_settings"`
	Position     int32        `json:"position"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// toNotebookResponse converts a db.Notebook to a NotebookResponse.
func toNotebookResponse(n db.Notebook) NotebookResponse {
	resp := NotebookResponse{
		ID:        n.ID,
		Name:      n.Name,
		Position:  n.Position,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}

	if n.Description.Valid {
		resp.Description = &n.Description.String
	}
	if n.Emoji.Valid {
		resp.Emoji = &n.Emoji.String
	}
	if n.Color.Valid {
		resp.Color = &n.Color.String
	}

	var settings FSRSSettings
	if err := json.Unmarshal(n.FsrsSettings, &settings); err != nil {
		settings = DefaultFSRSSettings()
	}
	resp.FSRSSettings = settings

	return resp
}

// CreateNotebookParams holds the parameters for creating a notebook.
type CreateNotebookParams struct {
	Name        string
	Description *string
	Emoji       *string
	Color       *string
	Position    *int32
}

// CreateNotebook creates a new notebook for the given user.
func (app *Application) CreateNotebook(ctx context.Context, userID uuid.UUID, params CreateNotebookParams) (db.Notebook, error) {
	// Generate default FSRS settings (Go is source of truth)
	fsrsJSON, err := json.Marshal(DefaultFSRSSettings())
	if err != nil {
		return db.Notebook{}, err
	}

	dbParams := db.CreateNotebookParams{
		UserID:       userID,
		Name:         params.Name,
		FsrsSettings: fsrsJSON,
	}

	if params.Description != nil {
		dbParams.Description = pgtype.Text{String: *params.Description, Valid: true}
	}

	if params.Emoji != nil {
		dbParams.Emoji = pgtype.Text{String: *params.Emoji, Valid: true}
	}

	if params.Color != nil {
		dbParams.Color = pgtype.Text{String: *params.Color, Valid: true}
	}

	if params.Position != nil {
		dbParams.Position = *params.Position
	}

	return app.Queries.CreateNotebook(ctx, dbParams)
}

// GetNotebook retrieves a notebook by ID for the given user.
// Returns errNotebookNotFound if the notebook doesn't exist or belongs to another user.
func (app *Application) GetNotebook(ctx context.Context, userID, notebookID uuid.UUID) (db.Notebook, error) {
	notebook, err := app.Queries.GetNotebook(ctx, db.GetNotebookParams{
		UserID: userID,
		ID:     notebookID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Notebook{}, errNotebookNotFound
		}
		return db.Notebook{}, err
	}
	return notebook, nil
}

// ListNotebooks retrieves all non-archived notebooks for the given user.
func (app *Application) ListNotebooks(ctx context.Context, userID uuid.UUID) ([]db.Notebook, error) {
	return app.Queries.ListNotebooks(ctx, userID)
}

// UpdateNotebookParams holds the optional fields for updating a notebook.
type UpdateNotebookParams struct {
	Name             *string
	Description      OptionalString // Supports clearing via explicit null
	Emoji            OptionalString // Supports clearing via explicit null
	Color            OptionalString // Supports clearing via explicit null
	Position         *int32
	DesiredRetention *float64
}

// UpdateNotebook updates a notebook with the provided fields.
// All updates happen atomically in a single query.
func (app *Application) UpdateNotebook(ctx context.Context, userID, notebookID uuid.UUID, params UpdateNotebookParams) (db.Notebook, error) {
	dbParams := db.UpdateNotebookParams{
		UserID: userID,
		ID:     notebookID,
	}

	if params.Name != nil {
		dbParams.Name = pgtype.Text{String: *params.Name, Valid: true}
	}

	// Handle description: can be updated or cleared
	if params.Description.Set {
		if params.Description.Value == nil {
			dbParams.ClearDescription = true
		} else {
			dbParams.Description = pgtype.Text{String: *params.Description.Value, Valid: true}
		}
	}

	if params.Emoji.Set {
		if params.Emoji.Value == nil {
			dbParams.ClearEmoji = true
		} else {
			dbParams.Emoji = pgtype.Text{String: *params.Emoji.Value, Valid: true}
		}
	}

	if params.Color.Set {
		if params.Color.Value == nil {
			dbParams.ClearColor = true
		} else {
			dbParams.Color = pgtype.Text{String: *params.Color.Value, Valid: true}
		}
	}

	if params.Position != nil {
		dbParams.Position = pgtype.Int4{Int32: *params.Position, Valid: true}
	}

	if params.DesiredRetention != nil {
		dbParams.UpdateRetention = true
		dbParams.DesiredRetention = *params.DesiredRetention
	}

	notebook, err := app.Queries.UpdateNotebook(ctx, dbParams)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Notebook{}, errNotebookNotFound
		}
		return db.Notebook{}, err
	}

	return notebook, nil
}

// ArchiveNotebook archives a notebook (soft delete).
// This operation is idempotent: archiving an already-archived notebook succeeds.
// Returns errNotebookNotFound only if the notebook doesn't exist.
func (app *Application) ArchiveNotebook(ctx context.Context, userID, notebookID uuid.UUID) error {
	rowsAffected, err := app.Queries.ArchiveNotebook(ctx, db.ArchiveNotebookParams{
		UserID: userID,
		ID:     notebookID,
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		// Check if notebook exists but is already archived (idempotent case)
		status, err := app.Queries.IsNotebookArchived(ctx, db.IsNotebookArchivedParams{
			UserID: userID,
			ID:     notebookID,
		})
		if err != nil {
			return err
		}
		if !status.Exists {
			return errNotebookNotFound
		}
		// Already archived - idempotent success
	}
	return nil
}

// errNotebookNotFound is a sentinel error for notebook not found conditions.
var errNotebookNotFound = errors.New("notebook not found")

// TODO(seed): Remove this function once auto-seed script is created.
// seedExampleNotebooks creates 4 example notebooks for new users.
func (app *Application) seedExampleNotebooks(ctx context.Context, userID uuid.UUID) ([]db.Notebook, error) {
	examples := []struct {
		emoji       string
		name        string
		description string
		position    int32
	}{
		{"🧬", "Biology 202", "Cell biology, genetics, and evolution fundamentals", 0},
		{"🇪🇸", "Spanish B2", "Intermediate Spanish vocabulary and grammar", 1},
		{"♾️", "Calculus", "Derivatives, integrals, and limits", 2},
		{"🇺🇸", "US History", "American history from colonial era to present", 3},
		{"🌍", "World History", "World history from colonial era to present", 4},
		{"🇺🇸", "Civil War", "American Civil War history", 5},
		{"🌡️", "Thermodynamics", "Thermodynamics and heat transfer", 6},
		// create
		{"🚀", "Quantum Mechanics", "Quantum mechanics and quantum field theory", 7},
		{"⚡", "Electrodynamics", "Electrodynamics and electromagnetism", 8},
		{"⚓", "Nautical Archaeology", "Nautical archaeology and shipwreck history", 9},
	}

	notebooks := make([]db.Notebook, 0, len(examples))
	for _, ex := range examples {
		desc := ex.description
		emoji := ex.emoji
		pos := ex.position
		nb, err := app.CreateNotebook(ctx, userID, CreateNotebookParams{
			Name:        ex.name,
			Description: &desc,
			Emoji:       &emoji,
			Position:    &pos,
		})
		if err != nil {
			return nil, err
		}
		notebooks = append(notebooks, nb)
	}
	return notebooks, nil
}
