package rest

import (
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

// GroupInviteDTO represents an invite link for joining a group
// swagger:model GroupInviteDTO
type GroupInviteDTO struct {
	// Unique invite identifier (UUID v7). Share this to invite users.
	// required: true
	// example: 018e1234-abcd-7000-8000-000000000001
	ID string `json:"id"`

	// ID of the group the invite grants access to
	// required: true
	// example: 018e1234-abcd-7000-8000-000000000002
	GroupID string `json:"group_id"`

	// When the invite expires (UTC)
	// required: true
	ExpiresAt time.Time `json:"expires_at"`

	// When the invite was created (UTC)
	// required: true
	CreatedAt time.Time `json:"created_at"`
}

func mapGroupInviteFromDomain(groupInvite domain.GroupInvite) (*GroupInviteDTO, error) {
	dto := &GroupInviteDTO{
		ID:        groupInvite.ID,
		GroupID:   groupInvite.GroupID,
		ExpiresAt: groupInvite.ExpiresAt,
		CreatedAt: groupInvite.CreatedAt,
	}

	return dto, nil
}
