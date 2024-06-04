package containers

import (
	"maps"
	"slices"
)

func RemoveKeys(m map[string]interface{}, keys ...string) {
	maps.DeleteFunc(m, func(s string, i interface{}) bool {
		return slices.Contains(keys, s)
	})
}

func KeepKeys(m map[string]interface{}, keys ...string) {
	maps.DeleteFunc(m, func(s string, i interface{}) bool {
		return !slices.Contains(keys, s)
	})
}
