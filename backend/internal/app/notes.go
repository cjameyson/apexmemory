package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"apexmemory.ai/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var (
	errNoteNotFound     = errors.New("note not found")
	errNoteTypeImmutable = errors.New("note type cannot be changed")
)

const maxElementsPerNote = 128

var clozePattern = regexp.MustCompile(`\{\{c(\d+)::`)
var clozeIDPattern = regexp.MustCompile(`^c[1-9][0-9]{0,2}$`)
var imageOcclusionIDPattern = regexp.MustCompile(`^m_[a-zA-Z0-9_-]{6,24}$`)

// NoteResponse is the API response for a note.
type NoteResponse struct {
	ID         uuid.UUID        `json:"id"`
	NotebookID uuid.UUID       `json:"notebook_id"`
	NoteType   string           `json:"note_type"`
	Content    json.RawMessage  `json:"content"`
	SourceID   *uuid.UUID       `json:"source_id"`
	CardCount  int64            `json:"card_count"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

// NoteDetailResponse includes the note and its cards.
type NoteDetailResponse struct {
	NoteResponse
	Cards []CardResponse `json:"cards"`
}

// CardResponse is the API response for a card.
type CardResponse struct {
	ID            uuid.UUID  `json:"id"`
	NoteID        uuid.UUID  `json:"note_id"`
	NotebookID    uuid.UUID  `json:"notebook_id"`
	ElementID     string     `json:"element_id"`
	State         string     `json:"state"`
	Stability     *float32   `json:"stability"`
	Difficulty    *float32   `json:"difficulty"`
	Due           *time.Time `json:"due"`
	Reps          int32      `json:"reps"`
	Lapses        int32      `json:"lapses"`
	SuspendedAt   *time.Time `json:"suspended_at"`
	BuriedUntil   *string    `json:"buried_until"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func toNoteResponse(n db.GetNoteRow) NoteResponse {
	resp := NoteResponse{
		ID:         n.ID,
		NotebookID: n.NotebookID,
		NoteType:   n.NoteType,
		Content:    json.RawMessage(n.Content),
		CardCount:  n.CardCount,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
	if n.SourceID.Valid {
		id := uuid.UUID(n.SourceID.Bytes)
		resp.SourceID = &id
	}
	return resp
}

func toNoteListResponse(n db.ListNotesByNotebookRow) NoteResponse {
	resp := NoteResponse{
		ID:         n.ID,
		NotebookID: n.NotebookID,
		NoteType:   n.NoteType,
		Content:    json.RawMessage(n.Content),
		CardCount:  n.CardCount,
		CreatedAt:  n.CreatedAt,
		UpdatedAt:  n.UpdatedAt,
	}
	if n.SourceID.Valid {
		id := uuid.UUID(n.SourceID.Bytes)
		resp.SourceID = &id
	}
	return resp
}

func toCardResponse(c db.Card) CardResponse {
	resp := CardResponse{
		ID:         c.ID,
		NoteID:     c.NoteID,
		NotebookID: c.NotebookID,
		ElementID:  c.ElementID,
		State:      string(c.State),
		Reps:       c.Reps,
		Lapses:     c.Lapses,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
	if c.Stability.Valid {
		resp.Stability = &c.Stability.Float32
	}
	if c.Difficulty.Valid {
		resp.Difficulty = &c.Difficulty.Float32
	}
	if c.Due.Valid {
		resp.Due = &c.Due.Time
	}
	if c.SuspendedAt.Valid {
		resp.SuspendedAt = &c.SuspendedAt.Time
	}
	if c.BuriedUntil.Valid {
		s := c.BuriedUntil.Time.Format("2006-01-02")
		resp.BuriedUntil = &s
	}
	return resp
}

// validateNoteContent validates the JSON structure and returns extracted element IDs.
func validateNoteContent(noteType string, content json.RawMessage) ([]string, error) {
	// Parse top-level structure
	var parsed struct {
		Version float64         `json:"version"`
		Fields  json.RawMessage `json:"fields"`
	}
	if err := json.Unmarshal(content, &parsed); err != nil {
		return nil, fmt.Errorf("invalid content JSON: %w", err)
	}
	if parsed.Version == 0 {
		return nil, errors.New("content must have a version")
	}
	if parsed.Fields == nil {
		return nil, errors.New("content must have fields")
	}

	elementIDs, err := extractElementIDs(noteType, content)
	if err != nil {
		return nil, err
	}

	if len(elementIDs) == 0 {
		return nil, errors.New("note must generate at least one card")
	}
	if len(elementIDs) > maxElementsPerNote {
		return nil, fmt.Errorf("note exceeds maximum of %d cards", maxElementsPerNote)
	}

	if err := validateElementIDs(noteType, elementIDs); err != nil {
		return nil, err
	}

	return elementIDs, nil
}

// extractElementIDs extracts element IDs from content based on note type.
func extractElementIDs(noteType string, content json.RawMessage) ([]string, error) {
	switch noteType {
	case "basic":
		return []string{""}, nil
	case "cloze":
		return extractClozeElementIDs(content)
	case "image_occlusion":
		return extractImageOcclusionElementIDs(content)
	default:
		return nil, fmt.Errorf("unsupported note type: %s", noteType)
	}
}

func extractClozeElementIDs(content json.RawMessage) ([]string, error) {
	var parsed struct {
		Fields []struct {
			ClozeText string `json:"cloze_text"`
		} `json:"fields"`
	}
	if err := json.Unmarshal(content, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse cloze fields: %w", err)
	}

	seen := make(map[string]bool)
	for _, field := range parsed.Fields {
		matches := clozePattern.FindAllStringSubmatch(field.ClozeText, -1)
		for _, m := range matches {
			num, err := strconv.Atoi(m[1])
			if err != nil || num < 1 || num > 999 {
				continue
			}
			id := fmt.Sprintf("c%d", num)
			seen[id] = true
		}
	}

	if len(seen) == 0 {
		return nil, errors.New("cloze note must contain at least one {{cN::...}} deletion")
	}

	ids := make([]string, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return ids, nil
}

func extractImageOcclusionElementIDs(content json.RawMessage) ([]string, error) {
	var parsed struct {
		Fields []struct {
			Type    string `json:"type"`
			Regions []struct {
				ID string `json:"id"`
			} `json:"regions"`
		} `json:"fields"`
	}
	if err := json.Unmarshal(content, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse image occlusion fields: %w", err)
	}

	seen := make(map[string]bool)
	var ids []string
	for _, field := range parsed.Fields {
		if field.Type != "image_occlusion" {
			continue
		}
		for _, region := range field.Regions {
			if region.ID == "" {
				return nil, errors.New("image occlusion region missing id")
			}
			if seen[region.ID] {
				return nil, fmt.Errorf("duplicate region id: %s", region.ID)
			}
			seen[region.ID] = true
			ids = append(ids, region.ID)
		}
	}

	if len(ids) == 0 {
		return nil, errors.New("image occlusion note must contain at least one region")
	}
	return ids, nil
}

func validateElementIDs(noteType string, ids []string) error {
	for _, id := range ids {
		switch noteType {
		case "basic":
			if id != "" {
				return fmt.Errorf("basic note element_id must be empty, got %q", id)
			}
		case "cloze":
			if !clozeIDPattern.MatchString(id) {
				return fmt.Errorf("invalid cloze element_id: %q", id)
			}
		case "image_occlusion":
			if !imageOcclusionIDPattern.MatchString(id) {
				return fmt.Errorf("invalid image occlusion element_id: %q", id)
			}
		}
	}
	return nil
}

// CreateNote creates a note and its derived cards in a transaction.
func (app *Application) CreateNote(ctx context.Context, userID, notebookID uuid.UUID, noteType string, content json.RawMessage) (db.GetNoteRow, []db.Card, error) {
	elementIDs, err := validateNoteContent(noteType, content)
	if err != nil {
		return db.GetNoteRow{}, nil, err
	}

	var note db.Note
	var cards []db.Card

	err = WithTx(ctx, app.DB, func(q *db.Queries) error {
		var txErr error
		note, txErr = q.CreateNote(ctx, db.CreateNoteParams{
			UserID:     userID,
			NotebookID: notebookID,
			NoteType:   noteType,
			Content:    content,
		})
		if txErr != nil {
			return txErr
		}

		cards = make([]db.Card, 0, len(elementIDs))
		for _, elemID := range elementIDs {
			card, txErr := q.CreateCard(ctx, db.CreateCardParams{
				UserID:     userID,
				NotebookID: notebookID,
				NoteID:     note.ID,
				ElementID:  elemID,
			})
			if txErr != nil {
				return txErr
			}
			cards = append(cards, card)
		}
		return nil
	})
	if err != nil {
		return db.GetNoteRow{}, nil, err
	}

	// Return as GetNoteRow with card count
	row := db.GetNoteRow{
		UserID:     note.UserID,
		ID:         note.ID,
		NotebookID: note.NotebookID,
		NoteType:   note.NoteType,
		Content:    note.Content,
		SourceID:   note.SourceID,
		CreatedAt:  note.CreatedAt,
		UpdatedAt:  note.UpdatedAt,
		CardCount:  int64(len(cards)),
	}
	return row, cards, nil
}

// GetNote retrieves a single note with card count.
func (app *Application) GetNote(ctx context.Context, userID, notebookID, noteID uuid.UUID) (db.GetNoteRow, error) {
	note, err := app.Queries.GetNote(ctx, db.GetNoteParams{
		UserID:     userID,
		ID:         noteID,
		NotebookID: notebookID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.GetNoteRow{}, errNoteNotFound
		}
		return db.GetNoteRow{}, err
	}
	return note, nil
}

// ListNotes retrieves paginated notes for a notebook.
func (app *Application) ListNotes(ctx context.Context, userID, notebookID uuid.UUID, limit, offset int32) ([]db.ListNotesByNotebookRow, int64, error) {
	notes, err := app.Queries.ListNotesByNotebook(ctx, db.ListNotesByNotebookParams{
		UserID:     userID,
		NotebookID: notebookID,
		RowLimit:   limit,
		RowOffset:  offset,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := app.Queries.CountNotesByNotebook(ctx, db.CountNotesByNotebookParams{
		UserID:     userID,
		NotebookID: notebookID,
	})
	if err != nil {
		return nil, 0, err
	}

	return notes, total, nil
}

// UpdateNoteResult contains the outcome of a note update.
type UpdateNoteResult struct {
	Note      db.Note
	Created   int
	Deleted   int
	Unchanged int
}

// UpdateNote updates note content and diffs cards atomically.
func (app *Application) UpdateNote(ctx context.Context, userID, notebookID, noteID uuid.UUID, content json.RawMessage) (UpdateNoteResult, error) {
	// Fetch existing note to verify ownership and get type
	existingNote, err := app.Queries.GetNote(ctx, db.GetNoteParams{
		UserID:     userID,
		ID:         noteID,
		NotebookID: notebookID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UpdateNoteResult{}, errNoteNotFound
		}
		return UpdateNoteResult{}, err
	}

	expectedIDs, err := validateNoteContent(existingNote.NoteType, content)
	if err != nil {
		return UpdateNoteResult{}, err
	}

	var result UpdateNoteResult

	err = WithTx(ctx, app.DB, func(q *db.Queries) error {
		// Update the note content
		note, txErr := q.UpdateNoteContent(ctx, db.UpdateNoteContentParams{
			Content:    content,
			UserID:     userID,
			ID:         noteID,
			NotebookID: notebookID,
		})
		if txErr != nil {
			return txErr
		}
		result.Note = note

		// Fetch existing cards to diff
		existingCards, txErr := q.ListCardsByNote(ctx, db.ListCardsByNoteParams{
			UserID: userID,
			NoteID: noteID,
		})
		if txErr != nil {
			return txErr
		}

		existingSet := make(map[string]bool, len(existingCards))
		for _, c := range existingCards {
			existingSet[c.ElementID] = true
		}

		expectedSet := make(map[string]bool, len(expectedIDs))
		for _, id := range expectedIDs {
			expectedSet[id] = true
		}

		// Compute diff
		var toCreate []string
		var toDelete []string
		for _, id := range expectedIDs {
			if !existingSet[id] {
				toCreate = append(toCreate, id)
			}
		}
		for _, c := range existingCards {
			if !expectedSet[c.ElementID] {
				toDelete = append(toDelete, c.ElementID)
			}
		}

		// Delete removed cards
		if len(toDelete) > 0 {
			if txErr := q.DeleteCardsByNoteAndElements(ctx, db.DeleteCardsByNoteAndElementsParams{
				UserID:     userID,
				NoteID:     noteID,
				ElementIds: toDelete,
			}); txErr != nil {
				return txErr
			}
		}

		// Create new cards
		for _, elemID := range toCreate {
			if _, txErr := q.CreateCard(ctx, db.CreateCardParams{
				UserID:     userID,
				NotebookID: notebookID,
				NoteID:     noteID,
				ElementID:  elemID,
			}); txErr != nil {
				return txErr
			}
		}

		result.Created = len(toCreate)
		result.Deleted = len(toDelete)
		result.Unchanged = len(expectedIDs) - len(toCreate)
		return nil
	})

	return result, err
}

// DeleteNote deletes a note (cascades to cards).
func (app *Application) DeleteNote(ctx context.Context, userID, noteID uuid.UUID) error {
	rows, err := app.Queries.DeleteNote(ctx, db.DeleteNoteParams{
		UserID: userID,
		ID:     noteID,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return errNoteNotFound
	}
	return nil
}
