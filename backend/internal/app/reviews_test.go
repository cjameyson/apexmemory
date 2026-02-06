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

func TestGetStudyCountsHandler_Empty(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Create notebook with no cards
	createTestNotebook(t, app, user.ID)

	req := httptest.NewRequest("GET", "/v1/reviews/study-counts", nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetStudyCountsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp StudyCountsResponse
	decodeResponse(t, rr, &resp)

	if resp.TotalDue != 0 {
		t.Errorf("expected total_due 0, got %d", resp.TotalDue)
	}
	if resp.TotalNew != 0 {
		t.Errorf("expected total_new 0, got %d", resp.TotalNew)
	}
	if len(resp.Counts) != 1 {
		t.Errorf("expected 1 notebook in counts, got %d", len(resp.Counts))
	}
}

func TestGetStudyCountsHandler_WithNewCards(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Create two notebooks with cards
	nbID1, _ := createTestCard(t, app, user.ID)
	nbID2, _ := createTestCard(t, app, user.ID)

	req := httptest.NewRequest("GET", "/v1/reviews/study-counts", nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetStudyCountsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp StudyCountsResponse
	decodeResponse(t, rr, &resp)

	// Both notebooks should have 1 new card each
	if resp.TotalNew != 2 {
		t.Errorf("expected total_new 2, got %d", resp.TotalNew)
	}
	if resp.TotalDue != 0 {
		t.Errorf("expected total_due 0, got %d", resp.TotalDue)
	}

	// Check individual notebook counts
	counts1, ok := resp.Counts[nbID1.String()]
	if !ok {
		t.Fatal("expected counts for notebook 1")
	}
	if counts1.New != 1 {
		t.Errorf("notebook 1: expected new 1, got %d", counts1.New)
	}
	if counts1.Total != 1 {
		t.Errorf("notebook 1: expected total 1, got %d", counts1.Total)
	}

	counts2, ok := resp.Counts[nbID2.String()]
	if !ok {
		t.Fatal("expected counts for notebook 2")
	}
	if counts2.New != 1 {
		t.Errorf("notebook 2: expected new 1, got %d", counts2.New)
	}
}

func TestGetStudyCountsHandler_AfterReview(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Create notebook with card
	nbID, cardID := createTestCard(t, app, user.ID)

	// Review the card (moves from new to learning)
	reviewReq := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	reviewReq = app.WithUser(reviewReq, user)
	rrReview := httptest.NewRecorder()
	app.SubmitReviewHandler(rrReview, reviewReq)

	if rrReview.Code != http.StatusOK {
		t.Fatalf("review failed: %d. Body: %s", rrReview.Code, rrReview.Body.String())
	}

	// Get counts - card is now not due (scheduled in future)
	req := httptest.NewRequest("GET", "/v1/reviews/study-counts", nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetStudyCountsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp StudyCountsResponse
	decodeResponse(t, rr, &resp)

	// After review, card is not new and not due
	if resp.TotalNew != 0 {
		t.Errorf("expected total_new 0 after review, got %d", resp.TotalNew)
	}
	if resp.TotalDue != 0 {
		t.Errorf("expected total_due 0 (card scheduled in future), got %d", resp.TotalDue)
	}

	// But total_cards should still be 1
	counts, ok := resp.Counts[nbID.String()]
	if !ok {
		t.Fatal("expected counts for notebook")
	}
	if counts.Total != 1 {
		t.Errorf("expected total 1, got %d", counts.Total)
	}
}

func TestGetStudyCountsHandler_TotalsMatchSum(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Create multiple notebooks with different card counts
	createTestCard(t, app, user.ID) // nb1: 1 new
	createTestCard(t, app, user.ID) // nb2: 1 new
	createTestCard(t, app, user.ID) // nb3: 1 new

	req := httptest.NewRequest("GET", "/v1/reviews/study-counts", nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetStudyCountsHandler(rr, req)

	var resp StudyCountsResponse
	decodeResponse(t, rr, &resp)

	// Sum individual counts
	var sumDue, sumNew int64
	for _, counts := range resp.Counts {
		sumDue += counts.Due
		sumNew += counts.New
	}

	if sumDue != resp.TotalDue {
		t.Errorf("sum of due counts (%d) != total_due (%d)", sumDue, resp.TotalDue)
	}
	if sumNew != resp.TotalNew {
		t.Errorf("sum of new counts (%d) != total_new (%d)", sumNew, resp.TotalNew)
	}
}

func TestUndoReviewHandler_Success(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)
	reviewID := uuid.New()

	// Submit a review
	submitReq := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      reviewID.String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	submitReq = app.WithUser(submitReq, user)
	rr := httptest.NewRecorder()
	app.SubmitReviewHandler(rr, submitReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("submit review failed: %d. Body: %s", rr.Code, rr.Body.String())
	}

	// Verify card is no longer in study queue
	studyReq := httptest.NewRequest("GET", "/v1/reviews/study?notebook_id="+nbID.String(), nil)
	studyReq = app.WithUser(studyReq, user)
	rr = httptest.NewRecorder()
	app.GetStudyCardsHandler(rr, studyReq)
	var studyCards []map[string]any
	decodeResponse(t, rr, &studyCards)
	if len(studyCards) != 0 {
		t.Fatalf("expected 0 study cards after review, got %d", len(studyCards))
	}

	// Undo the review
	undoReq := httptest.NewRequest("DELETE", "/v1/reviews/"+reviewID.String(), nil)
	undoReq.SetPathValue("id", reviewID.String())
	undoReq = app.WithUser(undoReq, user)
	rr = httptest.NewRecorder()
	app.UndoReviewHandler(rr, undoReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("undo review failed: %d. Body: %s", rr.Code, rr.Body.String())
	}

	var undoResp UndoReviewResponse
	decodeResponse(t, rr, &undoResp)

	// Card should be restored
	if undoResp.Card == nil {
		t.Fatal("expected card in undo response")
	}
	if undoResp.Card.State != "new" {
		t.Errorf("expected card state 'new' after undo, got %s", undoResp.Card.State)
	}

	// Card should be back in study queue
	studyReq = httptest.NewRequest("GET", "/v1/reviews/study?notebook_id="+nbID.String(), nil)
	studyReq = app.WithUser(studyReq, user)
	rr = httptest.NewRecorder()
	app.GetStudyCardsHandler(rr, studyReq)
	decodeResponse(t, rr, &studyCards)
	if len(studyCards) != 1 {
		t.Fatalf("expected 1 study card after undo, got %d", len(studyCards))
	}
}

func TestUndoReviewHandler_NotLatest(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	_, cardID := createTestCard(t, app, user.ID)
	reviewID1 := uuid.New()
	reviewID2 := uuid.New()

	// Submit first review
	submitReq := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      reviewID1.String(),
		"card_id": cardID.String(),
		"rating":  "again",
		"mode":    "scheduled",
	})
	submitReq = app.WithUser(submitReq, user)
	rr := httptest.NewRecorder()
	app.SubmitReviewHandler(rr, submitReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("first review failed: %d. Body: %s", rr.Code, rr.Body.String())
	}

	// Small delay to ensure different reviewed_at timestamp
	time.Sleep(10 * time.Millisecond)

	// Submit second review
	submitReq = jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      reviewID2.String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	submitReq = app.WithUser(submitReq, user)
	rr = httptest.NewRecorder()
	app.SubmitReviewHandler(rr, submitReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("second review failed: %d. Body: %s", rr.Code, rr.Body.String())
	}

	// Try to undo the first review (should fail with 409)
	undoReq := httptest.NewRequest("DELETE", "/v1/reviews/"+reviewID1.String(), nil)
	undoReq.SetPathValue("id", reviewID1.String())
	undoReq = app.WithUser(undoReq, user)
	rr = httptest.NewRecorder()
	app.UndoReviewHandler(rr, undoReq)
	if rr.Code != http.StatusConflict {
		t.Errorf("expected 409 Conflict, got %d. Body: %s", rr.Code, rr.Body.String())
	}
}

func TestUndoReviewHandler_NotFound(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Try to undo a non-existent review
	fakeID := uuid.New()
	undoReq := httptest.NewRequest("DELETE", "/v1/reviews/"+fakeID.String(), nil)
	undoReq.SetPathValue("id", fakeID.String())
	undoReq = app.WithUser(undoReq, user)
	rr := httptest.NewRecorder()
	app.UndoReviewHandler(rr, undoReq)
	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 Not Found, got %d. Body: %s", rr.Code, rr.Body.String())
	}
}

