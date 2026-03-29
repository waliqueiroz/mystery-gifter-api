# Research: Backend Gaps Fix for Group Management

**Branch**: `002-backend-gaps-fix` | **Date**: 2026-03-29

## 1. Codebase Audit

### Files to be changed

| File | Change type | Reason |
|------|-------------|--------|
| `internal/domain/group.go` | Modify | Make `Description` optional; add `IsMember`, `CanView` methods |
| `internal/domain/group_invite.go` | Modify | Add `GetActiveByGroupID` to `GroupInviteRepository` interface |
| `internal/application/group_service.go` | Modify | Add `requesterID` to `GetByID`; remove membership check from service (moved to domain) |
| `internal/application/group_invite_service.go` | Modify | Add `GetActive` method to interface and implementation |
| `internal/infra/entrypoint/rest/group_dto.go` | Modify | Make `description` optional in `CreateGroupDTO` and `GroupDTO`; remove `Matches` from `GroupDTO` |
| `internal/infra/entrypoint/rest/group_controller.go` | Modify | Pass `authUserID` to `GetByID` |
| `internal/infra/entrypoint/rest/group_invite_controller.go` | Modify | Add `GetActive` handler |
| `internal/infra/outgoing/postgres/group_invite_repository.go` | Modify | Add `GetActiveByGroupID` method |
| `internal/infra/entrypoint/routes.go` | Modify | Add new route; fix Swagger descriptions |
| `internal/infra/entrypoint/rest/build_rest/group_dto_builder.go` | Modify | Remove `Matches` from builder (if present) |
| `mock_application/group_service.go` | Regenerate | `GetByID` signature changed |
| `mock_application/group_invite_service.go` | Regenerate | New `GetActive` method |
| `mock_domain/group_invite_repository.go` | Regenerate | New `GetActiveByGroupID` method |

### No migrations needed

The `group_invites` table already has `id`, `group_id`, `expires_at`, `created_at`. The new repository query uses existing columns with a `WHERE expires_at > NOW()` filter — no schema change required.

---

## 2. Decision Log

### D-001: Where to enforce the membership check (BC-003)

**Decision**: Add `CanView(requesterID string) error` to `domain.Group`. The service calls it after fetching the group.

**Rationale**: Constitution III — "permission or state verification must be encapsulated in entity methods." ForbiddenError is already defined in `domain/errors.go` and is mapped to HTTP 403 by the error handler. Adding to domain makes the rule consistent with all other permission checks (`CanCreateInvite`, `GenerateMatches`, etc.).

**Implementation**: `group.CanView(requesterID)` iterates `group.Users` and returns `ForbiddenError` if requesterID is not found. The `GroupService.GetByID` signature changes to `GetByID(ctx, groupID, requesterID string)`.

---

### D-002: IsMember helper vs inline loop

**Decision**: Add private `IsMember(userID string) bool` to `domain.Group`, and build `CanView` on top of it.

**Rationale**: `IsMember` is also needed in `groupInviteService.GetActive` (to validate membership before returning invite). Sharing the logic avoids duplication without creating a premature abstraction — it solves an immediate, present need (2 callers).

---

### D-003: Active invite query strategy (BC-002)

**Decision**: New repository method `GetActiveByGroupID(ctx, groupID string) (*GroupInvite, error)` queries:
`SELECT * FROM group_invites WHERE group_id = $1 AND expires_at > NOW() ORDER BY created_at DESC LIMIT 1`

Returns `domain.NewResourceNotFoundError` when no rows found (maps to HTTP 404). The service checks membership before calling the repository.

**Rationale**: Single SQL query, no application-level filtering. `sql.ErrNoRows` → 404 follows the existing repository convention.

---

### D-004: Removing `Matches` from `GroupDTO` (BC-004)

**Decision**: Remove `Matches []MatchDTO` field from `GroupDTO` struct and `mapGroupFromDomain`. The `domain.Group.Matches` field remains intact (needed for `GetUserMatch` and `GenerateMatches` logic).

**Rationale**: The frontend must never see all draw pairs — not even the owner. The only way to reveal a match is `GET /api/v1/groups/{groupID}/matches/user`, which returns only the requester's recipient. Removing from DTO is the safest enforcement point.

**Side effect**: `mapGroupFromDomain` no longer calls `mapMatchesFromDomain`. The `match_dto.go` mapper function remains (still used internally if needed) but is not invoked by the group mapper.

---

### D-005: Description optionality — three layers (BC-001)

**Decision**: Remove `required` validation tag from `Description` in three places: `domain.Group`, `CreateGroupDTO`, and `GroupDTO`. No migration needed.

**Rationale**: `domain.Group.Description` currently has `validate:"required,max=255"`. Since `NewGroup` receives description from the service which receives it from the DTO, the domain must also allow empty. All three layers must be consistent to avoid "valid DTO rejected by domain".

---

### D-006: GroupController.GetByID — auth extraction

**Decision**: Extract `authUserID` in `GroupController.GetByID` using `c.AuthTokenManager.GetAuthUserID(ctx.Locals("user"))` and pass it to `groupService.GetByID`.

**Rationale**: The controller already does this pattern in `Reopen`, `Archive`, `GenerateMatches`, `GetUserMatch`. No new pattern introduced.
