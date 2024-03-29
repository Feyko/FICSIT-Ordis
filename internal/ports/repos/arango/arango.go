package arango

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"FICSIT-Ordis/internal/util"
	"context"
	"fmt"
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type Config struct {
	Username,
	Password,
	SuperUsername,
	SuperPassword,
	DBName string
	Endpoints []string
}

type Repository[T id.IDer] struct {
	client driver.Client
	db     driver.Database
}

func New[T id.IDer](conf Config) (repo.Repository[T], error) {
	repo := new(Repository[T])
	client, err := connectClient(conf)
	if err != nil {
		return nil, err
	}
	db, err := client.Database(nil, conf.DBName)
	if driver.IsNotFoundGeneral(err) {
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
		return nil, errors.Wrapf(err, "could not get the database '%v'", conf.DBName)
	}

	repo.client = client
	repo.db = db
	return repo, nil
}

func connectClient(conf Config) (driver.Client, error) {
	conn, err := http.NewConnection(http.ConnectionConfig{Endpoints: conf.Endpoints})
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to the endpoint")
	}
	client, err := driver.NewClient(driver.ClientConfig{Connection: conn, Authentication: driver.BasicAuthentication(conf.Username, conf.Password)})
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to the DB")
	}
	authed, err := authCheck(client, conf)
	if err != nil {
		return nil, errors.Wrap(err, "authentication test failed")
	}
	if !authed {
		err := superInit(client.Connection(), conf)
		if err != nil {
			return nil, errors.Wrap(err, "could not create the database and user")
		}
	}

	return client, nil
}

func authCheck(client driver.Client, conf Config) (bool, error) {
	_, err := client.DatabaseExists(nil, conf.DBName)
	if driver.IsUnauthorized(err) {
		return false, nil
	}
	if driver.IsForbidden(err) {
		return false, fmt.Errorf("the user '%v' has insufficient permissions to access the db '%v'", conf.Username, conf.DBName)
	}
	return true, nil
}

func superInit(conn driver.Connection, conf Config) error {
	client, err := driver.NewClient(driver.ClientConfig{
		Connection:     conn,
		Authentication: driver.BasicAuthentication(conf.SuperUsername, conf.SuperPassword),
	})
	if err != nil {
		return errors.Wrapf(err, "could not connect as the superuser '%v'", conf.SuperUsername)
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
		return errors.Wrapf(err, "could not create the db '%v'", conf.DBName)
	}
	return nil
}

func (r *Repository[T]) CreateCollection(name string) (any, error) {
	collection, err := r.db.CreateCollection(nil, name, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "could not create the collection '%v'", name)
	}
	return newCollection[T](collection, r.db), nil
}

func (r *Repository[T]) GetCollection(name string) (any, error) {
	collection, err := r.db.Collection(nil, name)
	if driver.IsNotFoundGeneral(err) {
		return nil, repo.ErrCollectionNotFound
	}
	if err != nil {
		return nil, errors.Wrapf(err, "could not get the collection '%v'", collection)
	}
	return repo.Collection[T](newCollection[T](collection, r.db)), nil
}

func (r *Repository[T]) DeleteCollection(name string) error {
	c, err := r.GetCollection(name)
	if err != nil {
		return errors.Wrap(err, "could not get the collection")
	}
	err = c.(*Collection[T]).c.Remove(nil)
	if err != nil {
		return err
	}
	return nil
}

type Collection[T id.IDer] struct {
	c  driver.Collection
	db driver.Database
}

func newCollection[T id.IDer](collection driver.Collection, db driver.Database) *Collection[T] {
	return &Collection[T]{
		c:  collection,
		db: db,
	}
}

func (c *Collection[T]) Get(ctx context.Context, ID string) (T, error) {
	query :=
		`filter doc.id == @id
return doc`
	elements, err := runQueryInCollection(ctx, c, query, map[string]interface{}{
		"id": ID,
	})
	if err != nil {
		return *new(T), errors.Wrap(err, "could not read the document")
	}
	if len(elements) == 0 {
		return *new(T), repo.ErrElementNotFound
	}
	return elements[0], nil
}

func (c *Collection[T]) GetAll(ctx context.Context) ([]T, error) {
	return runQueryInCollection(ctx, c, "return doc", nil)
}

