package repos

import (
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos/arango"
	"FICSIT-Ordis/internal/ports/repos/memrepo"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"reflect"
	"strings"
)

func GetCollection[T id.IDer, U id.IDer](repository repo.Repository[U], name string) (repo.Collection[T], error) {
	repoA, err := retype[T](repository)
	if err != nil {
		return nil, errors.Wrap(err, "could not retype repo")
	}
	collection, err := repoA.GetCollection(name)
	if err != nil {
		return nil, errors.Wrap(err, "could not get collection")
	}
	typed, ok := collection.(repo.Collection[T])
	if !ok {
		return nil, errors.New(fmt.Sprintf("collection '%v' does not hold type '%T'", name, *new(T)))
	}
	return typed, nil
}

func CreateCollection[T id.IDer, U id.IDer](repository repo.Repository[U], name string) (repo.Collection[T], error) {
	repoA, err := retype[T](repository)
	if err != nil {
		return nil, errors.Wrap(err, "could not retype repo")
	}
	collection, err := repoA.CreateCollection(name)
	if err != nil {
		return nil, errors.Wrap(err, "could not create collection")
	}
	typed, ok := collection.(repo.Collection[T])
	if !ok {
		return nil, errors.New(fmt.Sprintf("created collection '%v' does not hold type '%T'", name, *new(T)))
	}
	return typed, nil
}

func GetOrCreateCollection[T id.IDer, U id.IDer](repository repo.Repository[U], name string) (repo.Collection[T], error) {
	collection, err := GetCollection[T](repository, name)
	notFound := errors.Is(err, repo.ErrCollectionNotFound)
	if notFound {
		collection, err = CreateCollection[T](repository, name)
		if err != nil {
			return nil, errors.Wrap(err, "create")
		}
	}
	if err != nil && !notFound {
		return nil, errors.Wrap(err, "get")
	}
	return collection, nil
}

func retype[newT id.IDer, oldT id.IDer](repo repo.Repository[oldT]) (repo.Repository[newT], error) {
	switch typed := repo.(type) {
	case *memrepo.Repository[oldT]:
		return (*memrepo.Repository[newT])(typed), nil
	case *arango.Repository[oldT]:
		return (*arango.Repository[newT])(typed), nil
	default:
		return nil, errors.New("unsupported repository type")
	}
}

func GetTypeInfo(v any) (TypeInfo, error) {
	typ := reflect.TypeOf(v)

	numFields := typ.NumField()

	fields := make([]FieldInfo, numFields)

	for i := 0; i < numFields; i++ {
		field := typ.Field(i)

		info, err := fillFieldInfo(field)
		if err != nil {
			return TypeInfo{}, errors.Wrap(err, "error getting the tag values")
		}
		fields[i] = info
	}

	return TypeInfo{
		Fields: fields,
	}, nil
}

type TypeInfo struct {
	Fields []FieldInfo
}

type FieldInfo struct {
	reflect.StructField
	ToSearch bool
}

func fillFieldInfo(field reflect.StructField) (FieldInfo, error) {
	tag := field.Tag.Get("repos")
	stringValues := strings.Split(tag, ",")
	if len(stringValues) == 0 {
		return FieldInfo{}, nil
	}
	var info FieldInfo

	if slices.Contains(stringValues, "search") {
		info.ToSearch = true
	}

	return info, nil
}