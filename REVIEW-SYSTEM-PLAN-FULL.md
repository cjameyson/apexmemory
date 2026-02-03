# Review System Implementation Plan

## Progress Summary

| Phase | Status | Description |
|-------|--------|-------------|
| **Phase 1** | COMPLETED | Backend Core -- Due Cards + Submit Review |
| **Phase 2** | COMPLETED | Frontend Wiring -- Connect Focus Mode to Real API |
| **Phase 3** | COMPLETED | Due Counts + Review Launcher Polish |
| **Phase 4** | COMPLETED | Undo + Learning Card Queue Management |
| **Phase 5** | NOT STARTED | Session Stats + Review History |

**Last Updated:** 2026-02-03 (Phase 4 completed)

---

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

## Phase 2: Frontend Wiring -- Connect Focus Mode to Real API [COMPLETED]

**Goal:** Replace mock data with real API calls. Full end-to-end review session works for both scheduled and practice modes, both notebook-scoped and global.

### New Files
- [x] `frontend/src/routes/api/reviews/study/+server.ts` -- GET proxy (forwards `?notebook_id=&limit=`)
- [x] `frontend/src/routes/api/reviews/practice/+server.ts` -- GET proxy (forwards `?notebook_id=&limit=`)
- [x] `frontend/src/routes/api/reviews/+server.ts` -- POST proxy
- [x] `frontend/src/lib/services/reviews.ts` -- service layer for card fetching and display extraction
- [x] `frontend/src/lib/types/review.ts` -- `StudyCard`, `CardDisplay`, `ReviewMode` types

### Modified Files
- [x] `frontend/src/lib/components/overlays/focus-mode.svelte` -- major rewrite:
  - [x] Accept `cards` as prop (pre-fetched by parent), remove all mock imports
  - [x] Accept `mode: 'scheduled' | 'practice'` prop -- controls which mode is sent in POST
  - [x] Render card content from `fact.content` JSONB via `extractCardDisplay()` (handles basic and cloze)
  - [x] Pass real `intervals` to `RatingButtons`
  - [x] On rate: POST to `/api/reviews` with `crypto.randomUUID()`, track `duration_ms`, include `mode`
  - [x] Handle learning cards re-entering queue (insert at end when due within 10min) -- scheduled mode only
  - [x] Fire-and-forget POST with simple retry on failure (2 attempts, 2s delay)
  - [x] Visual indicator for practice mode ("Practice Mode" badge in header)
- [x] `frontend/src/lib/components/navigation/review-launcher.svelte` -- rewrite:
  - [x] "Review All" option first (global scheduled review)
  - [x] List of notebooks with due counts (uses `notebook.dueCount`)
  - [x] Practice mode option (per-notebook or global)
  - [x] Fetch real cards on selection before opening focus mode, with loading state
- [x] `frontend/src/lib/types/stats.ts` -- updated `ReviewScope` to include mode
- [x] `frontend/src/lib/api/types.ts` -- added `ApiStudyCard`, `ApiReviewRequest`, `ApiReviewResponse`
- [x] Mock data not imported in review flow (focus-mode.svelte and review-launcher.svelte use real API)

### Key Design Details
- **Card rendering:** `extractCardDisplay()` in `reviews.ts` extracts front/back from fact content based on `fact_type` (basic: field[0]=front, field[1]=back; cloze: replaces target cloze with `[...]` on front, shows filled on back).
- **UUID client-side:** Uses `crypto.randomUUID()` (standard UUID v4, not v7).
- **Optimistic transitions:** Show next card immediately; POST fires in background with retry.
- **Practice indicator:** "Practice Mode" amber badge in focus-mode header.

### Exit Criteria
- [x] Focus mode opens with real cards from DB (both scheduled and practice modes)
- [x] Each rating fires a real POST, card advances
- [x] Interval labels on buttons match server-computed values
- [x] Session completes, "all done" screen shows
- [x] Global review works (cards from multiple notebooks)
- [x] Practice mode: cards served regardless of due date, visual indicator present
- [x] No mock data imports remain in review flow

### Notes
- Mocks directory still exists but is not imported in review flow. Used elsewhere for non-review UI placeholders.
- `dueCount` on notebooks currently returns 0 (placeholder) -- real counts deferred to Phase 3.
- Client uses UUID v4 (`crypto.randomUUID()`) not v7 -- functionally equivalent for idempotency.

---

## Phase 3: Due Counts + Review Launcher Polish [COMPLETED]

**Goal:** Accurate due counts in the review launcher dropdown. Polished UX for selecting review scope and mode.

### Modified Files
- **Backend:**
  - [x] `backend/db/queries/reviews.sql` -- added `GetStudyCountsByNotebook`: returns `notebook_id`, `due_count` (non-new cards due now), `new_count` (new cards) grouped by notebook, excluding suspended/buried cards
  - [x] `backend/internal/app/reviews.go` -- added `getStudyCounts()` service method, `StudyCountsResponse` and `NotebookStudyCounts` types
  - [x] `backend/internal/app/reviews_handlers.go` -- added `GetStudyCountsHandler` for `GET /v1/reviews/study-counts`
  - [x] `backend/internal/app/routes.go` -- registered route
  - [x] `backend/internal/app/reviews_test.go` -- added 4 tests: empty counts, new cards, after review, totals match sum
