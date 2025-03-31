package containers

import (
	"fmt"
	"maps"
	"slices"

	"github.com/barkimedes/go-deepcopy"
)

// WalkMap recursively walks through a map[string]interface{}, applying a function to each key-value pair.
//
// It handles nested maps by recursively calling itself.
//
// Parameters:
func WalkMap(m map[string]interface{}, prefix string, f func(k string, v interface{})) {
	for k, v := range m {
		switch v.(type) {
		case map[string]interface{}: // Check if value is another map
			WalkMap(v.(map[string]interface{}), fmt.Sprintf("%s%s.", prefix, k), f) // Recursively call for nested map
		default:
			f(k, v)
		}
	}
}

// RemoveKeys removes all the keys from map m
func RemoveKeys[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	dc := DeepClone(m)
	maps.DeleteFunc(dc, func(key K, i V) bool {
		return slices.Contains(keys, key)
	})
	return dc
}

// KeepKeys remove all the keys that are NOT in keys from map m
func KeepKeys[K comparable, V any](m map[K]V, keys ...K) map[K]V {
	dc := DeepClone(m)
	maps.DeleteFunc(dc, func(key K, v V) bool {
		return !slices.Contains(keys, key)
	})
	return dc
}

func DeepClone[T any](m T) T {
	return deepcopy.MustAnything(m).(T)
}
