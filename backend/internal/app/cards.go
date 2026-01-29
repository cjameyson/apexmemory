package app

import (
	"context"
	"errors"

	"apexmemory.ai/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var errCardNotFound = errors.New("card not found")

// GetCard retrieves a single card by ID.
func (app *Application) GetCard(ctx context.Context, userID, notebookID, cardID uuid.UUID) (db.Card, error) {
	card, err := app.Queries.GetCard(ctx, db.GetCardParams{
		UserID:     userID,
		ID:         cardID,
		NotebookID: notebookID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.Card{}, errCardNotFound
		}
		return db.Card{}, err
	}
	return card, nil
}

// ListCards retrieves paginated cards for a notebook.
func (app *Application) ListCards(ctx context.Context, userID, notebookID uuid.UUID, state db.NullAppCardState, limit, offset int32) ([]db.Card, int64, error) {
	cards, err := app.Queries.ListCardsByNotebook(ctx, db.ListCardsByNotebookParams{
		UserID:     userID,
		NotebookID: notebookID,
		State:      state,
		RowLimit:   limit,
		RowOffset:  offset,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := app.Queries.CountCardsByNotebook(ctx, db.CountCardsByNotebookParams{
		UserID:     userID,
		NotebookID: notebookID,
	})
	if err != nil {
		return nil, 0, err
	}

	return cards, total, nil
}
