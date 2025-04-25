package build_rest

import "github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"

type CreateGroupDTOBuilder struct {
	createGroupDTO rest.CreateGroupDTO
}

func NewCreateGroupDTOBuilder() *CreateGroupDTOBuilder {
	return &CreateGroupDTOBuilder{
		createGroupDTO: rest.CreateGroupDTO{
			Name: "Test Group",
		},
	}
}

func (b *CreateGroupDTOBuilder) WithName(name string) *CreateGroupDTOBuilder {
	b.createGroupDTO.Name = name
	return b
}

func (b *CreateGroupDTOBuilder) Build() rest.CreateGroupDTO {
	return b.createGroupDTO
}
