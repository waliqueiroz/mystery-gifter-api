package build_domain

import (
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type CredentialsBuilder struct {
	credentials domain.Credentials
}

func NewCredentialsBuilder() *CredentialsBuilder {
	return &CredentialsBuilder{
		credentials: domain.Credentials{
			Email:    "default@example.com",
			Password: "defaultpassword",
		},
	}
}

func (b *CredentialsBuilder) WithEmail(email string) *CredentialsBuilder {
	b.credentials.Email = email
	return b
}

func (b *CredentialsBuilder) WithPassword(password string) *CredentialsBuilder {
	b.credentials.Password = password
	return b
}

func (b *CredentialsBuilder) Build() domain.Credentials {
	return b.credentials
}
