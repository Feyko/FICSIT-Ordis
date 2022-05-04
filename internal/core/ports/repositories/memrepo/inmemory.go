package memrepo

import (
	"FICSIT-Ordis/internal/storable"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func New[E storable.I]() *Repository[E] {
	return new(Repository[E])
}

type Repository[E storable.I] struct {
	elements []E
}

func (repo *Repository[E]) Find(ID string) (E, error) {
	elem, _, err := repo.findWithIndex(ID)
	return elem, err
}

func (repo *Repository[E]) findWithIndex(ID string) (E, int, error) {
	for i, elem := range repo.elements {
		if elem.ID() == ID {
			return elem, i, nil
		}
	}
	return *new(E), 0, fmt.Errorf("element with ID '%v' does not exist", ID)
}

func (repo *Repository[E]) GetAll() ([]E, error) {
	r := make([]E, len(repo.elements))
	copy(r, repo.elements)
	return r, nil
}

//Terrible code. Need to refactor this asap
func (repo *Repository[E]) Search(search string, fields []string) ([]E, error) {
	var r []E
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
			if field.IsZero() { // This can skip valid values.. let's hope it doesn't. I hate Go's zero value thing
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

func (repo *Repository[E]) Create(element E) error {
	_, err := repo.Find(element.ID())
	if err == nil {
		return fmt.Errorf("element with ID '%v' already exists", element.ID())
	}
	repo.elements = append(repo.elements, element)
	return nil
}

func (repo *Repository[E]) Update(ID string, newElement E) error {
	_, i, err := repo.findWithIndex(ID)
	if err != nil {
		return err
	}
	repo.elements[i] = newElement
	return nil
}

func (repo *Repository[E]) Delete(ID string) error {
	_, i, err := repo.findWithIndex(ID)
	if err != nil {
		return err
	}
	repo.elements = removeIndex(repo.elements, i)
	return nil
}

func removeIndex[E any](s []E, i int) []E {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
