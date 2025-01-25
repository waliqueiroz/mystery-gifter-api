package build_rest

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

type AuthSessionDTOBuilder struct {
	authSessionDTO rest.AuthSessionDTO
}

func NewAuthSessionDTOBuilder() *AuthSessionDTOBuilder {
	return &AuthSessionDTOBuilder{
		authSessionDTO: rest.AuthSessionDTO{
			User: rest.UserDTO{
				ID:        uuid.New().String(),
				Name:      "DefaultName",
				Surname:   "DefaultSurname",
				Email:     "default@example.com",
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			AccessToken: "DefaultAccessToken",
			TokenType:   "Bearer",
			ExpiresIn:   time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours by default
		},
	}
}

func (b *AuthSessionDTOBuilder) WithUser(user rest.UserDTO) *AuthSessionDTOBuilder {
	b.authSessionDTO.User = user
	return b
}

func (b *AuthSessionDTOBuilder) WithAccessToken(accessToken string) *AuthSessionDTOBuilder {
	b.authSessionDTO.AccessToken = accessToken
	return b
}

func (b *AuthSessionDTOBuilder) WithTokenType(tokenType string) *AuthSessionDTOBuilder {
	b.authSessionDTO.TokenType = tokenType
	return b
}

func (b *AuthSessionDTOBuilder) WithExpiresIn(expiresIn int64) *AuthSessionDTOBuilder {
	b.authSessionDTO.ExpiresIn = expiresIn
	return b
}

func (b *AuthSessionDTOBuilder) Build() rest.AuthSessionDTO {
	return b.authSessionDTO
}
