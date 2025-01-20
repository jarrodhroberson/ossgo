package containers

import (
	"fmt"
	"maps"
	"reflect"
	"slices"

	"github.com/barkimedes/go-deepcopy"
	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type lom[T any] struct {
	delegate T
}

func (l lom[T]) MarshalZerologObject(e *zerolog.Event) {
	WalkMap(must.MarshallMap(l.delegate), "", func(k string, v interface{}) {
		switch v.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			e.Int64(k, v.(int64))
		case float32, float64:
			e.Float64(k, v.(float64))
		case string:
			e.Str(k, v.(string))
		case bool:
			e.Bool(k, v.(bool))
		default:
			log.Warn().Msgf("unknown type %s:", reflect.TypeOf(v))
			e.Str(k, fmt.Sprintf("%s", v))
		}
	})
}

func DecorateWithLogObjectMarshaller[T any](s T) zerolog.LogObjectMarshaler {
	return lom[T]{
		delegate: s,
	}
}

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
