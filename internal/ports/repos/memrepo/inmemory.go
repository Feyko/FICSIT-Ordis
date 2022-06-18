package memrepo

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"reflect"
	"strings"
)

func New() Repository {
	return Repository{make(map[string]any)}
}

type Repository struct {
	collections map[string]any
}

func NewCollection[T id.IDer](repo *Repository, name string) (repos.TypedCollection[T], error) {
	collection := newCollection[T]()
	repo.collections[name] = collection
	return collection, nil
}

func (r Repository) GetCollection(name string) (any, error) {
	collection, ok := r.collections[name]
	if !ok {
		return nil, errors.New("collection not found")
	}
	return collection, nil
}

func newCollection[T id.IDer]() *Collection[T] {
	return new(Collection[T])
}

type Collection[T id.IDer] struct {
	elements []T
}

func (repo *Collection[T]) Get(ID string) (T, error) {
	elem, _, err := repo.findWithIndex(ID)
	return elem, err
}

func (repo *Collection[T]) findWithIndex(ID string) (T, int, error) {
	for i, elem := range repo.elements {
		if elem.ID() == ID {
			return elem, i, nil
		}
	}
	return *new(T), 0, fmt.Errorf("element with ID '%v' does not exist", ID)
}

func (repo *Collection[T]) GetAll() ([]T, error) {
	return slices.Clone(repo.elements), nil
}

//Terrible code. Need to refactor this asap
func (repo *Collection[T]) Search(search string, fields []string) ([]T, error) {
	var r []T
	for _, e := range repo.elements {
		reflected := reflect.ValueOf(e)
		if reflected.Kind() != reflect.Struct {
			if reflected.Kind() == reflect.String {
				if strings.Contains(reflected.String(), search) {
					r = append(r, e)
					continue
				}
			}
			return nil, errors.New("trying to search an invalid type. Search only supports structs and strings")
		}
		for _, fieldName := range fields {
			field := reflected.FieldByName(fieldName)
			if field.IsZero() { // This can skip valid values.. let's hope it doesn't. IDer hate Go's zero value thing
				continue
			}
			fieldString := field.String() // Might be too broad and a bad fix for "deep" search
			if strings.Contains(fieldString, search) {
				r = append(r, e)
			}
		}
	}
	return r, nil
}

func (repo *Collection[T]) Create(element T) error {
	_, err := repo.Get(element.ID())
	if err == nil {
		return fmt.Errorf("element with ID '%v' already exists", element.ID())
	}
	repo.elements = append(repo.elements, element)
	return nil
}

func (repo *Collection[T]) Update(ID string, updateElement repos.Updater[T]) error {
	current, err := repo.Get(ID)
	if err != nil {
		return errors.Wrap(err, "could not get the element")
	}
	err = id.Update(&current, updateElement)
	if err != nil {
		return err
	}

	return nil
}

func (repo *Collection[T]) Delete(ID string) error {
	_, i, err := repo.findWithIndex(ID)
	if err != nil {
		return err
	}
	repo.elements = removeIndex(repo.elements, i)
	return nil
}

func removeIndex[T id.IDer](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
