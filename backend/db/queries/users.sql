-- name: CreateUser :one
INSERT INTO app.users (email, username, display_name)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateUserWithPassword :one
-- Creates a user and their password auth identity in a single transaction-friendly call
WITH new_user AS (
    INSERT INTO app.users (email, username, display_name)
    VALUES (@email, @username, @display_name)
    RETURNING *
)
INSERT INTO app.auth_identities (user_id, provider, provider_user_id, email, password_hash)
SELECT id, 'password', id::text, @email, @password_hash
FROM new_user
RETURNING (SELECT id FROM new_user), (SELECT email FROM new_user), (SELECT username FROM new_user);

-- name: GetUserByID :one
SELECT * FROM app.users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByEmail :one
SELECT * FROM app.users
WHERE email = $1 AND deleted_at IS NULL;

-- name: GetUserByUsername :one
SELECT * FROM app.users
WHERE username = $1 AND deleted_at IS NULL;

-- name: GetUserByEmailPassword :one
-- Get user with password hash for authentication
SELECT u.id, u.email, u.username, ai.password_hash
FROM app.users u
JOIN app.auth_identities ai ON ai.user_id = u.id
WHERE ai.provider = 'password'
  AND ai.email = @email
  AND u.deleted_at IS NULL;

-- name: UpdateUser :one
UPDATE app.users
SET display_name = COALESCE(sqlc.narg('display_name'), display_name),
    avatar_url = COALESCE(sqlc.narg('avatar_url'), avatar_url),
    locale = COALESCE(sqlc.narg('locale'), locale)
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: SoftDeleteUser :exec
UPDATE app.users
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateAuthIdentity :one
INSERT INTO app.auth_identities (user_id, provider, provider_user_id, email, password_hash)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetAuthIdentityByEmail :one
SELECT ai.*, u.email as "user_email", u.username as "user_username"
FROM app.auth_identities ai
JOIN app.users u ON u.id = ai.user_id
WHERE ai.provider = 'password'
  AND ai.email = $1
  AND u.deleted_at IS NULL;

-- name: GetAuthIdentityByProviderID :one
SELECT ai.*, u.email as "user_email", u.username as "user_username"
FROM app.auth_identities ai
JOIN app.users u ON u.id = ai.user_id
WHERE ai.provider = $1
  AND ai.provider_user_id = $2
  AND u.deleted_at IS NULL;

-- name: UpsertTestUser :one
-- Idempotent test user creation for agent testing.
-- Uses fixed UUID to ensure consistent user_id across reseeds.
INSERT INTO app.users (id, email, username, display_name, created_at, updated_at)
VALUES (@id, @email, @username, @display_name, now(), now())
ON CONFLICT (email) DO UPDATE SET updated_at = now()
RETURNING *;

-- name: UpsertTestAuthIdentity :exec
-- Idempotent auth identity creation for test users.
-- Creates or updates the password auth identity.
INSERT INTO app.auth_identities (id, user_id, provider, provider_user_id, email, password_hash, created_at, updated_at)
VALUES (gen_random_uuid(), @user_id, 'password', @provider_user_id, @email, @password_hash, now(), now())
ON CONFLICT (user_id, provider) DO UPDATE SET password_hash = EXCLUDED.password_hash, updated_at = now();
