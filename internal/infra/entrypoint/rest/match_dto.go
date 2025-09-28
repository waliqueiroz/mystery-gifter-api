package rest

import (
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

// MatchDTO represents a match between two users in a group
// swagger:model MatchDTO
type MatchDTO struct {
	// ID of the user who gives the gift
	// required: true
	// example: 01234567-89ab-cdef-0123-456789abcdef
	GiverID string `json:"giver_id" validate:"required,uuid"`

	// ID of the user who receives the gift
	// required: true
	// example: 01234567-89ab-cdef-0123-456789abcdef
	ReceiverID string `json:"receiver_id" validate:"required,uuid"`
}

func (m *MatchDTO) Validate() error {
	if errs := validator.Validate(m); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

func mapMatchFromDomain(match domain.Match) (*MatchDTO, error) {
	matchDTO := MatchDTO{
		GiverID:    match.GiverID,
		ReceiverID: match.ReceiverID,
	}

	if err := matchDTO.Validate(); err != nil {
		return nil, err
	}

	return &matchDTO, nil
}

func mapMatchesFromDomain(matches []domain.Match) ([]MatchDTO, error) {
	matchDTOs := make([]MatchDTO, 0, len(matches))
	for _, match := range matches {
		matchDTO, err := mapMatchFromDomain(match)
		if err != nil {
			return nil, err
		}
		matchDTOs = append(matchDTOs, *matchDTO)
	}
	return matchDTOs, nil
}
