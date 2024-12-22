package build_postgres

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/repository/postgres"
)

type UserBuilder struct {
	user postgres.User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		user: postgres.User{
			ID:        uuid.New().String(),
			Name:      "DefaultName",
			Surname:   "DefaultSurname",
			Email:     "default@example.com",
			Password:  "defaultpassword",
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
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

func (b *UserBuilder) Build() postgres.User {
	return b.user
}
