package containers

import (
	"fmt"
	"iter"
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

func Seq2[K comparable, V any](m map[K]V) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k := range maps.Keys(m) {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

type ImmutableMap[K comparable, V any] interface {
	Get(key K) (V,bool)
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
	IterSeq2() iter.Seq2[K,V]
}

type immutableMap[K comparable, V any] struct {
	m map[K]V
}

func (i immutableMap[K, V]) Get(key K) (V,bool) {
	v, ok := i.m[key]
	return v, ok
}

func (i immutableMap[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(i.m)
}

func (i immutableMap[K, V]) Values() iter.Seq[V] {
	return maps.Values(i.m)
}

func (i immutableMap[K, V]) IterSeq2() iter.Seq2[K,V] {
	return Seq2(i.m)
}

func NewImmutableMap[K comparable, V any](m map[K]V) ImmutableMap[K, V] {
	return immutableMap[K, V]{m: DeepClone(m)}
}