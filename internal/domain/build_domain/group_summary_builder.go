package build_domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type GroupSummaryBuilder struct {
	groupSummary domain.GroupSummary
}

func NewGroupSummaryBuilder() *GroupSummaryBuilder {
	now := time.Now().UTC()
	return &GroupSummaryBuilder{
		groupSummary: domain.GroupSummary{
			ID:        uuid.New().String(),
			Name:      "Test Group",
			Status:    domain.GroupStatusOpen,
			OwnerID:   uuid.New().String(),
			UserCount: 1,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (b *GroupSummaryBuilder) WithID(id string) *GroupSummaryBuilder {
	b.groupSummary.ID = id
	return b
}

func (b *GroupSummaryBuilder) WithName(name string) *GroupSummaryBuilder {
	b.groupSummary.Name = name
	return b
}

func (b *GroupSummaryBuilder) WithStatus(status domain.GroupStatus) *GroupSummaryBuilder {
	b.groupSummary.Status = status
	return b
}

func (b *GroupSummaryBuilder) WithOwnerID(ownerID string) *GroupSummaryBuilder {
	b.groupSummary.OwnerID = ownerID
	return b
}

func (b *GroupSummaryBuilder) WithUserCount(userCount int) *GroupSummaryBuilder {
	b.groupSummary.UserCount = userCount
	return b
}

func (b *GroupSummaryBuilder) WithCreatedAt(createdAt time.Time) *GroupSummaryBuilder {
	b.groupSummary.CreatedAt = createdAt
	return b
}

func (b *GroupSummaryBuilder) WithUpdatedAt(updatedAt time.Time) *GroupSummaryBuilder {
	b.groupSummary.UpdatedAt = updatedAt
	return b
}

func (b *GroupSummaryBuilder) Build() domain.GroupSummary {
	return b.groupSummary
}
