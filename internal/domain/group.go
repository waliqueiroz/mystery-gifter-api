package domain

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/group_repository.go . GroupRepository

import (
	"context"
	"time"

	"slices"

	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type GroupRepository interface {
	Create(ctx context.Context, group Group) error
	Update(ctx context.Context, group Group) error
	GetByID(ctx context.Context, groupID string) (*Group, error)
}

type Group struct {
	ID        string    `validate:"required,uuid"`
	Name      string    `validate:"required"`
	Users     []User    `validate:"required,min=1"`
	OwnerID   string    `validate:"required,uuid"`
	CreatedAt time.Time `validate:"required"`
	UpdatedAt time.Time `validate:"required"`
}

func NewGroup(identityGenerator IdentityGenerator, name string, owner User) (*Group, error) {
	id, err := identityGenerator.Generate()
	if err != nil {
		return nil, err
	}

	now := time.Now()

	group := &Group{
		ID:        id,
		Name:      name,
		OwnerID:   owner.ID,
		Users:     []User{owner},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := group.Validate(); err != nil {
		return nil, err
	}

	return group, nil
}

func (g *Group) Validate() error {
	if errs := validator.Validate(g); len(errs) > 0 {
		return NewValidationError(errs)
	}

	return nil
}

func (g *Group) AddUser(requesterID string, targetUser User) error {
	if requesterID != g.OwnerID && requesterID != targetUser.ID {
		return NewForbiddenError("only the group owner can add other users")
	}

	for _, existingUser := range g.Users {
		if existingUser.ID == targetUser.ID {
			return nil
		}
	}

	g.Users = append(g.Users, targetUser)
	g.UpdatedAt = time.Now()

	return g.Validate()
}

func (g *Group) RemoveUser(requesterID, targetUserID string) error {
	if requesterID != g.OwnerID && requesterID != targetUserID {
		return NewForbiddenError("only the group owner can remove other users")
	}

	if g.OwnerID == targetUserID {
		return NewForbiddenError("cannot remove group owner")
	}

	for i, user := range g.Users {
		if user.ID == targetUserID {
			g.Users = slices.Delete(g.Users, i, i+1)
			g.UpdatedAt = time.Now()
			break
		}
	}

	return g.Validate()
}
