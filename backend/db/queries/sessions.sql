-- name: CreateSession :one
INSERT INTO app.user_sessions (user_id, token_hash, user_agent, ip_address, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetSessionByToken :one
SELECT s.*, u.email, u.username
FROM app.user_sessions s
JOIN app.users u ON u.id = s.user_id
WHERE s.token_hash = $1
  AND s.expires_at > now()
  AND u.deleted_at IS NULL;

-- name: UpdateSessionLastUsed :exec
UPDATE app.user_sessions
SET last_used_at = now()
WHERE token_hash = $1;

-- name: DeleteSession :exec
DELETE FROM app.user_sessions
WHERE token_hash = $1;

-- name: DeleteUserSessions :exec
DELETE FROM app.user_sessions
WHERE user_id = $1;

-- name: DeleteExpiredSessions :exec
DELETE FROM app.user_sessions
WHERE expires_at < now();
