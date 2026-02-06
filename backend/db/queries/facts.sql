-- name: CreateFact :one
INSERT INTO app.facts (user_id, notebook_id, fact_type, content, source_id)
VALUES (@user_id, @notebook_id, @fact_type, @content, sqlc.narg('source_id'))
RETURNING *;

-- name: GetFact :one
SELECT n.*, (SELECT count(*) FROM app.cards c WHERE c.user_id = n.user_id AND c.fact_id = n.id) AS card_count
FROM app.facts n
WHERE n.user_id = @user_id AND n.id = @id AND n.notebook_id = @notebook_id;

-- name: CountFactsByNotebook :one
SELECT count(*) FROM app.facts
WHERE user_id = @user_id AND notebook_id = @notebook_id;

-- name: ListFactsByNotebookFiltered :many
SELECT f.*,
  (SELECT count(*) FROM app.cards c WHERE c.user_id = f.user_id AND c.fact_id = f.id) AS card_count,
  (SELECT count(*) FROM app.cards c WHERE c.user_id = f.user_id AND c.fact_id = f.id
    AND c.due <= now() AND c.suspended_at IS NULL AND (c.buried_until IS NULL OR c.buried_until <= CURRENT_DATE)) AS due_count
FROM app.facts f
WHERE f.user_id = @user_id AND f.notebook_id = @notebook_id
  AND (sqlc.narg('fact_type')::text IS NULL OR f.fact_type = sqlc.narg('fact_type'))
  AND (sqlc.narg('search')::text IS NULL OR f.content::text ILIKE '%' || sqlc.narg('search')::text || '%')
ORDER BY
  CASE WHEN @sort_field::text = 'created'  AND @sort_asc::bool = true  THEN f.created_at END ASC,
  CASE WHEN @sort_field::text = 'created'  AND @sort_asc::bool = false THEN f.created_at END DESC,
  CASE WHEN @sort_field::text = 'updated'  AND @sort_asc::bool = true  THEN f.updated_at END ASC,
  CASE WHEN @sort_field::text = 'updated'  AND @sort_asc::bool = false THEN f.updated_at END DESC,
  f.updated_at DESC
LIMIT @row_limit OFFSET @row_offset;

-- name: CountFactsByNotebookFiltered :one
SELECT count(*) FROM app.facts f
WHERE f.user_id = @user_id AND f.notebook_id = @notebook_id
  AND (sqlc.narg('fact_type')::text IS NULL OR f.fact_type = sqlc.narg('fact_type'))
  AND (sqlc.narg('search')::text IS NULL OR f.content::text ILIKE '%' || sqlc.narg('search')::text || '%');

-- name: GetFactStatsByNotebook :one
SELECT
  count(DISTINCT f.id) AS total_facts,
  count(DISTINCT c.id) AS total_cards,
  count(DISTINCT c.id) FILTER (WHERE c.due <= now() AND c.suspended_at IS NULL
    AND (c.buried_until IS NULL OR c.buried_until <= CURRENT_DATE)) AS total_due,
  count(DISTINCT f.id) FILTER (WHERE f.fact_type = 'basic') AS basic_count,
  count(DISTINCT f.id) FILTER (WHERE f.fact_type = 'cloze') AS cloze_count,
  count(DISTINCT f.id) FILTER (WHERE f.fact_type = 'image_occlusion') AS image_occlusion_count
FROM app.facts f
LEFT JOIN app.cards c ON c.user_id = f.user_id AND c.fact_id = f.id
WHERE f.user_id = @user_id AND f.notebook_id = @notebook_id;

-- name: UpdateFactContent :one
UPDATE app.facts
SET content = @content
WHERE user_id = @user_id AND id = @id AND notebook_id = @notebook_id
RETURNING *;

-- name: DeleteFact :execrows
DELETE FROM app.facts
WHERE user_id = @user_id AND id = @id AND notebook_id = @notebook_id;
