# Tasks: Backend Gaps Fix for Group Management

**Input**: Design documents from `/specs/002-backend-gaps-fix/`
**Feature**: `002-backend-gaps-fix`

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel within the same branch (different files, no dependencies)
- **[Story]**: Which user story this task belongs to
- Each Phase = one PR / one branch. Move to the next Phase only after the current PR is merged and approved.

---

## Branching Strategy (Stacked Branches)

```
main
 ├── [CP1] fix/002-bc005-reopen-swagger    ← Independent PR, merge at any time
 └── feat/002-bc001-description-optional    ← Start after CP1 is in review or merged
      └── feat/002-bc003-domain-membership  ← Stacked on CP2 (same group.go file)
           └── feat/002-bc004-remove-matches ← Stacked on CP3 (same group_dto.go)
                └── feat/002-bc003-service-getbyid ← Stacked on CP4 (needs IsMember/CanView)
                     └── feat/002-bc002-active-invite ← Stacked on CP5 (needs IsMember domain + service)
```

**Parallel opportunities:**
- CP1 is fully independent — can be opened as a PR immediately.
- CP2–CP6 form a sequential stack — each waits for the previous to merge.

---

## Phase 2: CP1 — Fix Swagger Description for Reopen (US5) 🟢 Independent

**Branch**: `fix/002-bc005-reopen-swagger`
**Based on**: `002-backend-gaps-fix` (or `main`)
**Goal**: Correct the misleading Swagger description for `POST /api/v1/groups/{groupID}/reopen`.
**Independent Test**: Run `make generate-docs`; verify the generated spec describes the reopen endpoint correctly as reopening a MATCHED group, not an archived one.

- [ ] T001 [US5] Update the Swagger operation description and summary for `ReopenGroup` in `internal/infra/entrypoint/routes.go` (lines ~428–461): change description from "Reopen an archived group" / "This endpoint reopens an archived group" to "Reopen a group with MATCHED status, clearing all draw results and returning the group to OPEN status."
- [ ] T002 [US5] Run `make generate-docs` and verify the updated description appears in `docs/specs/swagger.yaml`

**Checkpoint**: `make generate-docs` passes; description is accurate. Open PR and request review. ✋

---

## Phase 3: CP2 — Make Description Optional (US1)

**Branch**: `feat/002-bc001-description-optional`
**Based on**: `002-backend-gaps-fix` (or `main`)
**Goal**: Allow group creation without providing a description.
**Independent Test**: Send `POST /api/v1/groups` with only `name` and verify 201 response with no description.

- [ ] T003 [P] [US1] Remove `required` from `Description` validation tag in `internal/domain/group.go` (line ~33): change `validate:"required,max=255"` to `validate:"omitempty,max=255"`
- [ ] T004 [P] [US1] Remove `required` from `CreateGroupDTO.Description` in `internal/infra/entrypoint/rest/group_dto.go` (line ~22): change `validate:"required,max=255"` to `validate:"omitempty,max=255"` and remove the `// required: true` Swagger annotation
- [ ] T005 [P] [US1] Remove `required` from `GroupDTO.Description` in `internal/infra/entrypoint/rest/group_dto.go` (line ~48): change `validate:"required"` to `validate:"omitempty,max=255"`
- [ ] T006 [US1] Update `internal/domain/build_domain/group_builder.go` if `Description` defaults to a non-empty value for required-validation reasons; ensure `NewGroupBuilder()` still works with empty description
- [ ] T007 [US1] Update `internal/infra/entrypoint/rest/build_rest/create_group_dto_builder.go` if needed — verify the default description value is still acceptable
- [ ] T008 [US1] Update unit tests in `internal/domain/group_test.go` to add a scenario: "should create group successfully when description is empty"
- [ ] T009 [US1] Update unit tests in `internal/application/group_service_test.go` to cover the no-description creation scenario
- [ ] T010 [US1] Update unit tests in `internal/infra/entrypoint/rest/group_controller_test.go` if any test builds a `CreateGroupDTO` with description marked required
- [ ] T011 [US1] Run `make test` — all tests must pass; run `make generate-docs` — Swagger spec must generate without errors

**Checkpoint**: `make test` and `make generate-docs` both pass. Open PR and request review. ✋

