package repos

import (
	"FICSIT-Ordis/internal/id"
	"fmt"
)

type Repository interface {
	GetCollection(name string) (any, error)
}

func GetCollection[T id.IDer](repository Repository, name string) (TypedCollection[T], error) {
	collection, err := repository.GetCollection(name)
	if err != nil {
		return nil, err
	}
	typed, ok := collection.(TypedCollection[T])
	if !ok {
		return nil, fmt.Errorf("collection %v does not hold the type %t", name, *new(T))
	}
	return typed, nil
}

//type UntypedCollection interface {
//	Get(ID string) (any, error)
//	GetAll() ([]any, error)
//	Search(search string, fields []string) ([]any, error)
//
//	Create(element id.IDer) error
//	Update(ID string, updateElement Updater[any]) error
//	Delete(ID string) error
//}

type TypedCollection[E id.IDer] interface {
	Get(ID string) (E, error)
	GetAll() ([]E, error)
	Search(search string, fields []string) ([]E, error)

	Create(element E) error
	Update(ID string, updateElement Updater[E]) error
	Delete(ID string) error
}

type Updater[T any] interface {
	Update(update T) T
}
