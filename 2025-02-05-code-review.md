# Backend Code Review - 2025-02-05

Four-agent parallel review covering correctness, performance, maintainability, and security.

---

## Correctness

### Medium

**1. `CountCardsByNotebook` missing state filter**
- `backend/internal/app/cards.go:43-46`
- `ListCards` passes the `state` filter to `ListCardsByNotebook` but NOT to `CountCardsByNotebook`. When filtering by state (e.g., `?state=review`), the list returns filtered cards but `total` reflects the unfiltered count. Pagination metadata (`has_more`, `total`) is wrong.
- Fix: Pass `State: state` to `CountCardsByNotebookParams`.

**2. Archived notebook cards appear in study queue**
- `backend/db/queries/reviews.sql:1-33` (GetStudyCards)
- `backend/db/queries/reviews.sql:92-106` (GetStudyCountsByNotebook)
- `backend/db/queries/reviews.sql:35-45` (GetPracticeCards)
- None of these queries join with `app.notebooks` to exclude archived notebooks (`archived_at IS NOT NULL`). After archiving a notebook, its cards still appear in study sessions and practice mode.
- Fix: Add `JOIN app.notebooks n ON n.user_id = c.user_id AND n.id = c.notebook_id AND n.archived_at IS NULL`.

**3. `buried_until` filter prevents auto-unbury**
- All study/practice/count queries use `buried_until IS NULL` (`reviews.sql:10,43,52,101`, `cards.sql:15,23`, `facts.sql:26,50`)
- Cards buried until a past date remain excluded indefinitely. The semantics of `buried_until` imply the card should become eligible once the date passes.
- Fix: Change `buried_until IS NULL` to `(buried_until IS NULL OR buried_until <= CURRENT_DATE)`.

### Low

**4. Study counts response includes stale archived notebook entries**
- `backend/internal/app/reviews.go:728-740`
- `GetStudyCountsByNotebook` returns counts for all notebooks including archived. The overlay loop adds entries with `Total: 0` but non-zero `Due`/`New` for archived notebooks.
- Fix: Skip notebook IDs not present in the `counts` map, or fix query per issue #2.

**5. `GetNotebookFSRSSettings` doesn't filter by `archived_at IS NULL`**
- `backend/db/queries/notebooks.sql:64-67`
- Inconsistent with `GetNotebook` which excludes archived. Minimal impact since fallback to defaults is safe.

**6. `GetCardForReview` doesn't verify notebook is not archived**
- `backend/db/queries/reviews.sql:54-57`
- Only checks `user_id` and `id`. A card from an archived notebook can still be reviewed if the client submits the review directly.

---

## Performance

### High

**1. Correlated subquery in `GetStudyCards` scans reviews table per candidate**
- `backend/db/queries/reviews.sql:17-22`
- The new-card cap subquery scans the reviews table for every candidate card row. No index on `(user_id, state_before, mode, reviewed_at)` covers this pattern.
- Fix: (a) Add partial index: `CREATE INDEX ix_reviews_new_today ON app.reviews(user_id, reviewed_at) WHERE state_before = 'new' AND mode = 'scheduled';` (b) Restructure as a CTE computed once.

**2. `GetFactStatsByNotebook` runs 6 sequential correlated subqueries**
- `backend/db/queries/facts.sql:46-53`
- 6 independent count subqueries in a single SELECT. Each is O(n) where n is facts/cards in the notebook.
- Fix: Consolidate with `count(*) FILTER (WHERE ...)` via JOINs.

**3. ILIKE search on JSONB cast to text**
- `backend/db/queries/facts.sql:30`
- `f.content::text ILIKE '%' || @search || '%'` forces sequential scan. No index can accelerate this.
- Fix (future): GIN trigram index with `pg_trgm`, or extract searchable text to a `tsvector` column.

**4. N+1 pattern in `cleanupOrphanedAssets`**
- `backend/internal/app/facts.go:201-240`
- Iterates over each removed asset ID with separate `CountFactsReferencingAsset` + `DeleteAsset` + `Storage.Delete` calls per asset.
- Fix: Batch query to find all truly orphaned assets in one query, then batch delete.

### Medium

**5. `ListCards` count query doesn't match list query filters**
- `backend/internal/app/cards.go:43-48`
- `CountCardsByNotebook` ignores the `state` filter. Wastes a query and returns incorrect pagination metadata.

**6. Redundant notebook lookup in `CreateFactHandler`**
- `backend/internal/app/facts_handlers.go:25-32`
- Pre-checks notebook existence before `CreateFact`, but FK constraint would catch it. Extra DB round-trip on every fact creation.

**7. Study counts handler does 2 queries where 1 suffices**
- `backend/internal/app/reviews.go:705-747`
- Calls `GetStudyCountsByNotebook` then `ListNotebooks` separately. The second returns all columns including JSONB `fsrs_settings` when only `id` and `total_cards` are needed.

