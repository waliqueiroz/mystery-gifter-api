package build_domain

import "github.com/waliqueiroz/mystery-gifter-api/internal/domain"

type UserFiltersBuilder struct {
	userFilters domain.UserFilters
}

func NewUserFiltersBuilder() *UserFiltersBuilder {
	return &UserFiltersBuilder{
		userFilters: domain.UserFilters{
			Name:          "",
			Surname:       "",
			Email:         "",
			Limit:         domain.DefaultUserLimit,
			Offset:        domain.DefaultUserOffset,
			SortDirection: domain.DefaultUserSortDirection,
			SortBy:        domain.DefaultUserSortBy,
		},
	}
}

func (b *UserFiltersBuilder) WithName(name string) *UserFiltersBuilder {
	b.userFilters.Name = name
	return b
}

func (b *UserFiltersBuilder) WithSurname(surname string) *UserFiltersBuilder {
	b.userFilters.Surname = surname
	return b
}

func (b *UserFiltersBuilder) WithEmail(email string) *UserFiltersBuilder {
	b.userFilters.Email = email
	return b
}

func (b *UserFiltersBuilder) WithLimit(limit int) *UserFiltersBuilder {
	b.userFilters.Limit = limit
	return b
}

func (b *UserFiltersBuilder) WithOffset(offset int) *UserFiltersBuilder {
	b.userFilters.Offset = offset
	return b
}

func (b *UserFiltersBuilder) WithSortDirection(sortDirection domain.SortDirectionType) *UserFiltersBuilder {
	b.userFilters.SortDirection = sortDirection
	return b
}

func (b *UserFiltersBuilder) WithSortBy(sortBy string) *UserFiltersBuilder {
	b.userFilters.SortBy = sortBy
	return b
}

func (b *UserFiltersBuilder) Build() domain.UserFilters {
	return b.userFilters
}
