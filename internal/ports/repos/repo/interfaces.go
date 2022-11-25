package repo

import (
	"FICSIT-Ordis/internal/id"
	"context"
)

type Repository[T id.IDer] interface {
	// Get a collection with a unique name. Use repos.GetCollection instead.
	GetCollection(name string) (any, error)

	// Create a collection that holds values of the current type parameter. Use repos.CreateCollection instead.
	CreateCollection(name string) (any, error)

	DeleteCollection(name string) error
}

type Collection[E id.IDer] interface {
	Get(ctx context.Context, ID string) (E, error)
	GetAll(ctx context.Context) ([]E, error)
	Search(ctx context.Context, search string) ([]E, error)

	Create(ctx context.Context, element E) error
	Update(ctx context.Context, ID string, updateElement any) (old E, new E, err error)
	Delete(ctx context.Context, ID string) error
}
