package domain

import (
	"context"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user User) error
	// GetByID(ctx context.Context, userID string)
	// Update(ctx context.Context, user User) error
}

type User struct {
	ID        string
	Name      string
	Surname   string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
