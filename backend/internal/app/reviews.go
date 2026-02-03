package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"apexmemory.ai/internal/db"
	"apexmemory.ai/internal/fsrs"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	errCardNotReviewable  = errors.New("card not found or not reviewable")
	errInvalidRating      = errors.New("invalid rating")
	errInvalidMode        = errors.New("invalid mode")
	errReviewNotFound     = errors.New("review not found")
	errReviewNotLatest    = errors.New("review is not the latest for this card")
	errCardAlreadyDeleted = errors.New("card has been deleted")
)

const defaultNewCardCap int64 = 20

// ReviewRequest is the JSON body for POST /v1/reviews.
type ReviewRequest struct {
	ID         uuid.UUID `json:"id"`
	CardID     uuid.UUID `json:"card_id"`
	Rating     string    `json:"rating"`
	DurationMs *int32    `json:"duration_ms"`
	Mode       string    `json:"mode"`
}

// ReviewResponse is the API response after submitting a review.
type ReviewResponse struct {
	Review ReviewData   `json:"review"`
	Card   CardResponse `json:"card"`
}

// StudyCountsResponse is the API response for GET /v1/reviews/study-counts.
type StudyCountsResponse struct {
	Counts   map[string]NotebookStudyCounts `json:"counts"`
	TotalDue int64                          `json:"total_due"`
	TotalNew int64                          `json:"total_new"`
}

// NotebookStudyCounts contains study counts for a single notebook.
type NotebookStudyCounts struct {
	Due   int64 `json:"due"`
	New   int64 `json:"new"`
	Total int32 `json:"total"`
}

// ReviewData is the review portion of the response.
type ReviewData struct {
	ID         uuid.UUID `json:"id"`
	CardID     uuid.UUID `json:"card_id"`
	Rating     string    `json:"rating"`
	Mode       string    `json:"mode"`
	ReviewedAt time.Time `json:"reviewed_at"`
}

// UndoSnapshot captures card state fields needed for undo restoration.
// This is stored as JSONB in the review record.
type UndoSnapshot struct {
	Step       *int16     `json:"step"`
	Due        *time.Time `json:"due"`
	LastReview *time.Time `json:"last_review"`
	Reps       int32      `json:"reps"`
	Lapses     int32      `json:"lapses"`
}

// UndoReviewResponse is returned after undoing a review.
type UndoReviewResponse struct {
	Card *CardResponse `json:"card"` // nil for practice mode
}

// StudyCard is a card with precomputed intervals for all 4 ratings.
type StudyCard struct {
	CardResponse
	FactType    string          `json:"fact_type"`
	FactContent json.RawMessage `json:"fact_content"`
	Intervals   IntervalPreview `json:"intervals"`
}

// IntervalPreview shows the next interval for each rating.
type IntervalPreview struct {
	Again string `json:"again"`
	Hard  string `json:"hard"`
	Good  string `json:"good"`
	Easy  string `json:"easy"`
}

