package memrepo

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/translators"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func New() Repository {
	return Repository{make(map[string]repos.UntypedCollection)}
}

type Repository struct {
	collections map[string]repos.UntypedCollection
}

func (r *Repository) newCollection(name string) (repos.UntypedCollection, error) {
	collection := newCollection()
	r.collections[name] = collection
	return collection, nil
}

func (r Repository) GetCollection(name string) (repos.UntypedCollection, error) {
	collection, ok := r.collections[name]
	if !ok {
		collection, _ = r.newCollection(name)
	}
	return collection, nil
}

func newCollection() *Collection {
	return new(Collection)
}

type Collection struct {
	elements []id.IDer
}

func (repo *Collection) Get(ID string) (any, error) {
	elem, _, err := repo.findWithIndex(ID)
	return elem, err
}

func (repo *Collection) findWithIndex(ID string) (id.IDer, int, error) {
	for i, elem := range repo.elements {
		if elem.ID() == ID {
			return elem, i, nil
		}
	}
	return *new(id.IDer), 0, fmt.Errorf("element with ID '%v' does not exist", ID)
}

func (repo *Collection) GetAll() ([]any, error) {
	r, err := translators.RetypeSlice[any](repo.elements)
	if err != nil {
		return nil, fmt.Errorf("could not retype the inner slice: %w", err)
	}
	return r, nil
}

//Terrible code. Need to refactor this asap
func (repo *Collection) Search(search string, fields []string) ([]any, error) {
	var r []any
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

func (repo *Collection) Create(element id.IDer) error {
	_, err := repo.Get(element.ID())
	if err == nil {
		return fmt.Errorf("element with ID '%v' already exists", element.ID())
	}
	repo.elements = append(repo.elements, element)
	return nil
}

func (repo *Collection) Update(ID string, updateElement id.IDer) error {
	current, err := repo.Get(ID)
	if err != nil {
		return errors.Wrap(err, "could not get the element")
	}
	asMap, err := id.ToMap(current.(id.IDer), "id")
	if err != nil {
		return err
	}
	err = id.Update(asMap, updateElement)
	if err != nil {
		return err
	}
}

func (repo *Collection) Delete(ID string) error {
	_, i, err := repo.findWithIndex(ID)
	if err != nil {
		return err
	}
	repo.elements = removeIndex(repo.elements, i)
	return nil
}

func removeIndex(s []id.IDer, i int) []id.IDer {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
