package build_postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres"
)

type GroupBuilder struct {
	group postgres.Group
}

func NewGroupBuilder() *GroupBuilder {
	return &GroupBuilder{
		group: postgres.Group{
			ID:        uuid.New().String(),
			Name:      "Test Group",
			OwnerID:   uuid.New().String(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}
}

func (b *GroupBuilder) WithID(id string) *GroupBuilder {
	b.group.ID = id
	return b
}

func (b *GroupBuilder) WithName(name string) *GroupBuilder {
	b.group.Name = name
	return b
}

func (b *GroupBuilder) WithOwnerID(ownerID string) *GroupBuilder {
	b.group.OwnerID = ownerID
	return b
}

func (b *GroupBuilder) WithCreatedAt(createdAt time.Time) *GroupBuilder {
	b.group.CreatedAt = createdAt
	return b
}

func (b *GroupBuilder) WithUpdatedAt(updatedAt time.Time) *GroupBuilder {
	b.group.UpdatedAt = updatedAt
	return b
}

func (b *GroupBuilder) Build() postgres.Group {
	return b.group
}
