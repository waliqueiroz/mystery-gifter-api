package application

//go:generate go run go.uber.org/mock/mockgen -destination mock_application/user_service.go . UserService

import (
	"context"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type UserService interface {
	Create(ctx context.Context, user domain.User) error
	GetByID(ctx context.Context, userID string) (*domain.User, error)
	Search(ctx context.Context, filters domain.UserFilters) (*domain.SearchResult[domain.User], error)
}

type userService struct {
	userRepository domain.UserRepository
}

func NewUserService(userRepository domain.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) Create(ctx context.Context, user domain.User) error {
	if err := user.Validate(); err != nil {
		return err
	}

	return s.userRepository.Create(ctx, user)
}

func (s *userService) GetByID(ctx context.Context, userID string) (*domain.User, error) {
	return s.userRepository.GetByID(ctx, userID)
}

func (s *userService) Search(ctx context.Context, filters domain.UserFilters) (*domain.SearchResult[domain.User], error) {
	if err := filters.Validate(); err != nil {
		return nil, err
	}

	return s.userRepository.Search(ctx, filters)
}
