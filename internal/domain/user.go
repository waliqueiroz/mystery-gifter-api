package domain

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/user_repository.go . UserRepository

import (
	"context"
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type UserRepository interface {
	Search(ctx context.Context, filters UserFilters) (*SearchResult[User], error)
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, userID string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	// Update(ctx context.Context, user User) error
}

type User struct {
	ID        string    `validate:"required,uuid"`
	Name      string    `validate:"required"`
	Surname   string    `validate:"required"`
	Email     string    `validate:"required,email"`
	Password  string    `validate:"required"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
}

func NewUser(identity IdentityGenerator, passwordManager PasswordManager, name, surname, email, password string) (*User, error) {
	hashedPassword, err := passwordManager.Hash(password)
	if err != nil {
		return nil, err
	}

	id, err := identity.Generate()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	user := User{
		ID:        id,
		Name:      name,
		Surname:   surname,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return &user, err
}

func (u *User) Validate() error {
	if errs := validator.Validate(u); len(errs) > 0 {
		return NewValidationError(errs)
	}
	return nil
}

type SortDirectionType string

const (
	SortDirectionTypeAsc  SortDirectionType = "ASC"
	SortDirectionTypeDesc SortDirectionType = "DESC"
)

type UserFilters struct {
	Name          *string
	Surname       *string
	Email         *string
	Limit         int               `validate:"required,min=1"`
	Offset        int               `validate:"min=0"`
	SortDirection SortDirectionType `validate:"required,oneof=ASC DESC"`
	SortBy        string            `validate:"required,oneof=name surname email created_at updated_at"`
}

func NewUserFilters(name, surname, email string, limit, offset int, sortDirection SortDirectionType, sortBy string) (*UserFilters, error) {
	userFilters := UserFilters{
		Name:          &name,
		Surname:       &surname,
		Email:         &email,
		Limit:         limit,
		Offset:        offset,
		SortDirection: sortDirection,
		SortBy:        sortBy,
	}

	if err := userFilters.Validate(); err != nil {
		return nil, err
	}

	return &userFilters, nil
}

func (u *UserFilters) Validate() error {
	if errs := validator.Validate(u); len(errs) > 0 {
		return NewValidationError(errs)
	}
	return nil
}
