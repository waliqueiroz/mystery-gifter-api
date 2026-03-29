# Research: Group Invite Links

**Feature**: 001-group-invite-links
**Date**: 2026-03-29

## Decision 1: Invite ID Strategy

**Decision**: Use UUID v7 (via the existing `IdentityGenerator` interface) as the `GroupInvite.ID`.

**Rationale**: UUID v7 is time-ordered, 128-bit random, and collision-resistant — sufficient
security for an invite identifier in this domain. The `IdentityGenerator` interface and its
UUID implementation already exist in `internal/infra/outgoing/identity/`. Injecting it into
`GroupInviteService` introduces zero new abstractions (YAGNI). A dedicated `TokenGenerator`
interface would add complexity with no material security benefit for this use case.

**Alternatives considered**:
- `crypto/rand` → base64url: Stronger randomness, but requires a new interface and
  implementation; overkill given UUID v7's entropy level.
- Signed JWT: No DB storage needed, but cannot be revoked (even theoretically in v2),
  and adds dependency on JWT signing for a non-auth purpose.
- Separate `token` field (non-PK): Rejected — the invite ID already serves as the shared
  identifier; a separate token field would be redundant (YAGNI).

---

## Decision 2: Expiration Configuration

**Decision**: Add `InviteConfig.LinkExpiration time.Duration` to `internal/infra/config/config.go`
using `caarlos0/env` tag `env:"INVITE_LINK_EXPIRATION" envDefault:"24h"`.

**Rationale**: `caarlos0/env` v11 natively parses `time.Duration` from strings like `"24h"`,
`"30m"`, `"1h30m"`. Using `envDefault:"24h"` satisfies the spec requirement of a fallback
default without any custom parsing code. The existing `RequiredIfNoDef: true` parse option
applies only to fields without a default, so this field will be optional in the environment.

**Alternatives considered**:
- Storing as integer (hours): Less expressive, requires custom parsing, not idiomatic Go.
- Storing as seconds: Same problem as above.

---

## Decision 3: Invite Storage

**Decision**: Store invite records in a new PostgreSQL table `group_invites` with `id` (UUID)
as the primary key.

**Rationale**: DB-backed invites can be queried by ID (O(1) PK lookup), support future
revocation (v2), and allow audit logging. The expiration timestamp is stored as `TIMESTAMPTZ`
and evaluated server-side at redemption time, which is immune to client clock drift.
ID existence check + expiration check are a single `GetContext` call.

**Alternatives considered**:
- Self-validating JWT: Would avoid DB reads at redemption, but prevents revocation,
  requires JWT signing infrastructure for non-auth purpose, and complicates domain purity.

---

## Decision 4: Entity Naming

**Decision**: Domain entity named `GroupInvite` (not `InviteLink`); primary key field named `ID`.

**Rationale**: `GroupInvite` is consistent with the table name `group_invites` and follows
the same naming pattern as other entities (`Group`, `User`, `Match`). Using `ID` as the
primary key field name is consistent with every other entity in the codebase (`User.ID`,
`Group.ID`) and avoids a special-case `Token` field. The UUID serves as both the database
PK and the identifier shared with users — no separate "token" concept is needed.

**Alternatives considered**:
- `InviteLink` + `Token` field: Descriptive of purpose but inconsistent with project conventions.

---

## Decision 5: Breaking Change Scope

**Decision**: Modify `Group.AddUser()` domain method to reject calls where
`requesterID != g.OwnerID`. The self-join path (`requesterID == targetUser.ID`) is removed.

**Rationale**: The spec (FR-013, FR-014) explicitly requires this change. Bundling it with
the invite-link feature is correct because the invite link becomes the replacement path
for self-service joining. Shipping the restriction without the invite link would break
existing workflows.

**Impact**:
- `internal/domain/group.go` — condition update in `AddUser`
- `internal/domain/group_test.go` — test cases for removed self-join path become negative
  test cases (self-join now returns `ForbiddenError`)
- Swagger description for `POST /api/v1/groups/{groupID}/users` must be updated to reflect
  owner-only restriction

---

## Decision 6: Expired Invite Error Type

**Decision**: Map expired invites to `domain.ConflictError` (HTTP 409).

**Rationale**: The invite exists in the database but the temporal state conflicts with the
redemption attempt. `ResourceNotFoundError` (404) is reserved for invite IDs that do not
exist at all. Using 409 for "invite exists but cannot be used" is consistent with how the
project uses `ConflictError` for other state-conflict scenarios (e.g., group already archived).

---

## Decision 7: Route Structure

**Decision**:
- `POST /api/v1/groups/:groupID/invites` — create invite (auth required, owner only)
- `POST /api/v1/invites/:inviteID/join` — redeem invite (auth required)

**Rationale**: The create endpoint is scoped under the group resource it belongs to.
The redeem endpoint is at the top-level `/invites` resource because the user only knows
the invite ID (not the group ID) when following a shared link. This mirrors common patterns
(e.g., Slack invite URLs do not embed the workspace ID in the user-visible path).