func TestUndoReviewHandler_PracticeMode(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)
	reviewID := uuid.New()

	// Submit a practice mode review
	submitReq := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      reviewID.String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "practice",
	})
	submitReq = app.WithUser(submitReq, user)
	rr := httptest.NewRecorder()
	app.SubmitReviewHandler(rr, submitReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("submit practice review failed: %d. Body: %s", rr.Code, rr.Body.String())
	}

	// Verify card is still in study queue (practice mode doesn't change card state)
	studyReq := httptest.NewRequest("GET", "/v1/reviews/study?notebook_id="+nbID.String(), nil)
	studyReq = app.WithUser(studyReq, user)
	rr = httptest.NewRecorder()
	app.GetStudyCardsHandler(rr, studyReq)
	var studyCards []map[string]any
	decodeResponse(t, rr, &studyCards)
	if len(studyCards) != 1 {
		t.Fatalf("expected 1 study card (practice doesn't change state), got %d", len(studyCards))
	}

	// Undo the practice review
	undoReq := httptest.NewRequest("DELETE", "/v1/reviews/"+reviewID.String(), nil)
	undoReq.SetPathValue("id", reviewID.String())
	undoReq = app.WithUser(undoReq, user)
	rr = httptest.NewRecorder()
	app.UndoReviewHandler(rr, undoReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("undo practice review failed: %d. Body: %s", rr.Code, rr.Body.String())
	}

	var undoResp UndoReviewResponse
	decodeResponse(t, rr, &undoResp)

	// Practice mode undo returns nil card (no state restoration needed)
	if undoResp.Card != nil {
		t.Error("expected nil card for practice mode undo")
	}

	// Card should still be in study queue
	studyReq = httptest.NewRequest("GET", "/v1/reviews/study?notebook_id="+nbID.String(), nil)
	studyReq = app.WithUser(studyReq, user)
	rr = httptest.NewRecorder()
	app.GetStudyCardsHandler(rr, studyReq)
	decodeResponse(t, rr, &studyCards)
	if len(studyCards) != 1 {
		t.Fatalf("expected 1 study card after practice undo, got %d", len(studyCards))
	}
}

