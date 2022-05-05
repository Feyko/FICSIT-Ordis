package base

import (
	"FICSIT-Ordis/internal/core/ports/repos"
	"FICSIT-Ordis/internal/id"
	"fmt"
)

func New[S id.IDer](collection repos.Collection) *BasicModule[S] {
	return &BasicModule[S]{
		Collection: collection,
	}
}

type BasicModule[S id.IDer] struct {
	Collection repos.Collection
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
