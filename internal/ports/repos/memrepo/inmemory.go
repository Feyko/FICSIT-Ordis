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
	"strings"
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
func (coll *Collection[T]) Search(ctx context.Context, search string) ([]T, error) {
	var found []T
	for _, elem := range coll.elements {
		matched, err := searchInElement(elem, search)
		if err != nil {
			return nil, errors.Wrap(err, "error searching an element")
		}
		if matched {
			found = append(found, elem)
		}
	}
	return found, nil
}

func searchInElement(elem any, search string) (bool, error) {
	typeInfo, err := repo.GetTypeInfo(elem)
	if err != nil {
		return false, errors.Wrap(err, "error getting type info")
	}

	value := reflect.ValueOf(elem)

	for _, field := range typeInfo.Fields {
		if field.ToSearch {
			matched := searchInField(value.FieldByIndex(field.Index), search)
			if matched {
				return true, nil
			}
		}
	}
	return false, nil
}

func searchInField(value reflect.Value, search string) bool {
	switch value.Kind() {
	case reflect.Pointer, reflect.Interface:
		return searchInField(value.Elem(), search)
	case reflect.Func, reflect.Uintptr, reflect.UnsafePointer, reflect.Chan, reflect.Invalid:
		return false
	case reflect.Array, reflect.Slice:
		length := value.Len()
		for i := 0; i < length; i++ {
			matched := searchInField(value.Index(i), search)
			if matched {
				return true
			}
		}
	case reflect.Map:
		iter := value.MapRange()
		for iter.Next() {
			matched := searchInField(iter.Value(), search)
			if matched {
				return true
			}
		}
	case reflect.Struct:
		return false // We search in all sub-structs too, no need for further handling here
	default:
		return strings.Contains(value.String(), search)
	}

	return false
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
