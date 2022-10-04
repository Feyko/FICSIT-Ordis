package util

import (
	"github.com/fatih/structs"
)

func ToMapNoNil(v any) map[string]any {
	asMap := structs.Map(v)
	for k, v := range asMap {
		if IsNil(v) {
			delete(asMap, k)
		}
	}
	return asMap
}