- **Frontend:**
  - [x] `frontend/src/lib/api/types.ts` -- added `ApiStudyCountsResponse` type
  - [x] `frontend/src/routes/api/reviews/study-counts/+server.ts` -- GET proxy
  - [x] `frontend/src/routes/(app)/+layout.server.ts` -- parallel fetch of notebooks and study counts
  - [x] `frontend/src/lib/services/notebooks.ts` -- added `toNotebooksWithCounts()` to merge API notebooks with study counts
  - [x] `frontend/src/routes/(app)/+layout.svelte` -- uses `toNotebooksWithCounts()` to populate real `dueCount` and `totalCards`

### Exit Criteria
- [x] Review launcher shows accurate due counts per notebook
- [x] Total due count shown on "Review All" option
- [x] Counts update on page navigation after completing a review session

### Key Design Details
- **Response format:** `{ counts: { [notebook_id]: { due, new, total } }, total_due, total_new }` -- map for efficient frontend lookups
- **Due count = due + new:** Frontend sums `due_count + new_count` for display (all reviewable cards)
- **Parallel fetch:** Layout loads notebooks and counts in `Promise.all()` to avoid waterfall
- **Graceful fallback:** If counts endpoint fails, notebooks show 0 counts (UI still works)

---

## Phase 4: Undo + Learning Card Queue Management [COMPLETED]

**Goal:** Undo last review within a session. Proper handling of learning/relearning cards that come due within the session.

### New Files
- N/A (undo logic integrated into `reviews.go` and `reviews_handlers.go`)

### Modified Files
- **Backend:**
  - [x] `backend/db/queries/reviews.sql` -- added `GetReviewByID`, `GetLatestReviewForCard`, `DeleteReview` (`:execrows`), `RestoreCardAfterUndo`
  - [x] `backend/internal/app/routes.go` -- added `DELETE /v1/reviews/{id}`
  - [x] `backend/internal/app/reviews.go` -- added `undoReview()` service method, `UndoReviewResponse` type
  - [x] `backend/internal/app/reviews_handlers.go` -- added `UndoReviewHandler`
  - [x] `backend/internal/app/reviews_test.go` -- added 4 undo tests: success, not-latest (409), not-found (404), practice mode
- **Frontend:**
  - [x] `frontend/src/routes/api/reviews/[reviewId]/+server.ts` -- DELETE proxy
  - [x] `frontend/src/lib/api/types.ts` -- added `ApiUndoReviewResponse` type
  - `frontend/src/lib/components/overlays/focus-mode.svelte`:
    - [x] Track `lastReview` (UndoState) in session state with `reviewId`, `cardBefore`, `insertPosition`, `wasRequeued`
    - [x] Show undo toast/snackbar for 8s after each rating with "Undo" action button
    - [x] On undo: DELETE, re-insert card at original position with restored FSRS state, reset reveal state
    - [x] Handle requeued card removal on undo (removes duplicate from end of queue)
    - [x] Z key shortcut for undo
    - [x] Learning card re-queue (scheduled mode only): when POST response shows card due within session window (<10min), insert it back into the queue at end *(implemented in Phase 2)*
    - [x] Practice mode: no re-queue logic needed (cards are already all served)

### Exit Criteria
- [x] Undo appears after rating, clicking it restores card for re-rating
- [x] Undo of non-latest review returns 409
- [x] Learning cards with short intervals reappear in scheduled session *(implemented in Phase 2)*
- [x] Undo works in both scheduled and practice modes (practice undo deletes the practice review log)

### Key Design Details
- **Undo window:** 8 seconds after each rating, toast with action button
- **State restoration:** For scheduled reviews, card FSRS state (stability, difficulty, due, reps, lapses) restored from previous review's `prev_*` columns. For practice reviews, only review log deleted (card state unchanged).
- **Conflict detection:** Only the most recent review for a card can be undone (verified via `GetLatestReviewForCard`). Attempting to undo an earlier review returns 409 Conflict.
- **Requeue handling:** If the undone card was requeued (learning card due soon), it's removed from the queue end before being re-inserted at original position.

---

## Phase 5: Session Stats + Review History [NOT STARTED]

**Goal:** Post-session summary and review history for analytics.

### New Files
- [ ] `backend/db/queries/review_stats.sql` -- aggregation queries

### Modified Files
- **Backend:**
  - [ ] `backend/internal/app/reviews_handlers.go` -- add:
    - `GET /v1/reviews/summary?date=YYYY-MM-DD` -- daily summary (total, by rating, by mode, duration, new cards seen)
    - `GET /v1/notebooks/{notebook_id}/reviews?date=&limit=&offset=` -- paginated review history
  - [ ] `backend/internal/app/routes.go` -- register
- **Frontend:**
  - [ ] `frontend/src/lib/components/overlays/focus-mode.svelte` -- session complete screen shows: cards reviewed, time spent, rating breakdown (tracked locally during session, no extra API call). Differentiate stats display for practice vs scheduled.

### Exit Criteria
- [ ] Session completion shows real stats (count, time, breakdown by rating)
- [ ] Summary endpoint returns correct daily aggregations (separates scheduled vs practice)
- [ ] Review list endpoint supports pagination + date filter

### Current State
- Session complete screen exists with basic "You reviewed N cards" message.
- No rating breakdown, time tracking display, or differentiation between practice/scheduled stats.
- No backend stats endpoints exist.

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
