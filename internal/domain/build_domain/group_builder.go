package build_domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type GroupBuilder struct {
	group domain.Group
}

func NewGroupBuilder() *GroupBuilder {
	user := NewUserBuilder().Build()
	now := time.Now()

	return &GroupBuilder{
		group: domain.Group{
			ID:        uuid.New().String(),
			Name:      "Test Group",
			Users:     []domain.User{user},
			OwnerID:   user.ID,
			CreatedAt: now,
			UpdatedAt: now,
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

func (b *GroupBuilder) WithUsers(users []domain.User) *GroupBuilder {
	b.group.Users = users
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

func (b *GroupBuilder) Build() domain.Group {
	return b.group
}
