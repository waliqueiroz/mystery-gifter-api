package build_rest

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

type GroupDTOBuilder struct {
	groupDTO rest.GroupDTO
}

func NewGroupDTOBuilder() *GroupDTOBuilder {
	user := NewUserDTOBuilder().Build()
	now := time.Now()

	return &GroupDTOBuilder{
		groupDTO: rest.GroupDTO{
			ID:        uuid.NewString(),
			Name:      "Default Group",
			Users:     []rest.UserDTO{user},
			OwnerID:   user.ID,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (b *GroupDTOBuilder) WithID(id string) *GroupDTOBuilder {
	b.groupDTO.ID = id
	return b
}

func (b *GroupDTOBuilder) WithName(name string) *GroupDTOBuilder {
	b.groupDTO.Name = name
	return b
}

func (b *GroupDTOBuilder) WithUsers(users []rest.UserDTO) *GroupDTOBuilder {
	b.groupDTO.Users = users
	return b
}

func (b *GroupDTOBuilder) WithOwnerID(ownerID string) *GroupDTOBuilder {
	b.groupDTO.OwnerID = ownerID
	return b
}

func (b *GroupDTOBuilder) WithCreatedAt(createdAt time.Time) *GroupDTOBuilder {
	b.groupDTO.CreatedAt = createdAt
	return b
}

func (b *GroupDTOBuilder) WithUpdatedAt(updatedAt time.Time) *GroupDTOBuilder {
	b.groupDTO.UpdatedAt = updatedAt
	return b
}

func (b *GroupDTOBuilder) Build() rest.GroupDTO {
	return b.groupDTO
}
