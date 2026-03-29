# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Language

Always respond in **Brazilian Portuguese (pt-BR)**.

## Commands

```bash
make build           # Compile to bin/mystery-gifter-api
make run             # Run locally
make test            # Run all unit tests
make generate-docs   # Generate Swagger/OpenAPI docs
make serve-docs      # Serve Swagger UI at http://localhost:8081
make install-tools   # Install swagger and mockgen
make clean           # Remove generated files

go test ./...                        # All tests
go test -cover ./...                 # With coverage
go test ./internal/application/...  # Specific package
```

## Architecture

Clean Architecture with three main layers:

- **`internal/domain/`** — Entities, interfaces (Repository, IdentityGenerator, PasswordManager, AuthTokenManager), domain errors, and validation logic. No external dependencies.
- **`internal/application/`** — Service interfaces and implementations. Orchestrates domain logic via injected repositories and managers. No infrastructure code.
- **`internal/infra/`** — Infrastructure:
  - `entrypoint/rest/` — Fiber controllers, DTOs, and mapper functions
  - `entrypoint/routes.go` — Route registration
  - `outgoing/postgres/` — Repository implementations using `sqlx` + `squirrel`
  - `outgoing/security/` — BCrypt password manager and JWT auth token manager
  - `outgoing/identity/` — UUID generator
  - `config/` — Environment variable loading via `caarlos0/env`

**Entry point:** `cmd/api/main.go` → `internal/infra/runner.go` (wires all dependencies and starts Fiber on port 8080)

**Key stack:** Fiber v2, PostgreSQL, sqlx, squirrel (query builder), golang-migrate, golang-jwt, go-playground/validator, gomock (uber), go-swagger.

## Coding Conventions

### Services (`internal/application/`)
- Define a public interface `XService` with a `//go:generate mockgen` directive
- Implement with unexported struct `xService`
- Constructor: `NewXService(...) XService`
- All methods take `context.Context` as the first parameter
- Call `entity.Validate()` before any business logic
- Mocks go in `mock_application/`

### Controllers (`internal/infra/entrypoint/rest/`)
- Named `XController` with constructor `NewXController`
- Handler signature: `func (c *XController) Method(ctx *fiber.Ctx) error`
- Parse body → validate DTO → map to domain → call service → map to response DTO → return JSON
- Return `fiber.NewError(fiber.StatusUnprocessableEntity)` on parse failure, propagate other errors directly
- No business logic in controllers

### Repositories (`internal/infra/outgoing/postgres/`)
- Unexported struct `xRepository`, constructor `NewXRepository(db DB) domain.XRepository`
- Use `squirrel` with `PlaceholderFormat(squirrel.Dollar)` for all queries
- `ExecContext` for writes, `GetContext` for single row, `SelectContext` for multiple rows
- Map `sql.ErrNoRows` → `domain.NewResourceNotFoundError`
- Map `pq.Error` unique violation → `domain.NewConflictError`
- Use transactions (`BeginTxx`/`defer tx.Rollback()`/`Commit`) for multi-step operations
- Mocks go in `mock_postgres/`

### Mapper functions
- `mapXTo<Dest>` or `mapXFrom<Origin>` naming
- Keep them private and in the same file as the type they serve

### Tests
- Pattern: `Test_<Type>_<Method>` with `t.Run("should ... when ...", ...)`
- Three phases annotated with comments: `// given`, `// when`, `// then`
- One `gomock.Controller` per subtest for isolation
- Use builders from `build_domain/`, `build_postgres/`, `build_rest/` subdirectories
- `testify/assert` for all assertions; always check both the error and the result

## Mock Generation

Mocks are generated with `go.uber.org/mock/mockgen`. To regenerate after interface changes:

```bash
go generate ./...
```

Each interface file has a `//go:generate` directive pointing to its mock destination.

## Environment Variables

Required: `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD`, `AUTH_SECRET_KEY`, `AUTH_SESSION_DURATION`. Copy `.env.example` to `.env` to get started.
