# Feature Specification: Group Invite Links

**Feature Branch**: `001-group-invite-links`
**Created**: 2026-03-28
**Status**: Draft
**Input**: User description: "Quero implementar o recurso na API para que usuários possam gerar links de convite para convidar outros usuários para grupos. O link deve expirar depois de um tempo que deve ser configurado como variável de ambiente"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Generate Invite Link (Priority: P1)

A group owner wants to invite someone to join their group. They generate an invite link that can be shared with anyone (e.g., via messaging apps). The link is valid for a configurable period and can be used by the recipient to join the group.

**Why this priority**: This is the foundational action — without generating an invite link, the entire feature is unusable. It directly enables the invitation flow.

**Independent Test**: Can be fully tested by calling the generate-invite endpoint for a group and verifying that a link with an expiration is returned. Delivers value as soon as a link can be produced.

**Acceptance Scenarios**:

1. **Given** an authenticated user who is the owner of an open group, **When** they request an invite link for that group, **Then** the system returns a unique invite link along with its expiration date/time.
2. **Given** an authenticated user who is NOT the owner of the group, **When** they request an invite link, **Then** the system returns a forbidden error.
3. **Given** a group that is not in OPEN status, **When** the owner requests an invite link, **Then** the system returns a conflict error indicating that the group must be open to generate invites.
4. **Given** a valid invite link generation request, **When** another invite link is generated for the same group, **Then** each request produces a new unique link (multiple active invites may coexist).

---

### User Story 2 - Join Group via Invite Link (Priority: P2)

A user receives an invite link and wants to join the group it refers to. They use the link/token to be added to the group as a member, without needing the group owner to add them manually.

**Why this priority**: Without the ability to use the invite link, generating it has no value. This story completes the invitation flow end-to-end.

**Independent Test**: Can be fully tested by generating a valid invite token (from US1), then using it to join the group and verifying the user appears in the group's member list.

**Acceptance Scenarios**:

1. **Given** an authenticated user who is not yet a member of the group, **When** they use a valid non-expired invite link/token, **Then** they are added to the group and the updated group data is returned.
2. **Given** an authenticated user who is already a member of the group, **When** they use a valid invite link/token, **Then** the system accepts the request gracefully (idempotent — no error, no duplicate membership).
3. **Given** an authenticated user, **When** they use an expired invite link/token, **Then** the system returns an error indicating the link has expired.
4. **Given** an authenticated user, **When** they use an invite link/token that does not exist or is invalid, **Then** the system returns a not-found error.
5. **Given** an authenticated user using a valid invite link, **When** the referenced group is no longer in OPEN status, **Then** the system returns a conflict error indicating the group is not accepting new members.

---

### Edge Cases

- What happens when the invite expiration environment variable is missing or not a valid Go duration string? The system uses `"24h"` as the fallback default and logs a warning at startup.
- What happens when a user tries to join a group they are already the owner of via an invite link? They are already a member; the request is treated as idempotent.
- What happens when an invite link is generated for a group that is subsequently archived or matched before it is used? The join attempt must fail with a conflict error (group not open).
- What happens when the system clock drifts? Expiration must be stored as an absolute timestamp and evaluated server-side at redemption time.
- What happens when a non-owner user attempts to add themselves (or any other user) directly via the `AddUser` endpoint without an invite link? The system returns a forbidden error — self-join is no longer permitted.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow only the group owner to generate invite links for their group.
- **FR-002**: System MUST generate a cryptographically unique, unguessable token for each invite link.
- **FR-003**: System MUST persist each invite token associated with its target group and an absolute expiration timestamp.
- **FR-004**: The invite expiration duration MUST be configurable via an environment variable accepting a Go duration string (e.g., `"24h"`, `"30m"`). If absent or invalid, the system MUST fall back to `"24h"` and log a warning at startup.
- **FR-005**: System MUST reject redemption attempts for expired invite tokens with a clear expiration error.
- **FR-006**: System MUST reject redemption attempts for tokens that do not exist or cannot be found.
- **FR-007**: System MUST add the authenticated user to the group when a valid, non-expired invite token is redeemed.
- **FR-008**: System MUST enforce that joining via invite is only possible when the group is in OPEN status.
- **FR-009**: System MUST handle duplicate join attempts idempotently — redeeming an invite when already a member must not produce an error.
- **FR-010**: System MUST return the full invite details (token, expiration) to the group owner upon invite creation.
- **FR-011**: Multiple active invite links for the same group MUST be supported simultaneously.
- **FR-012**: A single invite link MUST be usable by any number of authenticated users until it expires (multi-use behavior); no per-redemption invalidation occurs.
- **FR-013**: The "add user" operation MUST be restricted so that only the group owner can add members directly (by user ID). Self-join by third parties without an invite link MUST be rejected with a forbidden error.
- **FR-014**: The invite link redemption endpoint is the sole path for a non-owner user to join a group on their own initiative.

### Key Entities

- **InviteLink**: Represents a single invite token. Key attributes: unique token, target group, expiration timestamp, creation timestamp. Linked to the group it grants access to.
- **Group**: Existing entity. An OPEN group can have invite links generated for it and new members added via those links.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A group owner can generate an invite link in a single request and immediately share it.
- **SC-002**: A new user can join a group via an invite link in a single request without any additional action from the group owner. Direct self-join without an invite link is rejected in 100% of attempts.
- **SC-003**: Expired invite links are consistently rejected — 100% of redemption attempts after expiration must fail with an appropriate error message.
- **SC-004**: The invite expiration duration can be changed via environment variable without any code changes or redeployment of new code.
- **SC-005**: All invite-related operations complete within the same response-time expectations as other group operations (no perceptible latency increase for end users).

## Clarifications

### Session 2026-03-29

- Q: Should invite links be single-use (expire after first redemption) or multi-use (usable by any number of users until expiration)? → A: Multi-use — any number of users can join via the same link before it expires; no usage counter is tracked.
- Q: Should the group owner be able to list or revoke active invite links? → A: Out of scope for v1 — no listing or revocation endpoints.
- Q: What format should the invite expiration duration environment variable accept? → A: Go duration string format (e.g., `"24h"`, `"48h"`, `"30m"`, `"1h30m"`).
- Q: Should any authenticated user be able to self-add to a group (current behavior), or should third-party joining be restricted to invite links only? → A: Restrict — remove self-join from `AddUser`; the invite link becomes the only path for third parties. The owner retains the ability to add members directly.

## Assumptions

- Only the group owner can generate invite links (consistent with the existing permission model where owners manage group membership).
- An invite link is represented as a URL or token string; the API returns the token and the client constructs or uses the full URL as needed.
- Multiple invite links can be active at the same time for the same group (no single-active-invite constraint).
- Invite links are **multi-use**: any number of authenticated users can join via the same link before it expires. No per-link usage counter is tracked.
- The environment variable for expiration duration accepts a Go duration string (e.g., `"24h"`, `"30m"`, `"1h30m"`). If not set or invalid, the system applies a default of `"24h"` and logs a warning at startup.
- Listing and revocation of active invite links are out of scope for v1. No endpoints will be provided for these operations.
- Users must be authenticated to redeem an invite link — anonymous join is out of scope.
- The invite link redemption reuses the group's existing membership rules (group must be OPEN, no duplicate members).
- **Breaking change to existing behavior**: `AddUser` will be restricted — self-join (a user adding themselves directly) is removed. Only the group owner may use `AddUser` to add a specific member by user ID. This change ships as part of this feature.
