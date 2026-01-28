-- name: CreateNotebook :one
-- Note: fsrs_settings is always provided from Go code (source of truth for defaults)
INSERT INTO app.notebooks (user_id, name, description, emoji, color, position, fsrs_settings)
VALUES (
    @user_id,
    @name,
    sqlc.narg('description'),
    sqlc.narg('emoji'),
    sqlc.narg('color'),
    COALESCE(sqlc.narg('position'), 0),
    @fsrs_settings
)
RETURNING *;

-- name: GetNotebook :one
SELECT * FROM app.notebooks
WHERE user_id = @user_id AND id = @id AND archived_at IS NULL;

-- name: ListNotebooks :many
SELECT * FROM app.notebooks
WHERE user_id = @user_id AND archived_at IS NULL
ORDER BY position ASC, created_at DESC;

-- name: UpdateNotebook :one
UPDATE app.notebooks
SET name = COALESCE(sqlc.narg('name'), name),
    description = CASE
        WHEN sqlc.arg('clear_description')::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('description'), description)
    END,
    emoji = CASE
        WHEN sqlc.arg('clear_emoji')::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('emoji'), emoji)
    END,
    color = CASE
        WHEN sqlc.arg('clear_color')::boolean THEN NULL
        ELSE COALESCE(sqlc.narg('color'), color)
    END,
    position = COALESCE(sqlc.narg('position'), position),
    fsrs_settings = CASE
        WHEN @update_retention::boolean
        THEN jsonb_set(fsrs_settings, '{desired_retention}', to_jsonb(@desired_retention::float8))
        ELSE fsrs_settings
    END
WHERE user_id = @user_id AND id = @id AND archived_at IS NULL
RETURNING *;

-- name: ArchiveNotebook :execrows
UPDATE app.notebooks
SET archived_at = now()
WHERE user_id = @user_id AND id = @id AND archived_at IS NULL;

-- name: UnarchiveNotebook :exec
UPDATE app.notebooks
SET archived_at = NULL
WHERE user_id = @user_id AND id = @id AND archived_at IS NOT NULL;

-- name: IsNotebookArchived :one
-- Returns whether a notebook exists and its archived status.
-- Used for idempotent archive operations.
SELECT EXISTS(SELECT 1 FROM app.notebooks n WHERE n.user_id = @user_id AND n.id = @id) as exists,
       EXISTS(SELECT 1 FROM app.notebooks n WHERE n.user_id = @user_id AND n.id = @id AND n.archived_at IS NOT NULL) as is_archived;
