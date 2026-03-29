package postgres

import (
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type GroupInvite struct {
	ID        string    `db:"id"`
	GroupID   string    `db:"group_id"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

func mapGroupInviteToDomain(groupInvite GroupInvite) (*domain.GroupInvite, error) {
	domainGroupInvite := domain.GroupInvite{
		ID:        groupInvite.ID,
		GroupID:   groupInvite.GroupID,
		ExpiresAt: groupInvite.ExpiresAt,
		CreatedAt: groupInvite.CreatedAt,
	}

	if err := domainGroupInvite.Validate(); err != nil {
		return nil, err
	}

	return &domainGroupInvite, nil
}
