---
description: "Task list for Group Invite Links feature"
---

# Tasks: Group Invite Links

**Input**: Design documents from `/specs/001-group-invite-links/`
**Prerequisites**: plan.md ✅, spec.md ✅, data-model.md ✅, contracts/ ✅, research.md ✅

**Stacked Branches Strategy**: Each task is scoped to produce one small, focused PR.
Tasks T003 and T004 are marked [P] — they touch different files and can land in the
same branch commit or in parallel branches stacked on T002.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies between each other)
- **[Story]**: User story this task belongs to (US1 or US2)

---

## Phase 1: Setup

**Purpose**: Database schema — no application code, isolated migration PR.

- [X] T001 Add migration `000004_create_group_invites_table` in `internal/infra/outgoing/postgres/migrations/` (up: `CREATE TABLE group_invites` with `id UUID PK`, `group_id UUID FK groups`, `expires_at TIMESTAMPTZ`, `created_at TIMESTAMPTZ`; down: `DROP TABLE`)

**Checkpoint → PR1**: Migration only, reviewable standalone.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Breaking change + domain types + config. All user story work depends on this phase.

**⚠️ CRITICAL**: No user story work can begin until T002, T003, T004 are merged.

- [X] T002 Restrict `Group.AddUser()` to owner only in `internal/domain/group.go` (remove `|| requesterID == targetUser.ID` from the condition); update affected test cases in `internal/domain/group_test.go` (self-join scenarios must now assert `ForbiddenError`)
- [X] T003 [P] Create `GroupInvite` entity + `GroupInviteRepository` interface + `NewGroupInvite` factory + `Validate()` + `IsExpired()` in `internal/domain/group_invite.go`; add `//go:generate` directive and run `go generate ./...` to produce `internal/domain/mock_domain/group_invite_repository.go`; add test builder in `internal/domain/build_domain/group_invite_builder.go`; add unit tests in `internal/domain/group_invite_test.go`
- [X] T004 [P] Add `InviteConfig` struct with `LinkExpiration time.Duration` (`env:"INVITE_LINK_EXPIRATION" envDefault:"24h"`) to `internal/infra/config/config.go` and embed it in `Config`; document `INVITE_LINK_EXPIRATION=24h` in `.env.example`

**Checkpoint → PR2 (T002) + PR3 (T003 + T004)**: Breaking change and domain types reviewed separately.

---

## Phase 3: User Story 1 — Generate Invite (Priority: P1) 🎯 MVP

**Goal**: Group owner calls `POST /groups/:groupID/invites` and receives a `GroupInviteDTO` with the invite ID and expiration.

**Independent Test**: `POST /api/v1/groups/:groupID/invites` as owner returns 201 with `id` + `expires_at`. Non-owner returns 403. Group not OPEN returns 409.

### Infrastructure for User Story 1

- [X] T005 [US1] Create postgres model `GroupInvite` in `internal/infra/outgoing/postgres/group_invite.go`; implement `groupInviteRepository` (with `Create` + `GetByID`) in `internal/infra/outgoing/postgres/group_invite_repository.go` using `squirrel`; add builder in `internal/infra/outgoing/postgres/build_postgres/group_invite_builder.go`; add unit tests in `internal/infra/outgoing/postgres/group_invite_repository_test.go`

### Application for User Story 1

- [X] T006 [US1] Add `GroupInviteService` interface (with `Create` and `JoinGroup` signatures) + implement `groupInviteService.Create` in `internal/application/group_invite_service.go`; add `//go:generate` directive and run `go generate ./...` to produce `internal/application/mock_application/group_invite_service.go`; add unit tests for `Create` in `internal/application/group_invite_service_test.go`

### HTTP Layer for User Story 1

- [X] T007 [US1] Create `GroupInviteDTO` + `mapGroupInviteFromDomain` mapper in `internal/infra/entrypoint/rest/group_invite_dto.go`; add builder in `internal/infra/entrypoint/rest/build_rest/group_invite_dto_builder.go`
- [X] T008 [US1] Implement `GroupInviteController` struct + `NewGroupInviteController` + `Create` handler in `internal/infra/entrypoint/rest/group_invite_controller.go`; add unit tests in `internal/infra/entrypoint/rest/group_invite_controller_test.go`
- [X] T009 [US1] Register `POST /api/v1/groups/:groupID/invites` in `internal/infra/entrypoint/routes.go` (add `groupInviteController` param); instantiate `groupInviteRepository`, `groupInviteService`, `groupInviteController` and wire them in `internal/infra/runner.go`

