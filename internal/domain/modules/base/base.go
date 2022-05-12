package base

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/memrepo"
	"FICSIT-Ordis/internal/ports/repos/translators"
	"fmt"
	"log"
)

func newDefault[S id.IDer]() *Module[S] {
	repo := memrepo.New()
	collection, err := repo.GetCollection(fmt.Sprintf("%T", *new(S)))
	if err != nil {
		log.Fatalf("Something went horribly wrong and we could not create a new collection in the memrepo: %v", err)
	}
	translator := translators.Wrap[S](collection)
	return New[S](translator)
}

func New[S id.IDer](collection repos.TypedCollection[S]) *Module[S] {
	return &Module[S]{
		Collection: collection,
	}
}

type Module[S id.IDer] struct {
	Collection repos.TypedCollection[S]
}

func (mod *Module[S]) Create(cmd S) error {
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

func (mod *Module[S]) Get(id string) (S, error) {
	cmd, err := mod.Collection.Get(id)
	if err != nil {
		return *new(S), fmt.Errorf("could not get command: %v", err)
	}
	return cmd, nil
}

func (mod *Module[S]) List() ([]S, error) {
	elems, err := mod.Collection.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not get all the elements: %w", err)
	}
	return elems, nil
}

func (mod *Module[S]) Delete(name string) error {
	err := mod.Collection.Delete(name)
	if err != nil {
		return fmt.Errorf("could not delete the element: %w", err)
	}
	return nil
}

func (mod *Module[S]) Update(name string, newcmd S) error {
	err := mod.Collection.Update(name, newcmd)
	if err != nil {
		return fmt.Errorf("could not update the element: %w", err)
	}
	return nil
}
