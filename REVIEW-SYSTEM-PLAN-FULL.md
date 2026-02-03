# Review System Implementation Plan

## Overview

Wire up the complete review flow: backend FSRS scheduling, review submission, due-card queue, and connect the existing focus-mode UI to real data. Supports both **notebook-scoped** and **global** review from Phase 1. Includes a **practice mode** that tracks activity without affecting FSRS state. Delivered in 5 phases with human review gates between each.

### Review Modes

| Mode | Cards served | FSRS update | Review logged | Streaks/activity |
|------|-------------|-------------|---------------|-----------------|
| **Scheduled** | Due cards only (overdue + learning + new card cap) | Yes | Yes | Yes |
| **Practice** | All cards in scope (no due filter) | No | Yes (flagged as practice) | Yes |

### Scoping

Two GET endpoints with optional `notebook_id` query param. When omitted, they operate globally across all user notebooks.

- **Study:** `GET /v1/reviews/study?notebook_id=&limit=20` -- due cards (optional notebook scope)
- **Practice:** `GET /v1/reviews/practice?notebook_id=&limit=50` -- all cards (optional notebook scope)
- **Submit:** `POST /v1/reviews` -- single endpoint for both modes

---

## Phase 1: Backend Core -- Due Cards + Submit Review [COMPLETED]

**Goal:** Three endpoints: study (due cards), practice (all cards), and submit review. Optional `notebook_id` query param for scoping. Both modes work end-to-end.

### Schema Change
- [x] Add `mode` column to `app.reviews`: `TEXT NOT NULL DEFAULT 'scheduled' CHECK (mode IN ('scheduled', 'practice'))`.
- [x] Fold into migration `003` - do not create a new migration.

### New Files
- [x] `backend/db/queries/reviews.sql` -- sqlc queries:
  - [x] `GetStudyCards` -- cards where `due <= now()` or `state = 'new'`, not suspended/buried, joined with `facts` for content. New card cap via subquery counting today's new-card reviews. Optional `notebook_id` filter (use `sqlc.narg` for nullable param). `LIMIT` param.
  - [x] `GetPracticeCards` -- all cards (no due filter), joined with facts. Optional `notebook_id` filter. `LIMIT`/`OFFSET` for pagination.
  - [x] `CountPracticeCards` -- count query for pagination.
  - [x] `CreateReview` -- insert into `app.reviews` (includes `mode` column) with `ON CONFLICT (user_id, id) DO NOTHING` for idempotency, returning row.
  - [x] `UpdateCardAfterReview` -- update all FSRS columns on `app.cards`.
  - [x] `GetCardForReview` -- `SELECT ... FOR UPDATE` row lock.
- [x] `backend/internal/app/reviews.go` -- business logic:
  - [x] `submitReview()` -- in `WithTx`: lock card, build `fsrs.Scheduler`, call `ReviewCard()`, insert review row. **If mode=scheduled:** update card FSRS state. **If mode=practice:** skip card update, only insert review log. Handles idempotent re-submit via `pgx.ErrNoRows` on conflict.
  - [x] `getStudyCards()` -- call sqlc query, then for each card compute intervals for all 4 ratings via `ReviewCard()` (no fuzzing) and attach to response.
  - [x] `getPracticeCards()` -- call sqlc query, compute intervals same as study cards (informational only in practice mode).
  - [x] Helpers: `dbCardToFSRS`, `fsrsStateToDBState`, `ratingFromString`, `formatInterval`, row-to-card converters.
- [x] `backend/internal/app/reviews_handlers.go` -- HTTP handlers:
  - [x] `GET /v1/reviews/study?notebook_id=&limit=20`
  - [x] `GET /v1/reviews/practice?notebook_id=&limit=50`
  - [x] `POST /v1/reviews` -- request body includes `mode: "scheduled"|"practice"`. Single endpoint for both modes.
- [x] `backend/internal/app/reviews_test.go` -- 9 integration tests:
  - [x] Submit scheduled review on new card (FSRS state updates)
  - [x] Submit practice review (FSRS state unchanged, review logged)
  - [x] Idempotent re-submit (same review ID returns 200)
  - [x] Card not found returns 404
  - [x] Invalid rating returns 400
  - [x] Study cards with intervals for all 4 ratings
  - [x] No cards returned after scheduled review (card due in future)
  - [x] Practice returns all cards regardless of due date
  - [x] Global vs notebook-scoped queries
  - [x] `formatInterval` unit tests

### Modified Files
- [x] `backend/internal/app/routes.go` -- registered 3 new routes under `protected`.

