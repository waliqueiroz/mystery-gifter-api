package build_postgres

import (
	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/infra/outgoing/postgres"
)

type MatchBuilder struct {
	match postgres.Match
}

func NewMatchBuilder() *MatchBuilder {
	return &MatchBuilder{
		match: postgres.Match{
			GiverID:    uuid.New().String(),
			ReceiverID: uuid.New().String(),
		},
	}
}

func (b *MatchBuilder) WithGiverID(giverID string) *MatchBuilder {
	b.match.GiverID = giverID
	return b
}

func (b *MatchBuilder) WithReceiverID(receiverID string) *MatchBuilder {
	b.match.ReceiverID = receiverID
	return b
}

func (b *MatchBuilder) Build() postgres.Match {
	return b.match
}
