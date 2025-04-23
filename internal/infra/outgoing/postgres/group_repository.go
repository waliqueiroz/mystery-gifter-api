package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type groupRepository struct {
	db DB
}

func NewGroupRepository(db DB) domain.GroupRepository {
	return &groupRepository{
		db: db,
	}
}

func (r *groupRepository) Create(ctx context.Context, group domain.Group) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	query, args, err := squirrel.Insert("groups").
		Columns("id", "name", "owner_id", "created_at", "updated_at").
		Values(group.ID, group.Name, group.OwnerID, group.CreatedAt, group.UpdatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building group insert query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("error inserting group:", err)

		var currentError *pq.Error
		if errors.As(err, &currentError) && currentError.Code.Name() == POSTGRES_UNIQUE_VIOLATION {
			return domain.NewConflictError("you already have a group with this name")
		}

		return fmt.Errorf("error inserting group: %w", err)
	}

	groupUsersInsert := squirrel.Insert("group_users").
		Columns("group_id", "user_id", "created_at").
		PlaceholderFormat(squirrel.Dollar)

	for _, user := range group.Users {
		groupUsersInsert = groupUsersInsert.Values(group.ID, user.ID, group.CreatedAt)
	}

	query, args, err = groupUsersInsert.ToSql()
	if err != nil {
		return fmt.Errorf("error building group_users insert query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("error inserting group users:", err)
		return fmt.Errorf("error inserting group users: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}