// -----------------------------------------------------------------------------
// Review Summary Tests
// -----------------------------------------------------------------------------

func TestGetReviewSummaryHandler_Empty(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	req := httptest.NewRequest("GET", "/v1/reviews/summary", nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetReviewSummaryHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp ReviewSummaryResponse
	decodeResponse(t, rr, &resp)

	if resp.TotalReviews != 0 {
		t.Errorf("expected total_reviews 0, got %d", resp.TotalReviews)
	}
	if resp.RatingBreakdown.Again != 0 || resp.RatingBreakdown.Good != 0 {
		t.Errorf("expected zero rating counts, got %+v", resp.RatingBreakdown)
	}
	if resp.TotalDurationMs != 0 {
		t.Errorf("expected total_duration_ms 0, got %d", resp.TotalDurationMs)
	}
}

func TestGetReviewSummaryHandler_WithReviews(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	_, cardID := createTestCard(t, app, user.ID)

	// Submit multiple reviews with different ratings
	ratings := []string{"again", "hard", "good", "easy"}
	for i, rating := range ratings {
		reviewReq := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
			"id":          uuid.New().String(),
			"card_id":     cardID.String(),
			"rating":      rating,
			"duration_ms": 1000 + i*100,
			"mode":        "scheduled",
		})
		reviewReq = app.WithUser(reviewReq, user)
		rr := httptest.NewRecorder()
		app.SubmitReviewHandler(rr, reviewReq)
		if rr.Code != http.StatusOK {
			t.Fatalf("review %d failed: %d. Body: %s", i, rr.Code, rr.Body.String())
		}
		time.Sleep(10 * time.Millisecond)
	}

	req := httptest.NewRequest("GET", "/v1/reviews/summary", nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetReviewSummaryHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp ReviewSummaryResponse
	decodeResponse(t, rr, &resp)

	if resp.TotalReviews != 4 {
		t.Errorf("expected total_reviews 4, got %d", resp.TotalReviews)
	}
	if resp.RatingBreakdown.Again != 1 {
		t.Errorf("expected again_count 1, got %d", resp.RatingBreakdown.Again)
	}
	if resp.RatingBreakdown.Hard != 1 {
		t.Errorf("expected hard_count 1, got %d", resp.RatingBreakdown.Hard)
	}
	if resp.RatingBreakdown.Good != 1 {
		t.Errorf("expected good_count 1, got %d", resp.RatingBreakdown.Good)
	}
	if resp.RatingBreakdown.Easy != 1 {
		t.Errorf("expected easy_count 1, got %d", resp.RatingBreakdown.Easy)
	}
	if resp.ModeBreakdown.Scheduled != 4 {
		t.Errorf("expected scheduled_count 4, got %d", resp.ModeBreakdown.Scheduled)
	}
	// Duration: 1000 + 1100 + 1200 + 1300 = 4600
	if resp.TotalDurationMs != 4600 {
		t.Errorf("expected total_duration_ms 4600, got %d", resp.TotalDurationMs)
	}
	// First review was on a new card
	if resp.NewCardsSeen != 1 {
		t.Errorf("expected new_cards_seen 1, got %d", resp.NewCardsSeen)
	}
}

