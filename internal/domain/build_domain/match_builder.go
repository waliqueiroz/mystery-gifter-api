package build_domain

import (
	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type MatchBuilder struct {
	match domain.Match
}

func NewMatchBuilder() *MatchBuilder {
	return &MatchBuilder{
		match: domain.Match{
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

func (b *MatchBuilder) Build() domain.Match {
	return b.match
}
