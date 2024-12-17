package postgres

import (
	"context"
	"log"

	"github.com/Masterminds/squirrel"
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
	query, args, err := squirrel.Insert("users").
		Columns("id", "name", "surname", "email", "password", "created_at", "updated_at").
		Values(user.ID, user.Name, user.Surname, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("error creating user", err)
		return err
	}

	return nil
}
