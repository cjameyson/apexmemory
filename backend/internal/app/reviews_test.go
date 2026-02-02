//go:build integration

package app

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

// createTestCard creates a notebook, basic fact, and returns the notebook ID and card ID.
func createTestCard(t *testing.T, app *Application, userID uuid.UUID) (notebookID, cardID uuid.UUID) {
	t.Helper()
	nbID := createTestNotebook(t, app, userID)
	_, cards, err := app.CreateFact(t.Context(), userID, nbID, "basic", basicContent())
	if err != nil {
		t.Fatalf("failed to create fact: %v", err)
	}
	return nbID, cards[0].ID
}

func TestSubmitReviewHandler_Scheduled(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	_, cardID := createTestCard(t, app, user.ID)
	reviewID := uuid.New()

	req := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      reviewID.String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.SubmitReviewHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	decodeResponse(t, rr, &resp)

	// Verify review data
	review := resp["review"].(map[string]any)
	if review["rating"] != "good" {
		t.Errorf("expected rating 'good', got %v", review["rating"])
	}
	if review["mode"] != "scheduled" {
		t.Errorf("expected mode 'scheduled', got %v", review["mode"])
	}

	// Verify card state changed from 'new'
	card := resp["card"].(map[string]any)
	if card["state"] == "new" {
		t.Error("expected card state to change from 'new' after scheduled review")
	}
	if card["stability"] == nil {
		t.Error("expected stability to be set after review")
	}
	if card["difficulty"] == nil {
		t.Error("expected difficulty to be set after review")
	}
}

func TestSubmitReviewHandler_Practice(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	_, cardID := createTestCard(t, app, user.ID)
	reviewID := uuid.New()

	req := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      reviewID.String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "practice",
	})
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.SubmitReviewHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	decodeResponse(t, rr, &resp)

	// Card state should remain 'new' for practice mode
	card := resp["card"].(map[string]any)
	if card["state"] != "new" {
		t.Errorf("expected card state to remain 'new' for practice review, got %v", card["state"])
	}
}

func TestSubmitReviewHandler_Idempotent(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	_, cardID := createTestCard(t, app, user.ID)
	reviewID := uuid.New()

	body := map[string]any{
		"id":      reviewID.String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "scheduled",
	}

	// First submit
	req := jsonRequest(t, "POST", "/v1/reviews", body)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()
	app.SubmitReviewHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("first submit: expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	// Second submit with same ID
	req = jsonRequest(t, "POST", "/v1/reviews", body)
	req = app.WithUser(req, user)
	rr = httptest.NewRecorder()
	app.SubmitReviewHandler(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("second submit: expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}
}

func TestSubmitReviewHandler_CardNotFound(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	req := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": uuid.New().String(),
		"rating":  "good",
	})
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.SubmitReviewHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d. Body: %s", rr.Code, rr.Body.String())
	}
}

func TestSubmitReviewHandler_InvalidRating(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	_, cardID := createTestCard(t, app, user.ID)

	req := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": cardID.String(),
		"rating":  "invalid",
	})
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.SubmitReviewHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d. Body: %s", rr.Code, rr.Body.String())
	}
}

func TestGetStudyCardsHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, _ := createTestCard(t, app, user.ID)

	// New cards should appear in study queue
	req := httptest.NewRequest("GET", "/v1/reviews/study?notebook_id="+nbID.String(), nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetStudyCardsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp []map[string]any
	decodeResponse(t, rr, &resp)

	if len(resp) != 1 {
		t.Fatalf("expected 1 study card, got %d", len(resp))
	}

	// Verify intervals are present
	card := resp[0]
	intervals, ok := card["intervals"].(map[string]any)
	if !ok {
		t.Fatal("expected intervals object")
	}
	for _, key := range []string{"again", "hard", "good", "easy"} {
		if intervals[key] == nil || intervals[key] == "" {
			t.Errorf("expected non-empty interval for %s", key)
		}
	}
}

func TestGetStudyCardsHandler_NoCardsAfterReview(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)

	// Review the card
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

	// Study queue should be empty (card is now due in the future)
	req = httptest.NewRequest("GET", "/v1/reviews/study?notebook_id="+nbID.String(), nil)
	req = app.WithUser(req, user)
	rr = httptest.NewRecorder()

	app.GetStudyCardsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var resp []map[string]any
	decodeResponse(t, rr, &resp)

	if len(resp) != 0 {
		t.Errorf("expected 0 study cards after review, got %d", len(resp))
	}
}

func TestGetPracticeCardsHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)

	// Review the card so it's no longer due
	req := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()
	app.SubmitReviewHandler(rr, req)

	// Practice should still return the card
	req = httptest.NewRequest("GET", "/v1/reviews/practice?notebook_id="+nbID.String(), nil)
	req = app.WithUser(req, user)
	rr = httptest.NewRecorder()

	app.GetPracticeCardsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp map[string]any
	decodeResponse(t, rr, &resp)

	data, ok := resp["data"].([]any)
	if !ok {
		t.Fatal("expected data array")
	}
	if len(data) != 1 {
		t.Errorf("expected 1 practice card, got %d", len(data))
	}
}

func TestGetStudyCardsHandler_GlobalScope(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Create cards in two different notebooks
	createTestCard(t, app, user.ID)
	createTestCard(t, app, user.ID)

	// Global query (no notebook_id) should return both
	req := httptest.NewRequest("GET", "/v1/reviews/study", nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetStudyCardsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp []map[string]any
	decodeResponse(t, rr, &resp)

	if len(resp) != 2 {
		t.Errorf("expected 2 study cards across notebooks, got %d", len(resp))
	}
}

func TestFormatInterval(t *testing.T) {
	tests := []struct {
		name string
		d    int // seconds
		want string
	}{
		{"30 seconds", 30, "1m"},
		{"1 minute", 60, "1m"},
		{"10 minutes", 600, "10m"},
		{"59 minutes", 3540, "59m"},
		{"1 hour", 3600, "1h"},
		{"6 hours", 21600, "6h"},
		{"1 day", 86400, "1d"},
		{"3 days", 259200, "3d"},
		{"29 days", 2505600, "29d"},
		{"30 days", 2592000, "1mo"},
		{"60 days", 5184000, "2mo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatInterval(time.Duration(tt.d) * time.Second)
			if got != tt.want {
				t.Errorf("formatInterval(%ds) = %q, want %q", tt.d, got, tt.want)
			}
		})
	}
}