**Checkpoint → PR4 (T005) + PR5 (T006) + PR6 (T007–T009)**: US1 fully functional and independently testable.

---

## Phase 4: User Story 2 — Join via Invite (Priority: P2)

**Goal**: Authenticated user calls `POST /invites/:inviteID/join` and is added to the group.

**Independent Test**: Calling join with a valid invite ID adds the user to the group (200 + GroupDTO). Expired invite returns 409. Unknown ID returns 404. Idempotent: already-member returns 200.

### Application for User Story 2

- [X] T010 [US2] Implement `groupInviteService.JoinGroup` in `internal/application/group_invite_service.go`; add unit tests for `JoinGroup` in `internal/application/group_invite_service_test.go`

### HTTP Layer for User Story 2

- [X] T011 [US2] Implement `GroupInviteController.Join` handler in `internal/infra/entrypoint/rest/group_invite_controller.go`; add unit tests for `Join` in `internal/infra/entrypoint/rest/group_invite_controller_test.go`
- [X] T012 [US2] Register `POST /api/v1/invites/:inviteID/join` in `internal/infra/entrypoint/routes.go`

**Checkpoint → PR7 (T010) + PR8 (T011–T012)**: US2 fully functional. Full invite flow now works end-to-end.

---

## Phase 5: Polish & Cross-Cutting Concerns

- [X] T013 Add Swagger annotations for `POST /groups/:groupID/invites` and `POST /invites/:inviteID/join` in `internal/infra/entrypoint/routes.go`; update `AddUserToGroup` Swagger description to reflect owner-only restriction; run `make generate-docs` and confirm `docs/specs/swagger.yaml` is up to date

**Checkpoint → PR9 (T013)**: Swagger complete, `make generate-docs` passes.

---

## Dependencies & Execution Order

### Phase Dependencies

- **Phase 1** (T001): No dependencies — migration can land first.
- **Phase 2** (T002–T004): Depends on Phase 1. T003 and T004 are independent of each other [P].
- **Phase 3** (T005–T009): Depends on Phase 2 completion.
  - T005 → T006 → T007 → T008 → T009 (sequential within Phase 3)
- **Phase 4** (T010–T012): Depends on Phase 3 (needs repository, service interface, controller struct).
  - T010 → T011 → T012 (sequential within Phase 4)
- **Phase 5** (T013): Depends on Phase 4 completion.

### Stacked Branches Map

```
main
 └─ PR1: T001  (migration)
     └─ PR2: T002  (restrict AddUser)
         └─ PR3: T003 + T004  (domain entity + config)
             └─ PR4: T005  (repository)
                 └─ PR5: T006  (service Create)
                     └─ PR6: T007 + T008 + T009  (HTTP US1)
                         └─ PR7: T010  (service JoinGroup)
                             └─ PR8: T011 + T012  (HTTP US2)
                                 └─ PR9: T013  (swagger polish)
```

### Within Each Phase

- Models/entities before services
- Services before controllers
- Controllers before route registration
- Tests co-located in the same task as implementation (run `make test` after each task)

---

## Parallel Opportunities

Tasks T003 and T004 within Phase 2 can be executed simultaneously by two developers:

```bash
# Developer A: domain entity
Task: "GroupInvite entity + repository interface in group_invite.go"  # T003

# Developer B: config
Task: "InviteConfig in config.go + .env.example"                       # T004
```

---

## Implementation Strategy

### MVP (User Story 1 only — PR1 through PR6)

1. T001 — migration (PR1)
2. T002 — breaking change (PR2)
3. T003 + T004 — domain + config (PR3)
4. T005 — repository (PR4)
5. T006 — service Create (PR5)
6. T007 + T008 + T009 — HTTP US1 (PR6)
7. **STOP and validate**: `POST /groups/:groupID/invites` works end-to-end

### Full Delivery (add PR7–PR9)

8. T010 — service JoinGroup (PR7)
9. T011 + T012 — HTTP US2 (PR8)
10. T013 — swagger polish (PR9)

---

## Notes

- [P] tasks in Phase 2 = different files, no inter-dependency
- Run `go generate ./...` as part of T003 (domain mock) and T006 (service mock)
- Run `make test` at the end of every task before marking complete
- Run `make generate-docs` only in T013 (Swagger polish)
- Each PR should include its own test updates — no separate "add tests" PRs
