package application

//go:generate go run go.uber.org/mock/mockgen -destination mock_application/group_service.go . GroupService

import (
	"context"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type GroupService interface {
	Create(ctx context.Context, name, ownerID string) (*domain.Group, error)
	GetByID(ctx context.Context, groupID string) (*domain.Group, error)
	AddUser(ctx context.Context, groupID, requesterID, targetUserID string) (*domain.Group, error)
	RemoveUser(ctx context.Context, groupID, requesterID, targetUserID string) (*domain.Group, error)
}

type groupService struct {
	groupRepository   domain.GroupRepository
	userService       UserService
	identityGenerator domain.IdentityGenerator
}

func NewGroupService(
	groupRepository domain.GroupRepository,
	userService UserService,
	identityGenerator domain.IdentityGenerator,
) GroupService {
	return &groupService{
		groupRepository:   groupRepository,
		userService:       userService,
		identityGenerator: identityGenerator,
	}
}

func (s *groupService) Create(ctx context.Context, name, ownerID string) (*domain.Group, error) {
	owner, err := s.userService.GetByID(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	group, err := domain.NewGroup(s.identityGenerator, name, *owner)
	if err != nil {
		return nil, err
	}

	err = s.groupRepository.Create(ctx, *group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *groupService) GetByID(ctx context.Context, groupID string) (*domain.Group, error) {
	group, err := s.groupRepository.GetByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *groupService) AddUser(ctx context.Context, groupID, requesterID, targetUserID string) (*domain.Group, error) {
	group, err := s.groupRepository.GetByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	targetUser, err := s.userService.GetByID(ctx, targetUserID)
	if err != nil {
		return nil, err
	}

	if err := group.AddUser(requesterID, *targetUser); err != nil {
		return nil, err
	}

	err = s.groupRepository.Update(ctx, *group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (s *groupService) RemoveUser(ctx context.Context, groupID, requesterID, targetUserID string) (*domain.Group, error) {
	group, err := s.groupRepository.GetByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if err := group.RemoveUser(requesterID, targetUserID); err != nil {
		return nil, err
	}

	err = s.groupRepository.Update(ctx, *group)
	if err != nil {
		return nil, err
	}

	return group, nil
}
