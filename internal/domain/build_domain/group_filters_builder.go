package build_domain

import "github.com/waliqueiroz/mystery-gifter-api/internal/domain"

type GroupFiltersBuilder struct {
	groupFilters domain.GroupFilters
}

func NewGroupFiltersBuilder() *GroupFiltersBuilder {
	return &GroupFiltersBuilder{
		groupFilters: domain.GroupFilters{
			Name:          "",
			Status:        "",
			OwnerID:       "",
			UserID:        "",
			Limit:         domain.DefaultGroupLimit,
			Offset:        domain.DefaultGroupOffset,
			SortDirection: domain.DefaultGroupSortDirection,
			SortBy:        domain.DefaultGroupSortBy,
		},
	}
}

func (b *GroupFiltersBuilder) WithName(name string) *GroupFiltersBuilder {
	b.groupFilters.Name = name
	return b
}

func (b *GroupFiltersBuilder) WithStatus(status domain.GroupStatus) *GroupFiltersBuilder {
	b.groupFilters.Status = status
	return b
}

func (b *GroupFiltersBuilder) WithOwnerID(ownerID string) *GroupFiltersBuilder {
	b.groupFilters.OwnerID = ownerID
	return b
}

func (b *GroupFiltersBuilder) WithUserID(userID string) *GroupFiltersBuilder {
	b.groupFilters.UserID = userID
	return b
}

func (b *GroupFiltersBuilder) WithLimit(limit int) *GroupFiltersBuilder {
	b.groupFilters.Limit = limit
	return b
}

func (b *GroupFiltersBuilder) WithOffset(offset int) *GroupFiltersBuilder {
	b.groupFilters.Offset = offset
	return b
}

func (b *GroupFiltersBuilder) WithSortDirection(sortDirection domain.SortDirectionType) *GroupFiltersBuilder {
	b.groupFilters.SortDirection = sortDirection
	return b
}

func (b *GroupFiltersBuilder) WithSortBy(sortBy string) *GroupFiltersBuilder {
	b.groupFilters.SortBy = sortBy
	return b
}

func (b *GroupFiltersBuilder) Build() domain.GroupFilters {
	return b.groupFilters
}
