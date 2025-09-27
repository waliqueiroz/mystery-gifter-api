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
		Columns("id", "name", "status", "owner_id", "created_at", "updated_at").
		Values(group.ID, group.Name, group.Status, group.OwnerID, group.CreatedAt, group.UpdatedAt).
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

	if len(group.Matches) > 0 {
		groupMatchesInsert := squirrel.Insert("group_matches").
			Columns("group_id", "giver_id", "receiver_id", "created_at").
			PlaceholderFormat(squirrel.Dollar)

		for _, match := range group.Matches {
			groupMatchesInsert = groupMatchesInsert.Values(
				group.ID, match.GiverID, match.ReceiverID, group.CreatedAt,
			)
		}

		query, args, err = groupMatchesInsert.ToSql()
		if err != nil {
			return fmt.Errorf("error building group_matches insert query: %w", err)
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Println("error inserting group matches:", err)
			return fmt.Errorf("error inserting group matches: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (r *groupRepository) Update(ctx context.Context, group domain.Group) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	defer tx.Rollback()

	query, args, err := squirrel.Update("groups").
		Set("name", group.Name).
		Set("status", group.Status).
		Set("updated_at", group.UpdatedAt).
		Where(squirrel.Eq{"id": group.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building group update query: %w", err)
	}

	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("error updating group:", err)

		var currentError *pq.Error
		if errors.As(err, &currentError) && currentError.Code.Name() == POSTGRES_UNIQUE_VIOLATION {
			return domain.NewConflictError("you already have a group with this name")
		}

		return fmt.Errorf("error updating group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.NewResourceNotFoundError("group not found")
	}

	// Remove existing group users
	query, args, err = squirrel.Delete("group_users").
		Where(squirrel.Eq{"group_id": group.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building group_users delete query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error deleting group users: %w", err)
	}

	// Insert new group users
	groupUsersInsert := squirrel.Insert("group_users").
		Columns("group_id", "user_id", "created_at").
		PlaceholderFormat(squirrel.Dollar)

	for _, user := range group.Users {
		groupUsersInsert = groupUsersInsert.Values(group.ID, user.ID, group.UpdatedAt)
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

	// Remove existing group matches
	query, args, err = squirrel.Delete("group_matches").
		Where(squirrel.Eq{"group_id": group.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building group_matches delete query: %w", err)
	}

	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("error deleting group matches: %w", err)
	}

	// Insert new group matches if any
	if len(group.Matches) > 0 {
		groupMatchesInsert := squirrel.Insert("group_matches").
			Columns("group_id", "giver_id", "receiver_id", "created_at").
			PlaceholderFormat(squirrel.Dollar)

		for _, match := range group.Matches {
			groupMatchesInsert = groupMatchesInsert.Values(
				group.ID, match.GiverID, match.ReceiverID, group.UpdatedAt,
			)
		}

		query, args, err = groupMatchesInsert.ToSql()
		if err != nil {
			return fmt.Errorf("error building group_matches insert query: %w", err)
		}

		_, err = tx.ExecContext(ctx, query, args...)
		if err != nil {
			log.Println("error inserting group matches:", err)
			return fmt.Errorf("error inserting group matches: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func (r *groupRepository) GetByID(ctx context.Context, groupID string) (*domain.Group, error) {
	query, args, err := squirrel.Select("g.*").
		From("groups g").
		Where(squirrel.Eq{"g.id": groupID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building group select query: %w", err)
	}

	var group Group
	err = r.db.GetContext(ctx, &group, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("group not found")
		}
		return nil, fmt.Errorf("error getting group: %w", err)
	}

	// Get group users
	query, args, err = squirrel.Select("u.*").
		From("users u").
		Join("group_users gu ON gu.user_id = u.id").
		Where(squirrel.Eq{"gu.group_id": groupID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building group users select query: %w", err)
	}

	var users []User
	err = r.db.SelectContext(ctx, &users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting group users: %w", err)
	}

	// Get group matches
	query, args, err = squirrel.Select("giver_id", "receiver_id").
		From("group_matches").
		Where(squirrel.Eq{"group_id": groupID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building group matches select query: %w", err)
	}

	var matches []Match
	err = r.db.SelectContext(ctx, &matches, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error getting group matches: %w", err)
	}

	domainGroup, err := mapGroupToDomain(group, users, matches)
	if err != nil {
		return nil, err
	}

	return domainGroup, nil
}
