-- name: CreateCard :one
INSERT INTO app.cards (user_id, notebook_id, note_id, element_id)
VALUES (@user_id, @notebook_id, @note_id, @element_id)
RETURNING *;

-- name: ListCardsByNote :many
SELECT * FROM app.cards
WHERE user_id = @user_id AND note_id = @note_id
ORDER BY element_id ASC;

-- name: ListCardsByNotebook :many
SELECT * FROM app.cards
WHERE user_id = @user_id AND notebook_id = @notebook_id
  AND (sqlc.narg('state')::app.card_state IS NULL OR state = sqlc.narg('state'))
  AND suspended_at IS NULL AND buried_until IS NULL
ORDER BY due ASC NULLS FIRST, created_at ASC
LIMIT @row_limit OFFSET @row_offset;

-- name: CountCardsByNotebook :one
SELECT count(*) FROM app.cards
WHERE user_id = @user_id AND notebook_id = @notebook_id
  AND (sqlc.narg('state')::app.card_state IS NULL OR state = sqlc.narg('state'))
  AND suspended_at IS NULL AND buried_until IS NULL;

-- name: GetCard :one
SELECT * FROM app.cards
WHERE user_id = @user_id AND id = @id AND notebook_id = @notebook_id;

-- name: DeleteCardsByNoteAndElements :exec
DELETE FROM app.cards
WHERE user_id = @user_id AND note_id = @note_id AND element_id = ANY(@element_ids::text[]);
