package base

import (
	"FICSIT-Ordis/internal/core/ports/repos"
	"FICSIT-Ordis/internal/core/ports/repos/memrepo"
	"FICSIT-Ordis/internal/core/ports/repos/translators"
	"FICSIT-Ordis/internal/id"
	"fmt"
	"log"
)

func newDefaultSearchable[S id.Searchable]() *Searchable[S] {
	repo := memrepo.New()
	collection, err := repo.GetCollection(fmt.Sprintf("%T", *new(S)))
	if err != nil {
		log.Fatalf("Something went horribly wrong and we could not create a new collection in the memrepo: %v", err)
	}
	translator := translators.Wrap[S](collection)
	return NewSearchable(translator)
}

func NewSearchable[S id.Searchable](collection repos.TypedCollection[S]) *Searchable[S] {
	var defaultS S
	base := New(collection)
	return &Searchable[S]{
		base,
		defaultS.SearchFields(),
	}
}

type Searchable[S id.Searchable] struct {
	*Module[S]
	searchFields []string
}

func (s *Searchable[S]) Search(search string) ([]S, error) {
	return s.Collection.Search(search, s.searchFields)
}