---

## Phase 4: CP3 — Add IsMember + CanView to Domain (US3 — Part 1)

**Branch**: `feat/002-bc003-domain-membership`
**Based on**: CP2 branch (stacked; rebase on CP2 after it merges)
**Goal**: Add the domain-level membership check methods to `Group`. This is the foundational rule that all membership enforcement depends on.
**Independent Test**: Run `make test` for `internal/domain/...` — new test scenarios for `IsMember` and `CanView` must pass.

- [ ] T012 [P] [US3] Add `IsMember(userID string) bool` method to `Group` in `internal/domain/group.go`: iterates `g.Users` and returns true if `userID` matches any `user.ID`
- [ ] T013 [P] [US3] Add `CanView(requesterID string) error` method to `Group` in `internal/domain/group.go`: calls `g.IsMember(requesterID)` and returns `NewForbiddenError("user is not a member of this group")` if false
- [ ] T014 [US3] Add unit test scenarios in `internal/domain/group_test.go`:
  - `Test_Group_IsMember`: "should return true when user is a member", "should return false when user is not a member"
  - `Test_Group_CanView`: "should return nil when user is a member", "should return forbidden error when user is not a member"
- [ ] T015 [US3] Run `make test ./internal/domain/...` — all tests must pass

**Checkpoint**: Domain methods tested and passing. Open PR and request review. ✋

---

## Phase 5: CP4 — Remove Matches from GroupDTO (US4)

**Branch**: `feat/002-bc004-remove-matches`
**Based on**: CP3 branch (stacked; rebase on CP3 after it merges)
**Goal**: Eliminate draw pairs from all REST responses. `GroupDTO` must not expose `matches`.
**Independent Test**: After a draw is performed, `GET /api/v1/groups/{groupID}` response must not contain a `matches` field.

- [ ] T016 [US4] Remove the `Matches []MatchDTO` field from `GroupDTO` struct in `internal/infra/entrypoint/rest/group_dto.go` (line ~60) and remove its Swagger annotation comment
- [ ] T017 [US4] Update `mapGroupFromDomain` in `internal/infra/entrypoint/rest/group_dto.go`: remove the `matches, err := mapMatchesFromDomain(group.Matches)` call and the `Matches: matches` assignment from the `GroupDTO` construction
- [ ] T018 [US4] Remove `WithMatches` method and `Matches` initialization from `internal/infra/entrypoint/rest/build_rest/group_dto_builder.go`; remove `Matches: []rest.MatchDTO{}` from `NewGroupDTOBuilder()`
- [ ] T019 [US4] Update all test files that call `.WithMatches(...)` on `GroupDTOBuilder` or assert on `GroupDTO.Matches`:
  - `internal/infra/entrypoint/rest/group_controller_test.go` — remove `.WithMatches(...)` calls from expected DTO construction
  - `internal/infra/entrypoint/rest/group_invite_controller_test.go` — same
- [ ] T020 [US4] Run `make test` — all tests must pass; run `make generate-docs` — Swagger spec must generate without errors

**Checkpoint**: `make test` and `make generate-docs` both pass; no `matches` field in GroupDTO. Open PR and request review. ✋

---

## Phase 6: CP5 — Enforce Membership on GetByID (US3 — Part 2)

**Branch**: `feat/002-bc003-service-getbyid`
**Based on**: CP4 branch (stacked; rebase on CP4 after it merges)
**Goal**: `GET /api/v1/groups/{groupID}` returns 403 for non-members. Requires `IsMember`/`CanView` from CP3.
**Independent Test**: Authenticated non-member calls `GET /api/v1/groups/{groupID}` → 403 Forbidden.

