package base

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"context"
)

func NewSearchable[E id.IDer](conf Config, collection repo.Collection[E]) *Searchable[E] {
	base := New(conf, collection)
	return &Searchable[E]{
		base,
	}
}

type Searchable[E id.IDer] struct {
	*Module[E]
}

func (s *Searchable[E]) Search(ctx context.Context, search string) ([]E, error) {
	return s.Collection.Search(ctx, search)
}
