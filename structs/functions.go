package structs

import (
	"reflect"
	"strings"

	"github.com/jarrodhroberson/destruct/destruct"
	"github.com/rs/zerolog"
)

func parseTags(tag reflect.StructTag) []Tag {
	tagsWithValues := strings.Split(string(tag), " ")
	tags := make([]Tag, 0, len(tagsWithValues))
	for _, tag := range tagsWithValues {
		kv := strings.Split(tag, ":")
		name := kv[0]
		values := strings.Split(kv[1], ",")
		tags = append(tags, Tag{
			Name:   name,
			Values: values,
		})
	}
	return tags
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

func DecorateWithLogObjectMarshaller[T any](s any) zerolog.LogObjectMarshaler {

}
