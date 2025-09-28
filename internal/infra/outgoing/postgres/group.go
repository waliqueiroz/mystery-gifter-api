package postgres

import (
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
)

type Group struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	OwnerID   string    `db:"owner_id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type GroupSummary struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	OwnerID   string    `db:"owner_id"`
	Status    string    `db:"status"`
	UserCount int       `db:"user_count"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func mapGroupToDomain(group Group, groupUsers []User, matches []Match) (*domain.Group, error) {
	domainUsers, err := mapUsersToDomain(groupUsers)
	if err != nil {
		return nil, err
	}

	domainMatches, err := mapMatchesToDomain(matches)
	if err != nil {
		return nil, err
	}

	domainGroup := domain.Group{
		ID:        group.ID,
		Name:      group.Name,
		OwnerID:   group.OwnerID,
		Users:     domainUsers,
		Status:    domain.GroupStatus(group.Status),
		Matches:   domainMatches,
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}

	if err := domainGroup.Validate(); err != nil {
		return nil, err
	}

	return &domainGroup, nil
}

func mapGroupSummariesToDomain(groupSummaries []GroupSummary) ([]domain.GroupSummary, error) {
	domainGroupSummaries := make([]domain.GroupSummary, 0, len(groupSummaries))

	for _, groupSummary := range groupSummaries {
		domainGroupSummary, err := mapGroupSummaryToDomain(groupSummary)
		if err != nil {
			return nil, err
		}
		domainGroupSummaries = append(domainGroupSummaries, *domainGroupSummary)
	}

	return domainGroupSummaries, nil
}

func mapGroupSummaryToDomain(groupSummary GroupSummary) (*domain.GroupSummary, error) {
	domainGroupSummary := domain.GroupSummary{
		ID:        groupSummary.ID,
		Name:      groupSummary.Name,
		OwnerID:   groupSummary.OwnerID,
		Status:    domain.GroupStatus(groupSummary.Status),
		UserCount: groupSummary.UserCount,
		CreatedAt: groupSummary.CreatedAt,
		UpdatedAt: groupSummary.UpdatedAt,
	}

	if err := domainGroupSummary.Validate(); err != nil {
		return nil, err
	}

	return &domainGroupSummary, nil
}
