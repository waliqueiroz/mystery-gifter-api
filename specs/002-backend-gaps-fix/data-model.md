# Data Model: Backend Gaps Fix for Group Management

**Branch**: `002-backend-gaps-fix` | **Date**: 2026-03-29

## Domain Changes

### `domain.Group` — modified fields and new methods

```go
// Description: remove "required" tag
Description string `validate:"omitempty,max=255"`

// New method — checks if userID is in group.Users
func (g *Group) IsMember(userID string) bool

// New method — returns ForbiddenError if not a member
func (g *Group) CanView(requesterID string) error
```

No structural fields added or removed. `Matches []Match` remains in the domain entity (required by `GetUserMatch` and `GenerateMatches`).

### `domain.GroupInviteRepository` — new method

```go
type GroupInviteRepository interface {
    Create(ctx context.Context, groupInvite GroupInvite) error
    GetByID(ctx context.Context, id string) (*GroupInvite, error)
    GetActiveByGroupID(ctx context.Context, groupID string) (*GroupInvite, error) // NEW
}
```

`GetActiveByGroupID` returns the most recent non-expired invite for the given group, or `ResourceNotFoundError` if none exists.

---

## Application Layer Changes

### `application.GroupService` — modified interface

```go
type GroupService interface {
    Create(ctx context.Context, name, description, ownerID string) (*domain.Group, error)
    GetByID(ctx context.Context, groupID, requesterID string) (*domain.Group, error) // CHANGED: added requesterID
    Search(ctx context.Context, filters domain.GroupFilters) (*domain.SearchResult[domain.GroupSummary], error)
    AddUser(ctx context.Context, groupID, requesterID, targetUserID string) (*domain.Group, error)
    RemoveUser(ctx context.Context, groupID, requesterID, targetUserID string) (*domain.Group, error)
    GenerateMatches(ctx context.Context, groupID, requesterID string) (*domain.Group, error)
    Reopen(ctx context.Context, groupID, requesterID string) (*domain.Group, error)
    Archive(ctx context.Context, groupID, requesterID string) (*domain.Group, error)
    GetUserMatch(ctx context.Context, groupID, requesterID string) (*domain.User, error)
}
```

### `application.GroupInviteService` — new method

```go
type GroupInviteService interface {
    Create(ctx context.Context, groupID, requesterID string) (*domain.GroupInvite, error)
    GetActive(ctx context.Context, groupID, requesterID string) (*domain.GroupInvite, error) // NEW
    JoinGroup(ctx context.Context, inviteID, userID string) (*domain.Group, error)
}
```

---

## REST Layer Changes

### `CreateGroupDTO` — description optional

```go
type CreateGroupDTO struct {
    // Group name
    // required: true
    Name string `json:"name" validate:"required"`

    // Group description
    // max length: 255
    Description string `json:"description" validate:"omitempty,max=255"`
    // NOTE: removed "required: true" Swagger annotation and "required" validator tag
}
```

### `GroupDTO` — description optional + matches removed

```go
type GroupDTO struct {
    ID          string    `json:"id" validate:"required,uuid"`
    Name        string    `json:"name" validate:"required"`
    Description string    `json:"description" validate:"omitempty,max=255"` // removed "required"
    Users       []UserDTO `json:"users" validate:"required,min=1"`
    OwnerID     string    `json:"owner_id" validate:"required,uuid"`
    // Matches field REMOVED — draw pairs must not be exposed
    Status      string    `json:"status" validate:"required,oneof=OPEN MATCHED ARCHIVED"`
    CreatedAt   time.Time `json:"created_at" validate:"required"`
    UpdatedAt   time.Time `json:"updated_at" validate:"required"`
}
```

---

## State Transitions

No new state transitions. Existing group status machine (`OPEN → MATCHED → ARCHIVED`, `MATCHED → OPEN`) is unchanged.

The `IsMember` check is purely a read-time authorization guard, not a state mutation.

---

## Validation Rules Summary

| Field | Rule | Layer |
|-------|------|-------|
| `Group.Description` | `omitempty,max=255` | Domain |
| `CreateGroupDTO.Description` | `omitempty,max=255` | REST DTO |
| `GroupDTO.Description` | `omitempty,max=255` | REST DTO |
| `Group.CanView` | requesterID must be in `group.Users` | Domain |
| `GroupInvite` (active) | `expires_at > NOW()`, most recent by `created_at` | Repository query |
