# Quickstart: Backend Gaps Fix for Group Management

**Branch**: `002-backend-gaps-fix`

## Environment Setup

```bash
cp .env.example .env
# Fill in: DB_HOST, DB_PORT, DB_DATABASE, DB_USERNAME, DB_PASSWORD, AUTH_SECRET_KEY, AUTH_SESSION_DURATION
```

## Run & Test

```bash
make build          # Compile
make run            # Start on :8080
make test           # All unit tests
make generate-docs  # Regenerate Swagger spec (required after any DTO/route change)
```

## Mock Regeneration

After any interface change, regenerate mocks:

```bash
go generate ./...
```

Affected interfaces in this feature:
- `domain.GroupInviteRepository` — new `GetActiveByGroupID` method
- `application.GroupService` — `GetByID` signature changed
- `application.GroupInviteService` — new `GetActive` method

## New Endpoint to Test

```bash
# Create group invite (owner only)
curl -X POST http://localhost:8080/api/v1/groups/{groupID}/invites \
  -H "Authorization: Bearer <token>"

# Get active invite (any member)
curl http://localhost:8080/api/v1/groups/{groupID}/invites/active \
  -H "Authorization: Bearer <token>"
# → 200 GroupInviteDTO  (if active invite exists)
# → 404                 (if no active invite)
# → 403                 (if not a group member)

# Verify matches are NOT in group detail
curl http://localhost:8080/api/v1/groups/{groupID} \
  -H "Authorization: Bearer <token>"
# → GroupDTO without "matches" field

# Verify non-member gets 403
curl http://localhost:8080/api/v1/groups/{groupID} \
  -H "Authorization: Bearer <non-member-token>"
# → 403 Forbidden

# Create group without description
curl -X POST http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name": "My Group"}'
# → 201 Created
```

## Key INVITE_ID Note

The `id` field of `GroupInviteDTO` is the token used in the share URL:
`https://<host>/invite/{id}`

This token is also the `inviteID` parameter in `POST /api/v1/invites/{inviteID}/join`.
