package domain

//go:generate go run go.uber.org/mock/mockgen -destination mock_domain/group_repository.go . GroupRepository

import (
	"context"
	"math/rand"
	"time"

	"slices"

	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type GroupStatus string

const (
	GroupStatusOpen     GroupStatus = "OPEN"
	GroupStatusMatched  GroupStatus = "MATCHED"
	GroupStatusArchived GroupStatus = "ARCHIVED"
)

type GroupRepository interface {
	Create(ctx context.Context, group Group) error
	Update(ctx context.Context, group Group) error
	GetByID(ctx context.Context, groupID string) (*Group, error)
}

type Group struct {
	ID        string      `validate:"required,uuid"`
	Name      string      `validate:"required"`
	Users     []User      `validate:"required,min=1"`
	OwnerID   string      `validate:"required,uuid"`
	Matches   []Match     `validate:"dive,omitempty"`
	Status    GroupStatus `validate:"required,oneof=OPEN MATCHED ARCHIVED"`
	CreatedAt time.Time   `validate:"required"`
	UpdatedAt time.Time   `validate:"required"`
}

type Match struct {
	GiverID    string `validate:"required,uuid"`
	ReceiverID string `validate:"required,uuid"`
}

func (m *Match) Validate() error {
	if errs := validator.Validate(m); len(errs) > 0 {
		return NewValidationError(errs)
	}

	return nil
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
		Status:    GroupStatusOpen,
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

func (g *Group) IsOpen() bool {
	return g.Status == GroupStatusOpen
}

func (g *Group) IsMatched() bool {
	return g.Status == GroupStatusMatched
}

func (g *Group) IsArchived() bool {
	return g.Status == GroupStatusArchived
}

func (g *Group) AddUser(requesterID string, targetUser User) error {
	if !g.IsOpen() {
		return NewConflictError("group is not open for registration, contact the group owner to reopen the group")
	}

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
	if !g.IsOpen() {
		return NewConflictError("group is not open for removal, contact the group owner to reopen the group")
	}

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

func (g *Group) GenerateMatches(requesterID string) error {
	if requesterID != g.OwnerID {
		return NewForbiddenError("only the group owner can generate matches")
	}

	if !g.IsOpen() {
		return NewConflictError("group is not open for matches")
	}

	if len(g.Users) < 3 {
		return NewConflictError("group must have at least 3 users to generate matches")
	}

	userIDs := make([]string, len(g.Users))
	for i, user := range g.Users {
		userIDs[i] = user.ID
	}

	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	r.Shuffle(len(userIDs), func(i, j int) {
		userIDs[i], userIDs[j] = userIDs[j], userIDs[i]
	})

	currentMatches := make([]Match, len(userIDs))
	for i := range userIDs {
		giverID := userIDs[i]
		receiverID := userIDs[(i+1)%len(userIDs)]

		currentMatches[i] = Match{
			GiverID:    giverID,
			ReceiverID: receiverID,
		}
	}

	g.Matches = currentMatches
	g.Status = GroupStatusMatched
	g.UpdatedAt = time.Now()

	return g.Validate()
}

func (g *Group) Reopen(requesterID string) error {
	if requesterID != g.OwnerID {
		return NewForbiddenError("only the group owner can reopen the group")
	}

	if g.IsArchived() {
		return NewConflictError("group is archived and cannot be reopened")
	}

	if g.Status == GroupStatusOpen {
		return NewConflictError("group is already open")
	}

	g.Matches = []Match{}
	g.Status = GroupStatusOpen
	g.UpdatedAt = time.Now()

	return g.Validate()
}

func (g *Group) Archive(requesterID string) error {
	if requesterID != g.OwnerID {
		return NewForbiddenError("only the group owner can archive the group")
	}

	if g.Status == GroupStatusArchived {
		return NewConflictError("group is already archived")
	}

	g.Status = GroupStatusArchived
	g.UpdatedAt = time.Now()

	return g.Validate()
}

func (g *Group) GetUserMatch(requesterID string) (*User, error) {
	if !g.IsMatched() {
		return nil, NewConflictError("group is not matched")
	}

	var userMatch *Match
	for _, match := range g.Matches {
		if match.GiverID == requesterID {
			userMatch = &match
			break
		}
	}

	if userMatch == nil {
		return nil, NewConflictError("match not found for the given user")
	}

	for _, user := range g.Users {
		if user.ID == userMatch.ReceiverID {
			return &user, nil
		}
	}

	return nil, NewConflictError("receiver user not found for the identified match")
}
