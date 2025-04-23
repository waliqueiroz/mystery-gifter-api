package identity

import (
	"fmt"

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
		return "", fmt.Errorf("error generating UUID: %w", err)
	}
	return id.String(), nil
}
