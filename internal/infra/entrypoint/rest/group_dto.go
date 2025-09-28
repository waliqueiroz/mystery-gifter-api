package rest

import (
	"time"

	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

// CreateGroupDTO represents the data needed to create a new group
// swagger:model CreateGroupDTO
type CreateGroupDTO struct {
	// Group name
	// required: true
	// example: Secret Santa 2024
	Name string `json:"name" validate:"required"`
}

func (g *CreateGroupDTO) Validate() error {
	if errs := validator.Validate(g); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

// GroupDTO represents a complete group with all its information
// swagger:model GroupDTO
type GroupDTO struct {
	// Unique group identifier
	// required: true
	// example: 01234567-89ab-cdef-0123-456789abcdef
	ID string `json:"id" validate:"required,uuid"`

	// Group name
	// required: true
	// example: Secret Santa 2024
	Name string `json:"name" validate:"required"`

	// List of users in the group
	// required: true
	Users []UserDTO `json:"users" validate:"required,min=1"`

	// ID of the group owner
	// required: true
	// example: 01234567-89ab-cdef-0123-456789abcdef
	OwnerID string `json:"owner_id" validate:"required,uuid"`

	// List of matches in the group
	Matches []MatchDTO `json:"matches" validate:"dive,omitempty"`

	// Group status
	// required: true
	// example: OPEN
	// enum: OPEN,MATCHED,ARCHIVED
	Status string `json:"status" validate:"required,oneof=OPEN MATCHED ARCHIVED"`

	// Group creation timestamp
	// required: true
	// example: 2024-01-01T00:00:00Z
	CreatedAt time.Time `json:"created_at" validate:"required"`

	// Group last update timestamp
	// required: true
	// example: 2024-01-01T00:00:00Z
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

func (g *GroupDTO) Validate() error {
	if errs := validator.Validate(g); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

// GroupSummaryDTO represents a summary of a group (used in search results)
// swagger:model GroupSummaryDTO
type GroupSummaryDTO struct {
	// Unique group identifier
	// required: true
	// example: 01234567-89ab-cdef-0123-456789abcdef
	ID string `json:"id" validate:"required,uuid"`

	// Group name
	// required: true
	// example: Secret Santa 2024
	Name string `json:"name" validate:"required"`

	// Group status
	// required: true
	// example: OPEN
	// enum: OPEN,MATCHED,ARCHIVED
	Status string `json:"status" validate:"required,oneof=OPEN MATCHED ARCHIVED"`

	// ID of the group owner
	// required: true
	// example: 01234567-89ab-cdef-0123-456789abcdef
	OwnerID string `json:"owner_id" validate:"required,uuid"`

	// Number of users in the group
	// required: true
	// example: 5
	UserCount int `json:"user_count" validate:"required,min=0"`

	// Group creation timestamp
	// required: true
	// example: 2024-01-01T00:00:00Z
	CreatedAt time.Time `json:"created_at" validate:"required"`

	// Group last update timestamp
	// required: true
	// example: 2024-01-01T00:00:00Z
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

func (g *GroupSummaryDTO) Validate() error {
	if errs := validator.Validate(g); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

func mapGroupFromDomain(group domain.Group) (*GroupDTO, error) {
	users, err := mapUsersFromDomain(group.Users)
	if err != nil {
		return nil, err
	}

	matches, err := mapMatchesFromDomain(group.Matches)
	if err != nil {
		return nil, err
	}

	groupDTO := GroupDTO{
		ID:        group.ID,
		Name:      group.Name,
		Users:     users,
		OwnerID:   group.OwnerID,
		Matches:   matches,
		Status:    string(group.Status),
		CreatedAt: group.CreatedAt,
		UpdatedAt: group.UpdatedAt,
	}

	if err := groupDTO.Validate(); err != nil {
		return nil, err
	}

	return &groupDTO, nil
}

func mapGroupSummaryFromDomain(groupSummary domain.GroupSummary) (*GroupSummaryDTO, error) {
	groupSummaryDTO := GroupSummaryDTO{
		ID:        groupSummary.ID,
		Name:      groupSummary.Name,
		Status:    string(groupSummary.Status),
		OwnerID:   groupSummary.OwnerID,
		UserCount: groupSummary.UserCount,
		CreatedAt: groupSummary.CreatedAt,
		UpdatedAt: groupSummary.UpdatedAt,
	}

	if err := groupSummaryDTO.Validate(); err != nil {
		return nil, err
	}

	return &groupSummaryDTO, nil
}

// GroupFiltersDTO represents filters for searching groups
// swagger:model GroupFiltersDTO
type GroupFiltersDTO struct {
	// Filter by group name
	// example: Secret Santa 2024
	Name string `query:"name" json:"name"`

	// Filter by group owner ID
	// example: 01234567-89ab-cdef-0123-456789abcdef
	OwnerID string `query:"owner_id" json:"owner_id"`

	// Filter by group status
	// example: OPEN
	Status string `query:"status" json:"status"`

	// Maximum number of results to return
	// example: 10
	Limit int `query:"limit" json:"limit"`

	// Number of results to skip
	// example: 0
	Offset int `query:"offset" json:"offset"`

	// Sort direction (asc, desc)
	// example: asc
	Sort string `query:"sort" json:"sort"`

	// Field to sort by
	// example: name
	SortField string `query:"sort_field" json:"sort_field"`
}

func mapGroupFiltersDTOToDomain(filtersDTO GroupFiltersDTO) (*domain.GroupFilters, error) {
	var status domain.GroupStatus
	if filtersDTO.Status != "" {
		status = domain.GroupStatus(filtersDTO.Status)
	}

	var sortDirection domain.SortDirectionType
	if filtersDTO.Sort != "" {
		sortDirection = domain.SortDirectionType(filtersDTO.Sort)
	}

	filters := &domain.GroupFilters{
		Name:          filtersDTO.Name,
		Status:        status,
		OwnerID:       filtersDTO.OwnerID,
		Limit:         filtersDTO.Limit,
		Offset:        filtersDTO.Offset,
		SortDirection: sortDirection,
		SortBy:        filtersDTO.SortField,
	}

	return filters, nil
}
