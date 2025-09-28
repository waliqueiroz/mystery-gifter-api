package rest

import (
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type CreateUserDTO struct {
	Name            string `json:"name" validate:"required"`
	Surname         string `json:"surname" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8,eqfield=PasswordConfirm"`
	PasswordConfirm string `json:"password_confirm" validate:"required"`
}

func (u *CreateUserDTO) Validate() error {
	if errs := validator.Validate(u); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

func mapCreateUserDTOToDomain(identity domain.IdentityGenerator, passwordManager domain.PasswordManager, userDTO CreateUserDTO) (*domain.User, error) {
	if err := userDTO.Validate(); err != nil {
		return nil, err
	}

	user, err := domain.NewUser(identity, passwordManager, userDTO.Name, userDTO.Surname, userDTO.Email, userDTO.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

type UserDTO struct {
	ID        string    `json:"id" validate:"required"`
	Name      string    `json:"name" validate:"required"`
	Surname   string    `json:"surname" validate:"required"`
	Email     string    `json:"email" validate:"required,email"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

func (u *UserDTO) Validate() error {
	if errs := validator.Validate(u); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

func mapUserFromDomain(user domain.User) (*UserDTO, error) {
	userDTO := UserDTO{
		ID:        user.ID,
		Name:      user.Name,
		Surname:   user.Surname,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if err := userDTO.Validate(); err != nil {
		return nil, err
	}

	return &userDTO, nil
}

func mapUsersFromDomain(users []domain.User) ([]UserDTO, error) {
	userDTOs := make([]UserDTO, 0, len(users))

	for _, user := range users {
		userDTO, err := mapUserFromDomain(user)
		if err != nil {
			return nil, err
		}
		userDTOs = append(userDTOs, *userDTO)
	}

	return userDTOs, nil
}

type AddUserDTO struct {
	UserID string `json:"user_id" validate:"required,uuid"`
}

func (a *AddUserDTO) Validate() error {
	if errs := validator.Validate(a); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

type UserFiltersDTO struct {
	Name          string `query:"name" json:"name"`
	Surname       string `query:"surname" json:"surname"`
	Email         string `query:"email" json:"email"`
	Limit         int    `query:"limit" json:"limit"`
	Offset        int    `query:"offset" json:"offset"`
	SortDirection string `query:"sort_direction" json:"sort_direction" validate:"omitempty,oneof=ASC DESC"`
	SortBy        string `query:"sort_by" json:"sort_by" validate:"omitempty,oneof=name surname email created_at updated_at"`
}

func (s *UserFiltersDTO) Validate() error {
	if errs := validator.Validate(s); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

func mapUserFiltersDTOToDomain(filtersDTO UserFiltersDTO) (*domain.UserFilters, error) {
	if err := filtersDTO.Validate(); err != nil {
		return nil, err
	}

	return domain.NewUserFilters(filtersDTO.Name, filtersDTO.Surname, filtersDTO.Email, filtersDTO.Limit, filtersDTO.Offset, domain.SortDirectionType(filtersDTO.SortDirection), filtersDTO.SortBy)
}
