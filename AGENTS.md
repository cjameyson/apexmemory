# Apex Memory - Project Context

## Project Identity
**Name:** Apex Memory (apexmemory.ai)
**Mission:** Simple and effective spaced repetition, AI-enhanced card creation, and mobile-first PWA experience.
**Type:** SaaS $12/mo, Mobile-first PWA, Solo-dev (Closed Source).
**Core Values:** Progress > Perfection, Performance First, Clean/Composable Code (No over-abstraction).
**Solo Dev:** Progress and delightful user experiences are very motivating!

## Agent First Design
- How to think about "agent first" design?
- Can apex memory be used by Claude Code, etc? 
- Apex Memory handles the FSRS system and card management
- Agent could generate flashcards from sources via API calls

## Domain
Apex Memory is centered around the notebook concept. A notebooks is typically created by a user for a particular subject, class or area of study.  Notebooks contain sources (pdfs, slides, audio, links, etc), flashcards, notes, etc.  These sources can be viewed and searched in the product, but more importantly can be parsed, chunked, embedded and indexed. This enables the product to:
- generate flashcards with citations
- enable interactive chat grounded in the source material
- semantic search
- gap identification of flashcards based on source material

The core value proposition is a flashcard SRS system with AI-enhanced workflows.
- **Algorithm:** FSRS v6 (Free Spaced Repetition Scheduler)
- **Ratings:** Again (1), Hard (2), Good (3), Easy (4)
- **Fact Types:** `basic`, `cloze`, `image_occlusion`
- **Field Types:** `plain_text`, `rich_text`, `cloze_text`, `image_occlusion`


## Tech Stack & Architecture

### Backend (Go + PostgreSQL)
- **Runtime:** Go 1.25+ (Standard `net/http`, avoid heavy frameworks).
- **Database:** PostgreSQL 18+. **Primary Keys:** UUIDv7 (Time-sortable).
- **Data Access:** `pgx` driver. `sqlc` for type-safe query generation.
- **Migrations:** `tern` (Two-way migrations in `backend/db/migrations`).
- **Auth:** Session-based. `app.user_sessions` table. Token passed via `Authorization: Bearer`.
- **TERN:** Uses magic comment `---- create above / drop below ----` where up SQL is above the comment and down SQL is below.
- **PATCH APIs:** Use `OptionalString`/`OptionalUUID` types in `helpers.go` to distinguish missing vs explicit null in JSON requests.
- **DELETE APIs:** Use `:execrows` in sqlc and check `rowsAffected == 0` for 404 detection. Prefer this atomic single-query pattern over separate existence checks unless special handling (e.g., idempotent archive) is needed.

Refer to `references/guide-backend.md` for an overview of backend structure.

#### Pagination Convention
All list endpoints use `?limit=&offset=` query params with `PageResponse[T]`:
- **Defaults:** limit=50, max=100. Parsed by `parsePagination()` in `helpers.go`.
- **Response shape:** `{"data": [...], "total": N, "has_more": bool}`
- **Helper:** `NewPageResponse(data, total, limit, offset)` computes `has_more` automatically.
- **Client derives:** current page (`floor(offset/limit)+1`), total pages (`ceil(total/limit)`), showing range.

### Core Entities
**MVP**
- Sources: source material the user is studying (pdfs, slides, audio, links, etc)
- Notebooks: container for sources, flashcards, notes, etc.
- Facts: (like a Note+NoteType in Anki) a set of fields with values that generates cards
- Cards: (like a Card in Anki) a reviewable item derived from a Fact
- Reviews: a review of a card by a user
- Assets: files uploaded by the user for use in facts/cards (images, audio, etc)

**Post MVP**
- Source chunks: chunks of source files, split into smaller pieces with embeddings to enable better search and chat grounding
- Chat sessions, messages: user chats within a notebook

### Database
- Use JSONB wisely
- All multi-tenant tables use composite PKs `(user_id, id)` for efficient partitioning.
- Database code (functions, views, stored procedures, etc) live in `app_code` schema
- **Planned:** `reviews` table partitioning

