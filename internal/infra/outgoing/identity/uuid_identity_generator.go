package identity

import (
	"github.com/google/uuid"
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type IdentityFunc func() (uuid.UUID, error)

type UUIDIdentityGenerator struct {
	newUUID IdentityFunc
}

func NewUUIDIdentityGenerator(newUUID IdentityFunc) domain.IdentityGenerator {
	return &UUIDIdentityGenerator{
		newUUID: newUUID,
	}
}

func (g *UUIDIdentityGenerator) Generate() (string, error) {
	id, err := g.newUUID()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
