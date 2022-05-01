package repositories

import "FICSIT-Ordis/internal/id"

type Repository[E id.IDer] interface {
	Find(ID string) (E, error)
	GetAll() ([]E, error)

	Create(element E) error
	Update(ID string, newElement E) error
	Delete(ID string) error
}

type Searchable[E id.IDer] interface {
	Repository[E]
	Search(search string, fields []string) ([]E, error)
}
