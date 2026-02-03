-- name: GetStudyCards :many
-- Returns due cards for review: overdue first, then learning, then new.
-- New card cap limits how many new cards are introduced per day.
SELECT c.*, f.fact_type, f.content AS fact_content
FROM app.cards c
JOIN app.facts f ON f.user_id = c.user_id AND f.id = c.fact_id
WHERE c.user_id = @user_id
  AND (sqlc.narg('notebook_id')::uuid IS NULL OR c.notebook_id = sqlc.narg('notebook_id'))
  AND c.suspended_at IS NULL
  AND c.buried_until IS NULL
  AND (
    -- Due cards (overdue or currently in learning)
    (c.state != 'new' AND c.due <= now())
    OR
    -- New cards, subject to daily cap
    (c.state = 'new' AND (
      SELECT count(*) FROM app.reviews r
      WHERE r.user_id = @user_id
        AND r.state_before = 'new'
        AND r.mode = 'scheduled'
        AND r.reviewed_at >= date_trunc('day', now())
    ) < @new_card_cap::bigint
    )
  )
ORDER BY
  -- Overdue non-new cards first
  CASE WHEN c.state != 'new' AND c.due <= now() THEN 0 ELSE 1 END,
  -- Then learning/relearning before new
  CASE WHEN c.state IN ('learning', 'relearning') THEN 0 ELSE 1 END,
  -- Within each group, earliest due first
  c.due ASC NULLS LAST,
  c.created_at ASC
LIMIT @row_limit;

-- name: GetPracticeCards :many
-- Returns all non-suspended cards for practice mode (no due filter).
SELECT c.*, f.fact_type, f.content AS fact_content
FROM app.cards c
JOIN app.facts f ON f.user_id = c.user_id AND f.id = c.fact_id
WHERE c.user_id = @user_id
  AND (sqlc.narg('notebook_id')::uuid IS NULL OR c.notebook_id = sqlc.narg('notebook_id'))
  AND c.suspended_at IS NULL
  AND c.buried_until IS NULL
ORDER BY c.created_at ASC
LIMIT @row_limit OFFSET @row_offset;

-- name: CountPracticeCards :one
SELECT count(*) FROM app.cards
WHERE user_id = @user_id
  AND (sqlc.narg('notebook_id')::uuid IS NULL OR notebook_id = sqlc.narg('notebook_id'))
  AND suspended_at IS NULL
  AND buried_until IS NULL;

-- name: GetCardForReview :one
SELECT * FROM app.cards
WHERE user_id = @user_id AND id = @id
FOR UPDATE;

-- name: CreateReview :one
INSERT INTO app.reviews (
    user_id, id, card_id, notebook_id, fact_id, element_id,
    reviewed_at, rating, review_duration_ms, mode,
    state_before, stability_before, difficulty_before,
    elapsed_days, scheduled_days,
    state_after, stability_after, difficulty_after,
    interval_days, retrievability
) VALUES (
    @user_id, @id, @card_id, @notebook_id, @fact_id, @element_id,
    @reviewed_at, @rating, @review_duration_ms, @mode,
    @state_before, @stability_before, @difficulty_before,
    @elapsed_days, @scheduled_days,
    @state_after, @stability_after, @difficulty_after,
    @interval_days, @retrievability
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
