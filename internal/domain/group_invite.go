package domain

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/group_invite_repository.go . GroupInviteRepository

import (
	"context"
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type GroupInviteRepository interface {
	Create(ctx context.Context, groupInvite GroupInvite) error
	GetByID(ctx context.Context, id string) (*GroupInvite, error)
}

type GroupInvite struct {
	ID        string    `validate:"required,uuid"`
	GroupID   string    `validate:"required,uuid"`
	ExpiresAt time.Time `validate:"required"`
	CreatedAt time.Time `validate:"required"`
}

func NewGroupInvite(identityGenerator IdentityGenerator, groupID string, expiration time.Duration) (*GroupInvite, error) {
	id, err := identityGenerator.Generate()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	groupInvite := &GroupInvite{
		ID:        id,
		GroupID:   groupID,
		ExpiresAt: now.Add(expiration),
		CreatedAt: now,
	}

	if err := groupInvite.Validate(); err != nil {
		return nil, err
	}

	return groupInvite, nil
}

func (i *GroupInvite) Validate() error {
	if errs := validator.Validate(i); len(errs) > 0 {
		return NewValidationError(errs)
	}

	return nil
}

func (i *GroupInvite) IsExpired() bool {
	return time.Now().After(i.ExpiresAt)
}
