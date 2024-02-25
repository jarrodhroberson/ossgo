package structs

import (
	"fmt"
	"reflect"
	"strings"
)

const ReadOnly = "readonly"
const Immutable = "immutable"

type Tag struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func (t Tag) String() string {
	return fmt.Sprintf("%s:\"%s\"", t.Name, strings.Join(t.Values, ","))
}

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
