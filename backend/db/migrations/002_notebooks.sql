-- Notebooks table
CREATE TABLE app.notebooks (
    user_id uuid NOT NULL REFERENCES app.users(id) ON DELETE CASCADE,
    id uuid NOT NULL DEFAULT uuidv7(),
    name text NOT NULL,
    description text,
    emoji text,
    color text,
    fsrs_settings jsonb NOT NULL DEFAULT '{
        "desired_retention": 0.9,
        "version": "6",
        "params": [0.212, 1.2931, 2.3065, 8.2956, 6.4133, 0.8334, 3.0194, 0.001, 1.8722, 0.1666, 0.796, 1.4835, 0.0614, 0.2629, 1.6483, 0.6014, 1.8729, 0.5425, 0.0912, 0.0658, 0.1542],
        "learning_steps": [60, 600],
        "relearning_steps": [600],
        "maximum_interval": 36500,
        "enable_fuzzing": true
    }'::jsonb,
    position integer NOT NULL DEFAULT 0,
    archived_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, id),
    CONSTRAINT notebooks_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT notebooks_name_length CHECK (length(name) <= 255)
);

CREATE INDEX ix_notebooks_user_list ON app.notebooks (user_id, archived_at, position);

CREATE TRIGGER trg_notebooks_set_updated_at
BEFORE UPDATE ON app.notebooks
FOR EACH ROW EXECUTE FUNCTION app_code.tg_set_updated_at();

---- create above / drop below ----

DROP TRIGGER IF EXISTS trg_notebooks_set_updated_at ON app.notebooks;
DROP INDEX IF EXISTS app.ix_notebooks_user_list;
DROP TABLE IF EXISTS app.notebooks;
