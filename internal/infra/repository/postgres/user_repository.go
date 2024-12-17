package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

const POSTGRES_UNIQUE_VIOLATION = "unique_violation"

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

		var currentError *pq.Error
		if errors.As(err, &currentError) && currentError.Code.Name() == POSTGRES_UNIQUE_VIOLATION {
			return domain.NewConflictError("the email is already registered")
		}

		return err
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	query, args, err := squirrel.Select("*").
		From("users").
		Where(squirrel.Eq{"id": userID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	var user User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("user not found")
		}
		return nil, err
	}

	return mapUserToDomain(user)
}
