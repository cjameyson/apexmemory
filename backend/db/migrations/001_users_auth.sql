-- Create schemas (also created by init script, but sqlc needs them here)
CREATE SCHEMA IF NOT EXISTS app;
CREATE SCHEMA IF NOT EXISTS app_code;

-- Extensions
CREATE EXTENSION IF NOT EXISTS "citext";

-- Updated_at trigger function
CREATE OR REPLACE FUNCTION app_code.tg_set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END$$;

-- Auth provider enum
CREATE TYPE app.auth_provider AS ENUM ('password', 'google', 'apple');

-- Users table
CREATE TABLE app.users (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    email citext NOT NULL UNIQUE,
    email_verified_at timestamptz,
    username citext NOT NULL UNIQUE,
    display_name text,
    avatar_url text,
    locale text DEFAULT 'en-US',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    deleted_at timestamptz
);

CREATE TRIGGER trg_users_set_updated_at
BEFORE UPDATE ON app.users
FOR EACH ROW EXECUTE FUNCTION app_code.tg_set_updated_at();

-- Auth identities table (supports password, google, apple auth)
CREATE TABLE app.auth_identities (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    user_id uuid NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
    provider app.auth_provider NOT NULL,
    provider_user_id text NOT NULL,
    email citext,
    email_verified_at timestamptz,
    password_hash text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (provider, provider_user_id),
    UNIQUE (user_id, provider),
    CHECK ((provider <> 'password') OR (password_hash IS NOT NULL AND email IS NOT NULL))
);

CREATE INDEX ix_auth_identities_user ON app.auth_identities (user_id);

CREATE UNIQUE INDEX ux_auth_identities_password_email
  ON app.auth_identities (email)
  WHERE provider = 'password';

CREATE TRIGGER trg_auth_identities_set_updated_at
BEFORE UPDATE ON app.auth_identities
FOR EACH ROW EXECUTE FUNCTION app_code.tg_set_updated_at();

-- User sessions table
CREATE TABLE app.user_sessions (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    user_id uuid NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
    token_hash bytea NOT NULL UNIQUE,
    user_agent text,
    ip_address inet,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    last_used_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX ix_user_sessions_user_id ON app.user_sessions (user_id);
CREATE INDEX ix_user_sessions_expires_at ON app.user_sessions (expires_at);

---- create above / drop below ----

DROP TABLE IF EXISTS app.user_sessions;
DROP TABLE IF EXISTS app.auth_identities;
DROP TABLE IF EXISTS app.users;
DROP TYPE IF EXISTS app.auth_provider;
DROP FUNCTION IF EXISTS app_code.tg_set_updated_at();
DROP EXTENSION IF EXISTS "citext";
