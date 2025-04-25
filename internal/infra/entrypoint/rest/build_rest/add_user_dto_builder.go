package build_rest

import (
	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

type AddUserDTOBuilder struct {
	addUserDTO rest.AddUserDTO
}

func NewAddUserDTOBuilder() *AddUserDTOBuilder {
	return &AddUserDTOBuilder{
		addUserDTO: rest.AddUserDTO{
			UserID: uuid.NewString(),
		},
	}
}

func (b *AddUserDTOBuilder) WithUserID(userID string) *AddUserDTOBuilder {
	b.addUserDTO.UserID = userID
	return b
}

func (b *AddUserDTOBuilder) Build() rest.AddUserDTO {
	return b.addUserDTO
}
