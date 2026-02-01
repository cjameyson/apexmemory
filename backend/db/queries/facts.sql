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

-- name: ListFactsByNotebookFiltered :many
SELECT f.*,
  (SELECT count(*) FROM app.cards c WHERE c.user_id = f.user_id AND c.fact_id = f.id) AS card_count,
  (SELECT count(*) FROM app.cards c WHERE c.user_id = f.user_id AND c.fact_id = f.id
    AND c.due <= now() AND c.suspended_at IS NULL AND c.buried_until IS NULL) AS due_count
FROM app.facts f
WHERE f.user_id = @user_id AND f.notebook_id = @notebook_id
  AND (sqlc.narg('fact_type')::text IS NULL OR f.fact_type = sqlc.narg('fact_type'))
  AND (sqlc.narg('search')::text IS NULL OR f.content::text ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY f.updated_at DESC
LIMIT @row_limit OFFSET @row_offset;

-- name: CountFactsByNotebookFiltered :one
SELECT count(*) FROM app.facts f
WHERE f.user_id = @user_id AND f.notebook_id = @notebook_id
  AND (sqlc.narg('fact_type')::text IS NULL OR f.fact_type = sqlc.narg('fact_type'))
  AND (sqlc.narg('search')::text IS NULL OR f.content::text ILIKE '%' || sqlc.narg('search')::text || '%');

-- name: GetFactStatsByNotebook :one
SELECT
  (SELECT count(*) FROM app.facts f1 WHERE f1.user_id = @user_id AND f1.notebook_id = @notebook_id) AS total_facts,
  (SELECT count(*) FROM app.cards c1 WHERE c1.user_id = @user_id AND c1.notebook_id = @notebook_id) AS total_cards,
  (SELECT count(*) FROM app.cards c2 WHERE c2.user_id = @user_id AND c2.notebook_id = @notebook_id
    AND c2.due <= now() AND c2.suspended_at IS NULL AND c2.buried_until IS NULL) AS total_due,
  (SELECT count(*) FROM app.facts f2 WHERE f2.user_id = @user_id AND f2.notebook_id = @notebook_id AND f2.fact_type = 'basic') AS basic_count,
  (SELECT count(*) FROM app.facts f3 WHERE f3.user_id = @user_id AND f3.notebook_id = @notebook_id AND f3.fact_type = 'cloze') AS cloze_count,
  (SELECT count(*) FROM app.facts f4 WHERE f4.user_id = @user_id AND f4.notebook_id = @notebook_id AND f4.fact_type = 'image_occlusion') AS image_occlusion_count;

-- name: UpdateFactContent :one
UPDATE app.facts
SET content = @content
WHERE user_id = @user_id AND id = @id AND notebook_id = @notebook_id
RETURNING *;

-- name: DeleteFact :execrows
DELETE FROM app.facts
WHERE user_id = @user_id AND id = @id AND notebook_id = @notebook_id;