### Key Design Details
- **FSRS mapping:** DB `card_state` enum (new/learning/review/relearning) maps to fsrs package `State` (Learning=1, Review=2, Relearning=3). New cards have `state='new'` in DB and get initialized as fresh cards in the scheduler.
- **Intervals response:** `{ again: "10m", hard: "1d", good: "3d", easy: "7d" }` -- human-readable strings computed server-side.
- **Review request:** `{ id: UUIDv7, card_id: UUID, rating: 1-4, duration_ms: int, mode: "scheduled"|"practice" }`
- **Review response:** created review + updated card state (or unchanged state for practice) + next due date.
- **Single POST endpoint:** `/v1/reviews` (not nested under notebook). The `card_id` already implies the notebook. Simplifies the client.

### Exit Criteria
- [x] `make db.sqlc` succeeds
- [x] `make test.backend` passes all new tests (9 integration + 1 unit)
- [x] curl: create fact -> GET due (global) returns card with intervals -> POST scheduled review (good) -> card state=learning, stability/difficulty set -> GET due returns empty
- [x] curl: POST practice review (easy) -> card state remains 'new', review logged with `mode='practice'`
- [x] curl: GET practice cards returns all cards (2) regardless of due date, paginated with `total`/`has_more`

### Notes
- Due card ordering and new card cap tests not yet written (covered implicitly by study cards query logic). Can add if needed.
- Learning step re-entry test deferred to Phase 4 (undo + learning queue management).

---

## Phase 2: Frontend Wiring -- Connect Focus Mode to Real API

**Goal:** Replace mock data with real API calls. Full end-to-end review session works for both scheduled and practice modes, both notebook-scoped and global.

### New Files
- `frontend/src/routes/api/reviews/study/+server.ts` -- GET proxy (forwards `?notebook_id=&limit=`)
- `frontend/src/routes/api/reviews/practice/+server.ts` -- GET proxy (forwards `?notebook_id=&limit=`)
- `frontend/src/routes/api/reviews/+server.ts` -- POST proxy