func TestGetReviewSummaryHandler_WithDateFilter(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	_, cardID := createTestCard(t, app, user.ID)

	// Submit a review today
	reviewReq := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	reviewReq = app.WithUser(reviewReq, user)
	rr := httptest.NewRecorder()
	app.SubmitReviewHandler(rr, reviewReq)
	if rr.Code != http.StatusOK {
		t.Fatalf("review failed: %d. Body: %s", rr.Code, rr.Body.String())
	}

	// Query for yesterday - should return zero
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	req := httptest.NewRequest("GET", "/v1/reviews/summary?date="+yesterday, nil)
	req = app.WithUser(req, user)
	rr = httptest.NewRecorder()

	app.GetReviewSummaryHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp ReviewSummaryResponse
	decodeResponse(t, rr, &resp)

	if resp.TotalReviews != 0 {
		t.Errorf("expected total_reviews 0 for yesterday, got %d", resp.TotalReviews)
	}
}

func TestGetReviewSummaryHandler_NotebookFilter(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Create cards in two notebooks
	nbID1, cardID1 := createTestCard(t, app, user.ID)
	_, cardID2 := createTestCard(t, app, user.ID)

	// Submit review in notebook 1
	reviewReq := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": cardID1.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	reviewReq = app.WithUser(reviewReq, user)
	rr := httptest.NewRecorder()
	app.SubmitReviewHandler(rr, reviewReq)

	// Submit review in notebook 2
	reviewReq = jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": cardID2.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	reviewReq = app.WithUser(reviewReq, user)
	rr = httptest.NewRecorder()
	app.SubmitReviewHandler(rr, reviewReq)

	// Query with notebook filter
	req := httptest.NewRequest("GET", "/v1/reviews/summary?notebook_id="+nbID1.String(), nil)
	req = app.WithUser(req, user)
	rr = httptest.NewRecorder()

	app.GetReviewSummaryHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp ReviewSummaryResponse
	decodeResponse(t, rr, &resp)

	if resp.TotalReviews != 1 {
		t.Errorf("expected total_reviews 1 for notebook filter, got %d", resp.TotalReviews)
	}
}

