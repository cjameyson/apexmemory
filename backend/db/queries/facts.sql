-- name: CreateFact :one
INSERT INTO app.facts (user_id, notebook_id, fact_type, content, source_id)
VALUES (@user_id, @notebook_id, @fact_type, @content, sqlc.narg('source_id'))
RETURNING *;

-- name: GetFact :one
SELECT n.*, (SELECT count(*) FROM app.cards c WHERE c.user_id = n.user_id AND c.fact_id = n.id) AS card_count
FROM app.facts n
WHERE n.user_id = @user_id AND n.id = @id AND n.notebook_id = @notebook_id;

-- name: ListFactsByNotebook :many
SELECT n.*, (SELECT count(*) FROM app.cards c WHERE c.user_id = n.user_id AND c.fact_id = n.id) AS card_count
FROM app.facts n
WHERE n.user_id = @user_id AND n.notebook_id = @notebook_id
ORDER BY n.created_at DESC
LIMIT @row_limit OFFSET @row_offset;

-- name: CountFactsByNotebook :one
SELECT count(*) FROM app.facts
WHERE user_id = @user_id AND notebook_id = @notebook_id;

-- name: UpdateFactContent :one
UPDATE app.facts
SET content = @content
WHERE user_id = @user_id AND id = @id AND notebook_id = @notebook_id
RETURNING *;

-- name: DeleteFact :execrows
DELETE FROM app.facts
WHERE user_id = @user_id AND id = @id AND notebook_id = @notebook_id;
