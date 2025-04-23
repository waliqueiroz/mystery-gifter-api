package rest

import (
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type CreateGroupDTO struct {
	Name string `json:"name" validate:"required"`
}

func (g *CreateGroupDTO) Validate() error {
	if errs := validator.Validate(g); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}
