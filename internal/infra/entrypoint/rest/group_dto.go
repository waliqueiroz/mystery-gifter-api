package rest

import (
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type CreateGroupDTO struct {
	Name string `json:"name" validate:"required"`
}

func (g *CreateGroupDTO) Validate() error {
	if errs := validator.Validate(g); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

type GroupDTO struct {
	ID        string    `json:"id" validate:"required,uuid"`
	Name      string    `json:"name" validate:"required"`
	Users     []UserDTO `json:"users" validate:"required,min=1"`
	OwnerID   string    `json:"owner_id" validate:"required,uuid"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

func (g *GroupDTO) Validate() error {
	if errs := validator.Validate(g); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

func mapGroupFromDomain(group domain.Group) (*GroupDTO, error) {
	users, err := mapUsersFromDomain(group.Users)
	if err != nil {
		return nil, err
	}

	groupDTO := GroupDTO{
		ID:        group.ID,
		Name:      group.Name,
		Users:     users,
		OwnerID:   group.OwnerID,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}

	if err := groupDTO.Validate(); err != nil {
		return nil, err
	}

	return &groupDTO, nil
}
