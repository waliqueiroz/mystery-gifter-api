package application

import (
	"context"
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type AuthService interface {
	Login(ctx context.Context, credentials domain.Credentials) (*domain.AuthSession, error)
}

type authService struct {
	sessionDuration time.Duration
	userRepository  domain.UserRepository
	passwordManager domain.PasswordManager
	tokenManager    domain.TokenManager
}

func NewAuthService(sessionDuration time.Duration, userRepository domain.UserRepository, passwordManager domain.PasswordManager, tokenManager domain.TokenManager) AuthService {
	return &authService{
		sessionDuration: sessionDuration,
		userRepository:  userRepository,
		passwordManager: passwordManager,
		tokenManager:    tokenManager,
	}
}

func (s *authService) Login(ctx context.Context, credentials domain.Credentials) (*domain.AuthSession, error) {
	user, err := s.userRepository.GetByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	if err := s.passwordManager.Compare(user.Password, credentials.Password); err != nil {
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	expiresIn := time.Now().Add(s.sessionDuration).Unix()

	token, err := s.tokenManager.Create(user.ID, expiresIn)
	if err != nil {
		return nil, err
	}

	return domain.NewAuthSession(*user, token, s.tokenManager.GetTokenType(), expiresIn)
}
