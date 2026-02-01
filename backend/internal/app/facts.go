package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"apexmemory.ai/internal/db"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	errFactNotFound     = errors.New("fact not found")
	errFactTypeImmutable = errors.New("fact type cannot be changed")
)

const maxElementsPerFact = 128

var clozePattern = regexp.MustCompile(`\{\{c(\d+)::`)
var clozeIDPattern = regexp.MustCompile(`^c[1-9][0-9]{0,2}$`)
var imageOcclusionIDPattern = regexp.MustCompile(`^m_[a-zA-Z0-9_-]{6,24}$`)

// FactResponse is the API response for a fact.
type FactResponse struct {
	ID         uuid.UUID        `json:"id"`
	NotebookID uuid.UUID       `json:"notebook_id"`
	FactType   string           `json:"fact_type"`
	Content    json.RawMessage  `json:"content"`
	SourceID   *uuid.UUID       `json:"source_id"`
	CardCount  int64            `json:"card_count"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}

// FactDetailResponse includes the fact and its cards.
type FactDetailResponse struct {
	FactResponse
	Cards []CardResponse `json:"cards"`
}

// CardResponse is the API response for a card.
type CardResponse struct {
	ID            uuid.UUID  `json:"id"`
	FactID        uuid.UUID  `json:"fact_id"`
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

func toFactResponse(n db.GetFactRow) FactResponse {
	resp := FactResponse{
		ID:         n.ID,
		NotebookID: n.NotebookID,
		FactType:   n.FactType,
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

func toFactListResponse(n db.ListFactsByNotebookRow) FactResponse {
	resp := FactResponse{
		ID:         n.ID,
		NotebookID: n.NotebookID,
		FactType:   n.FactType,
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

// FactFilteredResponse extends FactResponse with due_count.
type FactFilteredResponse struct {
	FactResponse
	DueCount int64 `json:"due_count"`
}

// FactStatsResponse is the stats summary for a notebook's facts.
type FactStatsResponse struct {
	TotalFacts int64                    `json:"total_facts"`
	TotalCards int64                    `json:"total_cards"`
	TotalDue   int64                    `json:"total_due"`
	ByType     FactStatsTypeBreakdown   `json:"by_type"`
}

// FactStatsTypeBreakdown is the per-type breakdown.
type FactStatsTypeBreakdown struct {
	Basic          int64 `json:"basic"`
	Cloze          int64 `json:"cloze"`
	ImageOcclusion int64 `json:"image_occlusion"`
}

func toFactFilteredResponse(n db.ListFactsByNotebookFilteredRow) FactFilteredResponse {
	resp := FactFilteredResponse{
		FactResponse: FactResponse{
			ID:         n.ID,
			NotebookID: n.NotebookID,
			FactType:   n.FactType,
			Content:    json.RawMessage(n.Content),
			CardCount:  n.CardCount,
			CreatedAt:  n.CreatedAt,
			UpdatedAt:  n.UpdatedAt,
		},
		DueCount: n.DueCount,
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
		FactID:     c.FactID,
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

// validateFactContent validates the JSON structure and returns extracted element IDs.
func validateFactContent(factType string, content json.RawMessage) ([]string, error) {
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

	elementIDs, err := extractElementIDs(factType, content)
	if err != nil {
		return nil, err
	}

	if len(elementIDs) == 0 {
		return nil, errors.New("fact must generate at least one card")
	}
	if len(elementIDs) > maxElementsPerFact {
		return nil, fmt.Errorf("fact exceeds maximum of %d cards", maxElementsPerFact)
	}

	if err := validateElementIDs(factType, elementIDs); err != nil {
		return nil, err
	}

	return elementIDs, nil
}

// extractElementIDs extracts element IDs from content based on fact type.
func extractElementIDs(factType string, content json.RawMessage) ([]string, error) {
	switch factType {
	case "basic":
		return []string{""}, nil
	case "cloze":
		return extractClozeElementIDs(content)
	case "image_occlusion":
		return extractImageOcclusionElementIDs(content)
	default:
		return nil, fmt.Errorf("unsupported fact type: %s", factType)
	}
}

func extractClozeElementIDs(content json.RawMessage) ([]string, error) {
	var parsed struct {
		Fields []struct {
			Value string `json:"value"`
		} `json:"fields"`
	}
	if err := json.Unmarshal(content, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse cloze fields: %w", err)
	}

	seen := make(map[string]bool)
	for _, field := range parsed.Fields {
		matches := clozePattern.FindAllStringSubmatch(field.Value, -1)
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
		return nil, errors.New("cloze fact must contain at least one {{cN::...}} deletion")
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
			Name    string `json:"name"`
			Type    string `json:"type"`
			Value   string `json:"value"`
			Regions []struct {
				ID string `json:"id"`
			} `json:"regions"`
		} `json:"fields"`
	}
	if err := json.Unmarshal(content, &parsed); err != nil {
		return nil, fmt.Errorf("failed to parse image occlusion fields: %w", err)
	}

	// Validate required title field
	hasTitle := false
	for _, field := range parsed.Fields {
		if field.Name == "title" && field.Type == "plain_text" {
			if strings.TrimSpace(field.Value) == "" {
				return nil, errors.New("image occlusion title must not be empty")
			}
			hasTitle = true
			break
		}
	}
	if !hasTitle {
		return nil, errors.New("image occlusion fact must have a title field")
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
		return nil, errors.New("image occlusion fact must contain at least one region")
	}
	return ids, nil
}

func validateElementIDs(factType string, ids []string) error {
	for _, id := range ids {
		switch factType {
		case "basic":
			if id != "" {
				return fmt.Errorf("basic fact element_id must be empty, got %q", id)
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

// CreateFact creates a fact and its derived cards in a transaction.
func (app *Application) CreateFact(ctx context.Context, userID, notebookID uuid.UUID, factType string, content json.RawMessage) (db.GetFactRow, []db.Card, error) {
	elementIDs, err := validateFactContent(factType, content)
	if err != nil {
		return db.GetFactRow{}, nil, err
	}

	var fact db.AppFact
	var cards []db.Card

	err = WithTx(ctx, app.DB, func(q *db.Queries) error {
		var txErr error
		fact, txErr = q.CreateFact(ctx, db.CreateFactParams{
			UserID:     userID,
			NotebookID: notebookID,
			FactType:   factType,
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
				FactID:     fact.ID,
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
		return db.GetFactRow{}, nil, err
	}

	// Return as GetFactRow with card count
	row := db.GetFactRow{
		UserID:     fact.UserID,
		ID:         fact.ID,
		NotebookID: fact.NotebookID,
		FactType:   fact.FactType,
		Content:    fact.Content,
		SourceID:   fact.SourceID,
		CreatedAt:  fact.CreatedAt,
		UpdatedAt:  fact.UpdatedAt,
		CardCount:  int64(len(cards)),
	}
	return row, cards, nil
}

// GetFact retrieves a single fact with card count.
func (app *Application) GetFact(ctx context.Context, userID, notebookID, factID uuid.UUID) (db.GetFactRow, error) {
	fact, err := app.Queries.GetFact(ctx, db.GetFactParams{
		UserID:     userID,
		ID:         factID,
		NotebookID: notebookID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.GetFactRow{}, errFactNotFound
		}
		return db.GetFactRow{}, err
	}
	return fact, nil
}

// ListFacts retrieves paginated facts for a notebook.
func (app *Application) ListFacts(ctx context.Context, userID, notebookID uuid.UUID, limit, offset int32) ([]db.ListFactsByNotebookRow, int64, error) {
	facts, err := app.Queries.ListFactsByNotebook(ctx, db.ListFactsByNotebookParams{
		UserID:     userID,
		NotebookID: notebookID,
		RowLimit:   limit,
		RowOffset:  offset,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := app.Queries.CountFactsByNotebook(ctx, db.CountFactsByNotebookParams{
		UserID:     userID,
		NotebookID: notebookID,
	})
	if err != nil {
		return nil, 0, err
	}

	return facts, total, nil
}

// UpdateFactResult contains the outcome of a fact update.
type UpdateFactResult struct {
	Fact      db.AppFact
	Created   int
	Deleted   int
	Unchanged int
}

// UpdateFact updates fact content and diffs cards atomically.
func (app *Application) UpdateFact(ctx context.Context, userID, notebookID, factID uuid.UUID, content json.RawMessage) (UpdateFactResult, error) {
	// Fetch existing fact to verify ownership and get type
	existingFact, err := app.Queries.GetFact(ctx, db.GetFactParams{
		UserID:     userID,
		ID:         factID,
		NotebookID: notebookID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UpdateFactResult{}, errFactNotFound
		}
		return UpdateFactResult{}, err
	}

	expectedIDs, err := validateFactContent(existingFact.FactType, content)
	if err != nil {
		return UpdateFactResult{}, err
	}

	var result UpdateFactResult

	err = WithTx(ctx, app.DB, func(q *db.Queries) error {
		// Update the fact content
		fact, txErr := q.UpdateFactContent(ctx, db.UpdateFactContentParams{
			Content:    content,
			UserID:     userID,
			ID:         factID,
			NotebookID: notebookID,
		})
		if txErr != nil {
			return txErr
		}
		result.Fact = fact

		// Fetch existing cards to diff
		existingCards, txErr := q.ListCardsByFact(ctx, db.ListCardsByFactParams{
			UserID: userID,
			FactID: factID,
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
			if txErr := q.DeleteCardsByFactAndElements(ctx, db.DeleteCardsByFactAndElementsParams{
				UserID:     userID,
				FactID:     factID,
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
				FactID:     factID,
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

// FactListFilters holds optional filters for listing facts.
type FactListFilters struct {
	FactType string    // empty = all
	Search   string    // empty = no search
	Sort     SortParam // Field: "created", "updated"; empty = default (updated desc)
}

// ListFactsFiltered retrieves paginated facts with optional filters and due counts.
func (app *Application) ListFactsFiltered(ctx context.Context, userID, notebookID uuid.UUID, filters FactListFilters, limit, offset int32) ([]db.ListFactsByNotebookFilteredRow, int64, error) {
	params := db.ListFactsByNotebookFilteredParams{
		UserID:     userID,
		NotebookID: notebookID,
		RowLimit:   limit,
		RowOffset:  offset,
	}
	if filters.FactType != "" {
		params.FactType = pgtype.Text{String: filters.FactType, Valid: true}
	}
	if filters.Search != "" {
		params.Search = pgtype.Text{String: filters.Search, Valid: true}
	}
	params.SortField = filters.Sort.Field
	params.SortAsc = filters.Sort.Asc

	facts, err := app.Queries.ListFactsByNotebookFiltered(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	countParams := db.CountFactsByNotebookFilteredParams{
		UserID:     userID,
		NotebookID: notebookID,
		FactType:   params.FactType,
		Search:     params.Search,
	}
	total, err := app.Queries.CountFactsByNotebookFiltered(ctx, countParams)
	if err != nil {
		return nil, 0, err
	}

	return facts, total, nil
}

// GetFactStats returns aggregate stats for a notebook.
func (app *Application) GetFactStats(ctx context.Context, userID, notebookID uuid.UUID) (db.GetFactStatsByNotebookRow, error) {
	return app.Queries.GetFactStatsByNotebook(ctx, db.GetFactStatsByNotebookParams{
		UserID:     userID,
		NotebookID: notebookID,
	})
}

// DeleteFact deletes a fact (cascades to cards).
func (app *Application) DeleteFact(ctx context.Context, userID, notebookID, factID uuid.UUID) error {
	rows, err := app.Queries.DeleteFact(ctx, db.DeleteFactParams{
		UserID:     userID,
		ID:         factID,
		NotebookID: notebookID,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return errFactNotFound
	}
	return nil
}
