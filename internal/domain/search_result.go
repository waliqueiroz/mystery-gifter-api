package domain

import "github.com/waliqueiroz/mystery-gifter-api/pkg/validator"

type Paging struct {
	Total  int `validate:"min=0"`
	Limit  int `validate:"required,min=1"`
	Offset int `validate:"min=0"`
}

type SearchResult[T any] struct {
	Result []T    `validate:"required"`
	Paging Paging `validate:"required"`
}

func NewSearchResult[T any](result []T, limit, offset, total int) (*SearchResult[T], error) {
	searchResult := SearchResult[T]{
		Result: result,
		Paging: Paging{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	}

	if err := searchResult.Validate(); err != nil {
		return nil, err
	}

	return &searchResult, nil
}

func (s *SearchResult[T]) Validate() error {
	if errs := validator.Validate(s); len(errs) > 0 {
		return NewValidationError(errs)
	}
	return nil
}
