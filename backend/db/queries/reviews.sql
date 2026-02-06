-- name: GetStudyCards :many
-- Returns due cards for review: overdue first, then learning, then new.
-- New card cap limits how many new cards are introduced per day.
WITH new_today AS (
    SELECT count(*) AS cnt FROM app.reviews r
    WHERE r.user_id = @user_id
      AND r.state_before = 'new'
      AND r.mode = 'scheduled'
      AND r.reviewed_at >= date_trunc('day', now())
)
SELECT c.*, f.fact_type, f.content AS fact_content
FROM app.cards c
JOIN app.facts f ON f.user_id = c.user_id AND f.id = c.fact_id
JOIN app.notebooks n ON n.user_id = c.user_id AND n.id = c.notebook_id AND n.archived_at IS NULL
WHERE c.user_id = @user_id
  AND (sqlc.narg('notebook_id')::uuid IS NULL OR c.notebook_id = sqlc.narg('notebook_id'))
  AND c.suspended_at IS NULL
  AND (c.buried_until IS NULL OR c.buried_until <= CURRENT_DATE)
  AND (
    (c.state != 'new' AND c.due <= now())
    OR (c.state = 'new' AND (SELECT cnt FROM new_today) < @new_card_cap::bigint)
  )
ORDER BY
  CASE WHEN c.state != 'new' AND c.due <= now() THEN 0 ELSE 1 END,
  CASE WHEN c.state IN ('learning', 'relearning') THEN 0 ELSE 1 END,
  c.due ASC NULLS LAST,
  c.created_at ASC
LIMIT @row_limit;

-- name: GetPracticeCards :many
-- Returns all non-suspended cards for practice mode (no due filter).
SELECT c.*, f.fact_type, f.content AS fact_content
FROM app.cards c
JOIN app.facts f ON f.user_id = c.user_id AND f.id = c.fact_id
JOIN app.notebooks n ON n.user_id = c.user_id AND n.id = c.notebook_id AND n.archived_at IS NULL
WHERE c.user_id = @user_id
  AND (sqlc.narg('notebook_id')::uuid IS NULL OR c.notebook_id = sqlc.narg('notebook_id'))
  AND c.suspended_at IS NULL
  AND (c.buried_until IS NULL OR c.buried_until <= CURRENT_DATE)
ORDER BY c.created_at ASC
LIMIT @row_limit OFFSET @row_offset;

-- name: CountPracticeCards :one
SELECT count(*) FROM app.cards c
JOIN app.notebooks n ON n.user_id = c.user_id AND n.id = c.notebook_id AND n.archived_at IS NULL
WHERE c.user_id = @user_id
  AND (sqlc.narg('notebook_id')::uuid IS NULL OR c.notebook_id = sqlc.narg('notebook_id'))
  AND c.suspended_at IS NULL
  AND (c.buried_until IS NULL OR c.buried_until <= CURRENT_DATE);

-- name: GetCardForReview :one
SELECT c.* FROM app.cards c
JOIN app.notebooks n ON n.user_id = c.user_id AND n.id = c.notebook_id AND n.archived_at IS NULL
WHERE c.user_id = @user_id AND c.id = @id
FOR UPDATE OF c;

-- name: CreateReview :one
INSERT INTO app.reviews (
    user_id, id, card_id, notebook_id, fact_id, element_id,
    reviewed_at, rating, review_duration_ms, mode,
    state_before, stability_before, difficulty_before,
    elapsed_days, scheduled_days,
    state_after, stability_after, difficulty_after,
    interval_days, retrievability, undo_snapshot
) VALUES (
    @user_id, @id, @card_id, @notebook_id, @fact_id, @element_id,
    @reviewed_at, @rating, @review_duration_ms, @mode,
    @state_before, @stability_before, @difficulty_before,
    @elapsed_days, @scheduled_days,
    @state_after, @stability_after, @difficulty_after,
    @interval_days, @retrievability, @undo_snapshot
)
ON CONFLICT (user_id, id) DO NOTHING
RETURNING *;

-- name: UpdateCardAfterReview :exec
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
    lapses = CASE WHEN @add_lapse::bool THEN lapses + 1 ELSE lapses END
WHERE user_id = @user_id AND id = @id;

-- name: GetStudyCountsByNotebook :many
-- Returns due card counts per notebook for the review launcher.
SELECT
    c.notebook_id,
    count(*) FILTER (WHERE c.state != 'new' AND c.due <= now()) AS due_count,
    count(*) FILTER (WHERE c.state = 'new') AS new_count
