package postgres

import (
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type Match struct {
	GiverID    string `db:"giver_id"`
	ReceiverID string `db:"receiver_id"`
}

func mapMatchToDomain(match Match) (*domain.Match, error) {
	domainMatch := &domain.Match{
		GiverID:    match.GiverID,
		ReceiverID: match.ReceiverID,
	}

	if err := domainMatch.Validate(); err != nil {
		return nil, err
	}

	return domainMatch, nil
}

func mapMatchesToDomain(matches []Match) ([]domain.Match, error) {
	domainMatches := make([]domain.Match, 0, len(matches))
	for _, match := range matches {
		domainMatch, err := mapMatchToDomain(match)
		if err != nil {
			return nil, err
		}

		domainMatches = append(domainMatches, *domainMatch)
	}

	return domainMatches, nil
}
