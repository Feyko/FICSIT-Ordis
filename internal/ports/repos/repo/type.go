package repo

import (
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"reflect"
	"strings"
)

func GetTypeInfo(v any) (TypeInfo, error) {
	typ := reflect.TypeOf(v)

	switch typ.Kind() {
	case reflect.Interface, reflect.Pointer:
		return GetTypeInfo(reflect.ValueOf(v).Elem().Interface())
	}

	if typ.Kind() != reflect.Struct {
		return TypeInfo{}, errors.New("v needs to be or contain a struct")
	}
	value := reflect.ValueOf(v)

	numFields := typ.NumField()

	fields := make([]FieldInfo, 0, numFields)
	embeds := make(map[int]TypeInfo)

	for i := 0; i < numFields; i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}

		if field.Type.Kind() == reflect.Struct {
			embed, err := GetTypeInfo(value.Field(i).Interface())
			if err != nil {
				return TypeInfo{}, errors.Wrap(err, "error getting type info for embed struct")
			}
			embeds[field.Index[0]] = embed
			continue
		}

		info, err := fillFieldInfo(field)
		if err != nil {
			return TypeInfo{}, errors.Wrap(err, "error getting the tag values")
		}
		fields = append(fields, info)
	}

	return TypeInfo{
		Fields: fields,
		Embeds: embeds,
	}, nil
}

type TypeInfo struct {
	Fields []FieldInfo
	Embeds map[int]TypeInfo
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
