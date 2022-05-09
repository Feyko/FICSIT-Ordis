package base

import (
	"FICSIT-Ordis/internal/core/ports/repos"
	"FICSIT-Ordis/internal/core/ports/repos/memrepo"
	"FICSIT-Ordis/internal/core/ports/repos/translators"
	"FICSIT-Ordis/internal/id"
	"fmt"
	"log"
)

func newDefault[S id.IDer]() *BasicModule[S] {
	repo := memrepo.New()
	collection, err := repo.GetCollection(fmt.Sprintf("%T", *new(S)))
	if err != nil {
		log.Fatalf("Something went horribly wrong and we could not create a new collection in the memrepo: %v", err)
	}
	translator := translators.Wrap[S](collection)
	return New[S](translator)
}

func New[S id.IDer](collection repos.TypedCollection[S]) *BasicModule[S] {
	return &BasicModule[S]{
		Collection: collection,
	}
}

type BasicModule[S id.IDer] struct {
	Collection repos.TypedCollection[S]
}

func (mod *BasicModule[S]) Create(cmd S) error {
	_, err := mod.Get(cmd.ID())
	if err == nil {
		return fmt.Errorf("element with ID '%v' already exists", cmd.ID())
	}
	err = mod.Collection.Create(cmd)
	if err != nil {
		return fmt.Errorf("could not create a new element: %w", err)
	}
	return nil
}

func (mod *BasicModule[S]) Get(id string) (S, error) {
	cmd, err := mod.Collection.Get(id)
	if err != nil {
		return *new(S), fmt.Errorf("could not get command: %v", err)
	}
	return cmd, nil
}

func (mod *BasicModule[S]) List() ([]S, error) {
	elems, err := mod.Collection.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not get all the elements: %w", err)
	}
	return elems, nil
}

func (mod *BasicModule[S]) Delete(name string) error {
	err := mod.Collection.Delete(name)
	if err != nil {
		return fmt.Errorf("could not delete the element: %w", err)
	}
	return nil
}

func (mod *BasicModule[S]) Update(name string, newcmd S) error {
	err := mod.Collection.Update(name, newcmd)
	if err != nil {
		return fmt.Errorf("could not update the element: %w", err)
	}
	return nil
}
