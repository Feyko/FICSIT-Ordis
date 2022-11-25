package id

import (
	"FICSIT-Ordis/internal/util"
	"github.com/pkg/errors"
)

type IDer interface {
	ID() string
}

func ToMap(v IDer) (map[string]any, error) {
	asMap, err := AnyToMapNoID(v)
	if err != nil {
		return nil, err
	}
	asMap["id"] = v.ID()
	return asMap, nil
}

func AnyToMapNoID(v any) (map[string]any, error) {
	asMap := util.ToMapNoNil(v)
	if _, ok := asMap["id"]; ok {
		return nil, errors.Errorf("value of type %T already has field id", v)
	}
	return asMap, nil
}
