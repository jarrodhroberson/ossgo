package containers

import (
	"maps"
	"slices"

	"github.com/barkimedes/go-deepcopy"
)

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
