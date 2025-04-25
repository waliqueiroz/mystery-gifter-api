package build_domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type UserBuilder struct {
	user domain.User
}

func NewUserBuilder() *UserBuilder {
	now := time.Now().UTC()

	return &UserBuilder{
		user: domain.User{
			ID:        uuid.New().String(),
			Name:      "DefaultName",
			Surname:   "DefaultSurname",
			Email:     "default@example.com",
			Password:  "defaultpassword",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (b *UserBuilder) WithID(id string) *UserBuilder {
	b.user.ID = id
	return b
}

func (b *UserBuilder) WithName(name string) *UserBuilder {
	b.user.Name = name
	return b
}

func (b *UserBuilder) WithSurname(surname string) *UserBuilder {
	b.user.Surname = surname
	return b
}

func (b *UserBuilder) WithEmail(email string) *UserBuilder {
	b.user.Email = email
	return b
}

func (b *UserBuilder) WithPassword(password string) *UserBuilder {
	b.user.Password = password
	return b
}

func (b *UserBuilder) WithCreatedAt(createdAt time.Time) *UserBuilder {
	b.user.CreatedAt = createdAt
	return b
}

func (b *UserBuilder) WithUpdatedAt(updatedAt time.Time) *UserBuilder {
	b.user.UpdatedAt = updatedAt
	return b
}

func (b *UserBuilder) Build() domain.User {
	return b.user
}
