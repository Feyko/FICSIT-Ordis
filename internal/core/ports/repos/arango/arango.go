package arango

import (
	"FICSIT-Ordis/internal/core/config"
	"FICSIT-Ordis/internal/core/ports/repos"
	"FICSIT-Ordis/internal/id"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
)

type Repository struct {
	client driver.Client
	db     driver.Database
}

func New(conf config.ArangoConfig) (*Repository, error) {
	repo := new(Repository)
	client, err := connectClient(conf)
	if err != nil {
		return nil, err
	}
	db, err := client.Database(nil, conf.DBName)
	if err != nil {
		return nil, fmt.Errorf("could not get the database '%v': %w", conf.DBName, err)
	}

	repo.client = client
	repo.db = db
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
	authed, err := authCheck(client, conf)
	if err != nil {
		return nil, fmt.Errorf("authentication test failed: %w", err)
	}
	if !authed {
		err := superInit(client.Connection(), conf)
		if err != nil {
			return nil, fmt.Errorf("could not create the database and user: %w", err)
		}
	}

	return client, nil
}

func authCheck(client driver.Client, conf config.ArangoConfig) (bool, error) {
	_, err := client.DatabaseExists(nil, conf.DBName)
	if driver.IsUnauthorized(err) {
		return false, nil
	}
	if driver.IsForbidden(err) {
		return false, fmt.Errorf("the user '%v' has insufficient permissions to access the db '%v'", conf.Username, conf.DBName)
	}
	return true, nil
}

func superInit(conn driver.Connection, conf config.ArangoConfig) error {
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(conf.SuperUsername, conf.SuperPassword),
	})
	if err != nil {
		return fmt.Errorf("could not connect as the superuser '%v': %w", conf.SuperUsername, err)
	}
	_, err = client.CreateDatabase(nil, conf.DBName, &driver.CreateDatabaseOptions{Users: []driver.CreateDatabaseUserOptions{{
		UserName: conf.Username,
		Password: conf.Password,
	},
	}})

	if driver.IsUnauthorized(err) {
		return fmt.Errorf("invalid credentials for superuser '%v'", conf.SuperUsername)
	}

	if driver.IsForbidden(err) {
		return fmt.Errorf("superuser '%v' has insufficient permissions to create a new database", conf.SuperUsername)
	}
	if err != nil {
		return fmt.Errorf("could not create the db '%v': %w", conf.DBName, err)
	}
	return nil
}

func (r Repository) GetCollection(name string) (repos.UntypedCollection, error) {
	collection, err := getCollectionSafe(r.db, name)
	if err != nil {
		return nil, fmt.Errorf("could not get the collection '%v': %w", name, err)
	}
	return newCollection(collection), nil
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

type Collection struct {
	collection driver.Collection
}

func newCollection(collection driver.Collection) Collection {
	return Collection{
		collection: collection,
	}
}

func (r Collection) Get(ID string) (id.IDer, error) {
	//TODO implement me
	panic("implement me")
}

func (r Collection) GetAll() ([]id.IDer, error) {
	//TODO implement me
	panic("implement me")
}

func (r Collection) Create(element id.IDer) error {
	//TODO implement me
	panic("implement me")
}

func (r Collection) Update(ID string, newElement id.IDer) error {
	//TODO implement me
	panic("implement me")
}

func (r Collection) Delete(ID string) error {
	//TODO implement me
	panic("implement me")
}

func (r Collection) Search(search string, fields []string) ([]id.IDer, error) {
	//TODO implement me
	panic("implement me")
}
