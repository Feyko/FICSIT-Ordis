package id

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
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
