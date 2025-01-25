package build_rest

import "github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"

type CredentialsDTOBuilder struct {
	credentialsDTO rest.CredentialsDTO
}

func NewCredentialsDTOBuilder() *CredentialsDTOBuilder {
	return &CredentialsDTOBuilder{
		credentialsDTO: rest.CredentialsDTO{
			Email:    "default@example.com",
			Password: "defaultpassword",
		},
	}
}

func (b *CredentialsDTOBuilder) WithEmail(email string) *CredentialsDTOBuilder {
	b.credentialsDTO.Email = email
	return b
}

func (b *CredentialsDTOBuilder) WithPassword(password string) *CredentialsDTOBuilder {
	b.credentialsDTO.Password = password
	return b
}

func (b *CredentialsDTOBuilder) Build() rest.CredentialsDTO {
	return b.credentialsDTO
}