### Frontend (SvelteKit + Tailwind)
- **Framework:** SvelteKit 2 + Svelte 5 (Runes mode: `$state`, `$derived`). BFF pattern.
- **Styling:** TailwindCSS 4+ (Light/Dark mode support).
- **Rich Text (planned):** TipTap. Math: KaTeX + MathLive.
- **Routing:** File-based. `hooks.server.ts` validates sessions via `/v1/auth/me`.
- **LayerChart:** for charting (https://next.layerchart.com/). Ensure Svelte 5 version is used.
- **Shadcn Svelte:** for UI components (https://shadcn-svelte.com/docs).
- **Forms:** sveltekit-superforms + Zod for validation.
- **Svelte Documentation**: 
  - `references/svelte-llms-full.txt` : full examples with all details
  - `references/svelte-llms-small.txt` : examples and non-essential details removed

#### API Strategy
Go API hidden from browser. All requests flow through SvelteKit server.
- **Reads:** `+page.server.ts` load functions with `apiRequest()`.
- **Mutations (default):** Form actions with `use:enhance`. No client JS needed.
- **Mutations (edge cases):** `/api/[...path]` proxy routes for modals, drag-drop, real-time.
- **Client helper (planned):** `api()` in `$lib/api/client.ts` for proxy route calls.

#### Component Patterns
- Use Svelte 5 runes (`$state`, `$derived`, `$effect`) not legacy `let` reactivity.
- Props via `let { prop = default }: Props = $props()`.
- Prefer `$derived` over `$effect` for computed values.
- Colocate component-specific types in the same file or `types.ts` sibling.

#### Styling Guidelines
- Use Tailwind utilities; avoid custom CSS unless necessary.
- Use consistent naming conventions, colors, spacing, and structure
- Flexbox for 1D layouts, Grid for 2D layouts.
- Use shadcn components before building custom ones.
- Dark mode: use `dark:` variants, colors from shadcn theme.
- **Design tokens over hardcoded colors:** Always use semantic color tokens (`bg-card`, `text-foreground`, `border-border`, `bg-primary`, `text-muted-foreground`, etc.) defined in `app.css`. Only use raw Tailwind color families (`slate-*`, `sky-*`, `red-*`, etc.) in components when absolutely necessary, as they bypass theming and break dark mode. The only exceptions are shadcn upstream components and intentionally forced-color surfaces (e.g., focus-mode overlay).

#### File Structure
```
src/
├── lib/
│   ├── components/    # Shared components
│   ├── api/           # API helpers
│   └── stores/        # Global state (if needed)
├── routes/
│   └── (app)/         # Authenticated routes
│       └── notebooks/
│           └── [id]/
│               ├── +page.svelte
│               ├── +page.server.ts
│               └── +layout.server.ts
```

#### Anti-patterns to Avoid
- Don't use `onMount` for data fetching; use `+page.server.ts` load functions.
- Don't expose Go API URLs to browser; always proxy through SvelteKit.
- Don't mix Svelte 4 (`$:`) and Svelte 5 (`$derived`) reactive syntax.


## Project Structure
```text
apexmemory/
├── backend/
│   ├── cmd/api/            # Entry point (main.go)
│   ├── internal/
│   │   ├── app/            # HTTP handlers, middleware, routes
│   │   ├── db/             # sqlc generated code (DO NOT EDIT)
│   │   └── testutil/       # Test helpers
│   ├── db/
│   │   ├── init/           # Docker initdb scripts
│   │   ├── migrations/     # Tern migration files
│   │   ├── queries/        # SQL input for sqlc
│   │   ├── backups/        # Database dumps (gitignored)
│   │   └── schema-dump.sql # Current schema snapshot
│   ├── sqlc.yml
│   ├── tern.conf
│   └── go.mod
├── frontend/               # SvelteKit frontend
├── references/             # Architecture docs
├── Makefile                # Task runner
├── docker-compose.yml      # PostgreSQL 18
└── .env                    # Environment variables (gitignored)
```

## Development Workflow

### Setup
```bash
make docker.up     # Start PostgreSQL
make tern.migrate  # Apply migrations
make dev.backend   # Start API on :4000 (logs: /tmp/apexmemory-api.log)
```

### Key Commands
Refer to the Makefile for most dev commands
**Database:**
- `make tern.new name=x` - Create migration
- `make tern.migrate` - Apply migrations (also dumps schema)
- `make tern.rollback` - Roll back one migration
- `make tern.status` - Show migration status
- `make db.sqlc` - Generate Go code from SQL
- `make db.psql.claude SQL="..."` - Run a SQL query non-interactively (use this to inspect the database)

**Development:**
- `make dev.up` - Start backend + frontend
- `make dev.backend` - Restart backend only
- `make dev.mobile` - Expose to network for mobile testing
- `make dev.check` - Run go vet + svelte check

**Testing:**
- `make test.backend` - Run all backend tests
- `make test.backend.v` - Verbose test output
- `make test.backend.cover` - With coverage

**Seeding (planned):**
- `make seed.notebooks` - Seed example notebooks
- `make seed.reviews` - Generate review history

### Backend Dev Cycle
1. Write migration (if schema changes needed)
2. `make tern.migrate` to apply
3. Write queries in `db/queries/`
4. `make db.sqlc` to generate Go code
5. Write handlers in `internal/app/`
6. Write tests (`*_test.go`)
7. `make test.backend` to verify
8. `make dev.backend` to restart


### Logging

Optimized for human debugging **and** AI agents to diagnose and root-cause issues automatically.

#### Principles
- **Structured, not prose:** Stable field names; machine-parseable values.
- **Actionable errors:** Include `error_code`, message, and `remediation_hint`.
- **Log the "why":** Record intent + inputs at boundaries, not just outcomes.
- **Safe by default:** Never log secrets/PII; mask sensitive values.

#### Standard Pattern (Go)
```go
// Context fields (set in middleware)
slog.With("service", "api", "version", version.Commit, "env", cfg.Env, "request_id", reqID)

// Errors: always include error_code, retryable, remediation_hint
slog.Error("failed to create notebook",
    "error", err,
    "error_code", "NOTEBOOK_CREATE_FAILED",
    "user_id", userID,
    "retryable", false,
    "remediation_hint", "Check database connectivity",
)

// Operations: include IDs and duration
slog.Info("notebook created", "user_id", userID, "notebook_id", id, "duration_ms", elapsed)
```

#### Log Levels
| Level | Use for |
|-------|---------|
| `Debug` | Dev diagnostics (SQL, cache hits) |
| `Info` | Business events ("notebook created") |
| `Warn` | Recoverable issues (retries, slow queries) |
| `Error` | Failures requiring attention |

#### Correlation
- Generate `request_id` (UUIDv7) in middleware; propagate via context; include in all logs and error responses.

#### Never Log
Passwords, tokens, full request bodies, unmasked PII.

## Agent Testing

### Test Account
A dedicated test account for AI agents to perform interactive testing:

| Field | Value |
|-------|-------|
| **Email** | `test-agent@apexmemory.ai` |
| **Password** | `Always$InTheBananaStand!123` |
| **Username** | `test-agent` |
| **User ID** | `019501a0-0000-7000-8000-000000000001` |

### Setup
```bash
make seed.user      # Create test account only
make seed.all       # Create test account + sample notebooks
```

Run after: `docker.up && tern.migrate`, `db.truncate`, or `tern.reset`.

### Seed CLI Reference
```bash
# Direct CLI usage
cd backend && go run ./cmd/seed <command> [flags]

# Commands:
#   user       Create/verify test agent account
#   notebooks  Seed notebooks for a user
#   reviews    Generate review history (placeholder)
#   all        Run user + notebooks for test agent
#   help       Show help
```

### UI Testing (Browser)

**Clear any auto filled credentials**

1. Navigate to `/login`
2. Email: `test-agent@apexmemory.ai`
3. Password: `Always$InTheBananaStand!123`
4. Submit - session cookie set automatically

### API Testing (Go backend, port 4000)
```bash
# Login
curl -X POST http://localhost:4000/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test-agent@apexmemory.ai","password":"Always$InTheBananaStand!123"}'

# Use token from response
curl http://localhost:4000/v1/notebooks \
  -H "Authorization: Bearer <session_token>"
```


## Agent Behavior & Persona

**Default Role:** Senior Engineering Architect
**Default Tone:** Direct, Technical, Concise.

**STRICT GUIDELINES:**

1.  **NO SYCOPHANCY:** Be friendly, but do not use phrases like "Great idea," "You're right," or "Excellent point."
2.  **CRITICAL ANALYSIS:** If an approach has trade-offs, state them immediately.
3.  **CODE FIRST:** Focus on correctness, performance, simplicity, and maintainability. Consider wider architectural implications.
4.  **NO EMOJIS:** Do not add emojis to code, logs, or commands.
