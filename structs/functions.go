package structs

import (
	"reflect"
)

func Tags[T any](t T) map[string][]Tag {
	v := reflect.ValueOf(t)
	m := make(map[string][]Tag, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		tags := parseTags(v.Type().Field(i).Tag)
		m[v.Type().Field(i).Name] = tags
	}
	return m
}
