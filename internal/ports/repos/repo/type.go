package repo

import (
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"reflect"
	"strings"
)

func GetTypeInfo(v any) (TypeInfo, error) {
	typ := reflect.TypeOf(v)
	if typ.Kind() != reflect.Struct {
		return TypeInfo{}, errors.New("v needs to be a struct")
	}

	numFields := typ.NumField()

	fields := make([]FieldInfo, numFields)

	// TODO: Expand to sub-structs
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
	info := FieldInfo{StructField: field}

	if slices.Contains(stringValues, "search") {
		info.ToSearch = true
	}

	return info, nil
}
