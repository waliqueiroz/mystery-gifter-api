package rest

import (
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

type MatchDTO struct {
	GiverID    string `json:"giver_id" validate:"required,uuid"`
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
