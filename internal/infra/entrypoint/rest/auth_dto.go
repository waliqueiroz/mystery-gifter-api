package rest

import (
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

// CredentialsDTO represents the login credentials
// swagger:model CredentialsDTO
type CredentialsDTO struct {
	// Email address of the user
	// required: true
	// example: user@example.com
	Email string `json:"email" validate:"required,email"`

	// Password of the user
	// required: true
	// example: mypassword123
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

// AuthSessionDTO represents the authentication session response
// swagger:model AuthSessionDTO
type AuthSessionDTO struct {
	// User information
	// required: true
	User UserDTO `json:"user" validate:"required"`

	// JWT access token
	// required: true
	// example: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	AccessToken string `json:"access_token" validate:"required"`

	// Token type
	// required: true
	// example: Bearer
	TokenType string `json:"token_type" validate:"required"`

	// Token expiration time in seconds
	// required: true
	// example: 3600
	ExpiresIn int64 `json:"expires_in" validate:"required"`
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
