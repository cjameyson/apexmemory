package fsrs

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

var testRatings1 = []Rating{
	Good, Good, Good, Good, Good, Good,
	Again, Again,
	Good, Good, Good, Good, Good,
}

func mustScheduler(t *testing.T, opts ...Option) *Scheduler {
	t.Helper()
	s, err := NewScheduler(opts...)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func newCard() Card {
	return Card{
		State: Learning,
		Step:  intPtr(0),
		Due:   time.Now().UTC(),
	}
}

func ivlDays(c Card) int {
	return int(c.Due.Sub(*c.LastReview).Hours() / 24)
}

func TestReviewCard(t *testing.T) {
	s := mustScheduler(t, WithEnableFuzzing(false))
	card := newCard()
	now := time.Date(2022, 11, 29, 12, 30, 0, 0, time.UTC)

	var ivlHistory []int
	for _, rating := range testRatings1 {
		out := s.ReviewCard(card, rating, now)
		card = out.Card
		ivl := ivlDays(card)
		ivlHistory = append(ivlHistory, ivl)
		now = card.Due
	}

	expected := []int{0, 2, 11, 46, 163, 498, 0, 0, 2, 4, 7, 12, 21}
	if len(ivlHistory) != len(expected) {
		t.Fatalf("length mismatch: got %d, want %d", len(ivlHistory), len(expected))
	}
	for i := range expected {
		if ivlHistory[i] != expected[i] {
			t.Errorf("ivl[%d] = %d, want %d", i, ivlHistory[i], expected[i])
		}
	}
}

func TestRepeatedCorrectReviews(t *testing.T) {
	s := mustScheduler(t, WithEnableFuzzing(false))
	card := newCard()

	for i := 0; i < 10; i++ {
		now := time.Date(2022, 11, 29, 12, 30, 0, i, time.UTC)
		out := s.ReviewCard(card, Easy, now)
		card = out.Card
	}

	if *card.Difficulty != 1.0 {
		t.Errorf("difficulty = %f, want 1.0", *card.Difficulty)
	}
}

func TestMemoState(t *testing.T) {
	s := mustScheduler(t)

	ratings := []Rating{Again, Good, Good, Good, Good, Good}
	ivls := []int{0, 0, 1, 3, 8, 21}

	card := newCard()
	now := time.Date(2022, 11, 29, 12, 30, 0, 0, time.UTC)

	for i, rating := range ratings {
		now = now.Add(time.Duration(ivls[i]) * 24 * time.Hour)
		out := s.ReviewCard(card, rating, now)
		card = out.Card
	}

	if math.Abs(*card.Stability-53.62691) > 1e-4 {
		t.Errorf("stability = %f, want ~53.62691", *card.Stability)
	}
	if math.Abs(*card.Difficulty-6.3574867) > 1e-4 {
		t.Errorf("difficulty = %f, want ~6.3574867", *card.Difficulty)
	}
}

func TestRetrievability(t *testing.T) {
	s := mustScheduler(t)
	card := newCard()

	// new card
	r := s.Retrievability(card, time.Now().UTC())
	if r != 0 {
		t.Errorf("new card retrievability = %f, want 0", r)
	}

	// Learning
	out := s.ReviewCard(card, Good, time.Now().UTC())
	card = out.Card
	if card.State != Learning {
		t.Fatalf("expected Learning, got %d", card.State)
	}
	r = s.Retrievability(card, time.Now().UTC())
	if r < 0 || r > 1 {
		t.Errorf("learning R = %f, want [0,1]", r)
	}

	// Review
	out = s.ReviewCard(card, Good, card.Due)
	card = out.Card
	if card.State != Review {
		t.Fatalf("expected Review, got %d", card.State)
	}
	r = s.Retrievability(card, time.Now().UTC())
	if r < 0 || r > 1 {
		t.Errorf("review R = %f, want [0,1]", r)
	}

	// Relearning
	out = s.ReviewCard(card, Again, card.Due)
	card = out.Card
	if card.State != Relearning {
		t.Fatalf("expected Relearning, got %d", card.State)
	}
	r = s.Retrievability(card, time.Now().UTC())
	if r < 0 || r > 1 {
		t.Errorf("relearning R = %f, want [0,1]", r)
	}
}

func TestGoodLearningSteps(t *testing.T) {
	s := mustScheduler(t)
	card := newCard()

	if card.State != Learning || *card.Step != 0 {
		t.Fatal("bad initial state")
	}

	out := s.ReviewCard(card, Good, card.Due)
	card = out.Card
	if card.State != Learning || *card.Step != 1 {
		t.Errorf("after Good: state=%d step=%d, want Learning/1", card.State, *card.Step)
	}

	out = s.ReviewCard(card, Good, card.Due)
	card = out.Card
	if card.State != Review || card.Step != nil {
		t.Errorf("after 2nd Good: state=%d, want Review", card.State)
	}
}

func TestAgainLearningSteps(t *testing.T) {
	s := mustScheduler(t)
	card := newCard()

	out := s.ReviewCard(card, Again, card.Due)
	card = out.Card
	if card.State != Learning || *card.Step != 0 {
		t.Errorf("after Again: state=%d step=%d, want Learning/0", card.State, *card.Step)
	}
	// ~1 minute interval
	ivlSec := card.Due.Sub(*card.LastReview).Seconds()
	if ivlSec < 50 || ivlSec > 70 {
		t.Errorf("interval = %f sec, want ~60", ivlSec)
	}
}

func TestHardLearningSteps(t *testing.T) {
	s := mustScheduler(t)
	card := newCard()

	out := s.ReviewCard(card, Hard, card.Due)
	card = out.Card
	if card.State != Learning || *card.Step != 0 {
		t.Errorf("after Hard: state=%d step=%d", card.State, *card.Step)
	}
	// (1m + 10m) / 2 = 5.5m = 330s
	ivlSec := card.Due.Sub(*card.LastReview).Seconds()
	if ivlSec < 320 || ivlSec > 340 {
		t.Errorf("interval = %f sec, want ~330", ivlSec)
	}
}

func TestEasyLearningSteps(t *testing.T) {
	s := mustScheduler(t)
	card := newCard()

	out := s.ReviewCard(card, Easy, card.Due)
	card = out.Card
	if card.State != Review || card.Step != nil {
		t.Errorf("after Easy: state=%d, want Review", card.State)
	}
	ivlHrs := card.Due.Sub(*card.LastReview).Hours()
	if ivlHrs < 24 {
		t.Errorf("interval = %f hrs, want >= 24", ivlHrs)
	}
}

func TestReviewState(t *testing.T) {
	s := mustScheduler(t, WithEnableFuzzing(false))
	card := newCard()

	out := s.ReviewCard(card, Good, card.Due)
	card = out.Card
	out = s.ReviewCard(card, Good, card.Due)
	card = out.Card
	if card.State != Review {
		t.Fatalf("expected Review, got %d", card.State)
	}

	prevDue := card.Due
	out = s.ReviewCard(card, Good, card.Due)
	card = out.Card
	if card.State != Review {
		t.Errorf("expected Review after Good")
	}
	if card.Due.Sub(prevDue).Hours() < 24 {
		t.Errorf("interval too short")
	}

	prevDue = card.Due
	out = s.ReviewCard(card, Again, card.Due)
	card = out.Card
	if card.State != Relearning {
		t.Errorf("expected Relearning after Again")
	}
	ivlMin := int(card.Due.Sub(prevDue).Minutes())
	if ivlMin != 10 {
		t.Errorf("relearning interval = %d min, want 10", ivlMin)
	}
}

func TestRelearning(t *testing.T) {
	s := mustScheduler(t, WithEnableFuzzing(false))
	card := newCard()

	out := s.ReviewCard(card, Good, card.Due)
	card = out.Card
	out = s.ReviewCard(card, Good, card.Due)
	card = out.Card
	out = s.ReviewCard(card, Good, card.Due)
	card = out.Card

	prevDue := card.Due
	out = s.ReviewCard(card, Again, card.Due)
	card = out.Card
	if card.State != Relearning || *card.Step != 0 {
		t.Fatalf("expected Relearning/0")
	}
	if int(card.Due.Sub(prevDue).Minutes()) != 10 {
		t.Errorf("expected 10 min interval")
	}

	prevDue = card.Due
	out = s.ReviewCard(card, Again, card.Due)
	card = out.Card
	if card.State != Relearning || *card.Step != 0 {
		t.Fatalf("expected Relearning/0 after Again")
	}
	if int(card.Due.Sub(prevDue).Minutes()) != 10 {
		t.Errorf("expected 10 min interval")
	}

	prevDue = card.Due
	out = s.ReviewCard(card, Good, card.Due)
	card = out.Card
	if card.State != Review || card.Step != nil {
		t.Errorf("expected Review after Good in relearning")
	}
	if card.Due.Sub(prevDue).Hours() < 24 {
		t.Errorf("expected >= 1 day interval")
	}
}

func TestNoLearningSteps(t *testing.T) {
	s := mustScheduler(t, WithLearningSteps([]time.Duration{}))

	card := newCard()
	out := s.ReviewCard(card, Again, time.Now().UTC())
	card = out.Card
	if card.State != Review {
		t.Errorf("expected Review with no learning steps")
	}
	if ivlDays(card) < 1 {
		t.Errorf("expected >= 1 day interval")
	}
}

func TestNoRelearningSteps(t *testing.T) {
	s := mustScheduler(t, WithRelearningSteps([]time.Duration{}))
	card := newCard()

	out := s.ReviewCard(card, Good, time.Now().UTC())
	card = out.Card
	out = s.ReviewCard(card, Good, card.Due)
	card = out.Card
	if card.State != Review {
		t.Fatalf("expected Review")
	}
	out = s.ReviewCard(card, Again, card.Due)
	card = out.Card
	if card.State != Review {
		t.Errorf("expected Review with no relearning steps, got %d", card.State)
	}
	if ivlDays(card) < 1 {
		t.Errorf("expected >= 1 day interval")
	}
}

func TestOneCardMultipleSchedulers(t *testing.T) {
	s2 := mustScheduler(t, WithLearningSteps([]time.Duration{1 * time.Minute, 10 * time.Minute}))
	s1 := mustScheduler(t, WithLearningSteps([]time.Duration{1 * time.Minute}))
	s0 := mustScheduler(t, WithLearningSteps([]time.Duration{}))

	sr2 := mustScheduler(t, WithRelearningSteps([]time.Duration{1 * time.Minute, 10 * time.Minute}))
	sr1 := mustScheduler(t, WithRelearningSteps([]time.Duration{1 * time.Minute}))
	sr0 := mustScheduler(t, WithRelearningSteps([]time.Duration{}))

	card := newCard()
	now := time.Now().UTC()

	// learning with 2 steps
	out := s2.ReviewCard(card, Good, now)
	card = out.Card
	if card.State != Learning || *card.Step != 1 {
		t.Errorf("2-step: expected Learning/1")
	}

	out = s1.ReviewCard(card, Again, now)
	card = out.Card
	if card.State != Learning || *card.Step != 0 {
		t.Errorf("1-step: expected Learning/0")
	}

	out = s0.ReviewCard(card, Hard, now)
	card = out.Card
	if card.State != Review || card.Step != nil {
		t.Errorf("0-step: expected Review")
	}

	// relearning with 2 steps
	out = sr2.ReviewCard(card, Again, now)
	card = out.Card
	if card.State != Relearning || *card.Step != 0 {
		t.Errorf("2-rstep: expected Relearning/0")
	}

	out = sr2.ReviewCard(card, Good, now)
	card = out.Card
	if card.State != Relearning || *card.Step != 1 {
		t.Errorf("2-rstep good: expected Relearning/1")
	}

	out = sr1.ReviewCard(card, Again, now)
	card = out.Card
	if card.State != Relearning || *card.Step != 0 {
		t.Errorf("1-rstep: expected Relearning/0")
	}

	out = sr0.ReviewCard(card, Hard, now)
	card = out.Card
	if card.State != Review || card.Step != nil {
		t.Errorf("0-rstep: expected Review")
	}
}

func TestMaximumInterval(t *testing.T) {
	maxIvl := 100
	s := mustScheduler(t, WithMaximumInterval(maxIvl))
	card := newCard()

	ratings := []Rating{Easy, Good, Easy, Good}
	for _, r := range ratings {
		out := s.ReviewCard(card, r, card.Due)
		card = out.Card
		if ivlDays(card) > maxIvl {
			t.Errorf("interval %d exceeds max %d", ivlDays(card), maxIvl)
		}
	}
}

func TestStabilityLowerBound(t *testing.T) {
	s := mustScheduler(t)
	card := newCard()

	for i := 0; i < 1000; i++ {
		reviewTime := card.Due.Add(24 * time.Hour)
		out := s.ReviewCard(card, Again, reviewTime)
		card = out.Card
		if *card.Stability < StabilityMin {
			t.Fatalf("stability %f < %f at iteration %d", *card.Stability, StabilityMin, i)
		}
	}
}

func TestParameterValidation(t *testing.T) {
	// valid
	_, err := NewScheduler()
	if err != nil {
		t.Errorf("default params should be valid: %v", err)
	}

	// too high
	bad := DefaultParameters
	bad[6] = 100
	_, err = NewScheduler(WithParams(bad))
	if err == nil {
		t.Error("expected error for out-of-bounds param")
	}

	// too low
	bad = DefaultParameters
	bad[10] = -42
	_, err = NewScheduler(WithParams(bad))
	if err == nil {
		t.Error("expected error for out-of-bounds param")
	}
}

func TestHardRatingOneLearningStep(t *testing.T) {
	s := mustScheduler(t, WithLearningSteps([]time.Duration{10 * time.Minute}))
	card := newCard()

	out := s.ReviewCard(card, Hard, card.Due)
	card = out.Card
	if card.State != Learning {
		t.Fatal("expected Learning")
	}
	ivlSec := card.Due.Sub(*card.LastReview).Seconds()
	// 10m * 1.5 = 15m = 900s
	if math.Abs(ivlSec-900) > 1 {
		t.Errorf("interval = %f sec, want 900", ivlSec)
	}
}

func TestHardRatingTwoLearningSteps(t *testing.T) {
	s := mustScheduler(t, WithLearningSteps([]time.Duration{1 * time.Minute, 10 * time.Minute}))
	card := newCard()

	// advance to step 1
	out := s.ReviewCard(card, Good, card.Due)
	card = out.Card
	if *card.Step != 1 {
		t.Fatal("expected step 1")
	}

	prevDue := card.Due
	out = s.ReviewCard(card, Hard, card.Due)
	card = out.Card
	if card.State != Learning || *card.Step != 1 {
		t.Errorf("expected Learning/1")
	}
	ivlSec := card.Due.Sub(prevDue).Seconds()
	// step 1 = 10 min = 600s
	if math.Abs(ivlSec-600) > 1 {
		t.Errorf("interval = %f sec, want 600", ivlSec)
	}
}

func TestHardRatingOneRelearningStep(t *testing.T) {
	s := mustScheduler(t, WithRelearningSteps([]time.Duration{10 * time.Minute}))
	card := newCard()

	out := s.ReviewCard(card, Easy, card.Due)
	card = out.Card
	out = s.ReviewCard(card, Again, card.Due)
	card = out.Card
	if card.State != Relearning || *card.Step != 0 {
		t.Fatal("expected Relearning/0")
	}

	prevDue := card.Due
	out = s.ReviewCard(card, Hard, prevDue)
	card = out.Card
	ivlSec := card.Due.Sub(prevDue).Seconds()
	if math.Abs(ivlSec-900) > 1 {
		t.Errorf("interval = %f sec, want 900", ivlSec)
	}
}

func TestHardRatingTwoRelearningSteps(t *testing.T) {
	s := mustScheduler(t, WithRelearningSteps([]time.Duration{1 * time.Minute, 10 * time.Minute}))
	card := newCard()

	out := s.ReviewCard(card, Easy, card.Due)
	card = out.Card
	out = s.ReviewCard(card, Again, card.Due)
	card = out.Card
	if card.State != Relearning || *card.Step != 0 {
		t.Fatal("expected Relearning/0")
	}

	prevDue := card.Due
	out = s.ReviewCard(card, Hard, prevDue)
	card = out.Card
	ivlSec := card.Due.Sub(prevDue).Seconds()
	// (1m + 10m) / 2 = 330s
	if math.Abs(ivlSec-330) > 1 {
		t.Errorf("interval = %f sec, want 330", ivlSec)
	}
}

func TestLongTermStabilityRelearning(t *testing.T) {
	s := mustScheduler(t)
	card := newCard()

	out := s.ReviewCard(card, Easy, card.Due)
	card = out.Card
	if card.State != Review {
		t.Fatal("expected Review")
	}

	out = s.ReviewCard(card, Again, card.Due)
	card = out.Card
	if card.State != Relearning {
		t.Fatal("expected Relearning")
	}

	// review 1 day late
	lateReview := card.Due.Add(24 * time.Hour)
	out = s.ReviewCard(card, Good, lateReview)
	card = out.Card
	if card.State != Review {
		t.Errorf("expected Review after late relearning review")
	}
}

func TestFuzz(t *testing.T) {
	now := time.Now().UTC()

	// seed 1
	rng1 := rand.New(rand.NewSource(42))
	s1 := mustScheduler(t, WithRNG(rng1))
	card := newCard()
	out := s1.ReviewCard(card, Good, now)
	card = out.Card
	out = s1.ReviewCard(card, Good, card.Due)
	card = out.Card
	prevDue := card.Due
	out = s1.ReviewCard(card, Good, card.Due)
	card = out.Card
	ivl := int(card.Due.Sub(prevDue).Hours() / 24)
	// With fuzzing enabled, interval should differ from unfuzzed (11 days).
	// Just verify it's in a reasonable range around 11.
	if ivl < 9 || ivl > 13 {
		t.Errorf("seed 42: interval = %d, want ~11", ivl)
	}

	// seed 2: different seed should potentially give different fuzz
	rng2 := rand.New(rand.NewSource(12345))
	s2 := mustScheduler(t, WithRNG(rng2))
	card = newCard()
	out = s2.ReviewCard(card, Good, now)
	card = out.Card
	out = s2.ReviewCard(card, Good, card.Due)
	card = out.Card
	prevDue = card.Due
	out = s2.ReviewCard(card, Good, card.Due)
	card = out.Card
	ivl = int(card.Due.Sub(prevDue).Hours() / 24)
	if ivl < 9 || ivl > 13 {
		t.Errorf("seed 12345: interval = %d, want ~11", ivl)
	}
}

func TestRescheduleCardSameScheduler(t *testing.T) {
	s := mustScheduler(t, WithEnableFuzzing(false))
	card := newCard()

	var reviews []ReviewEntry
	now := card.Due
	for _, rating := range testRatings1 {
		out := s.ReviewCard(card, rating, now)
		reviews = append(reviews, ReviewEntry{Rating: rating, ReviewDatetime: now})
		card = out.Card
		now = card.Due
	}

	rescheduled := s.RescheduleCard(Card{State: Learning, Step: intPtr(0), Due: reviews[0].ReviewDatetime}, reviews)

	if *card.Stability != *rescheduled.Stability {
		t.Errorf("stability mismatch: %f vs %f", *card.Stability, *rescheduled.Stability)
	}
	if *card.Difficulty != *rescheduled.Difficulty {
		t.Errorf("difficulty mismatch: %f vs %f", *card.Difficulty, *rescheduled.Difficulty)
	}
	if card.State != rescheduled.State {
		t.Errorf("state mismatch: %d vs %d", card.State, rescheduled.State)
	}
}

func TestRescheduleCardDifferentParams(t *testing.T) {
	s := mustScheduler(t, WithEnableFuzzing(false))
	card := newCard()

	var reviews []ReviewEntry
	now := card.Due
	for _, rating := range testRatings1 {
		out := s.ReviewCard(card, rating, now)
		reviews = append(reviews, ReviewEntry{Rating: rating, ReviewDatetime: now})
		card = out.Card
		now = card.Due
	}

	diffParams := DefaultParameters
	diffParams[0] = 0.12340357383516173
	diffParams[2] = 2.397673571899466
	diffParams[4] = 6.686820427099132
	diffParams[5] = 0.45021679958387956
	diffParams[6] = 3.077875127553957
	diffParams[7] = 0.053520395733247045
	diffParams[8] = 1.6539992229052127
	diffParams[9] = 0.1466206769107436
	diffParams[10] = 0.6300772488850335
	diffParams[11] = 1.611965002575047
	diffParams[12] = 0.012840136810798864
	diffParams[13] = 0.34853762746216305
	diffParams[14] = 1.8878958285806287
	diffParams[15] = 0.8546376191171063
	diffParams[17] = 0.6748536823468675
	diffParams[18] = 0.20451266082721842
	diffParams[19] = 0.22622814695113844
	diffParams[20] = 0.46030603398979064

	s2 := mustScheduler(t, WithParams(diffParams), WithEnableFuzzing(false))
	rescheduled := s2.RescheduleCard(Card{State: Learning, Step: intPtr(0), Due: reviews[0].ReviewDatetime}, reviews)

	if *card.Stability == *rescheduled.Stability {
		t.Error("stability should differ with different params")
	}
	if *card.Difficulty == *rescheduled.Difficulty {
		t.Error("difficulty should differ with different params")
	}
	if *card.LastReview != *rescheduled.LastReview {
		t.Error("last_review should match")
	}
}
