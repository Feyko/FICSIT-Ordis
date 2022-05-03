package arango

import (
	"FICSIT-Ordis/internal/id"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func New[E id.IDer](endpoint string) (*Repository[E], error) {
	conn, err := http.NewConnection(http.ConnectionConfig{Endpoints: []string{endpoint}})
	if err != nil {
		return nil, fmt.Errorf("could not connect to the endpoint: %v", err)
	}
	client, err := driver.NewClient(driver.ClientConfig{Connection: conn})
	if err != nil {
		return nil, fmt.Errorf("could not connect to the DB: %v", err)
	}
	repo := new(Repository[E])
	repo.client = client
	return repo, nil
}

type Repository[E id.IDer] struct {
	client driver.Client
}

func (r Repository[E]) Find(ID string) (E, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository[E]) GetAll() ([]E, error) {
	//TODO implement me
	panic("implement me")
}

func (r Repository[E]) Create(element E) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository[E]) Update(ID string, newElement E) error {
	//TODO implement me
	panic("implement me")
}

func (r Repository[E]) Delete(ID string) error {
	//TODO implement me
	panic("implement me")
}
