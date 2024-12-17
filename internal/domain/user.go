package domain

import (
	"context"
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/pkg/identity"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/security"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type UserRepository interface {
	Create(ctx context.Context, user User) error
	// GetByID(ctx context.Context, userID string)
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

func NewUser(name, surname, email, password string) (*User, error) {
	hashedPassword, err := security.Hash(password)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	user := User{
		ID:        identity.NewUUID(),
		Name:      name,
		Surname:   surname,
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return &user, err

}

func (u *User) Validate() error {
	if err := validator.Validate(u); err != nil {
		return err
	}
	return nil
}
