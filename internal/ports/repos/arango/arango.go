package arango

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"strings"
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
		return nil, fmt.Errorf("could not get the c '%v': %w", name, err)
	}
	return newCollection(collection, r.db), nil
}

func getCollectionSafe(db driver.Database, collectionName string) (driver.Collection, error) {
	err := ensureCollectionExists(db, collectionName)
	if err != nil {
		return nil, fmt.Errorf("could not ensure the collection '%v' exists: %w", collectionName, err)
	}

	collection, err := db.Collection(nil, collectionName)
	if err != nil {
		return nil, fmt.Errorf("could not get the c '%v': %w", collection, err)
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
		return fmt.Errorf("could not create the c '%v': %w", collectionName, err)
	}
	return nil
}

type Collection struct {
	c  driver.Collection
	db driver.Database
}

func newCollection(collection driver.Collection, db driver.Database) Collection {
	return Collection{
		c:  collection,
		db: db,
	}
}

func (c Collection) Get(ID string) (any, error) {
	var r map[string]any
	_, err := c.c.ReadDocument(nil, ID, &r)
	if err != nil {
		return nil, fmt.Errorf("could not read the document: %w", err)
	}
	return r, nil
}

func (c Collection) GetAll() ([]any, error) {
	return runQueryInCollection(c, "return doc", nil)
}

func (c Collection) Create(element id.IDer) error {
	asMap, err := id.ToMap(element, "_key")
	if err != nil {
		return err
	}
	_, err = c.c.CreateDocument(nil, asMap)

	if err != nil {
		return fmt.Errorf("could not create the document: %w", err)
	}
	return nil
}

func (c Collection) Update(ID string, newElement id.IDer) error {
	// We remove them recreate the document to make sure the _key and the ID of our value are synced
	_, err := c.c.RemoveDocument(nil, ID)
	if err != nil {
		return fmt.Errorf("could not delete the old document: %w", err)
	}
	asMap, err := id.ToMap(newElement, "_key")
	if err != nil {
		return err
	}
	_, err = c.c.CreateDocument(nil, asMap)
	if err != nil {
		return fmt.Errorf("could not create the new document: %w", err)
	}
	return nil
}

func (c Collection) Delete(ID string) error {
	_, err := c.c.RemoveDocument(nil, ID)
	if err != nil {
		return fmt.Errorf("could not delete the document: %w", err)
	}
	return nil
}

func (c Collection) Search(search string, fields []string) ([]any, error) {
	params := map[string]any{}

	filters := make([]string, len(fields))
	for i, field := range fields {
		filters[i] = fmt.Sprintf(`doc.%v like "%%%v%%"`, field, search)
	}
	query := "\tfilter "
	query += strings.Join(filters, " || ")
	query += "\n\treturn doc"
	return runQueryInCollection(c, query, params)
}

func runQueryInCollection(coll Collection, query string, params map[string]any) ([]any, error) {
	if params == nil {
		params = make(map[string]any)
	}
	params["@coll"] = coll.c.Name()
	ctx := driver.WithQueryCount(context.Background(), true)
	query = "for doc in @@coll\n" + query
	cursor, err := coll.db.Query(ctx, query, params)
	if err != nil {
		return nil, fmt.Errorf("could not query the database: %w", err)
	}
	return flattenCursor(cursor)
}

func flattenCursor(cursor driver.Cursor) ([]any, error) {
	defer cursor.Close()
	count := int(cursor.Count())
	r := make([]any, count)
	for i := 0; i < count; i++ {
		_, err := cursor.ReadDocument(nil, &r[i])
		if err != nil {
			return nil, fmt.Errorf("could not read the document: %w", err)
		}
	}

	return r, nil
}
