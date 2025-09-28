package rest

import (
	"github.com/waliqueiroz/mystery-gifter-api/internal/domain"
	"github.com/waliqueiroz/mystery-gifter-api/pkg/validator"
)

// SearchResultDTO represents a paginated search result
type SearchResultDTO[T any] struct {
	// List of results
	Result []T `json:"result" validate:"required"`

	// Pagination information
	Paging PagingDTO `json:"paging" validate:"required"`
}

// PagingDTO represents pagination information
// swagger:model PagingDTO
type PagingDTO struct {
	// Total number of results
	// example: 100
	Total int `json:"total" validate:"omitempty,min=0"`

	// Maximum number of results per page
	// required: true
	// example: 10
	Limit int `json:"limit" validate:"required,min=1"`

	// Number of results to skip
	// example: 0
	Offset int `json:"offset" validate:"omitempty,min=0"`
}

type ResultMapper[T any, R any] func(T) (*R, error)

func (s *SearchResultDTO[T]) Validate() error {
	if errs := validator.Validate(s); len(errs) > 0 {
		return domain.NewValidationError(errs)
	}
	return nil
}

func mapSearchResultFromDomain[T any, R any](searchResult *domain.SearchResult[T], mapResultFromDomain ResultMapper[T, R]) (*SearchResultDTO[R], error) {
	results, err := mapResultsFromDomain(searchResult.Result, mapResultFromDomain)
	if err != nil {
		return nil, err
	}

	result := &SearchResultDTO[R]{
		Result: results,
		Paging: PagingDTO{
			Total:  searchResult.Paging.Total,
			Limit:  searchResult.Paging.Limit,
			Offset: searchResult.Paging.Offset,
		},
	}

	if err := result.Validate(); err != nil {
		return nil, err
	}

	return result, nil
}

func mapResultsFromDomain[T any, R any](results []T, mapResultFromDomain ResultMapper[T, R]) ([]R, error) {
	resultsDTOs := make([]R, 0, len(results))

	for _, result := range results {
		resultDTO, err := mapResultFromDomain(result)
		if err != nil {
			return nil, err
		}
		resultsDTOs = append(resultsDTOs, *resultDTO)
	}

	return resultsDTOs, nil
}

// UserSearchResultDTO represents a paginated search result for users
// swagger:model UserSearchResultDTO
type UserSearchResultDTO struct {
	// List of users
	// required: true
	Result []UserDTO `json:"result" validate:"required"`

	// Pagination information
	// required: true
	Paging PagingDTO `json:"paging" validate:"required"`
}

// GroupSearchResultDTO represents a paginated search result for groups
// swagger:model GroupSearchResultDTO
type GroupSearchResultDTO struct {
	// List of groups
	// required: true
	Result []GroupSummaryDTO `json:"result" validate:"required"`

	// Pagination information
	// required: true
	Paging PagingDTO `json:"paging" validate:"required"`
}
