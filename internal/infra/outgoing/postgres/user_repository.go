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

func (r *userRepository) Search(ctx context.Context, filters domain.UserFilters) (*domain.SearchResult[domain.User], error) {
	baseQuery := squirrel.Select("*").
		From("users").
		PlaceholderFormat(squirrel.Dollar)

	baseQuery = r.applyUserFilters(baseQuery, filters)
	baseQuery = r.applySorting(baseQuery, filters)
	baseQuery = baseQuery.Limit(uint64(filters.Limit)).Offset(uint64(filters.Offset))

	query, args, err := baseQuery.ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building users search query: %w", err)
	}

	var users []User
	err = r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		log.Println("error searching users:", err)
		return nil, fmt.Errorf("error searching users: %w", err)
	}

	total, err := r.countUsers(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("error counting users: %w", err)
	}

	domainUsers, err := mapUsersToDomain(users)
	if err != nil {
		return nil, fmt.Errorf("error mapping users to domain: %w", err)
	}

	return domain.NewSearchResult(domainUsers, filters.Limit, filters.Offset, total)
}

func (r *userRepository) applyUserFilters(query squirrel.SelectBuilder, filters domain.UserFilters) squirrel.SelectBuilder {
	if filters.Name != nil && *filters.Name != "" {
		query = query.Where(squirrel.ILike{"name": "%" + *filters.Name + "%"})
	}

	if filters.Surname != nil && *filters.Surname != "" {
		query = query.Where(squirrel.ILike{"surname": "%" + *filters.Surname + "%"})
	}

	if filters.Email != nil && *filters.Email != "" {
		query = query.Where(squirrel.ILike{"email": "%" + *filters.Email + "%"})
	}

	return query
}

func (r *userRepository) applySorting(query squirrel.SelectBuilder, filters domain.UserFilters) squirrel.SelectBuilder {
	orderBy := fmt.Sprintf("%s %s", filters.SortBy, filters.SortDirection)
	return query.OrderBy(orderBy)
}

func (r *userRepository) countUsers(ctx context.Context, filters domain.UserFilters) (int, error) {
	countQuery := squirrel.Select("COUNT(*)").
		From("users").
		PlaceholderFormat(squirrel.Dollar)

	countQuery = r.applyUserFilters(countQuery, filters)

	query, args, err := countQuery.ToSql()
	if err != nil {
		return 0, fmt.Errorf("error building users count query: %w", err)
	}

	var total int
	err = r.db.GetContext(ctx, &total, query, args...)
	if err != nil {
		log.Println("error counting users:", err)
		return 0, fmt.Errorf("error counting users: %w", err)
	}

	return total, nil
}
