package structs

import (
	"reflect"

	"github.com/jarrodhroberson/destruct/destruct"
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

func Hash[T any](t T) string {
	return destruct.MustHashIdentity(t)
}
