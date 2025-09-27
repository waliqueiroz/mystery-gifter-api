package build_rest

import (
	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"
)

type MatchDTOBuilder struct {
	matchDTO rest.MatchDTO
}

func NewMatchDTOBuilder() *MatchDTOBuilder {
	return &MatchDTOBuilder{
		matchDTO: rest.MatchDTO{
			GiverID:    uuid.New().String(),
			ReceiverID: uuid.New().String(),
		},
	}
}

func (b *MatchDTOBuilder) WithGiverID(giverID string) *MatchDTOBuilder {
	b.matchDTO.GiverID = giverID
	return b
}

func (b *MatchDTOBuilder) WithReceiverID(receiverID string) *MatchDTOBuilder {
	b.matchDTO.ReceiverID = receiverID
	return b
}

func (b *MatchDTOBuilder) Build() rest.MatchDTO {
	return b.matchDTO
}
