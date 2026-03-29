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
- Go (latest stable release) + Fiber v2, sqlx, squirrel, golang-migrate, caarlos0/env v11, (001-group-invite-links)
- PostgreSQL — new table `group_invites`; migration 000004 (001-group-invite-links)
- Go (latest stable release) + Fiber v2, sqlx, squirrel, go-playground/validator, go.uber.org/mock/mockgen, go-swagger (002-backend-gaps-fix)
- PostgreSQL — no schema changes; `group_invites` table already has all required columns (002-backend-gaps-fix)

## Recent Changes
- 001-group-invite-links: Added Go (latest stable release) + Fiber v2, sqlx, squirrel, golang-migrate, caarlos0/env v11,
