package build_domain

import "github.com/waliqueiroz/mystery-gifter-api/internal/domain"

type SearchResultBuilder[T any] struct {
	searchResult domain.SearchResult[T]
}

func NewSearchResultBuilder[T any]() *SearchResultBuilder[T] {
	return &SearchResultBuilder[T]{
		searchResult: domain.SearchResult[T]{
			Result: []T{},
			Paging: domain.Paging{
				Total:  0,
				Limit:  10,
				Offset: 0,
			},
		},
	}
}

func (b *SearchResultBuilder[T]) WithResult(result []T) *SearchResultBuilder[T] {
	b.searchResult.Result = result
	return b
}

func (b *SearchResultBuilder[T]) WithLimit(limit int) *SearchResultBuilder[T] {
	b.searchResult.Paging.Limit = limit
	return b
}

func (b *SearchResultBuilder[T]) WithOffset(offset int) *SearchResultBuilder[T] {
	b.searchResult.Paging.Offset = offset
	return b
}

func (b *SearchResultBuilder[T]) WithTotal(total int) *SearchResultBuilder[T] {
	b.searchResult.Paging.Total = total
	return b
}

func (b *SearchResultBuilder[T]) Build() domain.SearchResult[T] {
	return b.searchResult
}
