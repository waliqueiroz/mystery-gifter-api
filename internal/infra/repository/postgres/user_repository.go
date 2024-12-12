package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user domain.User) error {
	return nil
}
