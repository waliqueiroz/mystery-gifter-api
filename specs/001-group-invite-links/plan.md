# Implementation Plan: Group Invite Links

**Branch**: `001-group-invite-links` | **Date**: 2026-03-29 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-group-invite-links/spec.md`

## Summary

Add a group-invite system to Mystery Gifter API so that group owners can generate
time-limited, multi-use invite IDs that allow authenticated users to join a group without
requiring the owner to add them manually. As a security hardening step, self-join via
`AddUser` will be removed — the invite becomes the sole self-service path for third parties.
Invite lifetime is configurable via the `INVITE_LINK_EXPIRATION` env var (Go duration
string, defaulting to `"24h"`).

## Technical Context

**Language/Version**: Go (latest stable release)
**Primary Dependencies**: Fiber v2, sqlx, squirrel, golang-migrate, caarlos0/env v11,
go-playground/validator, go.uber.org/mock/mockgen, testify/assert
**Storage**: PostgreSQL — new table `group_invites`; migration 000004
**Testing**: `go test ./...` + `make test`; gomock for unit mocks; testify/assert
**Target Platform**: Linux server (Docker / bare metal)
**Project Type**: REST web service (Clean Architecture)
**Performance Goals**: Same p95 < 200ms target as existing CRUD endpoints
**Constraints**: Invite ID validated server-side; expiration stored as absolute `TIMESTAMPTZ`;
no JWT-style self-validating tokens (requires DB lookup)
**Scale/Scope**: Same scale assumptions as the rest of the API

## Constitution Check

*GATE: Must pass before implementation. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Clean Architecture | ✅ Pass | New entity in `domain/`, service in `application/`, repo in `infra/outgoing/postgres/`, controller in `infra/entrypoint/rest/` |
| II. Test-First Discipline | ✅ Pass | Tests planned for domain, service, repository, controller layers |
| III. Domain-Driven Validation | ✅ Pass | `GroupInvite` entity has `Validate()`, expiration check is domain logic |
| IV. Consistent API Contract | ✅ Pass | Controller follows BodyParser → Validate → service → map → response pattern |
| V. Infrastructure Abstraction | ✅ Pass | `GroupInviteRepository` interface injected; DB interface reused; squirrel queries |
| VI. Simplicity & YAGNI | ✅ Pass | UUID v7 reused as invite ID (no new abstraction); no revocation in v1 |
| VII. Performance & Observability | ✅ Pass | `context.Context` propagated; `log.Println` on infra errors |

No violations — Complexity Tracking table not required.

## Project Structure

### Documentation (this feature)

```text
specs/001-group-invite-links/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output
│   ├── POST_groups_groupID_invites.md
│   └── POST_invites_inviteID_join.md
└── tasks.md             # Phase 2 output (/speckit.tasks)
```

### Source Code (repository root)

```text
internal/
├── domain/
│   ├── group_invite.go                             # NEW — entity + GroupInviteRepository interface
│   ├── group.go                                    # MODIFIED — restrict AddUser to owner only
│   ├── mock_domain/
│   │   └── group_invite_repository.go              # NEW — generated mock
│   └── build_domain/
│       └── group_invite_builder.go                 # NEW — test builder
│
├── application/
│   ├── group_invite_service.go                     # NEW — GroupInviteService interface + impl
│   ├── mock_application/
│   │   └── group_invite_service.go                 # NEW — generated mock
│   └── (group_service.go unchanged)
│
└── infra/
    ├── config/
    │   └── config.go                               # MODIFIED — add InviteConfig
    ├── runner.go                                   # MODIFIED — wire new dependencies
    ├── entrypoint/
    │   ├── routes.go                               # MODIFIED — register invite routes
    │   └── rest/
    │       ├── group_invite_dto.go                 # NEW — DTOs + mapper functions
    │       ├── group_invite_controller.go           # NEW — GroupInviteController
    │       └── build_rest/
    │           └── group_invite_dto_builder.go      # NEW — test builder
    └── outgoing/
        └── postgres/
            ├── group_invite.go                     # NEW — postgres model struct
            ├── group_invite_repository.go           # NEW — repository implementation
            ├── build_postgres/
            │   └── group_invite_builder.go          # NEW — test builder
            └── migrations/
                ├── 000004_create_group_invites_table.up.sql   # NEW
                └── 000004_create_group_invites_table.down.sql # NEW

Test files (co-located with source):
internal/domain/group_invite_test.go
internal/application/group_invite_service_test.go
internal/infra/outgoing/postgres/group_invite_repository_test.go
internal/infra/entrypoint/rest/group_invite_controller_test.go
internal/domain/group_test.go                       # MODIFIED — update AddUser tests
```

**Structure Decision**: Single project (existing layout). All new code follows the
established Clean Architecture layer conventions. No new top-level directories required.

## Complexity Tracking

> No violations detected — table not required.
