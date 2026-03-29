# Quickstart: Group Invite Links

**Feature**: 001-group-invite-links
**Date**: 2026-03-29

This guide walks through the complete invite flow for manual validation after
implementation.

## Prerequisites

- API running locally (`make run`)
- Two registered users: **Alice** (group owner) and **Bob** (invitee)
- `INVITE_LINK_EXPIRATION` set in `.env` (e.g., `INVITE_LINK_EXPIRATION=24h`)

## Step 1 — Authenticate both users

```bash
# Login as Alice (owner)
curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret"}' | jq .token
# → ALICE_TOKEN

# Login as Bob (invitee)
curl -s -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{"email":"bob@example.com","password":"secret"}' | jq .token
# → BOB_TOKEN
```

## Step 2 — Alice creates a group

```bash
curl -s -X POST http://localhost:8080/api/v1/groups \
  -H "Authorization: Bearer $ALICE_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"Christmas 2026","description":"Our annual secret santa"}' | jq .id
# → GROUP_ID
```

## Step 3 — Alice generates an invite

```bash
curl -s -X POST http://localhost:8080/api/v1/groups/$GROUP_ID/invites \
  -H "Authorization: Bearer $ALICE_TOKEN" | jq .
```

Expected response (HTTP 201):
```json
{
  "id": "018e1234-abcd-7000-8000-000000000001",
  "group_id": "<GROUP_ID>",
  "expires_at": "2026-03-30T12:00:00Z",
  "created_at": "2026-03-29T12:00:00Z"
}
```

Save the invite ID: `INVITE_ID=018e1234-abcd-7000-8000-000000000001`

## Step 4 — Bob joins via invite

```bash
curl -s -X POST http://localhost:8080/api/v1/invites/$INVITE_ID/join \
  -H "Authorization: Bearer $BOB_TOKEN" | jq '.users[].name'
```

Expected: Bob's name appears in the group's user list. HTTP 200 with full `GroupDTO`.

## Step 5 — Validate idempotency

```bash
# Second join attempt by Bob — must return 200, not an error
curl -s -o /dev/null -w "%{http_code}" \
  -X POST http://localhost:8080/api/v1/invites/$INVITE_ID/join \
  -H "Authorization: Bearer $BOB_TOKEN"
# → 200
```

## Step 6 — Validate permission: non-owner cannot generate invite

```bash
curl -s -o /dev/null -w "%{http_code}" \
  -X POST http://localhost:8080/api/v1/groups/$GROUP_ID/invites \
  -H "Authorization: Bearer $BOB_TOKEN"
# → 403
```

## Step 7 — Validate self-join via AddUser is blocked

```bash
# Bob tries to add himself directly (old behavior — now forbidden)
curl -s -o /dev/null -w "%{http_code}" \
  -X POST http://localhost:8080/api/v1/groups/$GROUP_ID/users \
  -H "Authorization: Bearer $BOB_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"user_id\":\"$BOB_ID\"}"
# → 403
```

## Step 8 — Validate expired invite rejection

```bash
# Set INVITE_LINK_EXPIRATION=1ms, restart, generate a new invite, wait 1s, then:
curl -s -o /dev/null -w "%{http_code}" \
  -X POST http://localhost:8080/api/v1/invites/$EXPIRED_INVITE_ID/join \
  -H "Authorization: Bearer $BOB_TOKEN"
# → 409
```

## Step 9 — Run all tests

```bash
make test
# All tests must pass.

make generate-docs
# Swagger generation must succeed without errors.
```
