package repos

import "FICSIT-Ordis/internal/id"

type Repository interface {
	NewCollection(name string) (Collection, error)
	GetCollection(name string) (Collection, error)
}

type Collection interface {
	Get(ID string) (id.IDer, error)
	GetAll() ([]id.IDer, error)
	Search(search string, fields []string) ([]id.IDer, error)

	Create(element id.IDer) error
	Update(ID string, newElement id.IDer) error
	Delete(ID string) error
}
