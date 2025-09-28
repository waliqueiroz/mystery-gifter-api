package build_rest

import (
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

type GroupSummaryDTOBuilder struct {
	groupSummaryDTO rest.GroupSummaryDTO
}

func NewGroupSummaryDTOBuilder() *GroupSummaryDTOBuilder {
	return &GroupSummaryDTOBuilder{
		groupSummaryDTO: rest.GroupSummaryDTO{
			ID:        "550e8400-e29b-41d4-a716-446655440000",
			Name:      "Test Group",
			Status:    "OPEN",
			OwnerID:   "550e8400-e29b-41d4-a716-446655440001",
			UserCount: 5,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}
}

func (b *GroupSummaryDTOBuilder) WithID(id string) *GroupSummaryDTOBuilder {
	b.groupSummaryDTO.ID = id
	return b
}

func (b *GroupSummaryDTOBuilder) WithName(name string) *GroupSummaryDTOBuilder {
	b.groupSummaryDTO.Name = name
	return b
}

func (b *GroupSummaryDTOBuilder) WithStatus(status string) *GroupSummaryDTOBuilder {
	b.groupSummaryDTO.Status = status
	return b
}

func (b *GroupSummaryDTOBuilder) WithOwnerID(ownerID string) *GroupSummaryDTOBuilder {
	b.groupSummaryDTO.OwnerID = ownerID
	return b
}

func (b *GroupSummaryDTOBuilder) WithUserCount(userCount int) *GroupSummaryDTOBuilder {
	b.groupSummaryDTO.UserCount = userCount
	return b
}

func (b *GroupSummaryDTOBuilder) WithCreatedAt(createdAt time.Time) *GroupSummaryDTOBuilder {
	b.groupSummaryDTO.CreatedAt = createdAt
	return b
}

func (b *GroupSummaryDTOBuilder) WithUpdatedAt(updatedAt time.Time) *GroupSummaryDTOBuilder {
	b.groupSummaryDTO.UpdatedAt = updatedAt
	return b
}

func (b *GroupSummaryDTOBuilder) Build() rest.GroupSummaryDTO {
	return b.groupSummaryDTO
}
