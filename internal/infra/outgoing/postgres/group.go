package postgres

import (
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type Group struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	OwnerID   string    `db:"owner_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func mapGroupToDomain(group Group, groupUsers []User) (*domain.Group, error) {
	domainUsers, err := mapUsersToDomain(groupUsers)
	if err != nil {
		return nil, err
	}

	domainGroup := domain.Group{
		ID:        group.ID,
		Name:      group.Name,
		OwnerID:   group.OwnerID,
		Users:     domainUsers,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}

	if err := domainGroup.Validate(); err != nil {
		return nil, err
	}

	return &domainGroup, nil
}
