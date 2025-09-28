package rest

import (
	"time"

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

type GroupDTO struct {
	ID        string     `json:"id" validate:"required,uuid"`
	Name      string     `json:"name" validate:"required"`
	Users     []UserDTO  `json:"users" validate:"required,min=1"`
	OwnerID   string     `json:"owner_id" validate:"required,uuid"`
	Matches   []MatchDTO `json:"matches" validate:"dive,omitempty"`
	Status    string     `json:"status" validate:"required,oneof=OPEN MATCHED ARCHIVED"`
	CreatedAt time.Time  `json:"created_at" validate:"required"`
	UpdatedAt time.Time  `json:"updated_at" validate:"required"`
}

func (g *GroupDTO) Validate() error {
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

type GroupFiltersDTO struct {
	Name          string `query:"name" json:"name"`
	Status        string `query:"status" json:"status"`
	OwnerID       string `query:"owner_id" json:"owner_id"`
	UserID        string `query:"user_id" json:"user_id"`
	Limit         int    `query:"limit" json:"limit"`
	Offset        int    `query:"offset" json:"offset"`
	SortDirection string `query:"sort_direction" json:"sort_direction" validate:"omitempty,oneof=ASC DESC"`
	SortBy        string `query:"sort_by" json:"sort_by" validate:"omitempty,oneof=name status owner_id created_at updated_at"`
}

func (g *GroupFiltersDTO) Validate() error {
	if errs := validator.Validate(g); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

type GroupSummaryDTO struct {
	ID        string    `json:"id" validate:"required,uuid"`
	Name      string    `json:"name" validate:"required"`
	Status    string    `json:"status" validate:"required,oneof=OPEN MATCHED ARCHIVED"`
	OwnerID   string    `json:"owner_id" validate:"required,uuid"`
	UserCount int       `json:"user_count"`
	CreatedAt time.Time `json:"created_at" validate:"required"`
	UpdatedAt time.Time `json:"updated_at" validate:"required"`
}

func (g *GroupSummaryDTO) Validate() error {
	if errs := validator.Validate(g); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

func mapGroupFiltersDTOToDomain(dto GroupFiltersDTO) (*domain.GroupFilters, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}

	return domain.NewGroupFilters(
		dto.Name,
		dto.OwnerID,
		dto.UserID,
		domain.GroupStatus(dto.Status),
		dto.Limit,
		dto.Offset,
		domain.SortDirectionType(dto.SortDirection),
		dto.SortBy,
	)
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
