<!--
## Sync Impact Report

**Version Change**: [PROJECT_NAME] Constitution (template) → 1.0.0

**Modified Principles**: N/A (initial constitution from template)

**Added Sections**:
- Core Principles (7 principles: Clean Architecture, Test-First Discipline,
  Domain-Driven Validation, Consistent API Contract, Infrastructure Abstraction,
  Simplicity & YAGNI, Performance & Observability)
- Technology Standards
- Development Workflow
- Governance

**Removed Sections**: All placeholder template tokens replaced.

**Templates Requiring Updates**:
- ✅ `.specify/templates/plan-template.md` — Constitution Check gates align with principles
- ✅ `.specify/templates/spec-template.md` — Acceptance Scenarios align with given/when/then
- ✅ `.specify/templates/tasks-template.md` — Task categories reflect Clean Architecture layers
- ✅ `.specify/templates/constitution-template.md` — Source template unchanged (correct)

**Follow-up TODOs**: None — all fields resolved.
-->

# Mystery Gifter API Constitution

## Core Principles

### I. Clean Architecture (NON-NEGOTIABLE)

The codebase MUST maintain strict three-layer separation: **domain**, **application**,
and **infrastructure**. Dependency direction flows inward only — infrastructure depends
on application, application depends on domain, domain depends on nothing external.

- `internal/domain/` MUST contain only entities, interfaces, domain errors, and validation.
  No infrastructure imports permitted.
- `internal/application/` MUST orchestrate domain logic via injected interfaces only.
  No direct database, HTTP, or framework calls.
- `internal/infra/` MUST implement domain/application interfaces.
  Business logic MUST NOT appear in controllers or repositories.
- Every public service MUST follow: public interface (`XService`), private struct (`xService`),
  constructor returning the interface (`NewXService(...) XService`).
- Every repository MUST follow: private struct (`xRepository`), constructor returning
  `domain.XRepository` (`NewXRepository(db DB) domain.XRepository`).

**Rationale**: Violation of layer boundaries creates untestable code and hard coupling to
infrastructure choices. This rule was established from project inception and is never negotiable.

### II. Test-First Discipline

All business logic MUST have unit tests written in the same implementation cycle.
Tests MUST be run and pass before a feature is considered complete.

- Test function naming: `Test_<Type>_<Method>`
- Subtests MUST use `t.Run("should ... when ...", ...)` with descriptions in English
- Every test MUST have three explicit phases annotated with comments: `// given`, `// when`, `// then`
- Each subtest MUST create its own `gomock.Controller`: `mockCtrl := gomock.NewController(t)`
- Test data MUST be constructed via builders in `build_domain/`, `build_postgres/`,
  or `build_rest/` — never inline struct literals for complex objects
- All assertions MUST use `testify/assert`; both error AND result MUST always be verified
- `context.Background()` MUST be used for contexts in tests
- **Mandatory**: run tests after every implementation; identify and fix failures before proceeding

**Rationale**: Tests written after the fact are incomplete and miss edge cases. The
given/when/then pattern ensures traceability between test scenarios and acceptance criteria.

### III. Domain-Driven Validation

Entities and value objects MUST validate themselves. Validation logic belongs in the domain
layer, not in controllers or services.

- Every entity MUST implement a `Validate() error` method using `go-playground/validator` tags
- Factory functions (`NewX(...)`) MUST call `Validate()` before returning
- Mutating methods (`AddUser`, `Archive`, etc.) MUST call `Validate()` after mutation
- Domain errors MUST use the custom error types defined in `internal/domain/errors.go`:
  `ValidationError`, `ConflictError`, `ResourceNotFoundError`, `UnauthorizedError`, `ForbiddenError`
- Each error type MUST carry the correct HTTP status code as part of its definition
- Controllers MUST NOT define business error messages — they propagate domain errors directly

**Rationale**: Centralising validation in the domain ensures consistency across all entry
points (REST, CLI, queue consumers) and makes domain invariants explicit and testable.

### IV. Consistent API Contract

All REST controllers MUST follow a strict, uniform request-response flow to ensure
predictable client behavior.

- Handler signature: `func (c *XController) Method(ctx *fiber.Ctx) error`
- Request flow: `BodyParser` → `dto.Validate()` → `mapXToDomain` → service call →
  `mapXFromDomain` → response
- `BodyParser` failure MUST return `fiber.NewError(fiber.StatusUnprocessableEntity)`
- All other errors MUST be returned directly (error handler maps domain errors to HTTP status)
- Resource creation MUST return `ctx.Status(fiber.StatusCreated).JSON(...)`
- All other successful responses MUST return `ctx.JSON(...)`
- Route parameters MUST use `ctx.Params("paramName")`, never `ctx.Query` for resource IDs
- Auth user ID MUST be extracted via `c.AuthTokenManager.GetAuthUserID(ctx.Locals("user"))`
- Swagger annotations MUST be maintained for all endpoints; `make generate-docs` MUST pass
- Query parameters in Swagger MUST be defined individually (never `schema: "$ref"`)

**Rationale**: Uniform controller flow reduces cognitive load and prevents error-handling
inconsistencies that leak internal details to clients.

### V. Infrastructure Abstraction

All external dependencies MUST be accessed through interfaces, never through concrete types.

- Repositories MUST inject the `DB` interface, never `*sqlx.DB` directly
- All SQL queries MUST use `squirrel` with `PlaceholderFormat(squirrel.Dollar)`
- Write operations MUST use `ExecContext`, single-row reads `GetContext`,
  multi-row reads `SelectContext`
