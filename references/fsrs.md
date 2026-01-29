# FSRS v6 Schema Design: Implementation Guide

## Overview

This document accompanies the PostgreSQL schema design for implementing FSRS v6 (Free Spaced Repetition Scheduler) in Apex Memory. It covers the algorithm fundamentals, design decisions, and implementation guidance.

---

## Table of Contents

1. [FSRS v6 Algorithm Summary](#1-fsrs-v6-algorithm-summary)
2. [Data Model Mapping](#2-data-model-mapping)
3. [Schema Design Decisions](#3-schema-design-decisions)
4. [Implementation Notes](#4-implementation-notes)
5. [Go Backend Considerations](#5-go-backend-considerations)
6. [Optimizer Integration](#6-optimizer-integration)

---

## 1. FSRS v6 Algorithm Summary

### The Three-Component Memory Model (DSR)

FSRS is based on the DSR model where memory state is described by three variables:

| Variable | Symbol | Description | Range |
|----------|--------|-------------|-------|
| **Difficulty** | D | Inherent complexity of the memory item | [1, 10] |
| **Stability** | S | Time (days) for retrievability to decay from 100% to 90% | [0, ∞) |
| **Retrievability** | R | Probability of successful recall at a given moment | [0, 1] |

**Key insight:** D and S are stored on the card and updated after each review. R is computed dynamically based on S and time elapsed.

### Retrievability Formula (Forgetting Curve)

FSRS v6 uses a power-law forgetting curve:

```
R(t, S) = (1 + F × (t / S))^(-w₂₀)

where:
  t = elapsed time in days since last review
  S = stability
  F = 19/81 ≈ 0.2346 (constant)
  w₂₀ = decay parameter (0.1-0.8, default ~0.15)
```

At `t = S`, retrievability equals 90% (by definition of stability).

### Interval Calculation

The optimal interval is when R decays to the desired retention:

```
I = (S / F) × ((DR^(-1/w₂₀)) - 1)

where:
  I = interval in days
  S = stability
  DR = desired_retention (default 0.9)
  F = 19/81
  w₂₀ = decay parameter
```

When desired_retention = 0.9 and w₂₀ = 0.5, the interval equals stability (before fuzzing).

### Card States and Transitions

```
┌─────────────────────────────────────────────────────────────┐
│                         NEW                                  │
│              (never studied, no D/S values)                  │
└──────────────────────────┬──────────────────────────────────┘
                           │ First review
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                       LEARNING                               │
│          (stepping through learning_steps)                   │
│                                                              │
│  Again → restart step 0                                      │
│  Hard  → repeat current step                                 │
│  Good  → advance to next step                                │
│  Easy  → skip directly to Review                             │
└──────────────┬─────────────────────────────────┬────────────┘
               │ Complete final step              │ Easy rating
               ▼                                  ▼
┌─────────────────────────────────────────────────────────────┐
│                        REVIEW                                │
│           (scheduled by FSRS algorithm)                      │
│                                                              │
│  Hard/Good/Easy → update D, increase S, schedule             │
└──────────────────────────┬──────────────────────────────────┘
                           │ Again rating (lapse)
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                      RELEARNING                              │
│          (stepping through relearning_steps)                 │
│          (stability significantly reduced)                   │
│                                                              │
│  Again → restart step 0                                      │
│  Hard  → repeat current step                                 │
│  Good  → advance to next step                                │
│  Easy  → skip to Review with reduced S                       │
└──────────────────────────┬──────────────────────────────────┘
                           │ Complete final step
                           ▼
                    Back to REVIEW
```

### Rating Effects Summary

| Rating | Effect on D | Effect on S | Notes |
|--------|-------------|-------------|-------|
| Again (1) | Increases significantly | Decreases (lapse) | Card enters Relearning |
| Hard (2) | Increases moderately | Increases slightly or stays same | S×w₁₅ multiplier |
| Good (3) | May increase/decrease slightly | Increases | Standard progression |
| Easy (4) | Decreases moderately | Increases significantly | S×w₁₆ multiplier |

---

## 2. Data Model Mapping

### py-fsrs Card → PostgreSQL cards table

| py-fsrs Field | PostgreSQL Column | Type | Notes |
|---------------|-------------------|------|-------|
| `card_id` | `id` | UUID | Generated via UUIDv7 |
| `state` | `state` | ENUM | 'new', 'learning', 'review', 'relearning' |
| `step` | `step` | SMALLINT | NULL when state = 'review' |
| `stability` | `stability` | REAL | NULL for new cards |
| `difficulty` | `difficulty` | REAL | NULL for new cards, range [1,10] |
| `due` | `due` | TIMESTAMPTZ | NULL = due immediately |
| `last_review` | `last_review` | TIMESTAMPTZ | NULL if never reviewed |
| `elapsed_days` | `elapsed_days` | REAL | Days since last review |
| `scheduled_days` | `scheduled_days` | REAL | Assigned interval |
| `reps` | `reps` | INTEGER | Total review count |
| `lapses` | `lapses` | INTEGER | Lapse count (Again in Review) |

### py-fsrs ReviewLog → PostgreSQL reviews table

| py-fsrs Field | PostgreSQL Column | Type | Notes |
|---------------|-------------------|------|-------|
| `card_id` | `card_id` | UUID | Foreign key to cards |
| `rating` | `rating` | ENUM | 'again', 'hard', 'good', 'easy' |
| `review_datetime` | `reviewed_at` | TIMESTAMPTZ | When review occurred |
| `review_duration` | `review_duration_ms` | INTEGER | Milliseconds (optional) |
| `elapsed_days` | `elapsed_days` | REAL | Actual days elapsed |
| `scheduled_days` | `scheduled_days` | REAL | Days that were scheduled |
| — | `state_before` | ENUM | Card state before review |
| — | `stability_before` | REAL | S before review |
| — | `difficulty_before` | REAL | D before review |
| — | `state_after` | ENUM | Card state after review |
| — | `stability_after` | REAL | S after review |
| — | `difficulty_after` | REAL | D after review |
| — | `interval_days` | REAL | New scheduled interval |
| — | `retrievability` | REAL | R at review time |

**Note:** The schema stores both before/after states for comprehensive analytics and optimizer training.

### py-fsrs Scheduler → PostgreSQL notebooks.fsrs_params

| py-fsrs Parameter | JSONB Path | Notes |
|-------------------|------------|-------|
| `parameters` | `$.parameters` | Array of 21 floats (w₀-w₂₀) |
| `desired_retention` | `$.desired_retention` | Float 0.7-0.97 |
| `learning_steps` | `$.learning_steps` | Array of seconds [60, 600] |
| `relearning_steps` | `$.relearning_steps` | Array of seconds [600] |
| `maximum_interval` | `$.maximum_interval` | Days (default 36500) |
| `enable_fuzzing` | `$.enable_fuzzing` | Boolean |

---

## 3. Schema Design Decisions

### Why REAL for FSRS Values?
Human memory is noisy. The difference between stability of 45.123456789 and 45.123459 days is meaningless — the user might have been tired, distracted, or guessed.

### Why Composite Primary Keys?

From your Tech Stack: *"All multi-tenant tables use composite PKs (user_id, id) for efficient partitioning."*

Benefits:
- Future table partitioning by user_id becomes trivial
- Row-level security queries are efficient
- Related data is physically co-located

### Why Store state_before/state_after in reviews?

The optimizer needs to know the card's memory state at the time of each review. Storing snapshots means we don't need to replay history to reconstruct state.

### Why JSONB for fsrs_params?

From your Tech Stack: *"Avoid JSONB except for flexible config columns (e.g., fsrs_params)."*

FSRS parameters may evolve (v5 had 19, v6 has 21). JSONB allows:
- Adding parameters without migrations
- Storing user-specific optimized parameters
- Easy serialization to/from Go structs

### Why Denormalized notebook_id on reviews?

Reviews reference cards which reference notebooks. Denormalizing `notebook_id` onto reviews enables efficient notebook-level analytics without joins.

### Why Separate notes and cards Tables?

A single note can generate multiple cards:
- **Cloze notes:** "{{c1::Rome}} is the capital of {{c2::Italy}}" → 2 cards
- **Image occlusion:** One image with 5 hidden regions → 5 cards

Each card maintains its own independent FSRS state.

---

## 4. Implementation Notes

### Scheduling Algorithm (Pseudo-code)

```go
func (s *Scheduler) ReviewCard(card Card, rating Rating, now time.Time) (Card, ReviewLog) {
    // 1. Compute elapsed days
    elapsedDays := 0.0
    if card.LastReview != nil {
        elapsedDays = now.Sub(*card.LastReview).Hours() / 24
    }
    
    // 2. Record state before
    reviewLog := ReviewLog{
        CardID:           card.ID,
        Rating:           rating,
        ReviewedAt:       now,
        ElapsedDays:      elapsedDays,
        ScheduledDays:    card.ScheduledDays,
        StateBefore:      card.State,
        StabilityBefore:  card.Stability,
        DifficultyBefore: card.Difficulty,
    }
    
    // 3. Compute retrievability (if not new)
    if card.Stability != nil && *card.Stability > 0 {
        reviewLog.Retrievability = calculateRetrievability(
            *card.Stability, elapsedDays, s.params.W[20])
    }
    
    // 4. Update card based on current state and rating
    switch card.State {
    case StateNew:
        card = s.handleFirstReview(card, rating)
    case StateLearning:
        card = s.handleLearning(card, rating)
    case StateReview:
        card = s.handleReview(card, rating, elapsedDays)
    case StateRelearning:
        card = s.handleRelearning(card, rating)
    }
    
    // 5. Calculate next due date
    card.Due = s.calculateDue(card, now)
    card.LastReview = &now
    card.Reps++
    
    // 6. Record state after
    reviewLog.StateAfter = card.State
    reviewLog.StabilityAfter = *card.Stability
    reviewLog.DifficultyAfter = *card.Difficulty
    reviewLog.IntervalDays = card.ScheduledDays
    
    return card, reviewLog
}
```

### Fuzzing Implementation

Fuzzing adds random variance to intervals to spread reviews across days:

```go
func (s *Scheduler) applyFuzz(interval float64) float64 {
    if !s.params.EnableFuzzing || interval < 2.5 {
        return interval
    }
    
    // Fuzz range increases with interval length
    var fuzzRange float64
    switch {
    case interval < 7:
        fuzzRange = 1.0
    case interval < 30:
        fuzzRange = max(2.0, interval*0.15)
    default:
        fuzzRange = max(4.0, interval*0.05)
    }
    
    // Apply random offset
    offset := (rand.Float64() - 0.5) * 2 * fuzzRange
    return max(1.0, interval+offset)
}
```

### Due Query Optimization

The most critical query is fetching due cards:

```sql
-- Optimal query for due cards
SELECT c.*, n.fields
FROM app.cards c
JOIN app.notes n ON n.user_id = c.user_id AND n.id = c.note_id
WHERE c.user_id = $1
  AND c.notebook_id = $2
  AND c.suspended_at IS NULL
  AND (c.buried_until IS NULL OR c.buried_until <= CURRENT_DATE)
  AND (c.due IS NULL OR c.due <= $3)
ORDER BY 
  CASE c.state 
    WHEN 'learning' THEN 0 
    WHEN 'relearning' THEN 1 
    ELSE 2 
  END,
  c.due ASC NULLS FIRST
LIMIT 100;
```

The partial index `idx_cards_due` makes this query efficient.

---

## 5. Go Backend Considerations

### Type Definitions

```go
// types/fsrs.go

type CardState string

const (
    StateNew        CardState = "new"
    StateLearning   CardState = "learning"
    StateReview     CardState = "review"
    StateRelearning CardState = "relearning"
)

type Rating string

const (
    RatingAgain Rating = "again"
    RatingHard  Rating = "hard"
    RatingGood  Rating = "good"
    RatingEasy  Rating = "easy"
)

// FSRSParams stored in JSONB
type FSRSParams struct {
    Parameters       [21]float64 `json:"parameters"`
    DesiredRetention float64     `json:"desired_retention"`
    LearningSteps    []int       `json:"learning_steps"`    // seconds
    RelearningSteps  []int       `json:"relearning_steps"`  // seconds
    MaximumInterval  int         `json:"maximum_interval"`  // days
    EnableFuzzing    bool        `json:"enable_fuzzing"`
}

// Default FSRS v6 parameters
var DefaultFSRSParams = FSRSParams{
    Parameters: [21]float64{
        0.212, 1.2931, 2.3065, 8.2956, 6.4133, 0.8334, 3.0194, 0.001,
        1.8722, 0.1666, 0.796, 1.4835, 0.0614, 0.2629, 1.6483, 0.6014,
        1.8729, 0.5425, 0.0912, 0.0658, 0.1542,
    },
    DesiredRetention: 0.9,
    LearningSteps:    []int{60, 600},     // 1min, 10min
    RelearningSteps:  []int{600},          // 10min
    MaximumInterval:  36500,               // ~100 years
    EnableFuzzing:    true,
}
```

### Database Queries with sqlc

```sql
-- queries/cards.sql

-- name: GetDueCards :many
SELECT sqlc.embed(c), sqlc.embed(n)
FROM app.cards c
JOIN app.notes n ON n.user_id = c.user_id AND n.id = c.note_id
WHERE c.user_id = @user_id
  AND (@notebook_id::uuid IS NULL OR c.notebook_id = @notebook_id)
  AND c.suspended_at IS NULL
  AND (c.buried_until IS NULL OR c.buried_until <= CURRENT_DATE)
  AND (c.due IS NULL OR c.due <= @as_of)
ORDER BY 
  CASE c.state 
    WHEN 'learning' THEN 0 
    WHEN 'relearning' THEN 1 
    ELSE 2 
  END,
  c.due ASC NULLS FIRST
LIMIT @limit;

-- name: UpdateCardAfterReview :one
UPDATE app.cards SET
    state = @state,
    stability = @stability,
    difficulty = @difficulty,
    step = @step,
    due = @due,
    last_review = @last_review,
    elapsed_days = @elapsed_days,
    scheduled_days = @scheduled_days,
    reps = reps + 1,
    lapses = CASE WHEN @is_lapse THEN lapses + 1 ELSE lapses END,
    updated_at = now()
WHERE user_id = @user_id AND id = @id
RETURNING *;

-- name: InsertReview :one
INSERT INTO app.reviews (
    user_id, card_id, notebook_id, reviewed_at, rating,
    review_duration_ms, state_before, stability_before, difficulty_before,
    elapsed_days, scheduled_days, state_after, stability_after,
    difficulty_after, interval_days, retrievability
) VALUES (
    @user_id, @card_id, @notebook_id, @reviewed_at, @rating,
    @review_duration_ms, @state_before, @stability_before, @difficulty_before,
    @elapsed_days, @scheduled_days, @state_after, @stability_after,
    @difficulty_after, @interval_days, @retrievability
)
RETURNING *;
```

---

## 6. Optimizer Integration

### Data Export for Optimizer

The FSRS optimizer expects review logs in this format:

```python
# Required fields for optimizer
{
    "card_id": int,           # Unique card identifier
    "review_time": int,       # Unix timestamp (seconds)  
    "review_rating": int,     # 1-4
    "review_duration": int,   # milliseconds (optional)
}
```

Export query:

```sql
SELECT 
    c.id::text as card_id,
    EXTRACT(EPOCH FROM r.reviewed_at)::bigint as review_time,
    CASE r.rating 
        WHEN 'again' THEN 1
        WHEN 'hard' THEN 2
        WHEN 'good' THEN 3
        WHEN 'easy' THEN 4
    END as review_rating,
    r.review_duration_ms as review_duration
FROM app.reviews r
JOIN app.cards c ON c.user_id = r.user_id AND c.id = r.card_id
WHERE r.user_id = $1
  AND r.notebook_id = $2
  AND r.scheduled_days >= 1  -- Only long-term reviews
ORDER BY r.reviewed_at;
```

### Storing Optimized Parameters

After running the optimizer, store results at notebook level:

```sql
UPDATE app.notebooks
SET fsrs_params = jsonb_set(
    fsrs_params,
    '{parameters}',
    $2::jsonb  -- New 21-element array
)
WHERE user_id = $1 AND id = $3;
```

### Reschedule After Optimization

When parameters change, cards can be rescheduled to reflect new predictions:

```go
func (s *Service) RescheduleNotebook(ctx context.Context, userID, notebookID uuid.UUID) error {
    // Get new parameters
    notebook, _ := s.queries.GetNotebook(ctx, userID, notebookID)
    scheduler := NewScheduler(notebook.FSRSParams)
    
    // Get all cards with review history
    cards, _ := s.queries.GetCardsWithReviews(ctx, userID, notebookID)
    
    for _, card := range cards {
        // Replay review history with new parameters
        newCard := scheduler.RescheduleCard(card, card.Reviews)
        
        // Update card with new D, S, due values
        s.queries.UpdateCardFSRS(ctx, newCard)
    }
    
    return nil
}
```

---

## Quick Reference

### FSRS v6 Default Parameters

| Index | Symbol | Default | Description |
|-------|--------|---------|-------------|
| 0 | w₀ | 0.212 | Initial S after first Again |
| 1 | w₁ | 1.2931 | Initial S after first Hard |
| 2 | w₂ | 2.3065 | Initial S after first Good |
| 3 | w₃ | 8.2956 | Initial S after first Easy |
| 4 | w₄ | 6.4133 | Initial D (Good rating) |
| 5 | w₅ | 0.8334 | D scaling factor |
| 6 | w₆ | 3.0194 | D increase after Again |
| 7 | w₇ | 0.001 | D mean reversion |
| 8 | w₈ | 1.8722 | Base SInc factor |
| 9 | w₉ | 0.1666 | D effect on SInc |
| 10 | w₁₀ | 0.796 | S effect on SInc |
| 11 | w₁₁ | 1.4835 | Post-lapse S base |
| 12 | w₁₂ | 0.0614 | Post-lapse D effect |
| 13 | w₁₃ | 0.2629 | Post-lapse S effect |
| 14 | w₁₄ | 1.6483 | Post-lapse R effect |
| 15 | w₁₅ | 0.6014 | Hard penalty |
| 16 | w₁₆ | 1.8729 | Easy bonus |
| 17 | w₁₇ | 0.5425 | Same-day S base |
| 18 | w₁₈ | 0.0912 | Same-day grade offset |
| 19 | w₁₉ | 0.0658 | Same-day convergence |
| 20 | w₂₀ | 0.1542 | Forgetting curve decay |

### Key Formulas

```
# Retrievability (forgetting curve)
R = (1 + (19/81) × (t/S))^(-w₂₀)

# Interval from retention
I = (S × 81/19) × (DR^(-1/w₂₀) - 1)

# Stability increase (successful review)
SInc = e^(w₈) × (11 - D)^(w₉) × S^(-w₁₀) × (e^(w₁₅×(G-2)) - 1)  [simplified]

# Post-lapse stability (after Again in Review)
S' = w₁₁ × D^(-w₁₂) × ((S+1)^(w₁₃) - 1) × e^(w₁₄×(1-R))
```

---

## Resources

- [FSRS Algorithm Wiki](https://github.com/open-spaced-repetition/fsrs4anki/wiki/The-Algorithm)
- [py-fsrs Implementation](https://github.com/open-spaced-repetition/py-fsrs)
- [FSRS Visualizer](https://open-spaced-repetition.github.io/anki-fsrs-visualizer/)
- [Technical Explanation](https://expertium.github.io/Algorithm.html)
- [SRS Benchmark](https://github.com/open-spaced-repetition/srs-benchmark)