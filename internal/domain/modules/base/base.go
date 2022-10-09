package base

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"context"
	"fmt"
	"github.com/pkg/errors"
)

func New[E id.IDer](collection repo.Collection[E]) *Module[E] {
	return &Module[E]{
		Collection: collection,
	}
}

type Module[E id.IDer] struct {
	Collection repo.Collection[E]
}

func (mod *Module[E]) Create(ctx context.Context, element E) error {
	_, err := mod.Get(nil, element.ID())
	if err == nil {
		return fmt.Errorf("element with ID '%v' already exists", element.ID())
	}
	err = mod.Collection.Create(ctx, element)
	if err != nil {
		return errors.Wrap(err, "could not create a new element")
	}
	return nil
}

func (mod *Module[E]) Get(ctx context.Context, ID string) (E, error) {
	cmd, err := mod.Collection.Get(ctx, ID)
	if err != nil {
		return *new(E), errors.Wrapf(err, "could not get the command with ID '%v'", ID)
	}
	return cmd, nil
}

func (mod *Module[E]) List(ctx context.Context) ([]E, error) {
	elems, err := mod.Collection.GetAll(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "could not get all the elements")
	}
	return elems, nil
}

func (mod *Module[E]) Delete(ctx context.Context, ID string) error {
	err := mod.Collection.Delete(ctx, ID)
	if err != nil {
		return errors.Wrap(err, "could not delete the element")
	}
	return nil
}

func (mod *Module[E]) Update(ctx context.Context, ID string, updateElement any) (E, error) {
	elem, err := mod.Collection.Update(ctx, ID, updateElement)
	if err != nil {
		return *new(E), errors.Wrap(err, "could not update the element")
	}
	return elem, nil
}