- [ ] T021 [US3] Update `GroupService` interface in `internal/application/group_service.go`: change `GetByID(ctx context.Context, groupID string)` to `GetByID(ctx context.Context, groupID, requesterID string) (*domain.Group, error)`
- [ ] T022 [US3] Update `groupService.GetByID` implementation in `internal/application/group_service.go`: after fetching the group, call `group.CanView(requesterID)` and return the error if it fails
- [ ] T023 [US3] Regenerate `internal/application/mock_application/group_service.go` by running `go generate ./internal/application/...`
- [ ] T024 [US3] Update `GroupController.GetByID` in `internal/infra/entrypoint/rest/group_controller.go`: extract `authUserID` using `c.AuthTokenManager.GetAuthUserID(ctx.Locals("user"))` and pass it as second argument to `groupService.GetByID`
- [ ] T025 [US3] Update `internal/application/group_service_test.go` — add test scenarios for `Test_groupService_GetByID`:
  - "should return forbidden error when requester is not a member"
  - "should return group successfully when requester is a member"
  - Update existing scenario to pass `requesterID` argument
- [ ] T026 [US3] Update `internal/infra/entrypoint/rest/group_controller_test.go` — update `Test_GroupController_GetByID`:
  - Add mock expectation for `AuthTokenManager.GetAuthUserID`
  - Pass `requesterID` in `GetByID` mock expectation
  - Add scenario: "should return status 403 when user is not a member"
  - Update existing scenarios to match new mock signatures
- [ ] T027 [US3] Run `make test` — all tests must pass; run `make build` to confirm compilation

**Checkpoint**: `make test` and `make build` both pass; `GET /api/v1/groups/{groupID}` returns 403 for non-members. Open PR and request review. ✋

---

## Phase 7: CP6 — Active Invite Endpoint (US2)

**Branch**: `feat/002-bc002-active-invite`
**Based on**: CP5 branch (stacked; rebase on CP5 after it merges)
**Goal**: Any group member can retrieve the active invite link via `GET /api/v1/groups/{groupID}/invites/active`.
**Independent Test**: Member calls `GET /api/v1/groups/{groupID}/invites/active` after owner creates an invite → 200 with `GroupInviteDTO`. Same call when no invite exists → 404.

- [ ] T028 [P] [US2] Add `GetActiveByGroupID(ctx context.Context, groupID string) (*GroupInvite, error)` to `GroupInviteRepository` interface in `internal/domain/group_invite.go`
- [ ] T029 [P] [US2] Add `GetActive(ctx context.Context, groupID, requesterID string) (*domain.GroupInvite, error)` to `GroupInviteService` interface in `internal/application/group_invite_service.go`
- [ ] T030 [US2] Implement `GetActiveByGroupID` in `internal/infra/outgoing/postgres/group_invite_repository.go`: query `SELECT * FROM group_invites WHERE group_id = $1 AND expires_at > NOW() ORDER BY created_at DESC LIMIT 1` using squirrel; map `sql.ErrNoRows` to `domain.NewResourceNotFoundError("active group invite not found")`
- [ ] T031 [US2] Regenerate `internal/domain/mock_domain/group_invite_repository.go` by running `go generate ./internal/domain/...`
- [ ] T032 [US2] Implement `GetActive` in `internal/application/group_invite_service.go`:
  1. `groupRepository.GetByID(ctx, groupID)` — returns 404 if group not found
  2. `group.CanView(requesterID)` — returns 403 if not a member
  3. `groupInviteRepository.GetActiveByGroupID(ctx, groupID)` — returns 404 if no active invite
- [ ] T033 [US2] Regenerate `internal/application/mock_application/group_invite_service.go` by running `go generate ./internal/application/...`
- [ ] T034 [US2] Add `GetActive` handler to `GroupInviteController` in `internal/infra/entrypoint/rest/group_invite_controller.go`:
  - Extract `groupID` from `ctx.Params("groupID")`
  - Extract `authUserID` from `c.authTokenManager.GetAuthUserID(ctx.Locals("user"))`
  - Call `groupInviteService.GetActive(ctx.Context(), groupID, authUserID)`
  - Map result with `mapGroupInviteFromDomain` and return `ctx.JSON(...)`
- [ ] T035 [US2] Register the new route in `internal/infra/entrypoint/routes.go` with full Swagger annotation for `GET /api/v1/groups/{groupID}/invites/active` (operation name: `GetActiveGroupInvite`; tag: `groups`; responses: 200 `GroupInviteDTO`, 401, 403, 404)
- [ ] T036 [US2] Add unit tests in `internal/application/group_invite_service_test.go` — `Test_groupInviteService_GetActive`:
  - "should return active invite when requester is a member"
  - "should return not found error when group does not exist"
  - "should return forbidden error when requester is not a member"
  - "should return not found error when no active invite exists"