func (c *Collection[T]) Create(ctx context.Context, element T) error {
	_, err := c.Get(ctx, element.ID())
	if err == nil {
		return errors.Errorf("element with ID '%v' already exists", element.ID())
	}
	if !errors.Is(err, repo.ErrElementNotFound) {
		return errors.Wrap(err, "error checking existing element")
	}

	asMap, err := id.ToMap(element)
	if err != nil {
		return errors.Wrap(err, "could not turn the element into a map")
	}
	_, err = c.c.CreateDocument(ctx, asMap)
	if err != nil {
		return errors.Wrap(err, "could not create the document")
	}
	return nil
}

func (c *Collection[T]) Update(ctx context.Context, ID string, updateElement any) (T, T, error) {
	asMap, err := id.AnyToMapNoID(updateElement)
	if err != nil {
		return *new(T), *new(T), errors.Wrap(err, "could not turn the element into a map")
	}
	asMap["id"] = ID
	query := buildUpdateQuery(updateElement, "coll", "doc")

	cursor, err := runQueryInCollectionCursor(ctx, c, query, asMap)
	if err != nil {
		return *new(T), *new(T), errors.Wrap(err, "could not update the document")
	}

	var elems []T
	_, err = cursor.ReadDocument(ctx, &elems)
	if err != nil {
		return *new(T), *new(T), errors.Wrap(err, "could not read the old/new documents")
	}

	return elems[0], elems[1], nil
}

func (c *Collection[T]) Delete(ctx context.Context, ID string) error {
	query :=
		`filter doc.id == @id
remove doc in @@coll
return OLD`
	cur, err := runQueryInCollectionCursor(ctx, c, query, map[string]interface{}{
		"id": ID,
	})
	if err != nil {
		return errors.Wrap(err, "could not delete the element")
	}
	if cur.Count() == 0 {
		return errors.New("element does not exist")
	}
	return nil
}

func (c *Collection[T]) Search(ctx context.Context, search string) ([]T, error) {
	params := map[string]any{}

	//t := reflect.TypeOf(*new(T))

	//filters := make([]string, 0, len(fields))
	//
	//for _, fieldName := range fields {
	//	field, ok := t.FieldByName(fieldName)
	//	if !ok {
	//		continue
	//	}
	//
	//	var filter string
	//	kind := field.Type.Kind()
	//	switch kind {
	//	case reflect.Slice, reflect.Array:
	//		filter = fmt.Sprintf(`"%v" in doc.%v`, search, fieldName)
	//	default:
	//		filter = fmt.Sprintf(`doc.%v == "%v"`, fieldName, search)
	//	}
	//
	//	filters = append(filters, filter)
	//}

	query := "\tfilter "
	//query += strings.Join(filters, " || ")
	query += "\n\treturn doc"
	return runQueryInCollection(ctx, c, query, params)
}

func runQueryInCollectionCursor[T id.IDer](ctx context.Context, coll *Collection[T], query string, params map[string]any) (driver.Cursor, error) {
	if params == nil {
		params = make(map[string]any)
	}
	params["@coll"] = coll.c.Name()
	ctx = driver.WithQueryCount(ctx, true)
	query = "for doc in @@coll\n" + query
	cursor, err := coll.db.Query(ctx, query, params)
	if err != nil {
		return nil, errors.Wrap(err, "could not query the database")
	}
	return cursor, nil
}

func runQueryInCollection[T id.IDer](ctx context.Context, coll *Collection[T], query string, params map[string]any) ([]T, error) {
	cur, err := runQueryInCollectionCursor(ctx, coll, query, params)
	if err != nil {
		return nil, err
	}
	return flattenCursor[T](ctx, cur)
}

func flattenCursor[T id.IDer](ctx context.Context, cursor driver.Cursor) ([]T, error) {
	defer cursor.Close()
	count := int(cursor.Count())
	r := make([]T, count)
	for i := 0; i < count; i++ {
		_, err := cursor.ReadDocument(ctx, &r[i])
		if err != nil {
			return nil, errors.Wrap(err, "could not read the document")
		}
	}

	return r, nil
}

func buildUpdateQuery[T any](element T, collParam, elemParam string) string {
	query := "filter doc.id == @id\n"
	query += fmt.Sprintf("update %v with { ", elemParam)
	t := reflect.TypeOf(element)
	v := reflect.ValueOf(element)
	numField := t.NumField()
	fieldNames := make([]string, 0, numField)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		fieldV := v.Field(i)
		if util.IsNilable(fieldV) && fieldV.IsNil() {
			continue
		}
		fieldNames = append(fieldNames, field.Name)
	}
	for _, name := range fieldNames {
		query += fmt.Sprintf("%v: @%v, ", name, name)
	}
	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" } in @@%v\nreturn [OLD, NEW]", collParam)
	return query
}
