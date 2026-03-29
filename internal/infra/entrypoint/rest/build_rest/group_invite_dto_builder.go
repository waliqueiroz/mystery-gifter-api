package build_rest

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

type GroupInviteDTOBuilder struct {
	dto rest.GroupInviteDTO
}

func NewGroupInviteDTOBuilder() *GroupInviteDTOBuilder {
	now := time.Now().UTC()

	return &GroupInviteDTOBuilder{
		dto: rest.GroupInviteDTO{
			ID:        uuid.New().String(),
			GroupID:   uuid.New().String(),
			ExpiresAt: now.Add(24 * time.Hour),
			CreatedAt: now,
		},
	}
}

func (b *GroupInviteDTOBuilder) WithID(id string) *GroupInviteDTOBuilder {
	b.dto.ID = id
	return b
}

func (b *GroupInviteDTOBuilder) WithGroupID(groupID string) *GroupInviteDTOBuilder {
	b.dto.GroupID = groupID
	return b
}

func (b *GroupInviteDTOBuilder) WithExpiresAt(expiresAt time.Time) *GroupInviteDTOBuilder {
	b.dto.ExpiresAt = expiresAt
	return b
}

func (b *GroupInviteDTOBuilder) WithCreatedAt(createdAt time.Time) *GroupInviteDTOBuilder {
	b.dto.CreatedAt = createdAt
	return b
}

func (b *GroupInviteDTOBuilder) Build() rest.GroupInviteDTO {
	return b.dto
}
