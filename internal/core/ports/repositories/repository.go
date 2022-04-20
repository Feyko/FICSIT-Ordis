package repositories

import "FICSIT-Ordis/internal/identifiable"

type Repository[E identifiable.Identifiable] interface {
	Find(ID string) (E, error)
	GetAll() ([]E, error)

	Create(element E) error
	Update(ID string, newElement E) error
	Delete(ID string) error
}
