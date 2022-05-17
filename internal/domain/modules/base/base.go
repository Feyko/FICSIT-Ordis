package base

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"fmt"
)

func New[E id.IDer, U id.IDer](collection repos.TypedCollection[E, U]) *Module[E, U] {
	return &Module[E, U]{
		Collection: collection,
	}
}

type Module[E id.IDer, U id.IDer] struct {
	Collection repos.TypedCollection[E, U]
}

func (mod *Module[E, U]) Create(element E) error {
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

func (mod *Module[E, U]) Get(ID string) (E, error) {
	cmd, err := mod.Collection.Get(ID)
	if err != nil {
		return *new(E), fmt.Errorf("could not get the command with ID '%v': %v", ID, err)
	}
	return cmd, nil
}

func (mod *Module[E, U]) List() ([]E, error) {
	elems, err := mod.Collection.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not get all the elements: %w", err)
	}
	return elems, nil
}

func (mod *Module[E, U]) Delete(ID string) error {
	err := mod.Collection.Delete(ID)
	if err != nil {
		return fmt.Errorf("could not delete the element: %w", err)
	}
	return nil
}

func (mod *Module[E, U]) Update(ID string, updateElement U) error {
	err := mod.Collection.Update(ID, updateElement)
	if err != nil {
		return fmt.Errorf("could not update the element: %w", err)
	}
	return nil
}
