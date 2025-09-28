package build_rest

import "github.com/waliqueiroz/mystery-gifter-api/internal/infra/entrypoint/rest"

type SearchResultDTOBuilder[T any] struct {
	searchResultDTO rest.SearchResultDTO[T]
}

func NewSearchResultDTOBuilder[T any]() *SearchResultDTOBuilder[T] {
	return &SearchResultDTOBuilder[T]{
		searchResultDTO: rest.SearchResultDTO[T]{
			Result: []T{},
			Paging: rest.PagingDTO{
				Total:  0,
				Limit:  10,
				Offset: 0,
			},
		},
	}
}

func (b *SearchResultDTOBuilder[T]) WithResult(result []T) *SearchResultDTOBuilder[T] {
	b.searchResultDTO.Result = result
	return b
}

func (b *SearchResultDTOBuilder[T]) WithTotal(total int) *SearchResultDTOBuilder[T] {
	b.searchResultDTO.Paging.Total = total
	return b
}

func (b *SearchResultDTOBuilder[T]) WithLimit(limit int) *SearchResultDTOBuilder[T] {
	b.searchResultDTO.Paging.Limit = limit
	return b
}

func (b *SearchResultDTOBuilder[T]) WithOffset(offset int) *SearchResultDTOBuilder[T] {
	b.searchResultDTO.Paging.Offset = offset
	return b
}

func (b *SearchResultDTOBuilder[T]) Build() rest.SearchResultDTO[T] {
	return b.searchResultDTO
}
