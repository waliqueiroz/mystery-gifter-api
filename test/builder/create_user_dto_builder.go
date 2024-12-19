package builder

import "github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"

type CreateUserDTOBuilder struct {
	createUserDTO rest.CreateUserDTO
}

func NewCreateUserDTOBuilder() *CreateUserDTOBuilder {
	return &CreateUserDTOBuilder{
		createUserDTO: rest.CreateUserDTO{
			Name:            "Default Name",
			Surname:         "Default Surname",
			Email:           "default@mail.com",
			Password:        "defaultpassword",
			PasswordConfirm: "defaultpassword",
		},
	}
}

func (b *CreateUserDTOBuilder) WithName(name string) *CreateUserDTOBuilder {
	b.createUserDTO.Name = name
	return b
}

func (b *CreateUserDTOBuilder) WithSurname(surname string) *CreateUserDTOBuilder {
	b.createUserDTO.Surname = surname
	return b
}

func (b *CreateUserDTOBuilder) WithEmail(email string) *CreateUserDTOBuilder {
	b.createUserDTO.Email = email
	return b
}

func (b *CreateUserDTOBuilder) WithPassword(password string) *CreateUserDTOBuilder {
	b.createUserDTO.Password = password
	return b
}

func (b *CreateUserDTOBuilder) WithPasswordConfirm(passwordConfirm string) *CreateUserDTOBuilder {
	b.createUserDTO.PasswordConfirm = passwordConfirm
	return b
}

func (b *CreateUserDTOBuilder) Build() rest.CreateUserDTO {
	return b.createUserDTO
}
