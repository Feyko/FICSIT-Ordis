package repos

import "FICSIT-Ordis/internal/id"

type Repository interface {
	GetCollection(name string) (UntypedCollection, error)
}

type UntypedCollection interface {
	Get(ID string) (any, error)
	GetAll() ([]any, error)
	Search(search string, fields []string) ([]any, error)

	Create(element id.IDer) error
	Update(ID string, newElement id.IDer) error
	Delete(ID string) error
}

type TypedCollection[E id.IDer] interface {
	Get(ID string) (E, error)
	GetAll() ([]E, error)
	Search(search string, fields []string) ([]E, error)

	Create(element E) error
	Update(ID string, newElement E) error
	Delete(ID string) error
}
