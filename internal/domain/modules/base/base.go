package base

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"fmt"
)

func New[E id.IDer](collection repo.Collection[E]) *Module[E] {
	return &Module[E]{
		Collection: collection,
	}
}

type Module[E id.IDer] struct {
	Collection repo.Collection[E]
}

func (mod *Module[E]) Create(element E) error {
	_, err := mod.Get(element.ID())
	if err == nil {
		return fmt.Errorf("element with ID '%v' already exists", element.ID())
	}
	err = mod.Collection.Create(element)
	if err != nil {
		return fmt.Errorf("could not create a new element: %w", err)
	}
	return nil
}

func (mod *Module[E]) Get(ID string) (E, error) {
	cmd, err := mod.Collection.Get(ID)
	if err != nil {
		return *new(E), fmt.Errorf("could not get the command with ID '%v': %v", ID, err)
	}
	return cmd, nil
}

func (mod *Module[E]) List() ([]E, error) {
	elems, err := mod.Collection.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not get all the elements: %w", err)
	}
	return elems, nil
}

func (mod *Module[E]) Delete(ID string) error {
	err := mod.Collection.Delete(ID)
	if err != nil {
		return fmt.Errorf("could not delete the element: %w", err)
	}
	return nil
}

func (mod *Module[E]) Update(ID string, updateElement id.IDer) error {
	err := mod.Collection.Update(ID, updateElement)
	if err != nil {
		return fmt.Errorf("could not update the element: %w", err)
	}
	return nil
}
