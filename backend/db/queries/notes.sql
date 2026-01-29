-- name: CreateNote :one
INSERT INTO app.notes (user_id, notebook_id, note_type, content, source_id)
VALUES (@user_id, @notebook_id, @note_type, @content, sqlc.narg('source_id'))
RETURNING *;

-- name: GetNote :one
SELECT n.*, (SELECT count(*) FROM app.cards c WHERE c.user_id = n.user_id AND c.note_id = n.id) AS card_count
FROM app.notes n
WHERE n.user_id = @user_id AND n.id = @id AND n.notebook_id = @notebook_id;

-- name: ListNotesByNotebook :many
SELECT n.*, (SELECT count(*) FROM app.cards c WHERE c.user_id = n.user_id AND c.note_id = n.id) AS card_count
FROM app.notes n
WHERE n.user_id = @user_id AND n.notebook_id = @notebook_id
ORDER BY n.created_at DESC
LIMIT @row_limit OFFSET @row_offset;

-- name: CountNotesByNotebook :one
SELECT count(*) FROM app.notes
WHERE user_id = @user_id AND notebook_id = @notebook_id;

-- name: UpdateNoteContent :one
UPDATE app.notes
SET content = @content
WHERE user_id = @user_id AND id = @id AND notebook_id = @notebook_id
RETURNING *;

-- name: DeleteNote :execrows
DELETE FROM app.notes
WHERE user_id = @user_id AND id = @id;
