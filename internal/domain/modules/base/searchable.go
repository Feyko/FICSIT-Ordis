package base

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"context"
)

func NewSearchable[E id.Searchable](conf Config, collection repo.Collection[E]) *Searchable[E] {
	var defaultS E
	base := New(conf, collection)
	return &Searchable[E]{
		base,
		defaultS.SearchFields(),
	}
}

type Searchable[E id.Searchable] struct {
	*Module[E]
	searchFields []string
}

func (s *Searchable[E]) Search(ctx context.Context, search string) ([]E, error) {
	return s.Collection.Search(ctx, search, s.searchFields)
}
