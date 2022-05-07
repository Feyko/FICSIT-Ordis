package translators

import (
	"FICSIT-Ordis/internal/core/ports/repos"
	"FICSIT-Ordis/internal/id"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

func Wrap[E id.IDer](c repos.UntypedCollection) repos.TypedCollection[E] {
	translator := new(Translator[E])
	translator.c = c
	return translator
}

type Translator[E id.IDer] struct {
	c repos.UntypedCollection
}

func (t Translator[E]) Get(ID string) (E, error) {
	v, err := t.c.Get(ID)
	if err != nil {
		return *new(E), err
	}
	typed, err := retype[E](v)
	if err != nil {
		return typed, fmt.Errorf("could not translate the result: %w", err)
	}
	return typed, nil
}

func (t Translator[E]) GetAll() ([]E, error) {
	v, err := t.c.GetAll()
	if err != nil {
		return nil, err
	}
	typed, err := assertSliceTypes[E](v)
	if err != nil {
		return nil, fmt.Errorf("could not translate: %w", err)
	}
	return typed, nil
}

func (t Translator[E]) Search(search string, fields []string) ([]E, error) {
	v, err := t.c.Search(search, fields)
	if err != nil {
		return nil, err
	}
	typed, err := assertSliceTypes[E](v)
	if err != nil {
		return nil, fmt.Errorf("could not translate: %w", err)
	}
	return typed, nil
}

func (t Translator[E]) Create(element E) error {
	return t.c.Create(element)
}

func (t Translator[E]) Update(ID string, newElement E) error {
	return t.c.Update(ID, newElement)
}

func (t Translator[E]) Delete(ID string) error {
	return t.c.Delete(ID)
}

func retype[T any](v any) (T, error) {
	typed, ok := v.(T)
	if ok {
		return typed, nil
	}
	var r T
	err := mapstructure.Decode(v, &r)
	if err != nil {
		return *new(T), fmt.Errorf("could not decode into the struct: %w", err)
	}
	return r, nil
}

// This can result in quite a bit of "wasted" computation.
// But hopefully Search should result in short slices that don't take a lot of computation and GetAll should be used sparsly
func assertSliceTypes[E id.IDer](s []id.IDer) ([]E, error) {
	r := make([]E, len(s))
	for i := 0; i < len(s); i++ {
		elem := s[i]
		typed, err := retype[E](elem)
		if err != nil {
			return nil, err
		}
		r[i] = typed
	}
	return r, nil
}
