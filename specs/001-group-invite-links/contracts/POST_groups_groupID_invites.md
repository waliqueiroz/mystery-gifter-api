# Contract: Create Invite Link

**Method**: `POST`
**Path**: `/api/v1/groups/:groupID/invites`
**Auth**: Bearer token required
**Controller**: `GroupInviteController.Create`
**Tag**: `groups`

## Request

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `groupID` | UUID string | Yes | ID of the group to create an invite for |

### Headers

| Header | Value |
|--------|-------|
| `Authorization` | `Bearer <jwt-token>` |

### Body

None.

## Response

### 201 Created

```json
{
  "id": "018e1234-abcd-7000-8000-000000000001",
  "group_id": "018e1234-abcd-7000-8000-000000000002",
  "expires_at": "2026-03-30T12:00:00Z",
  "created_at": "2026-03-29T12:00:00Z"
}
```

**Schema**: `GroupInviteDTO`

| Field | Type | Description |
|-------|------|-------------|
| `id` | string | Unique invite ID (UUID v7). Share this to invite users. |
| `group_id` | string | ID of the group the invite grants access to |
| `expires_at` | ISO 8601 timestamp | When the invite expires (UTC) |
| `created_at` | ISO 8601 timestamp | When the invite was created (UTC) |

## Error Responses

| Status | Condition |
|--------|-----------|
| 401 | Missing or invalid authentication token |
| 403 | Authenticated user is not the group owner |
| 404 | Group not found |
| 409 | Group is not in OPEN status |
