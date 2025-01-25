package rest

import (
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type CredentialsDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (c *CredentialsDTO) Validate() error {
	if errs := validator.Validate(c); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

func mapCredentialsToDomain(credentialsDTO CredentialsDTO) (*domain.Credentials, error) {
	if err := credentialsDTO.Validate(); err != nil {
		return nil, err
	}

	credentials, err := domain.NewCredentials(credentialsDTO.Email, credentialsDTO.Password)
	if err != nil {
		return nil, err
	}

	return credentials, nil
}

type AuthSessionDTO struct {
	User        UserDTO `json:"user" validate:"required"`
	AccessToken string  `json:"access_token" validate:"required"`
	TokenType   string  `json:"token_type" validate:"required"`
	ExpiresIn   int64   `json:"expires_in" validate:"required"`
}

func (a *AuthSessionDTO) Validate() error {
	if errs := validator.Validate(a); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

func mapAuthSessionFromDomain(authSession domain.AuthSession) (*AuthSessionDTO, error) {
	userDTO, err := mapUserFromDomain(authSession.User)
	if err != nil {
		return nil, err
	}

	authSessionDTO := AuthSessionDTO{
		User:        *userDTO,
		AccessToken: authSession.AccessToken,
		TokenType:   authSession.TokenType,
		ExpiresIn:   authSession.ExpiresIn,
	}

	if err := authSessionDTO.Validate(); err != nil {
		return nil, err
	}

	return &authSessionDTO, nil
}
