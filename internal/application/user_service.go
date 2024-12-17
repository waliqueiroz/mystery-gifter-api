package application

import (
	"context"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type UserService interface {
	Create(ctx context.Context, user domain.User) error
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
	return s.userRepository.Create(ctx, user)
}
