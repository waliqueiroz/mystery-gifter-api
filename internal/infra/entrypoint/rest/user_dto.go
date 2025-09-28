package rest

import (
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

// CreateUserDTO represents the data needed to create a new user
// swagger:model CreateUserDTO
type CreateUserDTO struct {
	// User's first name
	// required: true
	// example: João
	Name string `json:"name" validate:"required"`

	// User's last name
	// required: true
	// example: Silva
	Surname string `json:"surname" validate:"required"`

	// User's email address
	// required: true
	// example: joao.silva@example.com
	Email string `json:"email" validate:"required,email"`

	// User's password
	// required: true
	// minLength: 8
	// example: mypassword123
	Password string `json:"password" validate:"required,min=8,eqfield=PasswordConfirm"`

	// Password confirmation (must match password)
	// required: true
	// example: mypassword123
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

// UserDTO represents a user in the system
// swagger:model UserDTO
type UserDTO struct {
	// Unique identifier for the user
	// required: true
	// example: 01234567-89ab-cdef-0123-456789abcdef
	ID string `json:"id" validate:"required"`

	// User's first name
	// required: true
	// example: João
	Name string `json:"name" validate:"required"`

	// User's last name
	// required: true
	// example: Silva
	Surname string `json:"surname" validate:"required"`

	// User's email address
	// required: true
	// example: joao.silva@example.com
	Email string `json:"email" validate:"required,email"`

	// When the user was created
	// required: true
	// example: 2023-12-01T10:00:00Z
	CreatedAt time.Time `json:"created_at" validate:"required"`

	// When the user was last updated
	// required: true
	// example: 2023-12-01T10:00:00Z
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

// AddUserDTO represents the data needed to add a user to a group
// swagger:model AddUserDTO
type AddUserDTO struct {
	// ID of the user to add to the group
	// required: true
	// example: 01234567-89ab-cdef-0123-456789abcdef
	UserID string `json:"user_id" validate:"required,uuid"`
}

func (a *AddUserDTO) Validate() error {
	if errs := validator.Validate(a); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

// UserFiltersDTO represents the filters for searching users
// swagger:model UserFiltersDTO
type UserFiltersDTO struct {
	// Filter by user's first name
	// example: João
	Name string `query:"name" json:"name"`

	// Filter by user's last name
	// example: Silva
	Surname string `query:"surname" json:"surname"`

	// Filter by user's email
	// example: joao@example.com
	Email string `query:"email" json:"email"`

	// Maximum number of results to return
	// example: 10
	Limit int `query:"limit" json:"limit"`

	// Number of results to skip
	// example: 0
	Offset int `query:"offset" json:"offset"`

	// Sort direction (ASC or DESC)
	// enum: ASC,DESC
	// example: ASC
	SortDirection string `query:"sort_direction" json:"sort_direction" validate:"omitempty,oneof=ASC DESC"`

	// Field to sort by
	// enum: name,surname,email,created_at,updated_at
	// example: name
	SortBy string `query:"sort_by" json:"sort_by" validate:"omitempty,oneof=name surname email created_at updated_at"`
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
