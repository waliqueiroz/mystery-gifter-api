# Feature Specification: Backend Gaps Fix for Group Management

**Feature Branch**: `002-backend-gaps-fix`
**Created**: 2026-03-29
**Status**: Draft
**Input**: Backend gaps identified during frontend implementation of group management (003-group-management)

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Create group without description (Priority: P1)

As an authenticated member, I want to create a group by providing only the name, without being required to fill in a description, because the description field is optional as defined in the product specification.

**Why this priority**: Directly blocks the group creation flow on the frontend. Without this fix, any attempt to create a group without a description returns a validation error.

**Independent Test**: Can be tested in isolation by sending a group creation request without the `description` field and verifying the group is created successfully.

**Acceptance Scenarios**:

1. **Given** an authenticated user, **When** they submit a group creation request with `name` filled and `description` absent, **Then** the group is created successfully and `description` is empty/null.
2. **Given** an authenticated user, **When** they submit a group creation request with both `name` and `description` filled, **Then** the group is created normally with the description saved.
3. **Given** an authenticated user, **When** they submit a group creation request with `description` containing more than 255 characters, **Then** the API returns a validation error.

---

### User Story 2 - Any member can view the active group invite (Priority: P1)

As a group member (not necessarily the owner), I want to be able to view the group's active invite link so I can share it with others.

**Why this priority**: Without this endpoint, only the owner can create invites but no member can retrieve the link to share. Blocks the invite flow for all non-owners.

**Independent Test**: Can be tested in isolation: the owner creates an invite, then a non-owner member queries the active invite and verifies they receive the invite data.

**Acceptance Scenarios**:

1. **Given** a group with an active invite (not expired), **When** any authenticated group member queries the active invite, **Then** the most recent non-expired invite data is returned (200).
2. **Given** a group with no active invite (none created or all expired), **When** any member queries the active invite, **Then** returns not found (404).
3. **Given** an authenticated user who is not a group member, **When** they try to query the active invite, **Then** returns forbidden (403).
4. **Given** an unauthenticated user, **When** they try to query the active invite, **Then** returns unauthorized (401).

---

### User Story 3 - Group detail access restricted to members (Priority: P1)

As the system, I want to ensure that only group members can view group details, preventing external users from accessing data such as member lists and group information.

**Why this priority**: Security vulnerability (data leakage). Any authenticated user can currently access full group data without being a member. Also blocks the correct redirect behavior on the frontend.

**Independent Test**: Can be tested in isolation: a user with no membership in the group attempts to access the group detail and receives a forbidden response.

**Acceptance Scenarios**:

1. **Given** an authenticated user who is a group member, **When** they query the group detail, **Then** returns 200 with the full group data.
2. **Given** an authenticated user who is not a group member, **When** they try to query the group detail, **Then** returns forbidden (403).
3. **Given** an unauthenticated user, **When** they try to query the group detail, **Then** returns unauthorized (401).

---

### User Story 4 - Matches not exposed in group responses (Priority: P1)

As a group member or owner, I do not want the complete draw pairs (who gifts whom) to be exposed in any API response — not in the group list, not in the group detail, and not even to the owner. The draw reveal must only happen through the dedicated individual endpoint.

**Why this priority**: Core security and business rule of the application. Exposing all draw pairs breaks the secret santa secrecy. Even the owner cannot see the full pairs; each user sees only their own recipient.

**Independent Test**: Can be tested by verifying that the `matches` field is absent from group listing and group detail responses, even after a draw has been performed.

**Acceptance Scenarios**:

1. **Given** a group with a completed draw, **When** the owner queries the group detail, **Then** the response does not contain the draw pairs.
2. **Given** a group with a completed draw, **When** any member queries the group listing, **Then** no item in the listing contains draw pairs.
3. **Given** any authenticated group member after a draw, **When** they want to know who they are gifting, **Then** they must use the individual match reveal endpoint to obtain only their recipient.

---

### User Story 5 - Correct documentation for the group reopen endpoint (Priority: P3)

As a developer consuming the API, I want the group reopen endpoint documentation to correctly describe its behavior: it reopens a group with MATCHED status (clearing all draw results), not "reopens an archived group".

**Why this priority**: Documentation-only impact. Does not block any functional flow; the domain behavior is already correct.

**Independent Test**: Verify in the generated API documentation that the reopen endpoint description is correct.

**Acceptance Scenarios**:

1. **Given** the generated API documentation, **When** a developer views the group reopen endpoint, **Then** the description states that it reopens groups with MATCHED status, clearing all draw results.

---

### Edge Cases

- What happens when the active invite is queried on a group with multiple invites, some expired and one active? Must return only the most recent non-expired one.
- What happens if the owner creates a new invite while one already exists and is still active? The query endpoint must return the most recent valid invite.
- What happens if `description` is sent as an empty string `""`? It should be treated as absent (omitempty behavior).
- What impact does removing `matches` from the group response have on existing clients? Must confirm no existing client depends on this field before removing it.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The system MUST accept group creation without the `description` field, treating it as optional with a maximum of 255 characters.
- **FR-002**: The system MUST provide an endpoint to retrieve the most recent active invite for a group, accessible to any authenticated member, returning the invite data or 404 if no valid invite exists.
- **FR-003**: The group detail endpoint MUST return 403 for any authenticated user who is not a member of the group.
- **FR-004**: The group detail endpoint MUST NOT include draw pairs in the response, regardless of the group status or the requester's role (including the owner).
- **FR-005**: The group listing endpoint MUST NOT include draw pairs in the returned items.
- **FR-006**: The active invite query endpoint MUST return 403 if the authenticated user is not a member of the group.
- **FR-007**: The API documentation for the group reopen endpoint MUST correctly describe that it reopens groups with MATCHED status, clearing all draw results.

### Key Entities

- **GroupInvite**: Represents a time-limited invite link for a group. Attributes: unique identifier (token), group identifier, expiration date, creation date. The active invite query filters by future expiration date and returns the most recent record.
- **Group (creation input)**: The `description` field becomes optional in group creation input.
- **Group (response)**: Draw pairs are removed from all public group responses (listing and detail).

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 100% of group creation requests without `description` are accepted successfully with no validation errors.
- **SC-002**: Non-owner members can view the active invite in 100% of cases where a valid invite exists.
- **SC-003**: 100% of attempts to access group detail by non-members result in a forbidden (403) response.
- **SC-004**: No response from the group listing or group detail endpoints contains draw pairs, regardless of group status or requester role.
- **SC-005**: The generated API documentation correctly describes all modified endpoints.

## Assumptions

- The `matches` field in the current `GroupDTO` is only consumed by the frontend under development; no external or legacy client depends on this field.
- The invite expiration check logic already exists in the domain or can be derived from the existing `GroupInvite` entity fields.
- The membership check for the group detail endpoint will be performed by verifying whether the authenticated user is present in the group's member list.
- The "most recent active invite" is defined as the record with the latest creation date whose expiration date is after the time of the query.
- No other endpoints need to be modified beyond those listed in the functional requirements.
