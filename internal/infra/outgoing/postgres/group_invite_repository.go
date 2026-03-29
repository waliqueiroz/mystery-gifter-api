package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type groupInviteRepository struct {
	db DB
}

func NewGroupInviteRepository(db DB) domain.GroupInviteRepository {
	return &groupInviteRepository{
		db: db,
	}
}

func (r *groupInviteRepository) Create(ctx context.Context, groupInvite domain.GroupInvite) error {
	query, args, err := squirrel.Insert("group_invites").
		Columns("id", "group_id", "expires_at", "created_at").
		Values(groupInvite.ID, groupInvite.GroupID, groupInvite.ExpiresAt, groupInvite.CreatedAt).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building group invite insert query: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Println("error inserting group invite:", err)
		return fmt.Errorf("error inserting group invite: %w", err)
	}

	return nil
}

func (r *groupInviteRepository) GetActiveByGroupID(ctx context.Context, groupID string) (*domain.GroupInvite, error) {
	query, args, err := squirrel.Select("*").
		From("group_invites").
		Where(squirrel.And{
			squirrel.Eq{"group_id": groupID},
			squirrel.Expr("expires_at > NOW()"),
		}).
		OrderBy("created_at DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building group invite active select query: %w", err)
	}

	var groupInvite GroupInvite
	err = r.db.GetContext(ctx, &groupInvite, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("no active invite found for this group")
		}
		return nil, fmt.Errorf("error getting active group invite: %w", err)
	}

	return mapGroupInviteToDomain(groupInvite)
}

func (r *groupInviteRepository) GetByID(ctx context.Context, id string) (*domain.GroupInvite, error) {
	query, args, err := squirrel.Select("*").
		From("group_invites").
		Where(squirrel.Eq{"id": id}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("error building group invite select query: %w", err)
	}

	var groupInvite GroupInvite
	err = r.db.GetContext(ctx, &groupInvite, query, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewResourceNotFoundError("group invite not found")
		}
		return nil, fmt.Errorf("error getting group invite: %w", err)
	}

	return mapGroupInviteToDomain(groupInvite)
}