func TestGetReviewSummaryHandler_InvalidDate(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	req := httptest.NewRequest("GET", "/v1/reviews/summary?date=invalid", nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetReviewSummaryHandler(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d. Body: %s", rr.Code, rr.Body.String())
	}
}

// -----------------------------------------------------------------------------
// Review History Tests
// -----------------------------------------------------------------------------

func TestGetReviewHistoryHandler(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)

	// Submit multiple reviews
	for i := 0; i < 3; i++ {
		reviewReq := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
			"id":      uuid.New().String(),
			"card_id": cardID.String(),
			"rating":  "good",
			"mode":    "scheduled",
		})
		reviewReq = app.WithUser(reviewReq, user)
		rr := httptest.NewRecorder()
		app.SubmitReviewHandler(rr, reviewReq)
		time.Sleep(10 * time.Millisecond)
	}

	req := httptest.NewRequest("GET", "/v1/notebooks/"+nbID.String()+"/reviews", nil)
	req.SetPathValue("notebook_id", nbID.String())
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetReviewHistoryHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp PageResponse[ReviewHistoryItem]
	decodeResponse(t, rr, &resp)

	if resp.Total != 3 {
		t.Errorf("expected total 3, got %d", resp.Total)
	}
	if len(resp.Data) != 3 {
		t.Errorf("expected 3 items, got %d", len(resp.Data))
	}

	// Verify descending order
	for i := 1; i < len(resp.Data); i++ {
		if resp.Data[i].ReviewedAt.After(resp.Data[i-1].ReviewedAt) {
			t.Error("expected reviews in descending order")
		}
	}
}

func TestGetReviewHistoryHandler_Pagination(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)

	// Submit 5 reviews
	for i := 0; i < 5; i++ {
		reviewReq := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
			"id":      uuid.New().String(),
			"card_id": cardID.String(),
			"rating":  "good",
			"mode":    "scheduled",
		})
		reviewReq = app.WithUser(reviewReq, user)
		rr := httptest.NewRecorder()
		app.SubmitReviewHandler(rr, reviewReq)
		time.Sleep(10 * time.Millisecond)
	}

	// Get first page
	req := httptest.NewRequest("GET", "/v1/notebooks/"+nbID.String()+"/reviews?limit=2&offset=0", nil)
	req.SetPathValue("notebook_id", nbID.String())
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetReviewHistoryHandler(rr, req)

	var resp PageResponse[ReviewHistoryItem]
	decodeResponse(t, rr, &resp)

	if resp.Total != 5 {
		t.Errorf("expected total 5, got %d", resp.Total)
	}
	if len(resp.Data) != 2 {
		t.Errorf("expected 2 items, got %d", len(resp.Data))
	}
	if !resp.HasMore {
		t.Error("expected has_more true")
	}

	// Get second page
	req = httptest.NewRequest("GET", "/v1/notebooks/"+nbID.String()+"/reviews?limit=2&offset=2", nil)
	req.SetPathValue("notebook_id", nbID.String())
	req = app.WithUser(req, user)
	rr = httptest.NewRecorder()

	app.GetReviewHistoryHandler(rr, req)

	decodeResponse(t, rr, &resp)

	if len(resp.Data) != 2 {
		t.Errorf("expected 2 items on page 2, got %d", len(resp.Data))
	}
}

func TestGetReviewHistoryHandler_DateFilter(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)

	// Submit a review today
	reviewReq := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	reviewReq = app.WithUser(reviewReq, user)
	rr := httptest.NewRecorder()
	app.SubmitReviewHandler(rr, reviewReq)

	// Query for yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	req := httptest.NewRequest("GET", "/v1/notebooks/"+nbID.String()+"/reviews?date="+yesterday, nil)
	req.SetPathValue("notebook_id", nbID.String())
	req = app.WithUser(req, user)
	rr = httptest.NewRecorder()

	app.GetReviewHistoryHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp PageResponse[ReviewHistoryItem]
	decodeResponse(t, rr, &resp)

	if resp.Total != 0 {
		t.Errorf("expected 0 reviews for yesterday, got %d", resp.Total)
	}
}

func TestParseDateValid(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"standard", "2026-02-03", "2026-02-03"},
		{"start of year", "2026-01-01", "2026-01-01"},
		{"end of year", "2026-12-31", "2026-12-31"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test?date="+tt.input, nil)
			got, err := parseDate(req, "date")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("expected non-nil date")
			}
			if got.Format("2006-01-02") != tt.want {
				t.Errorf("got %s, want %s", got.Format("2006-01-02"), tt.want)
			}
		})
	}
}

