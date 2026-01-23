# Monolith‑but‑Clean Guide (for AI agents)

This guide describes a **monolithic** Go backend style that stays clean and refactor‑friendly. The goal is to keep everything in one package (the “god object” Application), but structured so features are easy to move out later.

## Goals

- Keep a single `Application` object that owns dependencies and runtime state.
- Avoid fragmentation across many packages while preserving clean seams.
- Keep business logic testable without HTTP or global state.
- Make future extraction into separate packages easy.

## Core Principles

1. **One package, many files**: keep logic in `internal/app` but split by feature or concern.
2. **Application owns all state**: no global vars; inject dependencies via `Application` fields.
3. **Handlers are thin**: handlers only parse/validate input and call helper/business methods.
4. **Domain logic is HTTP‑agnostic**: helpers do not import `net/http`.
5. **Clear dependency direction**: handler → validation → domain logic → DB.
6. **Use small interfaces for seams**: DB, mailer, clock, token generator, etc.
7. **Routes are wiring only**: no logic or data access in `Routes()`.
8. **Errors are centralized**: sentinel errors and error helpers live in `errors.go`.

## Recommended Layout (single package)

```
cmd/api/main.go              # minimal bootstrap
internal/app/
  app.go                     # Application struct + lifecycle
  config.go                  # config parsing/validation
  routes.go                  # HTTP routing + middleware composition
  middleware.go              # auth, logging, recovery, rate limit
  auth.go                    # auth logic (non-HTTP)
  auth_handlers.go           # HTTP handlers for auth
  users.go                   # user logic (non-HTTP)
  users_handlers.go          # HTTP handlers for users
  sessions.go                # session logic
  response.go                # RespondJSON / RespondError helpers
  errors.go                  # sentinel errors + error wrappers
  logging.go                 # request logging helpers
  background.go              # background jobs + cleanup
  db.go                      # DB connection helpers
  types.go                   # request/response structs + shared types
  validate.go                # input validation helpers
```

## Application Struct (god object, but clean)

- All dependencies live on `Application`.
- Avoid implicit globals.
- Interfaces live in the same package to keep imports simple.

Example fields:

- `Logger` (structured logger)
- `DB` (pgx pool)
- `Queries` (sqlc or query layer)
- `Clock` (interface for time, makes tests deterministic)
- `Mailer` (interface)
- `TokenMaker` (interface)
- `RateLimiters` (stateful, owned here)

## Handler vs Logic Split

Keep logic in methods that **don’t import net/http**. Handlers should only:

- parse input (JSON/form)
- validate request
- call logic
- translate result to HTTP response

Example split:

- `func (app *Application) RegisterUser(ctx context.Context, req RegisterRequest) (*AppUser, error)`
- `func (app *Application) RegisterHandler(w http.ResponseWriter, r *http.Request)`

## Interfaces for Seams (even inside one package)

Define small interfaces to make extraction painless later:

- `type Mailer interface { Send(...)} `
- `type Clock interface { Now() time.Time }`
- `type TokenMaker interface { New() (string, error) }`

You can keep the concrete implementations in the same package for now.

## Routing Conventions

- `Routes()` wires handlers with middleware.
- Handlers are methods on `Application`.
- Middleware uses `Application` dependencies, but does not do domain work.

## Error Handling

- Define sentinel errors in `errors.go` (`ErrInvalidCredentials`, `ErrUnauthorized`, etc.).
- Convert domain errors to HTTP errors in a single helper (`RespondError`).

## Validation

- Keep input validation in `validate.go` or close to the handler.
- Avoid direct DB reads just to validate (unless required).

## Testing Strategy

- Test core logic methods directly (no HTTP) whenever possible.
- Handlers: only test routing and status codes for critical flows.
- Use interfaces to stub time, mailer, and token generation.

## Refactor‑Later Path (when the monolith grows)

When a feature grows, extract by **moving the logic first**, not the handlers:

1. Create an interface in `internal/app` (e.g., `AuthService`).
2. Move implementation to `internal/auth` package.
3. Handlers stay in `internal/app` until you’re ready to move them too.
4. `Application` depends on the interface, not the concrete type.

This makes extraction incremental and low‑risk.

## Agent Checklist (when adding new features)

- Add new logic in `internal/app/<feature>.go`.
- Add new handlers in `internal/app/<feature>_handlers.go`.
- Keep handler → logic separation.
- No globals. Use `Application` fields.
- Add tests for logic first; handler tests only if needed.
- Avoid importing `net/http` in logic files.
- Update `Routes()` only for wiring.

## “Don’t” List

- Don’t put DB calls directly in handlers.
- Don’t let helpers depend on `http.Request` or `http.ResponseWriter`.
- Don’t create new packages just for one file.
- Don’t introduce global singletons.
- Don’t mix two abstractions (e.g., direct SQL + repository) in the same feature.