**8. `reviewed_at::date` cast prevents index usage**
- `backend/db/queries/reviews.sql:160,171,180`
- Casting `reviewed_at` to date on the column side prevents the planner from using the `ix_reviews_notebook_time` index.
- Fix: Rewrite as range: `reviewed_at >= @date AND reviewed_at < @date + interval '1 day'`.

**9. `computeIntervalsWithScheduler` runs FSRS 4 times per study card**
- `backend/internal/app/reviews.go:442-451`
- 4 calls per card x 20 cards = 80 FSRS computations per study request. Likely acceptable (pure math), but `time.Now().UTC()` could be shared.

**10. Session validation on every request hits DB**
- `backend/internal/app/auth.go:175-202`
- No in-memory session cache. One DB round-trip per authenticated request.
- Fix: Short-lived LRU cache (30-60s TTL) for validated sessions.

**11. JSON `MarshalIndent` in non-production**
- `backend/internal/app/helpers.go:96-98`
- Allocates more memory than `json.Marshal`. Minor, dev/staging only.

### Low

**12. `ListNotebooks` returns all notebooks without pagination**
- `backend/db/queries/notebooks.sql:20-22`
- No LIMIT/OFFSET. Unbounded but unlikely to be large.

**13. `DeleteExpiredSessions` is unbounded**
- `backend/db/queries/sessions.sql:28-29`
- Deletes all expired sessions in a single statement. Could lock table if sessions accumulate.

**14. Asset upload reads entire file into memory**
- `backend/internal/app/assets.go:64`
- Up to 10MB per upload into memory. Acceptable at current scale.

**15. `db.New(tx)` allocates per transaction**
- `backend/internal/app/tx.go:22`
- Unavoidable sqlc pattern. Negligible cost.

### Missing Indexes

1. **HIGH**: Reviews new-card-cap query needs `(user_id, reviewed_at) WHERE state_before = 'new' AND mode = 'scheduled'`
2. **MEDIUM**: Facts search needs trigram GIN index if search is used at scale
3. **MEDIUM**: `reviewed_at::date` casts should be range-based to use existing indexes

---

## Maintainability

### High

**1. Redundant `IsAnonymous()` checks in every handler**
- ~20 occurrences across `notebooks_handlers.go`, `facts_handlers.go`, `cards_handlers.go`, `reviews_handlers.go`, `assets_handlers.go`, `users.go`
- `RequireAuth` middleware already rejects anonymous users. This is ~60-70 lines of identical boilerplate.
- Fix: Extract `MustUser(r.Context())` helper that panics (caught by RecoverPanic) or returns user directly.

**2. Response shape inconsistency**
- Auth responses (`auth_handlers.go:100-108,149-158`): `map[string]interface{}`
- User response (`users.go:30-34`): `map[string]any`
- Fact update (`facts_handlers.go:242-247`): inline `map[string]any`
- Other handlers use typed structs consistently.
- Fix: Define typed `AuthResponse` struct. Deduplicate login/register response construction.

### Medium

**3. Duplicate row-to-card mapper functions**
- `backend/internal/app/reviews.go:655-701`
- `studyRowToCard()` and `practiceRowToCard()` are field-by-field identical except input type. Same with `toFactResponse`/`toFactListResponse`/`toFactFilteredResponse` in `facts.go:78-110`.

**4. Global `slog` usage vs context logger**
- `facts.go:213-217,225-237`, `assets.go:123`, `reviews.go:559,569`
- These functions call `slog.Error()`/`slog.Warn()` directly instead of `GetLogger(ctx)`, losing `request_id`, `user_id`, and `trace_id` correlation.
- Fix: Use `GetLogger(ctx)` consistently. All affected functions already receive `ctx`.

**5. Inconsistent UUID parsing**
- `reviews_handlers.go:144-146` uses manual `uuid.Parse(r.PathValue("id"))` while other handlers use `app.PathUUID(w, r, "id")`.
- Fix: Use `PathUUID` everywhere.

**6. `GetStudyCardsHandler` has custom limit parsing**
- `reviews_handlers.go:20-25`
- Inline parsing with different default (20 vs 50) diverges from `parsePagination` pattern.

**7. Large files: `reviews.go` (984 LOC) and `facts.go` (743 LOC)**
- `reviews.go` mixes type definitions, business logic, helpers, and FSRS integration.
- Fix: Extract `fsrs_bridge.go` or similar for scheduler/cache logic.

**8. Auth handler uses `map[string]interface{}` instead of typed struct**
- `auth_handlers.go:100-108,149-157`
- Typos in map keys won't be caught at compile time.