- `sql.ErrNoRows` MUST be mapped to `domain.NewResourceNotFoundError`
- PostgreSQL unique violation (`pq.Error`) MUST be mapped to `domain.NewConflictError`
- All other errors MUST be wrapped: `fmt.Errorf("context message: %w", err)`
- Multi-step writes MUST use transactions: `BeginTxx` → `defer tx.Rollback()` →
  operations → `tx.Commit()`
- Mocks MUST be generated with `go.uber.org/mock/mockgen` via `//go:generate` directives;
  interface changes require `go generate ./...` to be re-run

**Rationale**: Interface-driven design enables isolated unit tests and decouples business
logic from specific infrastructure choices (database engine, storage backend).

### VI. Simplicity & YAGNI

Every abstraction MUST earn its place by solving a present problem, not a hypothetical
future one. Complexity requires explicit justification.

- No helpers, utilities, or abstractions for one-time operations
- No speculative features, extra configuration, or backwards-compatibility shims
- No docstrings, comments, or type annotations added to code that was not changed
- Error handling MUST NOT be added for scenarios that cannot happen given current invariants
- Mapper functions MUST be private and co-located with the type they serve
- Slices MUST be pre-allocated with `make([]T, 0, len(src))` before iteration
- Three similar lines of code is preferable to a premature abstraction

**Rationale**: Premature abstractions increase maintenance cost and obscure intent. The
right amount of complexity is exactly what the task requires.

### VII. Performance & Observability

The API MUST remain responsive under expected load. Errors MUST be observable without
exposing internal details to clients.

- All repository methods MUST pass `context.Context` to enable timeout propagation
- `log.Println` MUST be called for significant infrastructure errors before returning them
- Pagination MUST be applied to all list/search endpoints via `Limit` + `Offset` filters
- Default pagination values MUST be defined as constants in the domain layer
  (e.g., `DefaultGroupLimit = 15`)
- Sort direction and sort field MUST be validated via `go-playground/validator` `oneof` tags
- HTTP response times MUST remain suitable for interactive use (target p95 < 200ms for
  standard CRUD; no explicit SLO enforcement mechanism required at this stage)

**Rationale**: Context propagation and structured error logging are the minimum viable
observability baseline. Pagination prevents unbounded queries from degrading the database.

## Technology Standards

**Language**: Go (latest stable release)
**Web Framework**: Fiber v2 — use only idiomatic Fiber patterns
**Database**: PostgreSQL via `sqlx` + `squirrel` query builder
**Migrations**: `golang-migrate` — migration files MUST be committed alongside schema changes
**Authentication**: `golang-jwt` — JWT tokens; secret key and session duration via env vars
**Validation**: `go-playground/validator` — struct tags are the canonical validation source
**Testing**: `testify/assert` for assertions, `go.uber.org/mock/mockgen` for mocks
**Documentation**: `go-swagger` — `make generate-docs` MUST succeed at all times
**Identity**: UUID v4 via `internal/infra/outgoing/identity` (injected as `IdentityGenerator`)
**Configuration**: `caarlos0/env` — all config via environment variables; `.env.example` MUST
be kept up to date

Required environment variables: `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`,
`DB_PASSWORD`, `AUTH_SECRET_KEY`, `AUTH_SESSION_DURATION`.

## Development Workflow

- For files over 100 lines or complex changes: outline the plan and confirm approach before editing
- Alternative perspectives MUST be raised when there is opportunity to improve quality
- `context.Context` MUST be the first parameter of every service and repository method
- `error` MUST be the last return value of every method that can fail
- After every interface change, run `go generate ./...` to regenerate mocks
- After every DTO or route change, run `make generate-docs` to update Swagger spec
- All unit tests MUST pass before marking any task as complete (`make test`)
- Commit messages MUST be descriptive and reference the layer changed
  (e.g., `feat(domain): add group archiving logic`)

**Quality Gates** (all MUST pass before a feature is considered done):

1. `make build` — compiles without errors
2. `make test` — all unit tests pass
3. `make generate-docs` — Swagger spec generates without errors
4. No unexplained bracket tokens remain in any spec or plan document

## Governance

This constitution supersedes all other project practices and coding guidelines. Any rule
in CLAUDE.md or `.claude/rules/` that conflicts with this constitution MUST be reconciled
in favor of this document, and the conflicting rule updated accordingly.

**Amendment procedure**:
1. Propose change with rationale in a pull request description
2. Update `CONSTITUTION_VERSION` according to semantic versioning:
   - MAJOR: principle removal, redefinition, or backward-incompatible governance change
   - MINOR: new principle or section added, or materially expanded guidance
   - PATCH: clarification, wording fix, or non-semantic refinement
3. Update `LAST_AMENDED_DATE` to the amendment date (ISO format YYYY-MM-DD)
4. Run the consistency propagation checklist against all `.specify/templates/` files
5. Document impact in the Sync Impact Report comment at the top of this file

**Compliance review**: Every feature plan MUST include a "Constitution Check" gate
confirming that the proposed design does not violate any principle. Violations MUST be
justified in a Complexity Tracking table in the plan document.

Runtime development guidance: see `CLAUDE.md` and `.claude/rules/` for language-specific
conventions that complement these principles.

**Version**: 1.0.0 | **Ratified**: 2026-03-28 | **Last Amended**: 2026-03-28