func TestParseDateInvalid(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"invalid format", "02-03-2026"},
		{"not a date", "invalid"},
		{"partial date", "2026-02"},
		{"wrong separator", "2026/02/03"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test?date="+tt.input, nil)
			_, err := parseDate(req, "date")
			if err == nil {
				t.Error("expected error for invalid date format")
			}
		})
	}
}

// -----------------------------------------------------------------------------
// FSRS Settings Tests
// -----------------------------------------------------------------------------

func TestSubmitReview_UsesNotebookFSRSSettings(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Create notebook with default settings (retention=0.9)
	nbDefault := createTestNotebook(t, app, user.ID)
	_, cardsDefault, err := app.CreateFact(t.Context(), user.ID, nbDefault, "basic", basicContent())
	if err != nil {
		t.Fatalf("failed to create fact: %v", err)
	}
	cardDefault := cardsDefault[0].ID

	// Create notebook with custom retention (0.8 = easier retention = longer intervals)
	nbCustom := createTestNotebook(t, app, user.ID)
	customRetention := 0.8
	_, err = app.UpdateNotebook(t.Context(), user.ID, nbCustom, UpdateNotebookParams{
		DesiredRetention: &customRetention,
	})
	if err != nil {
		t.Fatalf("failed to update notebook retention: %v", err)
	}
	_, cardsCustom, err := app.CreateFact(t.Context(), user.ID, nbCustom, "basic", basicContent())
	if err != nil {
		t.Fatalf("failed to create fact: %v", err)
	}
	cardCustom := cardsCustom[0].ID

	// Review both cards with "good" rating
	reviewDefault := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": cardDefault.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	reviewDefault = app.WithUser(reviewDefault, user)
	rrDefault := httptest.NewRecorder()
	app.SubmitReviewHandler(rrDefault, reviewDefault)
	if rrDefault.Code != http.StatusOK {
		t.Fatalf("default review failed: %d. Body: %s", rrDefault.Code, rrDefault.Body.String())
	}

	reviewCustom := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": cardCustom.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	reviewCustom = app.WithUser(reviewCustom, user)
	rrCustom := httptest.NewRecorder()
	app.SubmitReviewHandler(rrCustom, reviewCustom)
	if rrCustom.Code != http.StatusOK {
		t.Fatalf("custom review failed: %d. Body: %s", rrCustom.Code, rrCustom.Body.String())
	}

	// Parse responses and compare due dates
	var respDefault, respCustom map[string]any
	decodeResponse(t, rrDefault, &respDefault)
	decodeResponse(t, rrCustom, &respCustom)

	cardRespDefault := respDefault["card"].(map[string]any)
	cardRespCustom := respCustom["card"].(map[string]any)

	dueDefaultStr := cardRespDefault["due"].(string)
	dueCustomStr := cardRespCustom["due"].(string)

	dueDefault, _ := time.Parse(time.RFC3339, dueDefaultStr)
	dueCustom, _ := time.Parse(time.RFC3339, dueCustomStr)

	// Lower retention (0.8) means the user accepts more forgetting,
	// so the interval should be LONGER than default (0.9)
	if !dueCustom.After(dueDefault) {
		t.Errorf("expected custom retention (0.8) to produce longer interval than default (0.9). Default due: %v, Custom due: %v",
			dueDefault, dueCustom)
	}
}

func TestGetStudyCards_UsesNotebookFSRSSettings(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Create notebook with custom retention
	nbID := createTestNotebook(t, app, user.ID)
	customRetention := 0.8
	_, err := app.UpdateNotebook(t.Context(), user.ID, nbID, UpdateNotebookParams{
		DesiredRetention: &customRetention,
	})
	if err != nil {
		t.Fatalf("failed to update notebook retention: %v", err)
	}

	// Create a card
	_, _, err = app.CreateFact(t.Context(), user.ID, nbID, "basic", basicContent())
	if err != nil {
		t.Fatalf("failed to create fact: %v", err)
	}

	// Get study cards - should include interval previews based on custom settings
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
	intervals, ok := resp[0]["intervals"].(map[string]any)
	if !ok {
		t.Fatal("expected intervals object")
	}

	// Intervals should be non-empty (basic sanity check)
	for _, key := range []string{"again", "hard", "good", "easy"} {
		if intervals[key] == nil || intervals[key] == "" {
			t.Errorf("expected non-empty interval for %s", key)
		}
	}
}

