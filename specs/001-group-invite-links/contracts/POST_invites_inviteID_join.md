# Contract: Join Group via Invite

**Method**: `POST`
**Path**: `/api/v1/invites/:inviteID/join`
**Auth**: Bearer token required
**Controller**: `GroupInviteController.Join`
**Tag**: `invites`

## Request

### Path Parameters

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `inviteID` | UUID string | Yes | Invite ID received from the group owner |

### Headers

| Header | Value |
|--------|-------|
| `Authorization` | `Bearer <jwt-token>` |

### Body

None.

## Response

### 200 OK

Full `GroupDTO` of the group the user has just joined.

```json
{
  "id": "018e1234-abcd-7000-8000-000000000002",
  "name": "Christmas 2026",
  "description": "Our annual secret santa",
  "status": "OPEN",
  "owner_id": "018e1234-abcd-7000-8000-000000000003",
  "users": [
    {
      "id": "018e1234-abcd-7000-8000-000000000003",
      "name": "Alice",
      "surname": "Smith",
      "email": "alice@example.com"
    },
    {
      "id": "018e1234-abcd-7000-8000-000000000004",
      "name": "Bob",
      "surname": "Jones",
      "email": "bob@example.com"
    }
  ],
  "matches": [],
  "created_at": "2026-03-28T10:00:00Z",
  "updated_at": "2026-03-29T12:00:00Z"
}
```

**Schema**: `GroupDTO` (existing schema — no new DTO required)

**Idempotency**: If the authenticated user is already a member of the group, the request
succeeds with 200 and returns the current group data (no duplicate membership created).

## Error Responses

| Status | Condition |
|--------|-----------|
| 401 | Missing or invalid authentication token |
| 404 | Invite ID does not exist |
| 409 | Invite has expired, or the group is no longer in OPEN status |
