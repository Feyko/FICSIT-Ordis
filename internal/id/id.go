package id

import (
	"FICSIT-Ordis/internal/util"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
)

type IDer interface {
	ID() string
}

type Searchable interface {
	IDer
	SearchFields() []string
}

func ToMap(v IDer) map[string]any {
	m, _ := toMap(v, false)
	return m
}

func ToMapNoOverwrite(v IDer) (map[string]any, error) {
	return toMap(v, true)
}

func toMap(v IDer, checkForIDField bool) (map[string]any, error) {
	asMap := structs.Map(v)
	if checkForIDField {
		_, ok := asMap["id"]
		if ok {
			return nil, errors.Errorf("value of type %T already has field id", v)
		}
	}
	for k, v := range asMap {
		if util.IsNil(v) {
			delete(asMap, k)
		}
	}
	asMap["id"] = v.ID()
	return asMap, nil
}
