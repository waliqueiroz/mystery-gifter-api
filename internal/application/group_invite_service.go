package application

//go:generate go run go.uber.org/mock/mockgen -destination mock_application/group_invite_service.go . GroupInviteService

import (
	"context"
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type GroupInviteService interface {
	Create(ctx context.Context, groupID, requesterID string) (*domain.GroupInvite, error)
	JoinGroup(ctx context.Context, inviteID, userID string) (*domain.Group, error)
}

type groupInviteService struct {
	groupInviteRepository domain.GroupInviteRepository
	groupRepository       domain.GroupRepository
	userService           UserService
	identityGenerator     domain.IdentityGenerator
	linkExpiration        time.Duration
}

func NewGroupInviteService(
	groupInviteRepository domain.GroupInviteRepository,
	groupRepository domain.GroupRepository,
	userService UserService,
	identityGenerator domain.IdentityGenerator,
	linkExpiration time.Duration,
) GroupInviteService {
	return &groupInviteService{
		groupInviteRepository: groupInviteRepository,
		groupRepository:       groupRepository,
		userService:           userService,
		identityGenerator:     identityGenerator,
		linkExpiration:        linkExpiration,
	}
}

func (s *groupInviteService) Create(ctx context.Context, groupID, requesterID string) (*domain.GroupInvite, error) {
	group, err := s.groupRepository.GetByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if err := group.CanCreateInvite(requesterID); err != nil {
		return nil, err
	}

	groupInvite, err := domain.NewGroupInvite(s.identityGenerator, groupID, s.linkExpiration)
	if err != nil {
		return nil, err
	}

	if err := s.groupInviteRepository.Create(ctx, *groupInvite); err != nil {
		return nil, err
	}

	return groupInvite, nil
}

func (s *groupInviteService) JoinGroup(ctx context.Context, inviteID, userID string) (*domain.Group, error) {
	groupInvite, err := s.groupInviteRepository.GetByID(ctx, inviteID)
	if err != nil {
		return nil, err
	}

	if groupInvite.IsExpired() {
		return nil, domain.NewConflictError("invite has expired")
	}

	group, err := s.groupRepository.GetByID(ctx, groupInvite.GroupID)
	if err != nil {
		return nil, err
	}

	targetUser, err := s.userService.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := group.AddUser(group.OwnerID, *targetUser); err != nil {
		return nil, err
	}

	if err := s.groupRepository.Update(ctx, *group); err != nil {
		return nil, err
	}

	return group, nil
}
