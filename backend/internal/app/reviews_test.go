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
