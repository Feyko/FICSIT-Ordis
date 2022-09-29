package arango

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/id"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type Repository[T id.IDer] struct {
	client driver.Client
	db     driver.Database
}

func New[T id.IDer](conf config.ArangoConfig) (*Repository[T], error) {
	repo := new(Repository[T])
	client, err := connectClient(conf)
	if err != nil {
		return nil, err
	}
	db, err := client.Database(nil, conf.DBName)
	if driver.IsNotFound(err) {
		db, err = client.CreateDatabase(nil, conf.DBName, &driver.CreateDatabaseOptions{
			Users: []driver.CreateDatabaseUserOptions{{
				UserName: conf.Username,
				Password: conf.Password,
			}},
		})
		if err != nil {
			return nil, errors.Wrap(err, "could not create the database")
		}
	}
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

func (r Repository[T]) GetCollection(name string) (any, error) {
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

type Collection[T id.IDer] struct {
	c  driver.Collection
	db driver.Database
}

func newCollection[T id.IDer](collection driver.Collection, db driver.Database) Collection[T] {
	return Collection[T]{
		c:  collection,
		db: db,
	}
}

func (c Collection[T]) Get(ID string) (T, error) {
	query :=
		`filter doc.id == @id
return doc`
	elements, err := runQueryInCollection(c, query, map[string]interface{}{
		"id": ID,
	})
	if err != nil {
		return nil, fmt.Errorf("could not read the document: %w", err)
	}
	if len(elements) == 0 {
		return nil, errors.Errorf("could not find the element with ID %v", ID)
	}
	return elements[0], nil
}

func (c Collection[T]) GetAll() ([]T, error) {
	return runQueryInCollection(c, "return doc", nil)
}

func (c Collection[T]) Create(element id.IDer) error {
	_, err := c.c.CreateDocument(nil, element)
	if err != nil {
		return fmt.Errorf("could not create the document: %w", err)
	}
	return nil
}

func (c Collection[T]) Update(ID string, updateElement id.IDer) error {
	asMap, err := id.ToMapNoOverwrite(updateElement)
	if err != nil {
		return errors.Wrap(err, "could not turn the element into a map")
	}
	asMap["id"] = ID
	query := updateQueryForType(updateElement, "coll", "doc")
	_, err = runQueryInCollection(c, query, asMap)
	if err != nil {
		return errors.Wrap(err, "could not update the document")
	}
	return nil
}

func (c Collection[T]) Delete(ID string) error {
	query :=
		`filter doc.id == @id
remove doc in @@coll`
	_, err := runQueryInCollection(c, query, map[string]interface{}{
		"id": ID,
	})
	if err != nil {
		return errors.Wrap(err, "could not delete the element")
	}
	return nil
}

func (c Collection[T]) Search(search string, fields []string) ([]T, error) {
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

func runQueryInCollection[T id.IDer](coll Collection[T], query string, params map[string]any) ([]T, error) {
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
	return flattenCursor[T](cursor)
}

func flattenCursor[T id.IDer](cursor driver.Cursor) ([]T, error) {
	defer cursor.Close()
	count := int(cursor.Count())
	r := make([]T, count)
	for i := 0; i < count; i++ {
		_, err := cursor.ReadDocument(nil, &r[i])
		if err != nil {
			return nil, fmt.Errorf("could not read the document: %w", err)
		}
	}

	return r, nil
}

func updateQueryForType[T any](v T, collParam, elemParam string) string {
	query := fmt.Sprintf("update @@%v with { ", elemParam)
	t := reflect.TypeOf(v)
	numField := t.NumField()
	fieldNames := make([]string, 0, numField)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldNames = append(fieldNames, field.Name)
	}
	for _, name := range fieldNames {
		query += fmt.Sprintf("%v: @%v", name, name)
	}
	query += fmt.Sprintf(" } in @@%v", collParam)
	return query
}
