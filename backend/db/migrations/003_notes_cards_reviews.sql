-- =============================================================================
-- NOTES, CARDS, REVIEWS - FSRS v6 Core Schema
-- =============================================================================

-- -----------------------------------------------------------------------------
-- TYPES
-- -----------------------------------------------------------------------------

CREATE TYPE app.card_state AS ENUM ('new', 'learning', 'review', 'relearning');
CREATE TYPE app.rating AS ENUM ('again', 'hard', 'good', 'easy');

-- -----------------------------------------------------------------------------
-- NOTES
-- -----------------------------------------------------------------------------
-- One note can generate one or more cards:
--   - basic:           1 note → 1 card (element_id = '')
--   - cloze:           1 note → N cards (element_id = 'c1', 'c2', ...)
--   - image_occlusion: 1 note → N cards (element_id = 'm_<nanoid>' per mask region)
--
-- Note type changes are NOT allowed after creation.
-- Maximum 128 cards per note (enforced in application layer).

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
-- element_id identifies which part of the note this card tests:
--   - basic:           '' (empty string, single card per note)
--   - cloze:           'c1', 'c2', etc. (matches c1::... in content)
--   - image_occlusion: 'm_<nanoid>' (matches region id in content)
--
-- element_id is stable: deleting one element doesn't affect others' IDs.
-- This preserves review history for unchanged cards when notes are edited.

CREATE TABLE app.cards (
    user_id           UUID NOT NULL,
    id                UUID NOT NULL DEFAULT uuidv7(),
    -- Denormalized from note; avoids join on the hot "fetch due cards" query.
    -- Notes are not expected to move between notebooks.
    notebook_id       UUID NOT NULL,
    note_id           UUID NOT NULL,
    -- Identifies which element of the note this card represents.
    -- Empty string for basic notes, 'c1'/'c2'/etc for cloze, 'm_xxx' for image masks.
    element_id        TEXT NOT NULL DEFAULT '',
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

    UNIQUE (user_id, note_id, element_id),

    CONSTRAINT cards_valid_difficulty CHECK (
        difficulty IS NULL OR (difficulty >= 1.0 AND difficulty <= 10.0)
    ),
    CONSTRAINT cards_valid_stability CHECK (
        stability IS NULL OR stability >= 0
    ),
    CONSTRAINT cards_valid_step_state CHECK (
        (state IN ('learning', 'relearning') AND step IS NOT NULL AND step >= 0) OR
        (state IN ('new', 'review') AND step IS NULL)
    ),
    -- element_id format validation:
    --   '' for basic, 'c1'-'c999' for cloze, 'm_' + alphanumeric for image masks
    CONSTRAINT cards_valid_element_id CHECK (
        element_id = '' OR                           -- basic
        element_id ~ '^c[1-9][0-9]{0,2}$' OR         -- cloze: c1 to c999
        element_id ~ '^m_[a-zA-Z0-9_-]{6,24}$'       -- image mask: m_ + nanoid
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
-- Reviews are immutable learning history used for:
--   1. User statistics and progress tracking
--   2. FSRS optimizer training data
--   3. Retention analysis
--
-- Denormalization strategy:
--   - notebook_id: Snapshot of which notebook at review time. No FK.
--   - note_id: Snapshot for context. Nullable, SET NULL on note delete.
--   - element_id: Snapshot of which element was reviewed.
--   - card_id: Reference to card. Nullable, SET NULL on card delete.
--
-- When cards/notes are deleted, reviews are preserved with NULL references.
-- This retains training data for the optimizer while allowing cleanup.
--
-- TODO: Consider adding an index for orphaned reviews if cleanup queries needed:
--   CREATE INDEX ix_reviews_orphaned ON app.reviews(user_id, reviewed_at)
--       WHERE card_id IS NULL;

CREATE TABLE app.reviews (
    user_id            UUID NOT NULL,
    id                 UUID NOT NULL DEFAULT uuidv7(),
    -- Nullable: preserved when card is deleted for optimizer training data
    card_id            UUID,
    -- Denormalized snapshots: record context at review time.
    -- These survive deletes and enable stats/optimizer even for orphaned reviews.
    notebook_id        UUID NOT NULL,
    note_id            UUID,
    element_id         TEXT,
    -- Review data
    reviewed_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    rating             app.rating NOT NULL,
    review_duration_ms INTEGER,
    -- State before review (for optimizer training)
    state_before       app.card_state NOT NULL,
    stability_before   REAL,
    difficulty_before  REAL,
    elapsed_days       REAL NOT NULL,
    scheduled_days     REAL NOT NULL,
    -- State after review
    state_after        app.card_state NOT NULL,
    stability_after    REAL NOT NULL,
    difficulty_after   REAL NOT NULL,
    interval_days      REAL NOT NULL,
    retrievability     REAL,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (user_id, id),

    -- SET NULL preserves review history when card is deleted
    FOREIGN KEY (user_id, card_id)
        REFERENCES app.cards(user_id, id) ON DELETE SET NULL,
    -- SET NULL preserves review history when note is deleted
    FOREIGN KEY (user_id, note_id)
        REFERENCES app.notes(user_id, id) ON DELETE SET NULL
);

-- For fetching review history of a specific card
CREATE INDEX ix_reviews_card ON app.reviews(user_id, card_id, reviewed_at DESC)
    WHERE card_id IS NOT NULL;

-- For notebook-level analytics and time-based queries
CREATE INDEX ix_reviews_notebook_time ON app.reviews(user_id, notebook_id, reviewed_at DESC);

-- For FSRS optimizer: needs state_before and rating for graduated cards
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
