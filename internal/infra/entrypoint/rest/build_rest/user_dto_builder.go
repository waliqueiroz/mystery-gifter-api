package build_rest

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

type UserDTOBuilder struct {
	userDTO rest.UserDTO
}

func NewUserDTOBuilder() *UserDTOBuilder {
	now := time.Now().UTC()
	return &UserDTOBuilder{
		userDTO: rest.UserDTO{
			ID:        uuid.NewString(),
			Name:      "DefaultName",
			Surname:   "DefaultSurname",
			Email:     "default@example.com",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func (b *UserDTOBuilder) WithID(id string) *UserDTOBuilder {
	b.userDTO.ID = id
	return b
}

func (b *UserDTOBuilder) WithName(name string) *UserDTOBuilder {
	b.userDTO.Name = name
	return b
}

func (b *UserDTOBuilder) WithSurname(surname string) *UserDTOBuilder {
	b.userDTO.Surname = surname
	return b
}

func (b *UserDTOBuilder) WithEmail(email string) *UserDTOBuilder {
	b.userDTO.Email = email
	return b
}

func (b *UserDTOBuilder) WithCreatedAt(createdAt time.Time) *UserDTOBuilder {
	b.userDTO.CreatedAt = createdAt
	return b
}

func (b *UserDTOBuilder) WithUpdatedAt(updatedAt time.Time) *UserDTOBuilder {
	b.userDTO.UpdatedAt = updatedAt
	return b
}

func (b *UserDTOBuilder) Build() rest.UserDTO {
	return b.userDTO
}