// submitReview processes a review submission within a transaction.
func (app *Application) submitReview(ctx context.Context, userID uuid.UUID, req ReviewRequest) (ReviewResponse, error) {
	rating, err := ratingFromString(req.Rating)
	if err != nil {
		return ReviewResponse{}, err
	}

	if req.Mode == "" {
		req.Mode = "scheduled"
	}
	if req.Mode != "scheduled" && req.Mode != "practice" {
		return ReviewResponse{}, fmt.Errorf("invalid mode %q: %w", req.Mode, errInvalidMode)
	}

	now := time.Now().UTC()
	var review db.AppReview
	var updatedCard db.Card

	err = WithTx(ctx, app.DB, func(q *db.Queries) error {
		// Lock the card
		card, txErr := q.GetCardForReview(ctx, db.GetCardForReviewParams{
			UserID: userID,
			ID:     req.CardID,
		})
		if txErr != nil {
			if errors.Is(txErr, pgx.ErrNoRows) {
				return errCardNotReviewable
			}
			return txErr
		}

		// Capture undo snapshot before processing
		undoSnapshot := UndoSnapshot{
			Reps:   card.Reps,
			Lapses: card.Lapses,
		}
		if card.Step.Valid {
			step := card.Step.Int16
			undoSnapshot.Step = &step
		}
		if card.Due.Valid {
			due := card.Due.Time
			undoSnapshot.Due = &due
		}
		if card.LastReview.Valid {
			lr := card.LastReview.Time
			undoSnapshot.LastReview = &lr
		}

		// Map DB card to FSRS card
		fsrsCard := dbCardToFSRS(card)

		// Build scheduler (no fuzzing for actual review to keep deterministic for now;
		// we can enable later)
		scheduler, txErr := fsrs.NewScheduler(fsrs.WithEnableFuzzing(true))
		if txErr != nil {
			return txErr
		}

		// Run FSRS
		output := scheduler.ReviewCard(fsrsCard, rating, now)

		// Capture before state
		stateBefore := card.State
		stabilityBefore := card.Stability
		difficultyBefore := card.Difficulty

		// Build after state
		stateAfter := fsrsStateToDBState(output.Card.State)
		stabilityAfter := float32(*output.Card.Stability)
		difficultyAfter := float32(*output.Card.Difficulty)
		intervalDays := float32(output.ScheduledDays)

		var retrievability pgtype.Float4
		if output.Retrievability > 0 {
			retrievability = pgtype.Float4{Float32: float32(output.Retrievability), Valid: true}
		}

		var durationMs pgtype.Int4
		if req.DurationMs != nil {
			durationMs = pgtype.Int4{Int32: *req.DurationMs, Valid: true}
		}

		// For practice mode, before=after (card state doesn't change)
		actualStateAfter := stateAfter
		actualStabilityAfter := stabilityAfter
		actualDifficultyAfter := difficultyAfter
		actualIntervalDays := intervalDays
		if req.Mode == "practice" {
			actualStateAfter = stateBefore
			actualStabilityAfter = stabilityBefore.Float32
			actualDifficultyAfter = difficultyBefore.Float32
			actualIntervalDays = card.ScheduledDays
		}

		// Serialize undo snapshot to JSON
		undoSnapshotJSON, txErr := json.Marshal(undoSnapshot)
		if txErr != nil {
			return fmt.Errorf("marshal undo snapshot: %w", txErr)
		}

		// Insert review (ON CONFLICT DO NOTHING — returns pgx.ErrNoRows on conflict)
		review, txErr = q.CreateReview(ctx, db.CreateReviewParams{
			UserID:           userID,
			ID:               req.ID,
			CardID:           pgtype.UUID{Bytes: card.ID, Valid: true},
			NotebookID:       card.NotebookID,
			FactID:           pgtype.UUID{Bytes: card.FactID, Valid: true},
			ElementID:        pgtype.Text{String: card.ElementID, Valid: true},
			ReviewedAt:       now,
			Rating:           db.AppRating(req.Rating),
			ReviewDurationMs: durationMs,
			Mode:             req.Mode,
			StateBefore:      stateBefore,
			StabilityBefore:  stabilityBefore,
			DifficultyBefore: difficultyBefore,
			ElapsedDays:      float32(output.ElapsedDays),
			ScheduledDays:    float32(output.ScheduledDays),
			StateAfter:       actualStateAfter,
			StabilityAfter:   actualStabilityAfter,
			DifficultyAfter:  actualDifficultyAfter,
			IntervalDays:     actualIntervalDays,
			Retrievability:   retrievability,
			UndoSnapshot:     undoSnapshotJSON,
		})
		if txErr != nil {
			if errors.Is(txErr, pgx.ErrNoRows) {
				// Idempotent re-submit: conflict on (user_id, id)
				updatedCard = card
				review.ID = req.ID
				review.Rating = db.AppRating(req.Rating)
				review.Mode = req.Mode
				review.ReviewedAt = now
				review.CardID = pgtype.UUID{Bytes: card.ID, Valid: true}
				return nil
			}
			return txErr
		}

		// Should not happen after the above, but guard anyway
		if review.ID == uuid.Nil {
			updatedCard = card
			// Reconstruct minimal review data for response
			review.ID = req.ID
			review.Rating = db.AppRating(req.Rating)
			review.Mode = req.Mode
			review.ReviewedAt = now
			review.CardID = pgtype.UUID{Bytes: card.ID, Valid: true}
			return nil
		}

		// Update card state (only for scheduled mode)
		if req.Mode == "scheduled" {
			isLapse := rating == fsrs.Again && (card.State == db.AppCardStateReview || card.State == db.AppCardStateRelearning)

			var step pgtype.Int2
			if output.Card.Step != nil {
				step = pgtype.Int2{Int16: int16(*output.Card.Step), Valid: true}
			}

			txErr = q.UpdateCardAfterReview(ctx, db.UpdateCardAfterReviewParams{
				State:         stateAfter,
				Stability:     pgtype.Float4{Float32: stabilityAfter, Valid: true},
				Difficulty:    pgtype.Float4{Float32: difficultyAfter, Valid: true},
				Step:          step,
				Due:           pgtype.Timestamptz{Time: output.Card.Due, Valid: true},
				LastReview:    pgtype.Timestamptz{Time: now, Valid: true},
				ElapsedDays:   float32(output.ElapsedDays),
				ScheduledDays: float32(output.ScheduledDays),
				AddLapse:      isLapse,
				UserID:        userID,
				ID:            card.ID,
			})
			if txErr != nil {
				return txErr
			}

			// Build updated card for response
			updatedCard = card
			updatedCard.State = stateAfter
			updatedCard.Stability = pgtype.Float4{Float32: stabilityAfter, Valid: true}
			updatedCard.Difficulty = pgtype.Float4{Float32: difficultyAfter, Valid: true}
			updatedCard.Step = step
			updatedCard.Due = pgtype.Timestamptz{Time: output.Card.Due, Valid: true}
			updatedCard.LastReview = pgtype.Timestamptz{Time: now, Valid: true}
			updatedCard.ElapsedDays = float32(output.ElapsedDays)
			updatedCard.ScheduledDays = float32(output.ScheduledDays)
			updatedCard.Reps = card.Reps + 1
			if isLapse {
				updatedCard.Lapses = card.Lapses + 1
			}
		} else {
			updatedCard = card
		}

		return nil
	})
	if err != nil {
		return ReviewResponse{}, err
	}

	return ReviewResponse{
		Review: ReviewData{
			ID:         review.ID,
			CardID:     req.CardID,
			Rating:     req.Rating,
			Mode:       req.Mode,
			ReviewedAt: review.ReviewedAt,
		},
		Card: toCardResponse(updatedCard),
	}, nil
}

