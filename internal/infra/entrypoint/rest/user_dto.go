package rest

import (
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type CreateUserDTO struct {
	Name            string `json:"name" validate:"required"`
	Surname         string `json:"surname" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,eqfield=PasswordConfirm"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}

func (u *CreateUserDTO) Validate() error {
	if err := validator.Validate(u); err != nil {
		return err
	}
	return nil
}

func mapCreateUserDTOToUser(userDTO CreateUserDTO) (*domain.User, error) {
	if err := userDTO.Validate(); err != nil {
		return nil, err
	}

	user, err := domain.NewUser(userDTO.Name, userDTO.Surname, userDTO.Email, userDTO.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}
