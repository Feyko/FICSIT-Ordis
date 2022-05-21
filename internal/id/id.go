package id

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"reflect"
)

type IDer interface {
	ID() string
}

func ToMap(e IDer, IDField string) (map[string]any, error) {
	var mapform map[string]any
	err := mapstructure.Decode(e, &mapform)
	if err != nil {
		return nil, fmt.Errorf("could not decode the input into a map: %w", err)
	}
	mapform[IDField] = e.ID()
	return mapform, nil
}

func Update(elem map[string]any, update IDer) error {
	v := reflect.ValueOf(update)
	t := reflect.TypeOf(update)
	if v.Kind() != reflect.Struct {
		return errors.New("update object must be a struct")
	}
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldInfo := t.Field(i)
		isNil, err := safeIsNil(field)
		if err != nil && isNil {
			continue
		}
		if field.Kind() == reflect.Pointer {
			field = field.Elem()
		}
		elem[fieldInfo.Name] = field.Interface()
	}
	return nil
}

func safeIsNil(field reflect.Value) (bool, error) {
	switch field.Kind() {
	case reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan, reflect.Interface:
		return field.IsNil(), nil
	default:
		return false, errors.New("field not nullable")
	}
}

type Searchable interface {
	IDer
	SearchFields() []string
}
