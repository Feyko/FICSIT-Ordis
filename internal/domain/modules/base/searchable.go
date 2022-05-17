package base

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/memrepo"
	"FICSIT-Ordis/internal/ports/repos/translators"
	"fmt"
	"log"
)

func newDefaultSearchable[E id.Searchable, U id.IDer]() *Searchable[E, U] {
	repo := memrepo.New()
	collection, err := repo.GetCollection(fmt.Sprintf("%T", *new(E)))
	if err != nil {
		log.Fatalf("Something went horribly wrong and we could not create a new collection in the memrepo: %v", err)
	}
	translator := translators.Wrap[E, U](collection)
	return NewSearchable(translator)
}

func NewSearchable[E id.Searchable, U id.IDer](collection repos.TypedCollection[E, U]) *Searchable[E, U] {
	var defaultS E
	base := New(collection)
	return &Searchable[E, U]{
		base,
		defaultS.SearchFields(),
	}
}

type Searchable[E id.Searchable, U id.IDer] struct {
	*Module[E, U]
	searchFields []string
}

func (s *Searchable[E, U]) Search(search string) ([]E, error) {
	return s.Collection.Search(search, s.searchFields)
}
