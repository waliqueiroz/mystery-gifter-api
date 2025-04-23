package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type userRepository struct {
	db DB
}

func NewUserRepository(db DB) domain.UserRepository {
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
		return fmt.Errorf("error building users insert query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("error creating user:", err)

		var currentError *pq.Error
		if errors.As(err, &currentError) && currentError.Code.Name() == POSTGRES_UNIQUE_VIOLATION {
			return domain.NewConflictError("the email is already registered")
		}

		return fmt.Errorf("error inserting user: %w", err)
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
		return nil, fmt.Errorf("error building users select query: %w", err)
	}

	var user User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("user not found")
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return mapUserToDomain(user)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query, args, err := squirrel.Select("*").
		From("users").
		Where(squirrel.Eq{"email": email}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building users select query: %w", err)
	}

	var user User
	err = r.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("user not found")
		}
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}

	return mapUserToDomain(user)
}
