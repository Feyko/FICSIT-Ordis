package repo

import (
	"FICSIT-Ordis/internal/id"
)

type Repository[T id.IDer] interface {
	// Get a collection with a unique name. Use repos.GetCollection instead.
	GetCollection(name string) (any, error)

	// Create a collection that holds values of the current type parameter. Use repos.CreateCollection instead.
	CreateCollection(name string) (any, error)

	DeleteCollection(name string) error
}

type Collection[E id.IDer] interface {
	Get(ID string) (E, error)
	GetAll() ([]E, error)
	Search(search string, fields []string) ([]E, error)

	Create(element E) error
	Update(ID string, updateElement any) (E, error)
	Delete(ID string) error
}

type Updater[T any] interface {
	Update(update T) T
}
