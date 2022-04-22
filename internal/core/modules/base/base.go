package base

import (
	"FICSIT-Ordis/internal/core/ports/repositories"
	"FICSIT-Ordis/internal/identifiable"
	"fmt"
)

func NewDefault[S identifiable.Identifiable]() *BasicModule[S] {
	return New[S](new(repositories.MemoryRepository[S]))
}

func New[S identifiable.Identifiable](repo repositories.Repository[S]) *BasicModule[S] {
	return &BasicModule[S]{
		Repository: repo,
	}
}

type BasicModule[S identifiable.Identifiable] struct {
	Repository repositories.Repository[S]
}

func (mod *BasicModule[S]) Create(cmd S) error {
	_, err := mod.Get(cmd.ID())
	if err == nil {
		return fmt.Errorf("element with ID '%v' already exists", cmd.ID())
	}
	err = mod.Repository.Create(cmd)
	if err != nil {
		return fmt.Errorf("could not create a new element: %w", err)
	}
	return nil
}

func (mod *BasicModule[S]) Get(id string) (S, error) {
	cmd, err := mod.Repository.Find(id)
	if err != nil {
		return *new(S), fmt.Errorf("could not get command: %v", err)
	}
	return cmd, nil
}

func (mod *BasicModule[S]) List() ([]S, error) {
	elems, err := mod.Repository.GetAll()
	if err != nil {
		return nil, fmt.Errorf("could not get all the elements: %w", err)
	}
	return elems, nil
}

func (mod *BasicModule[S]) Delete(name string) error {
	err := mod.Repository.Delete(name)
	if err != nil {
		return fmt.Errorf("could not delete the element: %w", err)
	}
	return nil
}

func (mod *BasicModule[S]) Update(name string, newcmd S) error {
	err := mod.Repository.Update(name, newcmd)
	if err != nil {
		return fmt.Errorf("could not update the element: %w", err)
	}
	return nil
}
