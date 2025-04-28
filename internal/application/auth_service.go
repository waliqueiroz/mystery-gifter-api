package application

//go:generate go run go.uber.org/mock/mockgen -destination mock_application/auth_service.go . AuthService

import (
	"context"
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type AuthService interface {
	Login(ctx context.Context, credentials domain.Credentials) (*domain.AuthSession, error)
}

type authService struct {
	sessionDuration  time.Duration
	userRepository   domain.UserRepository
	passwordManager  domain.PasswordManager
	authTokenManager domain.AuthTokenManager
}

func NewAuthService(sessionDuration time.Duration, userRepository domain.UserRepository, passwordManager domain.PasswordManager, authTokenManager domain.AuthTokenManager) AuthService {
	return &authService{
		sessionDuration:  sessionDuration,
		userRepository:   userRepository,
		passwordManager:  passwordManager,
		authTokenManager: authTokenManager,
	}
}

func (s *authService) Login(ctx context.Context, credentials domain.Credentials) (*domain.AuthSession, error) {
	if err := credentials.Validate(); err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetByEmail(ctx, credentials.Email)
	if err != nil {
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	if err := s.passwordManager.Compare(user.Password, credentials.Password); err != nil {
		return nil, domain.NewUnauthorizedError("invalid credentials")
	}

	expiresIn := time.Now().Add(s.sessionDuration).Unix()

	token, err := s.authTokenManager.Create(user.ID, expiresIn)
	if err != nil {
		return nil, err
	}

	return domain.NewAuthSession(*user, token, s.authTokenManager.GetTokenType(), expiresIn)
}
