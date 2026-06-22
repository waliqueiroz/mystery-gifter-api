# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Language

There is a strict language separation between human language and machine language:

**Brazilian Portuguese (pt-BR) — mandatory for**:
- All responses and explanations
- Code comments (the "why", not the "what")
- Documentation: Swagger `summary`/`description` annotations, README files, speckit artifacts (`spec.md`, `plan.md`, `tasks.md`, checklists)
- Commit messages and pull request descriptions

**English — mandatory for**:
- All source code: function names, variables, types, structs, constants, packages
- Error message strings (`errors.New`, `fmt.Errorf`, `fiber.NewError`, etc.)
- Subtest descriptions (`t.Run("should ... when ...", ...)`)
- JSON/query field names, HTTP headers, validation tags

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

## Working Protocol
- For files >100 lines or complex changes: outline the plan before editing
- Challenge decisions and bring alternative perspectives when there's opportunity to improve quality

## Mock Generation

Mocks are generated with `go.uber.org/mock/mockgen`. To regenerate after interface changes:

```bash
go generate ./...
```

Each interface file has a `//go:generate` directive pointing to its mock destination.

## Environment Variables

Required: `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD`, `AUTH_SECRET_KEY`, `AUTH_SESSION_DURATION`. Copy `.env.example` to `.env` to get started.

## Active Technologies
- Go (latest stable release) + Fiber v3, sqlx, squirrel, golang-migrate, caarlos0/env v11, go-playground/validator, go.uber.org/mock/mockgen, go-swagger
- PostgreSQL — tables: `users`, `groups`, `group_users`, `group_matches`, `group_invites` (migration 000004)
- Go 1.26.4 + Fiber v3 (v3.3.0), gofiber/contrib/v3/jwt (v1.1.6), fiber/v3/extractors (transitivo), golang-jwt/jwt v5.3.1 (004-jwt-cookie-auth)
- PostgreSQL — sem migrações (feature é puramente HTTP) (004-jwt-cookie-auth)

## Recent Changes
- 002-backend-gaps-fix: description optional on group creation; membership check on GET /groups/:groupID (403 for non-members); matches removed from GroupDTO; GET /groups/:groupID/invites/active endpoint added; Swagger ReopenGroup description corrected
- 003-fiber-v3-upgrade: Go 1.25.1 → 1.26.4; Fiber v2 → v3; handler signatures *fiber.Ctx → fiber.Ctx; BodyParser/QueryParser → Bind().Body()/Bind().Query(); ctx.Locals("user") → jwtware.FromContext(ctx)
