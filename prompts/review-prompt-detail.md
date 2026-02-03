<review_request>
Perform a thorough review of all changes made in this session. Apply Apex Memory project standards and best practices throughout.

<project_context>
- **Backend**: Go 1.25+ with net/http, PostgreSQL 18+, pgx driver, sqlc for type-safe queries, tern migrations
- **Frontend**: SvelteKit 2 + Svelte 5 (runes: $state, $derived, $effect), TailwindCSS 4, shadcn-svelte
- **Architecture**: BFF pattern (Go API hidden from browser), composite PKs (user_id, id), session-based auth
- **Domain**: FSRS v6 spaced repetition, notebook-centric model, fact→card relationships
</project_context>

<review_dimensions>

<correctness priority="critical">
- Verify FSRS scheduling logic matches v6 spec (stability, difficulty, retrievability calculations)
- Check SQL queries for correctness with composite PKs and proper user_id scoping
- Validate Svelte 5 reactivity ($state mutations, $derived dependencies, $effect cleanup)
- Test edge cases: empty states, concurrent updates, timezone handling, null/undefined
- Ensure form validation matches backend constraints
</correctness>

<simplicity priority="high">
- Prefer Go standard library over external dependencies
- Use sqlc-generated code rather than manual query building
- Avoid premature abstraction—extract only when pattern repeats 3+ times
- Question any solution requiring more than ~50 lines; can it be simpler?
- Leverage PostgreSQL features (JSONB, CTEs, window functions) over application logic
</simplicity>

<consistency priority="high">
- SQL: snake_case, app schema, composite PKs, TEXT+CHECK over ENUMs for extensible types
- Go: standard net/http patterns, error wrapping with context, structured logging
- Svelte: runes mode only, form actions with use:enhance by default, +page.server.ts for loads
- API: RESTful endpoints under /v1/, consistent error response shapes
- Match existing code style in the file/module being modified
</consistency>

<maintainability priority="medium">
- Functions should do one thing; name should describe that thing precisely
- Prefer explicit over clever—future agents need to understand the code
- Keep related code together; separate concerns at natural boundaries
- Use TypeScript types/Go structs that document the domain model
- Structure for testability even if tests aren't written yet
</maintainability>

<documentation priority="medium">
- Add comments explaining *intent* and *why*, not *what* the code does
- Document non-obvious decisions: "Chose X over Y because..."
- Note tradeoffs considered, especially for performance vs. simplicity choices
- Mark TODOs with context: what needs doing and why it's deferred
- Update CLAUDE.md or relevant docs if patterns/conventions change
</documentation>

</review_dimensions>

<output_instructions>
1. List each issue found with file:line reference
2. Categorize by dimension (correctness/simplicity/consistency/maintainability/documentation)
3. For critical issues: fix directly
4. For judgment calls: explain tradeoffs, recommend approach, ask if uncertain
5. Summarize changes made at the end
</output_instructions>
</review_request>