FROM app.cards c
JOIN app.notebooks n ON n.user_id = c.user_id AND n.id = c.notebook_id AND n.archived_at IS NULL
WHERE c.user_id = @user_id
  AND c.suspended_at IS NULL
  AND (c.buried_until IS NULL OR c.buried_until <= CURRENT_DATE)
  AND (
    (c.state != 'new' AND c.due <= now())
    OR c.state = 'new'
  )
GROUP BY c.notebook_id;

-- name: GetReviewByID :one
-- Fetch a review for undo validation.
SELECT * FROM app.reviews
WHERE user_id = @user_id AND id = @id;

-- name: GetLatestReviewForCard :one
-- Get the most recent review for a card to verify undo is for latest.
SELECT id FROM app.reviews
WHERE user_id = @user_id AND card_id = @card_id
ORDER BY reviewed_at DESC
LIMIT 1;

-- name: DeleteReview :execrows
-- Delete a review (for undo).
DELETE FROM app.reviews
WHERE user_id = @user_id AND id = @id;

-- name: RestoreCardAfterUndo :exec
-- Restore card state from review's before columns + undo_snapshot.
-- The review record stores state_before/stability_before/difficulty_before directly.
-- The undo_snapshot JSONB captures remaining card fields that aren't stored as columns:
--   step, due, last_review: scheduling position in learning/relearning steps
--   reps, lapses: lifetime counters incremented by reviews
--   elapsed_days, scheduled_days: FSRS interval tracking for retrievability calc
UPDATE app.cards SET
    state = @state,
    stability = @stability,
    difficulty = @difficulty,
    step = @step,
    due = @due,
    last_review = @last_review,
    elapsed_days = @elapsed_days,
    scheduled_days = @scheduled_days,
    reps = @reps,
    lapses = @lapses
WHERE user_id = @user_id AND id = @id;

-- name: GetReviewSummaryByDate :one
-- Daily review summary with breakdown by rating, mode, and new cards.
SELECT
    count(*) AS total_reviews,
    count(*) FILTER (WHERE rating = 'again') AS again_count,
    count(*) FILTER (WHERE rating = 'hard') AS hard_count,
    count(*) FILTER (WHERE rating = 'good') AS good_count,
    count(*) FILTER (WHERE rating = 'easy') AS easy_count,
    count(*) FILTER (WHERE mode = 'scheduled') AS scheduled_count,
    count(*) FILTER (WHERE mode = 'practice') AS practice_count,
    COALESCE(sum(review_duration_ms), 0)::bigint AS total_duration_ms,
    count(*) FILTER (WHERE state_before = 'new' AND mode = 'scheduled') AS new_cards_seen
FROM app.reviews
WHERE user_id = @user_id
  AND (sqlc.narg('notebook_id')::uuid IS NULL OR notebook_id = sqlc.narg('notebook_id'))
  AND reviewed_at >= COALESCE(sqlc.narg('date')::date, CURRENT_DATE)::timestamptz
  AND reviewed_at < (COALESCE(sqlc.narg('date')::date, CURRENT_DATE) + interval '1 day')::timestamptz;

-- name: GetReviewHistory :many
-- Paginated review history for a notebook, optionally filtered by date.
SELECT
    r.id, r.card_id, r.notebook_id, r.fact_id, r.element_id,
    r.reviewed_at, r.rating, r.review_duration_ms, r.mode,
    r.state_before, r.state_after
FROM app.reviews r
WHERE r.user_id = @user_id
  AND r.notebook_id = @notebook_id
  AND (sqlc.narg('date')::date IS NULL
       OR (r.reviewed_at >= sqlc.narg('date')::timestamptz
           AND r.reviewed_at < (sqlc.narg('date')::date + interval '1 day')::timestamptz))
ORDER BY r.reviewed_at DESC
LIMIT @row_limit OFFSET @row_offset;

-- name: CountReviewHistory :one
SELECT count(*)
FROM app.reviews
WHERE user_id = @user_id
  AND notebook_id = @notebook_id
  AND (sqlc.narg('date')::date IS NULL
       OR (reviewed_at >= sqlc.narg('date')::timestamptz
           AND reviewed_at < (sqlc.narg('date')::date + interval '1 day')::timestamptz));
