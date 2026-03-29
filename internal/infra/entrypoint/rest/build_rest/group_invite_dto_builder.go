package build_rest

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

type GroupInviteDTOBuilder struct {
	groupInviteDTO rest.GroupInviteDTO
}

func NewGroupInviteDTOBuilder() *GroupInviteDTOBuilder {
	now := time.Now().UTC()

	return &GroupInviteDTOBuilder{
		groupInviteDTO: rest.GroupInviteDTO{
			ID:        uuid.New().String(),
			GroupID:   uuid.New().String(),
			ExpiresAt: now.Add(24 * time.Hour),
			CreatedAt: now,
		},
	}
}

func (b *GroupInviteDTOBuilder) WithID(id string) *GroupInviteDTOBuilder {
	b.groupInviteDTO.ID = id
	return b
}

func (b *GroupInviteDTOBuilder) WithGroupID(groupID string) *GroupInviteDTOBuilder {
	b.groupInviteDTO.GroupID = groupID
	return b
}

func (b *GroupInviteDTOBuilder) WithExpiresAt(expiresAt time.Time) *GroupInviteDTOBuilder {
	b.groupInviteDTO.ExpiresAt = expiresAt
	return b
}

func (b *GroupInviteDTOBuilder) WithCreatedAt(createdAt time.Time) *GroupInviteDTOBuilder {
	b.groupInviteDTO.CreatedAt = createdAt
	return b
}

func (b *GroupInviteDTOBuilder) Build() rest.GroupInviteDTO {
	return b.groupInviteDTO
}
