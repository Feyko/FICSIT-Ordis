package arango

import (
	"FICSIT-Ordis/internal/core/config"
	"FICSIT-Ordis/internal/storable"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

func New[E storable.I](conf config.ArangoConfig, collectionName string) (*Repository[E], error) {
	repo := new(Repository[E])
	client, err := connectClient(conf)
	if err != nil {
		return nil, err
	}
	db, err := getDatabaseSafe(client, conf)
	if err != nil {
		return nil, fmt.Errorf("could not get the database '%v': %w", conf.DBName, err)
	}

	collection, err := getCollectionSafe(db, collectionName)
	if err != nil {
		return nil, fmt.Errorf("could not get the collection '%v': %w", collectionName, err)
	}

	repo.client = client
	repo.db = db
	repo.collection = collection
	return repo, nil
}

func connectClient(conf config.ArangoConfig) (driver.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{Endpoints: conf.Endpoints})
	if err != nil {
		return nil, fmt.Errorf("could not connect to the endpoint: %w", err)
	}
	client, err := driver.NewClient(driver.ClientConfig{Connection: conn, Authentication: driver.BasicAuthentication(conf.Username, conf.Password)})
	if err != nil {
		return nil, fmt.Errorf("could not connect to the DB: %w", err)
	}
	return client, nil
}

func getDatabaseSafe(client driver.Client, conf config.ArangoConfig) (driver.Database, error) {
	err := ensureDatabaseExists(client, conf)
	if err != nil {
		return nil, fmt.Errorf("could not ensure the database '%v' exists: %w", conf.DBName, err)
	}

	db, err := client.Database(nil, conf.DBName)
	if err != nil {
		return nil, fmt.Errorf("could not get the database '%v': %w", conf.DBName, err)
	}
	return db, nil
}

func ensureDatabaseExists(client driver.Client, conf config.ArangoConfig) error {
	exists, err := client.DatabaseExists(nil, conf.DBName)
	if err != nil {
		return fmt.Errorf("could not check if the database '%v' exists: %w", conf.DBName, err)
	}

	if exists {
		return nil
	}

	client, err = driver.NewClient(driver.ClientConfig{Connection: client.Connection(), Authentication: driver.BasicAuthentication("root", conf.RootPassword)})
	if err != nil {
		return fmt.Errorf("could not connect to the DB as root: %w", err)
	}

	_, err = client.CreateDatabase(nil, conf.DBName, &driver.CreateDatabaseOptions{
		Users:   []driver.CreateDatabaseUserOptions{{UserName: conf.Username, Password: conf.Password}},
		Options: driver.CreateDatabaseDefaultOptions{},
	})
	if err != nil {
		return fmt.Errorf("could not create the database '%v': %w", conf.DBName, err)
	}
	return nil
}

func getCollectionSafe(db driver.Database, collectionName string) (driver.Collection, error) {
	err := ensureCollectionExists(db, collectionName)
	if err != nil {
		return nil, fmt.Errorf("could not ensure the collection '%v' exists: %w", err)
	}

	collection, err := db.Collection(nil, collectionName)
	if err != nil {
		return nil, fmt.Errorf("could not get the collection '%v': %w", collection, err)
	}
	return collection, nil
}

func ensureCollectionExists(db driver.Database, collectionName string) error {
	exists, err := db.CollectionExists(nil, collectionName)

	if exists {
		return nil
	}
	_, err = db.CreateCollection(nil, collectionName, nil)
	if err != nil {
		return fmt.Errorf("could not create the collection '%v': %w", collectionName, err)
	}
	return nil
}

type Repository[E storable.I] struct {
	client     driver.Client
	db         driver.Database
	collection driver.Collection
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

func (r Repository[E]) Search(search string, fields []string) ([]E, error) {
	//TODO implement me
	panic("implement me")
}
