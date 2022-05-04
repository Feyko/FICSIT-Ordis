package repositories

import "FICSIT-Ordis/internal/storable"

type Repository[E storable.I] interface {
	Find(ID string) (E, error)
	GetAll() ([]E, error)

	Create(element E) error
	Update(ID string, newElement E) error
	Delete(ID string) error
}

type Searchable[E storable.I] interface {
	Repository[E]
	Search(search string, fields []string) ([]E, error)
}