### Modified Files
- `frontend/src/lib/components/overlays/focus-mode.svelte` -- major rewrite:
  - Accept `cards` as prop (pre-fetched by parent), remove all mock imports
  - Accept `mode: 'scheduled' | 'practice'` prop -- controls which mode is sent in POST
  - Render card content from `fact.content` JSONB (handle basic front/back and cloze)
  - Pass real `intervals` to `RatingButtons`
  - On rate: POST to `/api/reviews` with client-generated UUIDv7, track `duration_ms` (timestamp delta), include `mode`
  - Handle learning cards re-entering queue (insert back at appropriate position when due within session) -- only in scheduled mode
  - Fire-and-forget POST with simple retry on failure
  - Visual indicator for practice mode (subtle badge/label so user knows FSRS isn't updating)
- `frontend/src/lib/components/navigation/review-launcher.svelte` -- rewrite:
  - "Review All" option first (global scheduled review)
  - List of active (non-archived) notebooks with due counts
  - Practice mode option (per-notebook or global)
  - Fetch real due cards on selection before opening focus mode, with loading state
- `frontend/src/lib/types/` -- add `DueCard`, `ReviewMode` types matching backend response
- `frontend/src/lib/types/stats.ts` -- update `ReviewScope` to include mode:
  ```ts
  export type ReviewScope =
    | { type: 'all'; mode: ReviewMode }
    | { type: 'notebook'; notebook: Notebook; mode: ReviewMode }
    | { type: 'source'; notebook: Notebook; source: Source; mode: ReviewMode };
  ```
- Remove or gut mock card data (`$lib/mocks/`)

### Key Design Details
- **Card rendering:** `focus-mode.svelte` currently uses `currentCard.front`/`.back`. Real cards have `content.fields` JSONB array. Need a small rendering function that extracts front/back from fact content based on `fact_type` (basic: field[0]=front, field[1]=back; cloze: render with blanks on front, filled on back).
- **UUIDv7 client-side:** Use `crypto.randomUUID()` with timestamp prefix, or install `uuid` package.
- **Optimistic transitions:** Show next card immediately; POST fires in background.
- **Practice indicator:** Subtle "Practice Mode" label in focus-mode header so user has clear context.

### Exit Criteria
- Focus mode opens with real cards from DB (both scheduled and practice modes)
- Each rating fires a real POST, card advances
- Interval labels on buttons match server-computed values
- Session completes, "all done" screen shows
- Global review works (cards from multiple notebooks)
- Practice mode: cards served regardless of due date, visual indicator present
- No mock data imports remain in review flow

---

## Phase 3: Due Counts + Review Launcher Polish

**Goal:** Accurate due counts in the review launcher dropdown. Polished UX for selecting review scope and mode.

### Modified Files
- **Backend:**
  - `backend/db/queries/cards.sql` -- add `GetDueCounts`: `SELECT notebook_id, count(*) FROM app.cards WHERE user_id = @user_id AND (due <= now() OR state = 'new') AND suspended_at IS NULL AND (buried_until IS NULL OR buried_until <= now()) GROUP BY notebook_id`
  - `backend/internal/app/reviews_handlers.go` -- add `GET /v1/reviews/study-counts` handler
  - `backend/internal/app/routes.go` -- register
- **Frontend:**
  - `frontend/src/routes/api/reviews/study-counts/+server.ts` -- GET proxy
  - `frontend/src/routes/(app)/+layout.server.ts` -- fetch due counts in layout load, pass to children
  - `frontend/src/lib/components/navigation/review-launcher.svelte` -- display real due counts per notebook, total due count on "Review All"

### Exit Criteria
- Review launcher shows accurate due counts per notebook
- Total due count shown on "Review All" option
- Counts update on page navigation after completing a review session

---

## Phase 4: Undo + Learning Card Queue Management

**Goal:** Undo last review within a session. Proper handling of learning/relearning cards that come due within the session.

### New Files
- `backend/internal/app/reviews_undo.go` -- undo handler logic

### Modified Files
- **Backend:**
  - `backend/db/queries/reviews.sql` -- add `GetReview`, `GetLatestReviewForCard`, `DeleteReview` (`:execrows`), `RestoreCardState`
  - `backend/internal/app/routes.go` -- add `DELETE /v1/reviews/{id}`
- **Frontend:**
  - `frontend/src/routes/api/reviews/[reviewId]/+server.ts` -- DELETE proxy
  - `frontend/src/lib/components/overlays/focus-mode.svelte`:
    - Track `lastReviewId` in session state
    - Show undo toast/snackbar for ~8s after each rating
    - On undo: DELETE, re-insert card at current position, reset reveal state
    - Learning card re-queue (scheduled mode only): when POST response shows card due within session window (e.g., <10min), insert it back into the queue at the appropriate position
    - Practice mode: no re-queue logic needed (cards are already all served)

### Exit Criteria
- Undo appears after rating, clicking it restores card for re-rating
- Undo of non-latest review returns 409
- Learning cards with short intervals reappear in scheduled session
- Undo works in both scheduled and practice modes (practice undo deletes the practice review log)

---

## Phase 5: Session Stats + Review History

**Goal:** Post-session summary and review history for analytics.

### New Files
- `backend/db/queries/review_stats.sql` -- aggregation queries

### Modified Files
- **Backend:**
  - `backend/internal/app/reviews_handlers.go` -- add:
    - `GET /v1/reviews/summary?date=YYYY-MM-DD` -- daily summary (total, by rating, by mode, duration, new cards seen)
    - `GET /v1/notebooks/{notebook_id}/reviews?date=&limit=&offset=` -- paginated review history
  - `backend/internal/app/routes.go` -- register
- **Frontend:**
  - `frontend/src/lib/components/overlays/focus-mode.svelte` -- session complete screen shows: cards reviewed, time spent, rating breakdown (tracked locally during session, no extra API call). Differentiate stats display for practice vs scheduled.

### Exit Criteria
- Session completion shows real stats (count, time, breakdown by rating)
- Summary endpoint returns correct daily aggregations (separates scheduled vs practice)
- Review list endpoint supports pagination + date filter

---

## Phase Dependencies

```
Phase 1 (Backend core: due + practice + submit)
  └─> Phase 2 (Frontend wiring: both modes, both scopes)
        ├─> Phase 3 (Due counts + launcher polish)
        └─> Phase 4 (Undo + learning queue)
              └─> Phase 5 (Stats + history)
```

Phases 3 and 4 are independent of each other and can be swapped.

## Critical Files Reference
- `backend/internal/fsrs/fsrs.go` -- FSRS v6 scheduler (complete, tested, unused)
- `backend/internal/app/tx.go` -- `WithTx` transaction helper
- `backend/internal/app/cards_handlers.go` -- handler patterns to follow
- `backend/internal/app/helpers.go` -- pagination, OptionalString, etc.
- `backend/db/migrations/003_fact_cards_reviews.sql` -- schema (cards, facts, reviews)
- `frontend/src/lib/components/overlays/focus-mode.svelte` -- review UI (mock-driven)
- `frontend/src/lib/components/cards/rating-buttons.svelte` -- already accepts `intervals` prop
- `frontend/src/lib/components/navigation/review-launcher.svelte` -- scope selector
- `frontend/src/lib/types/stats.ts` -- `ReviewScope` type definition