// getStudyCards returns due cards with precomputed intervals.
func (app *Application) getStudyCards(ctx context.Context, userID uuid.UUID, notebookID *uuid.UUID, limit int32) ([]StudyCard, error) {
	var nbID pgtype.UUID
	if notebookID != nil {
		nbID = pgtype.UUID{Bytes: *notebookID, Valid: true}
	}

	rows, err := app.Queries.GetStudyCards(ctx, db.GetStudyCardsParams{
		UserID:     userID,
		NotebookID: nbID,
		NewCardCap: defaultNewCardCap,
		RowLimit:   limit,
	})
	if err != nil {
		return nil, err
	}

	scheduler, err := fsrs.NewScheduler(fsrs.WithEnableFuzzing(false))
	if err != nil {
		return nil, err
	}

	result := make([]StudyCard, len(rows))
	for i, row := range rows {
		card := studyRowToCard(row)
		result[i] = StudyCard{
			CardResponse: toCardResponse(card),
			FactType:     row.FactType,
			FactContent:  json.RawMessage(row.FactContent),
			Intervals:    computeIntervalsWithScheduler(card, scheduler),
		}
	}
	return result, nil
}

// getPracticeCards returns all cards for practice mode with intervals.
func (app *Application) getPracticeCards(ctx context.Context, userID uuid.UUID, notebookID *uuid.UUID, limit, offset int32) ([]StudyCard, int64, error) {
	var nbID pgtype.UUID
	if notebookID != nil {
		nbID = pgtype.UUID{Bytes: *notebookID, Valid: true}
	}

	rows, err := app.Queries.GetPracticeCards(ctx, db.GetPracticeCardsParams{
		UserID:     userID,
		NotebookID: nbID,
		RowLimit:   limit,
		RowOffset:  offset,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := app.Queries.CountPracticeCards(ctx, db.CountPracticeCardsParams{
		UserID:     userID,
		NotebookID: nbID,
	})
	if err != nil {
		return nil, 0, err
	}

	scheduler, err := fsrs.NewScheduler(fsrs.WithEnableFuzzing(false))
	if err != nil {
		return nil, 0, err
	}

	result := make([]StudyCard, len(rows))
	for i, row := range rows {
		card := practiceRowToCard(row)
		result[i] = StudyCard{
			CardResponse: toCardResponse(card),
			FactType:     row.FactType,
			FactContent:  json.RawMessage(row.FactContent),
			Intervals:    computeIntervalsWithScheduler(card, scheduler),
		}
	}
	return result, total, nil
}

// computeIntervalsWithScheduler computes interval previews using a provided scheduler.
func computeIntervalsWithScheduler(card db.Card, scheduler *fsrs.Scheduler) IntervalPreview {
	fsrsCard := dbCardToFSRS(card)
	now := time.Now().UTC()
	return IntervalPreview{
		Again: formatInterval(scheduler.ReviewCard(fsrsCard, fsrs.Again, now).Card.Due.Sub(now)),
		Hard:  formatInterval(scheduler.ReviewCard(fsrsCard, fsrs.Hard, now).Card.Due.Sub(now)),
		Good:  formatInterval(scheduler.ReviewCard(fsrsCard, fsrs.Good, now).Card.Due.Sub(now)),
		Easy:  formatInterval(scheduler.ReviewCard(fsrsCard, fsrs.Easy, now).Card.Due.Sub(now)),
	}
}

// formatInterval converts a duration to a human-readable string.
func formatInterval(d time.Duration) string {
	if d < time.Hour {
		m := int(math.Round(d.Minutes()))
		if m < 1 {
			m = 1
		}
		return fmt.Sprintf("%dm", m)
	}
	if d < 24*time.Hour {
		h := int(math.Round(d.Hours()))
		if h < 1 {
			h = 1
		}
		return fmt.Sprintf("%dh", h)
	}
	days := d.Hours() / 24
	if days < 30 {
		d := int(math.Round(days))
		if d < 1 {
			d = 1
		}
		return fmt.Sprintf("%dd", d)
	}
	months := int(math.Round(days / 30))
	if months < 1 {
		months = 1
	}
	return fmt.Sprintf("%dmo", months)
}

// dbCardToFSRS converts a database Card to an FSRS Card.
func dbCardToFSRS(card db.Card) fsrs.Card {
	fc := fsrs.Card{
		Due: card.Due.Time,
	}

	if card.LastReview.Valid {
		t := card.LastReview.Time
		fc.LastReview = &t
	}

	switch card.State {
	case db.AppCardStateNew:
		fc.State = fsrs.Learning
		step := 0
		fc.Step = &step
	case db.AppCardStateLearning:
		fc.State = fsrs.Learning
		if card.Step.Valid {
			step := int(card.Step.Int16)
			fc.Step = &step
		}
	case db.AppCardStateReview:
		fc.State = fsrs.Review
	case db.AppCardStateRelearning:
		fc.State = fsrs.Relearning
		if card.Step.Valid {
			step := int(card.Step.Int16)
			fc.Step = &step
		}
	}

	if card.Stability.Valid {
		s := float64(card.Stability.Float32)
		fc.Stability = &s
	}
	if card.Difficulty.Valid {
		d := float64(card.Difficulty.Float32)
		fc.Difficulty = &d
	}

	return fc
}

// fsrsStateToDBState converts FSRS state to database card state.
func fsrsStateToDBState(state fsrs.State) db.AppCardState {
	switch state {
	case fsrs.Learning:
		return db.AppCardStateLearning
	case fsrs.Review:
		return db.AppCardStateReview
	case fsrs.Relearning:
		return db.AppCardStateRelearning
	default:
		return db.AppCardStateLearning
	}
}

// ratingFromString converts a string rating to FSRS Rating.
func ratingFromString(s string) (fsrs.Rating, error) {
	switch s {
	case "again":
		return fsrs.Again, nil
	case "hard":
		return fsrs.Hard, nil
	case "good":
		return fsrs.Good, nil
	case "easy":
		return fsrs.Easy, nil
	default:
		return 0, fmt.Errorf("invalid rating %q: %w", s, errInvalidRating)
	}
}

// studyRowToCard converts a GetStudyCardsRow to a db.Card.
func studyRowToCard(row db.GetStudyCardsRow) db.Card {
	return db.Card{
		UserID:        row.UserID,
		ID:            row.ID,
		NotebookID:    row.NotebookID,
		FactID:        row.FactID,
		ElementID:     row.ElementID,
		State:         row.State,
		Stability:     row.Stability,
		Difficulty:    row.Difficulty,
		Step:          row.Step,
		Due:           row.Due,
		LastReview:    row.LastReview,
		ElapsedDays:   row.ElapsedDays,
		ScheduledDays: row.ScheduledDays,
		Reps:          row.Reps,
		Lapses:        row.Lapses,
		SuspendedAt:   row.SuspendedAt,
		BuriedUntil:   row.BuriedUntil,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

// practiceRowToCard converts a GetPracticeCardsRow to a db.Card.
func practiceRowToCard(row db.GetPracticeCardsRow) db.Card {
	return db.Card{
		UserID:        row.UserID,
		ID:            row.ID,
		NotebookID:    row.NotebookID,
		FactID:        row.FactID,
		ElementID:     row.ElementID,
		State:         row.State,
		Stability:     row.Stability,
		Difficulty:    row.Difficulty,
		Step:          row.Step,
		Due:           row.Due,
		LastReview:    row.LastReview,
		ElapsedDays:   row.ElapsedDays,
		ScheduledDays: row.ScheduledDays,
		Reps:          row.Reps,
		Lapses:        row.Lapses,
		SuspendedAt:   row.SuspendedAt,
		BuriedUntil:   row.BuriedUntil,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

// getStudyCounts returns due and new card counts grouped by notebook.
func (app *Application) getStudyCounts(ctx context.Context, userID uuid.UUID) (StudyCountsResponse, error) {
	// Get counts per notebook from cards table
	countRows, err := app.Queries.GetStudyCountsByNotebook(ctx, userID)
	if err != nil {
		return StudyCountsResponse{}, err
	}

	// Get all notebooks to include total_cards and ensure all notebooks appear
	notebooks, err := app.ListNotebooks(ctx, userID)
	if err != nil {
		return StudyCountsResponse{}, err
	}

	// Build map with all notebooks (zero counts for those with no due/new cards)
	counts := make(map[string]NotebookStudyCounts, len(notebooks))
	for _, nb := range notebooks {
		counts[nb.ID.String()] = NotebookStudyCounts{
			Due:   0,
			New:   0,
			Total: nb.TotalCards,
		}
	}

	// Overlay actual counts from query results
	var totalDue, totalNew int64
	for _, row := range countRows {
		nbID := row.NotebookID.String()
		existing := counts[nbID]
		counts[nbID] = NotebookStudyCounts{
			Due:   row.DueCount,
			New:   row.NewCount,
			Total: existing.Total,
		}
		totalDue += row.DueCount
		totalNew += row.NewCount
	}

	return StudyCountsResponse{
		Counts:   counts,
		TotalDue: totalDue,
		TotalNew: totalNew,
	}, nil
}

// undoReview deletes a review and restores the card to its pre-review state.
// For practice mode: just deletes the review (card state wasn't changed).
// For scheduled mode: verifies this is the latest review, restores card, deletes review.
func (app *Application) undoReview(ctx context.Context, userID uuid.UUID, reviewID uuid.UUID) (UndoReviewResponse, error) {
	var result UndoReviewResponse

	err := WithTx(ctx, app.DB, func(q *db.Queries) error {
		// Fetch the review
		review, txErr := q.GetReviewByID(ctx, db.GetReviewByIDParams{
			UserID: userID,
			ID:     reviewID,
		})
		if txErr != nil {
			if errors.Is(txErr, pgx.ErrNoRows) {
				return errReviewNotFound
			}
			return txErr
		}

		// Practice mode: just delete the review, no card state change
		if review.Mode == "practice" {
			_, txErr = q.DeleteReview(ctx, db.DeleteReviewParams{
				UserID: userID,
				ID:     reviewID,
			})
			return txErr
		}

		// Scheduled mode: need to restore card state
		if !review.CardID.Valid {
			return errCardAlreadyDeleted
		}
		cardID := review.CardID.Bytes

		// Verify this is the latest review for the card
		latestID, txErr := q.GetLatestReviewForCard(ctx, db.GetLatestReviewForCardParams{
			UserID: userID,
			CardID: pgtype.UUID{Bytes: cardID, Valid: true},
		})
		if txErr != nil {
			return fmt.Errorf("get latest review: %w", txErr)
		}
		if latestID != reviewID {
			return errReviewNotLatest
		}

		// Lock and fetch the card
		card, txErr := q.GetCardForReview(ctx, db.GetCardForReviewParams{
			UserID: userID,
			ID:     cardID,
		})
		if txErr != nil {
			if errors.Is(txErr, pgx.ErrNoRows) {
				return errCardAlreadyDeleted
			}
			return txErr
		}

		// Parse undo snapshot
		var snapshot UndoSnapshot
		if len(review.UndoSnapshot) > 0 {
			if txErr = json.Unmarshal(review.UndoSnapshot, &snapshot); txErr != nil {
				return fmt.Errorf("unmarshal undo snapshot: %w", txErr)
			}
		}

		// Restore card state from before columns + undo snapshot
		var step pgtype.Int2
		if snapshot.Step != nil {
			step = pgtype.Int2{Int16: *snapshot.Step, Valid: true}
		}
		var due pgtype.Timestamptz
		if snapshot.Due != nil {
			due = pgtype.Timestamptz{Time: *snapshot.Due, Valid: true}
		}
		var lastReview pgtype.Timestamptz
		if snapshot.LastReview != nil {
			lastReview = pgtype.Timestamptz{Time: *snapshot.LastReview, Valid: true}
		}

		txErr = q.RestoreCardAfterUndo(ctx, db.RestoreCardAfterUndoParams{
			State:         review.StateBefore,
			Stability:     review.StabilityBefore,
			Difficulty:    review.DifficultyBefore,
			Step:          step,
			Due:           due,
			LastReview:    lastReview,
			ElapsedDays:   review.ElapsedDays,
			ScheduledDays: review.ScheduledDays,
			Reps:          snapshot.Reps,
			Lapses:        snapshot.Lapses,
			UserID:        userID,
			ID:            cardID,
		})
		if txErr != nil {
			return fmt.Errorf("restore card: %w", txErr)
		}

		// Delete the review
		_, txErr = q.DeleteReview(ctx, db.DeleteReviewParams{
			UserID: userID,
			ID:     reviewID,
		})
		if txErr != nil {
			return fmt.Errorf("delete review: %w", txErr)
		}

		// Build restored card for response
		restoredCard := card
		restoredCard.State = review.StateBefore
		restoredCard.Stability = review.StabilityBefore
		restoredCard.Difficulty = review.DifficultyBefore
		restoredCard.Step = step
		restoredCard.Due = due
		restoredCard.LastReview = lastReview
		restoredCard.ElapsedDays = review.ElapsedDays
		restoredCard.ScheduledDays = review.ScheduledDays
		restoredCard.Reps = snapshot.Reps
		restoredCard.Lapses = snapshot.Lapses

		cardResp := toCardResponse(restoredCard)
		result.Card = &cardResp
		return nil
	})

	return result, err
}