// -----------------------------------------------------------------------------
// Archived Notebook + Buried Until Tests
// -----------------------------------------------------------------------------

func TestGetStudyCards_ExcludesArchivedNotebook(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, _ := createTestCard(t, app, user.ID)

	// Archive the notebook
	if err := app.ArchiveNotebook(t.Context(), user.ID, nbID); err != nil {
		t.Fatalf("failed to archive notebook: %v", err)
	}

	req := httptest.NewRequest("GET", "/v1/reviews/study", nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetStudyCardsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp []map[string]any
	decodeResponse(t, rr, &resp)

	if len(resp) != 0 {
		t.Errorf("expected 0 study cards from archived notebook, got %d", len(resp))
	}
}

func TestGetPracticeCards_ExcludesArchivedNotebook(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, _ := createTestCard(t, app, user.ID)

	// Archive the notebook
	if err := app.ArchiveNotebook(t.Context(), user.ID, nbID); err != nil {
		t.Fatalf("failed to archive notebook: %v", err)
	}

	req := httptest.NewRequest("GET", "/v1/reviews/practice", nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

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
	if len(data) != 0 {
		t.Errorf("expected 0 practice cards from archived notebook, got %d", len(data))
	}
}

func TestGetStudyCounts_ExcludesArchivedNotebook(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Create two notebooks with cards
	nbID1, _ := createTestCard(t, app, user.ID)
	nbID2, _ := createTestCard(t, app, user.ID)

	// Archive notebook 2
	if err := app.ArchiveNotebook(t.Context(), user.ID, nbID2); err != nil {
		t.Fatalf("failed to archive notebook: %v", err)
	}

	req := httptest.NewRequest("GET", "/v1/reviews/study-counts", nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetStudyCountsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp StudyCountsResponse
	decodeResponse(t, rr, &resp)

	// Only non-archived notebook should appear
	if _, ok := resp.Counts[nbID2.String()]; ok {
		t.Error("archived notebook should not appear in study counts")
	}
	if _, ok := resp.Counts[nbID1.String()]; !ok {
		t.Error("non-archived notebook should appear in study counts")
	}
	if resp.TotalNew != 1 {
		t.Errorf("expected total_new 1 (only non-archived), got %d", resp.TotalNew)
	}
}

func TestGetStudyCards_BuriedUntilPastDate_Unburries(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)

	// Set buried_until to yesterday
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	_, err := app.DB.Exec(t.Context(),
		"UPDATE app.cards SET buried_until = $1 WHERE user_id = $2 AND id = $3",
		yesterday, user.ID, cardID)
	if err != nil {
		t.Fatalf("failed to set buried_until: %v", err)
	}

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
		t.Errorf("expected 1 study card (past buried_until should auto-unbury), got %d", len(resp))
	}
}

func TestGetStudyCards_BuriedUntilFutureDate_Excluded(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)

	// Set buried_until to 2 days from now (avoids UTC/local timezone boundary issues)
	tomorrow := time.Now().UTC().AddDate(0, 0, 2).Format("2006-01-02")
	_, err := app.DB.Exec(t.Context(),
		"UPDATE app.cards SET buried_until = $1 WHERE user_id = $2 AND id = $3",
		tomorrow, user.ID, cardID)
	if err != nil {
		t.Fatalf("failed to set buried_until: %v", err)
	}

	req := httptest.NewRequest("GET", "/v1/reviews/study?notebook_id="+nbID.String(), nil)
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.GetStudyCardsHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d. Body: %s", rr.Code, rr.Body.String())
	}

	var resp []map[string]any
	decodeResponse(t, rr, &resp)

	if len(resp) != 0 {
		t.Errorf("expected 0 study cards (future buried_until should exclude), got %d", len(resp))
	}
}

