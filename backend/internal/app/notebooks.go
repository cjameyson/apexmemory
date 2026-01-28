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
	Weights          []float64 `json:"weights"`
}

// DefaultFSRSSettings returns the default FSRS v6 settings.
// This is the canonical source of truth for FSRS defaults.
func DefaultFSRSSettings() FSRSSettings {
	return FSRSSettings{
		DesiredRetention: 0.9,
		Version:          "6",
		Weights: []float64{
			0.212, 1.2931, 2.3065, 8.2956, 6.4133, 0.8334, 3.0194, 0.001,
			1.8722, 0.1666, 0.796, 1.4835, 0.0614, 0.2629, 1.6483, 0.6014,
			1.8729, 0.5425, 0.0912, 0.0658, 0.1542,
		},
	}
}

// NotebookResponse is the API response representation of a notebook.
type NotebookResponse struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Description      *string   `json:"description"`
	DesiredRetention float64   `json:"desired_retention"`
	Position         int32     `json:"position"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
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

	// Parse FSRS settings to extract desired_retention
	// Note: For high-volume lists, consider storing desired_retention as a separate column
	var settings FSRSSettings
	if err := json.Unmarshal(n.FsrsSettings, &settings); err == nil {
		resp.DesiredRetention = settings.DesiredRetention
	} else {
		resp.DesiredRetention = DefaultFSRSSettings().DesiredRetention
	}

	return resp
}

// CreateNotebookParams holds the parameters for creating a notebook.
type CreateNotebookParams struct {
	Name        string
	Description *string
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
