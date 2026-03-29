# Data Model: Group Invite Links

**Feature**: 001-group-invite-links
**Date**: 2026-03-29

## New Entity: GroupInvite

### Domain Struct (`internal/domain/group_invite.go`)

```go
type GroupInvite struct {
    ID        string    `validate:"required,uuid"`
    GroupID   string    `validate:"required,uuid"`
    ExpiresAt time.Time `validate:"required"`
    CreatedAt time.Time `validate:"required"`
}
```

**Fields**:

| Field | Type | Validation | Description |
|-------|------|-----------|-------------|
| `ID` | `string` | `required,uuid` | UUID v7 — unique invite identifier, used as PK and shared token |
| `GroupID` | `string` | `required,uuid` | ID of the group this invite grants access to |
| `ExpiresAt` | `time.Time` | `required` | Absolute UTC timestamp when the invite expires |
| `CreatedAt` | `time.Time` | `required` | Creation timestamp |

**Factory function**:

```go
func NewGroupInvite(identityGenerator IdentityGenerator, groupID string, expiration time.Duration) (*GroupInvite, error)
```

Creates a new `GroupInvite`, generating the ID via `identityGenerator.Generate()` and
computing `ExpiresAt = time.Now().Add(expiration)`. Calls `Validate()` before returning.

**Domain methods**:

```go
func (i *GroupInvite) Validate() error
func (i *GroupInvite) IsExpired() bool  // returns time.Now().After(i.ExpiresAt)
```

### Repository Interface (`internal/domain/group_invite.go`)

```go
//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/group_invite_repository.go . GroupInviteRepository

type GroupInviteRepository interface {
    Create(ctx context.Context, groupInvite GroupInvite) error
    GetByID(ctx context.Context, id string) (*GroupInvite, error)
}
```

---

## Database Schema

### New Table: `group_invites` (migration `000004`)

```sql
CREATE TABLE IF NOT EXISTS group_invites (
    id         UUID         NOT NULL PRIMARY KEY,
    group_id   UUID         NOT NULL REFERENCES groups(id),
    expires_at TIMESTAMPTZ  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

**Notes**:
- `id` is UUID — consistent with all other entity primary keys in this project.
- No FK to `users` — invites are not scoped to a specific inviter at the DB level
  (owner check happens in the service layer using the auth token).
- No `updated_at` — `GroupInvite` records are immutable once created.

### Postgres Model Struct (`internal/infra/outgoing/postgres/group_invite.go`)

```go
type GroupInvite struct {
    ID        string    `db:"id"`
    GroupID   string    `db:"group_id"`
    ExpiresAt time.Time `db:"expires_at"`
    CreatedAt time.Time `db:"created_at"`
}
```

---

## Modified Entity: Group

### Change to `AddUser` (`internal/domain/group.go`)

**Before** (current):
```go
if requesterID != g.OwnerID && requesterID != targetUser.ID {
    return NewForbiddenError("only the group owner can add other users")
}
```

**After** (new):
```go
if requesterID != g.OwnerID {
    return NewForbiddenError("only the group owner can add other users")
}
```

Self-join (where `requesterID == targetUser.ID` and `requesterID != g.OwnerID`) now returns
`ForbiddenError`. The invite redemption path handles self-service joining instead.

---

## Config Addition

### `internal/infra/config/config.go`

New `InviteConfig` struct added to `Config`:

```go
type InviteConfig struct {
    LinkExpiration time.Duration `env:"INVITE_LINK_EXPIRATION" envDefault:"24h"`
}

type Config struct {
    Database DatabaseConfig
    Auth     AuthConfig
    Invite   InviteConfig      // NEW
}
```

New env var: `INVITE_LINK_EXPIRATION` (optional — defaults to `"24h"`).
`.env.example` must be updated to document this variable.

---

## New Application Service: GroupInviteService

### Interface (`internal/application/group_invite_service.go`)

```go
//go:generate go run go.uber.org/mock/mockgen -destination mock_application/group_invite_service.go . GroupInviteService

type GroupInviteService interface {
    Create(ctx context.Context, groupID, requesterID string) (*domain.GroupInvite, error)
    JoinGroup(ctx context.Context, inviteID, userID string) (*domain.Group, error)
}
```

### Implementation logic

**Create**:
1. `groupRepository.GetByID(ctx, groupID)` — returns 404 if not found
2. Check `requesterID == group.OwnerID` — returns 403 if not owner
3. Check `group.IsOpen()` — returns 409 if not open
4. `domain.NewGroupInvite(identityGenerator, groupID, expiration)` — generates invite
5. `groupInviteRepository.Create(ctx, *groupInvite)` — persists
6. Return `groupInvite`

**JoinGroup**:
1. `groupInviteRepository.GetByID(ctx, inviteID)` — returns 404 if not found
2. `groupInvite.IsExpired()` — returns 409 ("invite has expired") if true
3. `groupRepository.GetByID(ctx, groupInvite.GroupID)` — returns 404 if not found
4. `userService.GetByID(ctx, userID)` — returns 404 if not found
5. `group.AddUser(userID, *targetUser)` — uses existing domain method (checks OPEN, idempotent)
6. `groupRepository.Update(ctx, *group)` — persists membership change
7. Return updated `group`

**Dependencies injected**: `domain.GroupInviteRepository`, `domain.GroupRepository`,
`application.UserService`, `domain.IdentityGenerator`, `time.Duration` (expiration)

---

## DTOs

### `internal/infra/entrypoint/rest/group_invite_dto.go`

```go
// swagger:model GroupInviteDTO
type GroupInviteDTO struct {
    ID        string    `json:"id"`
    GroupID   string    `json:"group_id"`
    ExpiresAt time.Time `json:"expires_at"`
    CreatedAt time.Time `json:"created_at"`
}
```

No request body DTOs required — both endpoints receive all needed data from path params
and the auth token.
