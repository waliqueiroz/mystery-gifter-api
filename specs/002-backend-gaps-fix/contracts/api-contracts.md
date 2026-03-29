# API Contracts: Backend Gaps Fix for Group Management

**Branch**: `002-backend-gaps-fix` | **Date**: 2026-03-29

---

## Modified Endpoints

### POST /api/v1/groups — CreateGroup

**Change**: `description` field is now optional.

**Request body** (`CreateGroupDTO`):
```json
{
  "name": "Secret Santa 2024",
  "description": "Optional description"
}
```

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| `name` | string | yes | non-empty |
| `description` | string | **no** | max 255 chars |

**Response** (`GroupDTO`) — unchanged shape except `matches` field removed (see GET below).

---

### GET /api/v1/groups/{groupID} — GetGroupByID

**Changes**:
1. Returns `403 Forbidden` if authenticated user is not a group member.
2. Response no longer includes the `matches` field.

**Authorization**: Bearer token required. User must be a member of the requested group.

**Response 200** (`GroupDTO`):
```json
{
  "id": "uuid",
  "name": "Secret Santa 2024",
  "description": "Optional description",
  "users": [{ "id": "uuid", "name": "...", "surname": "...", "email": "..." }],
  "owner_id": "uuid",
  "status": "OPEN",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

Note: `matches` field is **absent** from the response.

**Error responses**:
| Status | Condition |
|--------|-----------|
| 401 | Missing or invalid Bearer token |
| 403 | Authenticated user is not a group member |
| 404 | Group not found |

---

### GET /api/v1/groups — SearchGroups

**Change**: `GroupSummaryDTO` items in the result never contained `matches` (already correct). No functional change; documentation updated for clarity.

---

### POST /api/v1/groups/{groupID}/reopen — ReopenGroup

**Change**: Swagger description corrected only. No functional change.

**Corrected description**: "Reopens a group with MATCHED status, clearing all draw results and returning the group to OPEN status."

---

## New Endpoint

### GET /api/v1/groups/{groupID}/invites/active — GetActiveGroupInvite

**Authorization**: Bearer token required. User must be a member of the group.

**Path parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `groupID` | string (UUID) | yes | Unique group identifier |

**Response 200** (`GroupInviteDTO`):
```json
{
  "id": "uuid",
  "group_id": "uuid",
  "expires_at": "2024-01-08T00:00:00Z",
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Error responses**:
| Status | Condition |
|--------|-----------|
| 401 | Missing or invalid Bearer token |
| 403 | Authenticated user is not a group member |
| 404 | Group not found OR no active (non-expired) invite exists |

**Behavior**:
- Returns the most recent non-expired invite (`expires_at > NOW()`, ordered by `created_at DESC`).
- Returns 404 (not a distinct status code) for both "no invite created" and "all invites expired" — the client treats both identically ("no active invite").
- Read-only; does not create or modify any invite.
