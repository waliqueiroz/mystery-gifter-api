package build_domain

import (
	"time"

	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type AuthSessionBuilder struct {
	authSession domain.AuthSession
}

func NewAuthSessionBuilder() *AuthSessionBuilder {
	return &AuthSessionBuilder{
		authSession: domain.AuthSession{
			User: domain.User{
				ID:        uuid.New().String(),
				Name:      "DefaultName",
				Surname:   "DefaultSurname",
				Email:     "default@example.com",
				Password:  "defaultpassword",
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			AccessToken: "DefaultAccessToken",
			TokenType:   "Bearer",
			ExpiresIn:   time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours by default
		},
	}
}

func (b *AuthSessionBuilder) WithUser(user domain.User) *AuthSessionBuilder {
	b.authSession.User = user
	return b
}

func (b *AuthSessionBuilder) WithAccessToken(accessToken string) *AuthSessionBuilder {
	b.authSession.AccessToken = accessToken
	return b
}

func (b *AuthSessionBuilder) WithTokenType(tokenType string) *AuthSessionBuilder {
	b.authSession.TokenType = tokenType
	return b
}

func (b *AuthSessionBuilder) WithExpiresIn(expiresIn int64) *AuthSessionBuilder {
	b.authSession.ExpiresIn = expiresIn
	return b
}

func (b *AuthSessionBuilder) Build() domain.AuthSession {
	return b.authSession
}
