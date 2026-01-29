-- Notebooks table
CREATE TABLE app.notebooks (
    user_id uuid NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
    id uuid NOT NULL DEFAULT uuidv7(),
    name text NOT NULL,
    description text,
    emoji text,
    color text,
    position integer NOT NULL DEFAULT 0,
    total_cards integer NOT NULL DEFAULT 0, -- denormalized for performance
    fsrs_settings jsonb NOT NULL DEFAULT '{}'::jsonb,
    archived_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, id),
    CONSTRAINT notebooks_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT notebooks_name_length CHECK (length(name) <= 255)
);

ALTER TABLE app.notebooks
    ADD CONSTRAINT notebooks_total_cards_non_negative CHECK (total_cards >= 0);

CREATE INDEX ix_notebooks_user_list ON app.notebooks (user_id, archived_at, position);

CREATE TRIGGER trg_notebooks_set_updated_at
BEFORE UPDATE ON app.notebooks
FOR EACH ROW EXECUTE FUNCTION app_code.tg_set_updated_at();

---- create above / drop below ----

DROP TRIGGER IF EXISTS trg_notebooks_set_updated_at ON app.notebooks;
DROP INDEX IF EXISTS app.ix_notebooks_user_list;
DROP TABLE IF EXISTS app.notebooks;
