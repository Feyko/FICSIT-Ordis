package base

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/memrepo"
	"fmt"
	"log"
)

func newDefaultSearchable[E id.Searchable]() *Searchable[E] {
	repo := memrepo.New()
	collection, err := repos.GetCollection(repo, fmt.Sprintf("%T", *new(E)))
	if err != nil {
		log.Fatalf("Something went horribly wrong and we could not create a new collection in the memrepo: %v", err)
	}
	return NewSearchable(collection)
}

func NewSearchable[E id.Searchable](collection repos.TypedCollection[E]) *Searchable[E] {
	var defaultS E
	base := New(collection)
	return &Searchable[E]{
		base,
		defaultS.SearchFields(),
	}
}

type Searchable[E id.Searchable] struct {
	*Module[E]
	searchFields []string
}

func (s *Searchable[E]) Search(search string) ([]E, error) {
	return s.Collection.Search(search, s.searchFields)
}
