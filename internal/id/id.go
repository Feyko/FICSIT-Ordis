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

func Update(elem any, update IDer) error {
	v := reflect.ValueOf(update)
	t := reflect.TypeOf(update)
	elemV := reflect.ValueOf(elem).Elem().Elem()
	if v.Kind() != reflect.Struct /* || elemV.Kind() != reflect.Struct */ {
		return errors.New("inputs must be structs")
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
		elemField := elemV.FieldByName(fieldInfo.Name)
		elemField.SetString(field.String())
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
