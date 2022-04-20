package repositories

import (
	"FICSIT-Ordis/internal/identifiable"
	"fmt"
)

type MemoryRepository[E identifiable.Identifiable] struct {
	elements []E
}

func (repo *MemoryRepository[E]) Find(ID string) (E, error) {
	elem, _, err := repo.findWithIndex(ID)
	return elem, err
}

func (repo *MemoryRepository[E]) findWithIndex(ID string) (E, int, error) {
	for i, elem := range repo.elements {
		if elem.ID() == ID {
			return elem, i, nil
		}
	}
	return *new(E), 0, fmt.Errorf("element with ID '%v' does not exist", ID)
}

func (repo *MemoryRepository[E]) GetAll() ([]E, error) {
	r := make([]E, len(repo.elements))
	copy(r, repo.elements)
	return r, nil
}

func (repo *MemoryRepository[E]) Create(element E) error {
	_, err := repo.Find(element.ID())
	if err == nil {
		return fmt.Errorf("element with ID '%v' already exists", element.ID())
	}
	repo.elements = append(repo.elements, element)
	return nil
}

func (repo *MemoryRepository[E]) Update(ID string, newElement E) error {
	_, i, err := repo.findWithIndex(ID)
	if err != nil {
		return err
	}
	repo.elements[i] = newElement
	return nil
}

func (repo *MemoryRepository[E]) Delete(ID string) error {
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
