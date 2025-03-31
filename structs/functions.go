package structs

import (
	"reflect"
	"strings"

	"github.com/jarrodhroberson/destruct/destruct"
)

func parseTags(t reflect.StructTag) []Tag {
	tagsWithValues := strings.Split(string(t), " ")
	tags := make([]Tag, 0, len(tagsWithValues))
	for _, tag := range tagsWithValues {
		kv := strings.Split(tag, ":")
		name := kv[0]
		values := strings.Split(kv[1], ",")
		tags = append(tags, NewTag(name, values...))
	}
	return tags
}

func NewTag(name string, values ...string) Tag {
	return Tag{
		Name:   name,
		Values: values,
	}
}

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