- [ ] T037 [US2] Add unit tests in `internal/infra/entrypoint/rest/group_invite_controller_test.go` — `Test_GroupInviteController_GetActive`:
  - "should return status 200 and the active invite when found"
  - "should return status 403 when user is not a member"
  - "should return status 404 when no active invite exists"
- [ ] T038 [US2] Run `make test` — all tests pass; run `make generate-docs` — Swagger generates without errors; run `make build` — compiles successfully

**Checkpoint**: All quality gates pass. `GET .../invites/active` works correctly. Open PR and request review. ✋

---

## Final Phase: Polish & Quality Gates

**Branch**: can be part of CP6 or a separate cleanup PR after all merges
**Goal**: Final validation that all changes work together end-to-end.

- [ ] T039 Run `make build` — binary compiles without errors
- [ ] T040 Run `make test` — all unit tests pass across all packages (`go test -cover ./...`)
- [ ] T041 Run `make generate-docs` — Swagger spec generates without errors; verify `docs/specs/swagger.yaml` reflects all changes (no `matches` in GroupDTO, new route, corrected reopen description, description optional)
- [ ] T042 Verify `CLAUDE.md` Active Technologies section is up to date (no new dependencies added in this feature)

---

## Dependencies & Execution Order

### Phase Dependencies

```
CP1 (Phase 2)   ─────────────────────────────────────────────┐
                                                               ↓ merge any time
CP2 (Phase 3) → CP3 (Phase 4) → CP4 (Phase 5) → CP5 (Phase 6) → CP6 (Phase 7)
```

- **CP1**: No dependencies — can be opened and merged at any time
- **CP2**: Based on `002-backend-gaps-fix`; start immediately (parallel with CP1)
- **CP3**: Stacked on CP2 — wait for CP2 to merge
- **CP4**: Stacked on CP3 — wait for CP3 to merge
- **CP5**: Stacked on CP4 — wait for CP4 to merge
- **CP6**: Stacked on CP5 — wait for CP5 to merge

### Within Each Checkpoint

- Tasks marked [P] touch different files and can be coded simultaneously
- Run `make test` after each checkpoint before opening the PR

### Parallel Opportunities

- CP1 can be merged completely independently and in parallel with the CP2–CP6 stack
- T003, T004, T005 (CP2) are in the same file — implement in one edit pass
- T012, T013 (CP3) add separate methods to the same file — implement in one edit pass
- T028, T029 (CP6) add to different files — can be coded simultaneously

---

## Parallel Execution Example: CP6 (Active Invite Endpoint)

```
# Simultaneously:
Task T028: Add GetActiveByGroupID to domain/group_invite.go
Task T029: Add GetActive to application/group_invite_service.go (interface only)

# Then sequentially:
Task T030: Implement GetActiveByGroupID in postgres repo
Task T031: Regenerate mock_domain/group_invite_repository.go
Task T032: Implement GetActive in groupInviteService
Task T033: Regenerate mock_application/group_invite_service.go
Task T034: Add GetActive handler to controller
Task T035: Add route to routes.go
Tasks T036, T037: Tests
Task T038: Quality gates
```

---

## Implementation Strategy

### MVP First (CP1 + CP2)

1. Merge CP1 (trivial docs fix — unblocks PR review workflow)
2. Merge CP2 (description optional — unblocks frontend group creation)
3. **STOP and VALIDATE**: Frontend can now create groups without description

### Incremental Delivery

1. CP1 → Quick win, no risk
2. CP2 → Unblocks frontend group creation (BC-001)
3. CP3 → Domain foundation for all membership checks
4. CP4 → Removes security leak from GroupDTO (BC-004)
5. CP5 → Enforces access control on group detail (BC-003)
6. CP6 → Unblocks frontend invite flow for non-owners (BC-002)

---

## Notes

- [P] tasks = different files, can be coded simultaneously within the same branch
- Each checkpoint is a standalone PR — keep diffs small and focused
- Always rebase the next stack branch after the parent merges to keep diffs clean
- Run `go generate ./...` immediately after any interface change (before writing tests)
- `make test` must pass locally before opening a PR
