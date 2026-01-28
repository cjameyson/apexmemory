-- =============================================================================
-- NOTES, CARDS, REVIEWS - FSRS v6 Core Schema
-- =============================================================================

-- -----------------------------------------------------------------------------
-- NOTEBOOKS: add guard against negative card count
-- -----------------------------------------------------------------------------

ALTER TABLE app.notebooks
    ADD CONSTRAINT notebooks_total_cards_non_negative CHECK (total_cards >= 0);

-- -----------------------------------------------------------------------------
-- TYPES
-- -----------------------------------------------------------------------------

CREATE TYPE app.card_state AS ENUM ('new', 'learning', 'review', 'relearning');
CREATE TYPE app.rating AS ENUM ('again', 'hard', 'good', 'easy');

-- -----------------------------------------------------------------------------
-- NOTES
-- -----------------------------------------------------------------------------

CREATE TABLE app.notes (
    user_id       UUID NOT NULL,
    id            UUID NOT NULL DEFAULT uuidv7(),
    notebook_id   UUID NOT NULL,
    -- TEXT+CHECK rather than ENUM: note_type may evolve; enums are hard to alter in PG.
    note_type     TEXT NOT NULL DEFAULT 'basic',
    content       JSONB NOT NULL,
    source_id     UUID, -- TODO: FK to sources table once it exists
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, id),

    FOREIGN KEY (user_id, notebook_id)
        REFERENCES app.notebooks(user_id, id) ON DELETE CASCADE,

    CONSTRAINT notes_valid_type CHECK (
        note_type IN ('basic', 'cloze', 'image_occlusion')
    ),
    CONSTRAINT notes_valid_content CHECK (
        jsonb_typeof(content->'version') = 'number' AND
        jsonb_typeof(content->'fields') = 'array'
    )
);

CREATE INDEX ix_notes_notebook ON app.notes(user_id, notebook_id, created_at DESC);

CREATE TRIGGER trg_notes_set_updated_at
    BEFORE UPDATE ON app.notes
    FOR EACH ROW EXECUTE FUNCTION app_code.tg_set_updated_at();

-- -----------------------------------------------------------------------------
-- CARDS
-- -----------------------------------------------------------------------------

CREATE TABLE app.cards (
    user_id           UUID NOT NULL,
    id                UUID NOT NULL DEFAULT uuidv7(),
    -- Denormalized from note; avoids join on the hot "fetch due cards" query.
    -- Notes are not expected to move between notebooks.
    notebook_id       UUID NOT NULL,
    note_id           UUID NOT NULL,
    ordinal           SMALLINT NOT NULL DEFAULT 0,
    state             app.card_state NOT NULL DEFAULT 'new',
    stability         REAL,
    difficulty        REAL,
    step              SMALLINT,
    due               TIMESTAMPTZ,
    last_review       TIMESTAMPTZ,
    elapsed_days      REAL NOT NULL DEFAULT 0,
    scheduled_days    REAL NOT NULL DEFAULT 0,
    reps              INTEGER NOT NULL DEFAULT 0,
    lapses            INTEGER NOT NULL DEFAULT 0,
    suspended_at      TIMESTAMPTZ,
    buried_until      DATE,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, id),

    FOREIGN KEY (user_id, notebook_id)
        REFERENCES app.notebooks(user_id, id) ON DELETE CASCADE,
    FOREIGN KEY (user_id, note_id)
        REFERENCES app.notes(user_id, id) ON DELETE CASCADE,

    UNIQUE (user_id, note_id, ordinal),

    CONSTRAINT cards_valid_difficulty CHECK (
        difficulty IS NULL OR (difficulty >= 1.0 AND difficulty <= 10.0)
    ),
    CONSTRAINT cards_valid_stability CHECK (
        stability IS NULL OR stability >= 0
    ),
    CONSTRAINT cards_valid_step_state CHECK (
        (state IN ('learning', 'relearning') AND step IS NOT NULL AND step >= 0) OR
        (state IN ('new', 'review') AND step IS NULL)
    )
);

CREATE INDEX ix_cards_due ON app.cards(user_id, notebook_id, due)
    WHERE suspended_at IS NULL AND buried_until IS NULL;

CREATE INDEX ix_cards_state ON app.cards(user_id, notebook_id, state);

CREATE INDEX ix_cards_note ON app.cards(user_id, note_id);

CREATE TRIGGER trg_cards_set_updated_at
    BEFORE UPDATE ON app.cards
    FOR EACH ROW EXECUTE FUNCTION app_code.tg_set_updated_at();

CREATE TRIGGER trg_cards_sync_notebook_count
    AFTER INSERT OR DELETE OR UPDATE OF notebook_id ON app.cards
    FOR EACH ROW EXECUTE FUNCTION app_code.tg_sync_notebook_card_count();

-- -----------------------------------------------------------------------------
-- REVIEWS
-- -----------------------------------------------------------------------------

CREATE TABLE app.reviews (
    user_id            UUID NOT NULL,
    id                 UUID NOT NULL DEFAULT uuidv7(),
    card_id            UUID NOT NULL,
    -- Denormalized snapshot: records which notebook the card belonged to at review
    -- time. No FK intentionally -- reviews are immutable history and must survive
    -- notebook deletes or card moves.
    notebook_id        UUID NOT NULL,
    reviewed_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    rating             app.rating NOT NULL,
    review_duration_ms INTEGER,
    state_before       app.card_state NOT NULL,
    stability_before   REAL,
    difficulty_before  REAL,
    elapsed_days       REAL NOT NULL,
    scheduled_days     REAL NOT NULL,
    state_after        app.card_state NOT NULL,
    stability_after    REAL NOT NULL,
    difficulty_after   REAL NOT NULL,
    interval_days      REAL NOT NULL,
    retrievability     REAL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, id),

    FOREIGN KEY (user_id, card_id)
        REFERENCES app.cards(user_id, id) ON DELETE CASCADE
);

CREATE INDEX ix_reviews_card ON app.reviews(user_id, card_id, reviewed_at DESC);

CREATE INDEX ix_reviews_notebook_time ON app.reviews(user_id, notebook_id, reviewed_at DESC);

CREATE INDEX ix_reviews_optimizer ON app.reviews(user_id, notebook_id, state_before, rating)
    WHERE scheduled_days >= 1;

---- create above / drop below ----

DROP TRIGGER IF EXISTS trg_cards_sync_notebook_count ON app.cards;
DROP TRIGGER IF EXISTS trg_cards_set_updated_at ON app.cards;
DROP TRIGGER IF EXISTS trg_notes_set_updated_at ON app.notes;

DROP INDEX IF EXISTS app.ix_reviews_optimizer;
DROP INDEX IF EXISTS app.ix_reviews_notebook_time;
DROP INDEX IF EXISTS app.ix_reviews_card;
DROP TABLE IF EXISTS app.reviews;

DROP INDEX IF EXISTS app.ix_cards_note;
DROP INDEX IF EXISTS app.ix_cards_state;
DROP INDEX IF EXISTS app.ix_cards_due;
DROP TABLE IF EXISTS app.cards;

DROP INDEX IF EXISTS app.ix_notes_notebook;
DROP TABLE IF EXISTS app.notes;

DROP TYPE IF EXISTS app.rating;
DROP TYPE IF EXISTS app.card_state;

ALTER TABLE app.notebooks DROP CONSTRAINT IF EXISTS notebooks_total_cards_non_negative;
