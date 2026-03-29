package build_postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres"
)

type GroupInviteBuilder struct {
	groupInvite postgres.GroupInvite
}

func NewGroupInviteBuilder() *GroupInviteBuilder {
	now := time.Now().UTC()

	return &GroupInviteBuilder{
		groupInvite: postgres.GroupInvite{
			ID:        uuid.New().String(),
			GroupID:   uuid.New().String(),
			ExpiresAt: now.Add(24 * time.Hour),
			CreatedAt: now,
		},
	}
}

func (b *GroupInviteBuilder) WithID(id string) *GroupInviteBuilder {
	b.groupInvite.ID = id
	return b
}

func (b *GroupInviteBuilder) WithGroupID(groupID string) *GroupInviteBuilder {
	b.groupInvite.GroupID = groupID
	return b
}

func (b *GroupInviteBuilder) WithExpiresAt(expiresAt time.Time) *GroupInviteBuilder {
	b.groupInvite.ExpiresAt = expiresAt
	return b
}

func (b *GroupInviteBuilder) WithCreatedAt(createdAt time.Time) *GroupInviteBuilder {
	b.groupInvite.CreatedAt = createdAt
	return b
}

func (b *GroupInviteBuilder) Build() postgres.GroupInvite {
	return b.groupInvite
}
