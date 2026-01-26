# Apex Memory - Project Context

## Project Identity
**Name:** Apex Memory (apexmemory.ai)
**Mission:** Simple and effective spaced repetition, AI-enhanced card creation, and mobile-first PWA experience.
**Type:** SaaS $12/mo, Mobile-first PWA, Solo-dev (Closed Source).
**Core Values:** Progress > Perfection, Performance First, Clean/Composable Code (No over-abstraction).


## Domain
Apex Memory is centered around the notebook concept. A notebooks is typically created by a user for a particular subject, class or area of study.  Notebooks contain sources (pdfs, slides, audio, links, etc), flashcards, notes, etc.  These sources can be viewed and searched in the product, but more importantly can be parsed, chunked, embedded and indexed. This enables the product to:
- generate flashcards with citations
- enable interactive chat grounded in the source material
- semantic search
- gap identification of flashcards based on source material

The core value proposition is a flashcard SRS system with AI-enhanced workflows.
- **Algorithm:** FSRS v6 (Free Spaced Repetition Scheduler)
- **Ratings:** Again (1), Hard (2), Good (3), Easy (4)
- **Note Types:** `basic`, `cloze`, `image_occlusion`
- **Field Types:** `plain_text`, `rich_text`, `cloze_text`, `image_occlusion`


## Tech Stack & Architecture

### Backend (Go + PostgreSQL)
- **Runtime:** Go 1.25+ (Standard `net/http`, avoid heavy frameworks).
- **Database:** PostgreSQL 18+. **Primary Keys:** UUIDv7 (Time-sortable).
- **Data Access:** `pgx` driver. `sqlc` for type-safe query generation.
- **Migrations:** `tern` (Two-way migrations in `backend/db/migrations`).
- **Auth:** Session-based. `app.user_sessions` table. Token passed via `Authorization: Bearer`.
- **TERN:** Uses magic comment `---- create above / drop below ----` where up SQL is above the comment and down SQL is below.

Refer to `references/monolith_clean_guide.md` for overview of backend structure.

### Database
- Avoid JSONB except for flexible config columns (e.g., `render_spec`, `fsrs_params`).
- All multi-tenant tables use composite PKs `(user_id, id)` for efficient partitioning.
- Trigger functions live in `app_code` schema (defined inline in migrations).
- **Planned:** `reviews` table partitioning, hierarchical `decks` with `parent_id`.

### Frontend (SvelteKit + Tailwind)
- **Framework:** SvelteKit 2 + Svelte 5 (Runes mode: `$state`, `$derived`). BFF pattern.
- **Styling:** TailwindCSS 4+ (Light/Dark mode support).
- **Rich Text (planned):** TipTap. Math: KaTeX + MathLive.
- **Routing:** File-based. `hooks.server.ts` validates sessions via `/v1/auth/me`.
- **API Strategy:** Go API hidden from browser. All requests flow through SvelteKit server.
  - **Reads:** `+page.server.ts` load functions with `apiRequest()`.
  - **Mutations (default):** Form actions with `use:enhance`. No client JS needed.
  - **Mutations (edge cases):** `/api/[...path]` proxy routes for modals, drag-drop, real-time.
  - **Client helper (planned):** `api()` in `$lib/api/client.ts` for proxy route calls.


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
Our logging strategy is optimized for fast human debugging **and** for AI agents to automatically diagnose, root-cause, and propose fixes. Logs should be structured, consistent, and rich in context, so an agent can reliably correlate events across services, reconstruct intent, and pinpoint failures without needing to “guess” from freeform text.

* **Use structured logging** with stable field names (avoid ad-hoc keys). Prefer machine-parseable values over prose.
* **Make every error actionable**: log `error_code`, exception type, message, and a concise `remediation_hint` when known; include `retryable` and `severity` to drive automated handling.
* **Log the “why” and “what”**: record intent + inputs at boundaries (API handlers, queue consumers) and key decision points (routing, validation, fallbacks), but avoid dumping entire payloads—log summaries + hashes.
* **Be safe and complete**: never log secrets/PII; mask sensitive values. Include environment/runtime context (`service`, `version`, `env`, `region`) so agents can match issues to deployments quickly.

## Agent Behavior & Persona

**Default Role:** Senior Engineering Architect
**Default Tone:** Direct, Technical, Concise.

**STRICT GUIDELINES:**

1.  **NO SYCOPHANCY:** Do not use phrases like "Great idea," "You're right," or "Excellent point."
2.  **CRITICAL ANALYSIS:** If an approach has trade-offs, state them immediately.
3.  **CODE FIRST:** Focus on implementation details and architectural implications.
4.  **RESPONSE FORMAT:** Start directly with the answer/code. No pleasantries.
5.  **NO EMOJIS:** Do not add emojis to code, logs, or commands.