**9. Test coverage gaps**
- Missing test files for: `cards.go`/`cards_handlers.go`, `ratelimit.go` (265 LOC), `middleware.go` (237 LOC), `background.go`, `response.go`, `apperror.go`.
- Cards endpoints and rate limiting logic are completely untested.

### Low

**10. `toFactListResponse` and `ListFacts` appear unused (dead code)**
- `facts.go:95-110` and `facts.go:513-533`
- `ListFactsHandler` uses `ListFactsFiltered` instead.

**11. `errAssetNotFound` declared but never used**
- `assets.go:32`

**12. Handler-service file naming inconsistency**
- `users.go` contains `GetCurrentUserHandler` — should be in `users_handlers.go`.
- `app.go` contains `HealthcheckHandler`.

**13. `Version` constant is hardcoded**
- `app.go:16` — `const Version = "1.0.0"`. Should use `-ldflags` injection at build time.

**14. No `Unwrap()` on `FactValidationError` and `AssetValidationError`**
- `facts.go:28-34`, `assets.go:35-41`
- Deviates from `AppError` pattern which implements `Unwrap()`.

**15. Test helper `createTestNotebook` always uses same name**
- `facts_test.go:15-24` — Makes debugging test failures harder.

**16. `TruncateTables` doesn't include `assets` table**
- `testutil/testutil.go:90-101` — Could leak state between tests.

---

## Security

Overall: Strong security posture. Zero critical or high findings.

### Medium

**1. Request ID spoofing via `X-Request-ID` header**
- `backend/internal/app/middleware.go:23-25`
- Accepts externally-provided `X-Request-ID` without validation. Attacker can inject arbitrary strings that could corrupt structured log parsers.
- Fix: Validate against UUID pattern before accepting. Reject malformed values.
- CVSS: 4.3

**2. Content-Type declared vs detected mismatch on asset uploads**
- `backend/internal/app/assets.go:57-83`
- The declared content type from the multipart header is stored and served back, but may differ from the detected type. `http.DetectContentType` has known limitations.
- Fix: Override declared type with detected type, or require exact match. Set `X-Content-Type-Options: nosniff`.
- CVSS: 4.0

**3. Missing `X-Content-Type-Options: nosniff` on asset serving**
- `backend/internal/app/assets_handlers.go:126-129`
- Browsers could MIME-sniff the response and execute content unexpectedly.
- Fix: Add `w.Header().Set("X-Content-Type-Options", "nosniff")`.
- CVSS: 4.0

**4. No `Content-Disposition` header on asset serving**
- `backend/internal/app/assets_handlers.go:126-132`
- Assets served inline without Content-Disposition. If SVG is ever added to the allowlist, inline rendering becomes an XSS vector.
- Fix: Set `Content-Disposition: inline` for images, `attachment` as default for future types. Never add SVG without sanitization.
- CVSS: 3.5

### Low

**5. 30-day session duration without absolute lifetime cap**
- `backend/internal/app/auth_handlers.go:15`
- Active sessions kept alive indefinitely through `UpdateSessionLastUsed`. Stolen tokens valid for 30 days.
- Fix: Consider shorter duration (7 days) or sliding window with absolute maximum.
- CVSS: 3.0

**6. No CORS configuration**
- `backend/internal/app/routes.go`
- Not currently needed (API proxied through SvelteKit). Would need explicit config if API is ever exposed directly.
- CVSS: 2.0

**7. Missing security headers on API responses**
- `backend/internal/app/middleware.go`
- No global middleware sets `X-Content-Type-Options`, `X-Frame-Options`, `Strict-Transport-Security`, or `Cache-Control`.
- Fix: Add security headers middleware.
- CVSS: 2.0

**8. Email not normalized before storage**
- `backend/internal/app/auth_handlers.go:44-47`
- Emails validated with `mail.ParseAddress()` but not lowercased. Schema uses `citext` which handles uniqueness, but Go string comparisons outside DB would be case-sensitive.
- Fix: `strings.ToLower()` before storing.
- CVSS: 1.5

### Positive Security Findings (no action required)

- Argon2id password hashing (64MB memory, 3 iterations, 2 parallelism)
- `crypto/subtle.ConstantTimeCompare` for password comparison
- 256-bit session tokens from `crypto/rand`, stored as SHA-256 hashes
- All SQL parameterized via sqlc, no dynamic SQL construction
- Every query scopes by `user_id` with composite PKs `(user_id, id)`
- Rate limiting: login (5/min), register (5/hour), API (120/min)
- File upload: 10MB limit, content-type allowlist, byte detection, `safePath()` traversal prevention
- JSON body limited to 1MB, pagination capped at 100, password length bounded (8-128)
- Generic auth error messages prevent user enumeration
- Trusted proxy configuration with `RemoteAddr` fallback
- DSN masking in startup logs