func TestSubmitReview_ArchivedNotebook_Rejected(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)
	nbID, cardID := createTestCard(t, app, user.ID)

	// Archive the notebook
	if err := app.ArchiveNotebook(t.Context(), user.ID, nbID); err != nil {
		t.Fatalf("failed to archive notebook: %v", err)
	}

	req := jsonRequest(t, "POST", "/v1/reviews", map[string]any{
		"id":      uuid.New().String(),
		"card_id": cardID.String(),
		"rating":  "good",
		"mode":    "scheduled",
	})
	req = app.WithUser(req, user)
	rr := httptest.NewRecorder()

	app.SubmitReviewHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 for review on archived notebook card, got %d. Body: %s", rr.Code, rr.Body.String())
	}
}

func TestGetStudyCards_MultiNotebook_UsesDifferentSettings(t *testing.T) {
	app := testApp(t)
	user := createTestUser(t, app)

	// Create two notebooks with different retention settings
	nb1 := createTestNotebook(t, app, user.ID)
	retention1 := 0.9 // higher retention = shorter intervals
	_, err := app.UpdateNotebook(t.Context(), user.ID, nb1, UpdateNotebookParams{
		DesiredRetention: &retention1,
	})
	if err != nil {
		t.Fatalf("failed to update notebook 1: %v", err)
	}

	nb2 := createTestNotebook(t, app, user.ID)
	retention2 := 0.7 // lower retention = longer intervals
	_, err = app.UpdateNotebook(t.Context(), user.ID, nb2, UpdateNotebookParams{
		DesiredRetention: &retention2,
	})
	if err != nil {
		t.Fatalf("failed to update notebook 2: %v", err)
	}

	// Create cards in both notebooks
	_, _, err = app.CreateFact(t.Context(), user.ID, nb1, "basic", basicContent())
	if err != nil {
		t.Fatalf("failed to create fact in nb1: %v", err)
	}
	_, _, err = app.CreateFact(t.Context(), user.ID, nb2, "basic", basicContent())
	if err != nil {
		t.Fatalf("failed to create fact in nb2: %v", err)
	}

	// Get study cards across all notebooks (global query)
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
		t.Fatalf("expected 2 study cards, got %d", len(resp))
	}

	// Map cards by notebook ID to compare their intervals
	cardsByNotebook := make(map[string]map[string]any)
	for _, card := range resp {
		nbID := card["notebook_id"].(string)
		cardsByNotebook[nbID] = card
	}

	card1 := cardsByNotebook[nb1.String()]
	card2 := cardsByNotebook[nb2.String()]

	if card1 == nil || card2 == nil {
		t.Fatalf("expected cards from both notebooks")
	}

	intervals1, ok := card1["intervals"].(map[string]any)
	if !ok {
		t.Fatal("expected intervals object for card 1")
	}
	intervals2, ok := card2["intervals"].(map[string]any)
	if !ok {
		t.Fatal("expected intervals object for card 2")
	}

	// Verify all intervals are present
	for _, key := range []string{"again", "hard", "good", "easy"} {
		if intervals1[key] == nil || intervals1[key] == "" {
			t.Errorf("card 1: expected non-empty interval for %s", key)
		}
		if intervals2[key] == nil || intervals2[key] == "" {
			t.Errorf("card 2: expected non-empty interval for %s", key)
		}
	}

	// Lower retention (0.7) should produce LONGER intervals than higher retention (0.9)
	// for at least one rating. We check "good" as it's the most common rating.
	good1 := intervals1["good"].(string)
	good2 := intervals2["good"].(string)

	// For new cards in learning, the intervals might be the same (both use learning steps).
	// But the "easy" rating graduates immediately to review, so it should differ.
	easy1 := intervals1["easy"].(string)
	easy2 := intervals2["easy"].(string)

	// At least one of good/easy should differ between the two retention levels
	if good1 == good2 && easy1 == easy2 {
		t.Errorf("expected different intervals for different retention settings. "+
			"nb1 (0.9): good=%s easy=%s, nb2 (0.7): good=%s easy=%s",
			good1, easy1, good2, easy2)
	}
}
