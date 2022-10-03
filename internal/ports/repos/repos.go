package repos

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/arango"
	"FICSIT-Ordis/internal/ports/repos/memrepo"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"fmt"
	"github.com/pkg/errors"
)

func GetCollection[T id.IDer, U id.IDer](repository repo.Repository[U], name string) (repo.Collection[T], error) {
	repoA, err := Retype[T](repository)
	if err != nil {
		return nil, errors.Wrap(err, "could not retype repo")
	}
	collection, err := repoA.GetCollection(name)
	if err != nil {
		return nil, errors.Wrap(err, "could not get collection")
	}
	typed, ok := collection.(repo.Collection[T])
	if !ok {
		return nil, errors.New(fmt.Sprintf("collection '%v' does not hold type '%T'", name, *new(T)))
	}
	return typed, nil
}

func CreateCollection[T id.IDer, U id.IDer](repository repo.Repository[U], name string) (repo.Collection[T], error) {
	repoA, err := Retype[T](repository)
	if err != nil {
		return nil, errors.Wrap(err, "could not retype repo")
	}
	collection, err := repoA.CreateCollection(name)
	if err != nil {
		return nil, errors.Wrap(err, "could not create collection")
	}
	typed, ok := collection.(repo.Collection[T])
	if !ok {
		return nil, errors.New(fmt.Sprintf("created collection '%v' does not hold type '%T'", name, *new(T)))
	}
	return typed, nil
}

func Retype[newT id.IDer, oldT id.IDer](repo repo.Repository[oldT]) (repo.Repository[newT], error) {
	switch typed := repo.(type) {
	case *memrepo.Repository[oldT]:
		return (*memrepo.Repository[newT])(typed), nil
	case *arango.Repository[oldT]:
		return (*arango.Repository[newT])(typed), nil
	default:
		return nil, errors.New("unsupported repository type")
	}
}
