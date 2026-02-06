//go:build integration

package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"apexmemory.ai/internal/db"
	"github.com/google/uuid"
)

func TestListCards_StateFilterAffectsTotal(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)

	// Review the card to move it out of 'new' state
	req := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()
	app.SubmitReviewHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("review failed: %d. Body: %s", rr.Code, rr.Body.String())
	}

	// Create another new card in the same notebook
	_, _, err := app.CreateFact(t.Context(), user.ID, nbID, "basic", basicContent())
	if err != nil {
		t.Fatalf("failed to create second fact: %v", err)
	}

	// Filter by state=new should return only 1 card and total=1
	newState := db.NullAppCardState{AppCardState: db.AppCardStateNew, Valid: true}
	cards, total, err := app.ListCards(t.Context(), user.ID, nbID, newState, 50, 0)
	if err != nil {
		t.Fatalf("ListCards failed: %v", err)
	}

	if total != 1 {
		t.Errorf("expected total 1 for state=new filter, got %d", total)
	}
	if len(cards) != 1 {
		t.Errorf("expected 1 card for state=new filter, got %d", len(cards))
	}

	// Without filter should return 2
	noFilter := db.NullAppCardState{}
	_, totalAll, err := app.ListCards(t.Context(), user.ID, nbID, noFilter, 50, 0)
	if err != nil {
		t.Fatalf("ListCards (all) failed: %v", err)
	}
	if totalAll != 2 {
		t.Errorf("expected total 2 without filter, got %d", totalAll)
	}
}

func TestListCards_BuriedUntilPastDate(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)

	// Set buried_until to yesterday via raw SQL
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err := app.DB.Exec(t.Context(),
		"UPDATE app.cards SET buried_until = $1 WHERE user_id = $2 AND id = $3",
		yesterday, user.ID, cardID)
	if err != nil {
		t.Fatalf("failed to set buried_until: %v", err)
	}

	// Card should appear in list (auto-unburied)
	noFilter := db.NullAppCardState{}
	cards, total, err := app.ListCards(t.Context(), user.ID, nbID, noFilter, 50, 0)
	if err != nil {
		t.Fatalf("ListCards failed: %v", err)
	}
	if total != 1 {
		t.Errorf("expected total 1 (past buried_until should unbury), got %d", total)
	}
	if len(cards) != 1 {
		t.Errorf("expected 1 card, got %d", len(cards))
	}
}

func TestListCards_BuriedUntilFutureDate(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)

	// Set buried_until to 2 days from now (avoids UTC/local timezone boundary issues)
	future := time.Now().UTC().AddDate(0, 0, 2).Format("2006-01-02")
	_, err := app.DB.Exec(t.Context(),
		"UPDATE app.cards SET buried_until = $1 WHERE user_id = $2 AND id = $3",
		future, user.ID, cardID)
	if err != nil {
		t.Fatalf("failed to set buried_until: %v", err)
	}

	// Card should be excluded
	noFilter := db.NullAppCardState{}
	cards, total, err := app.ListCards(t.Context(), user.ID, nbID, noFilter, 50, 0)
	if err != nil {
		t.Fatalf("ListCards failed: %v", err)
	}
	if total != 0 {
		t.Errorf("expected total 0 (future buried_until should exclude), got %d", total)
	}
	if len(cards) != 0 {
		t.Errorf("expected 0 cards, got %d", len(cards))
	}
}
