package memrepo

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"FICSIT-Ordis/internal/util"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"reflect"
)

func New() repo.Repository[id.IDer] {
	return &Repository[id.IDer]{make(map[string]any)}
}

type Repository[T id.IDer] struct {
	collections map[string]any
}

func (r *Repository[T]) CreateCollection(name string) (any, error) {
	_, ok := r.collections[name]
	if ok {
		return nil, errors.New(fmt.Sprintf("collection '%v' already exists", name))
	}
	collection := repo.Collection[T](newCollection[T]())
	r.collections[name] = collection
	return collection, nil
}

func (r *Repository[T]) GetCollection(name string) (any, error) {
	collection, ok := r.collections[name]
	if !ok {
		return nil, repo.ErrCollectionNotFound
	}
	typed, ok := collection.(repo.Collection[T])
	if !ok {
		return nil, errors.New(fmt.Sprintf("collection '%v' does not hold type '%T'", name, *new(T)))
	}
	return typed, nil
}

func (r *Repository[T]) DeleteCollection(name string) error {
	if _, ok := r.collections[name]; !ok {
		return errors.Errorf("No collection named '%v'", name)
	}
	delete(r.collections, name)
	return nil
}

func newCollection[T id.IDer]() *Collection[T] {
	return &Collection[T]{elements: make([]T, 0)}
}

type Collection[T id.IDer] struct {
	elements []T
}

func (coll *Collection[T]) Get(ctx context.Context, ID string) (T, error) {
	elem, _, err := coll.findWithIndex(ID)
	return elem, err
}

func (coll *Collection[T]) findWithIndex(ID string) (T, int, error) {
	for i, elem := range coll.elements {
		if elem.ID() == ID {
			return elem, i, nil
		}
	}
	return *new(T), 0, repo.ErrElementNotFound
}

func (coll *Collection[T]) GetAll(ctx context.Context) ([]T, error) {
	return slices.Clone(coll.elements), nil
}

//Terrible code. Need to refactor this asap
func (coll *Collection[T]) Search(ctx context.Context, search string, fields []string) ([]T, error) {
	var r []T
	for _, e := range coll.elements {
		reflected := reflect.ValueOf(e)
		if reflected.Kind() != reflect.Struct {
			if reflected.Kind() == reflect.String {
				if reflected.String() == search {
					r = append(r, e)
					continue
				}
			}
			return nil, errors.New("trying to search an invalid type. Search only supports structs and strings")
		}

		for _, fieldName := range fields {
			field := reflected.FieldByName(fieldName)
			if field.Kind() == reflect.Array || field.Kind() == reflect.Slice {
				fieldLen := field.Len()
				for i := 0; i < fieldLen; i++ {
					fieldValue := field.Index(i)
					fieldString := fieldValue.String()
					if fieldString == search {
						r = append(r, e)
						continue
					}
				}
			}
			fieldString := field.String() // Might be too broad
			if fieldString == search {
				r = append(r, e)
			}
		}
	}
	return r, nil
}

func (coll *Collection[T]) Create(ctx context.Context, element T) error {
	_, err := coll.Get(ctx, element.ID())
	if err == nil {
		return fmt.Errorf("element with ID '%v' already exists", element.ID())
	}
	coll.elements = append(coll.elements, element)
	return nil
}

func (coll *Collection[T]) Update(ctx context.Context, ID string, updateElement any) (oldElem T, newElem T, err error) {
	found, i, err := coll.findWithIndex(ID)
	if err != nil {
		return *new(T), *new(T), errors.Wrap(err, "could not get the element")
	}
	oldElem = found
	err = util.PatchStruct(&coll.elements[i], updateElement)
	if err != nil {
		return *new(T), *new(T), errors.Wrap(err, "could not update the element")
	}
	return oldElem, coll.elements[i], nil
}

func (coll *Collection[T]) Delete(ctx context.Context, ID string) error {
	_, i, err := coll.findWithIndex(ID)
	if err != nil {
		return err
	}
	coll.elements = removeIndex(coll.elements, i)
	return nil
}

func removeIndex[T id.IDer](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
