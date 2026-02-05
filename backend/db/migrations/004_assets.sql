CREATE TABLE app.assets (
    user_id      UUID NOT NULL REFERENCES app.users(id),
    id           UUID NOT NULL DEFAULT uuidv7(),
    content_type TEXT NOT NULL,
    filename     TEXT NOT NULL,
    size_bytes   BIGINT NOT NULL,
    sha256       TEXT NOT NULL,
    metadata     JSONB NOT NULL DEFAULT '{}',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, id)
);

CREATE INDEX ix_assets_user ON app.assets(user_id, created_at DESC);
CREATE INDEX ix_assets_sha256 ON app.assets(user_id, sha256);
CREATE INDEX ix_facts_asset_ids ON app.facts USING GIN ((content->'asset_ids'));

CREATE TRIGGER trg_assets_set_updated_at
    BEFORE UPDATE ON app.assets
    FOR EACH ROW
    EXECUTE FUNCTION app_code.tg_set_updated_at();

---- create above / drop below ----

DROP TRIGGER IF EXISTS trg_assets_set_updated_at ON app.assets;
DROP INDEX IF EXISTS ix_facts_asset_ids;
DROP INDEX IF EXISTS ix_assets_sha256;
DROP INDEX IF EXISTS ix_assets_user;
DROP TABLE IF EXISTS app.assets;
