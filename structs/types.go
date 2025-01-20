package structs

import (
	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/rs/zerolog"
)

const ReadOnly = "readonly"
const Immutable = "immutable"
const Ignore = "ignore"

type Tag struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func (t Tag) String() string {
	return string(must.MarshalJson(must.MarshallMap(t)))
}

type logObjectMarshaller[T any] struct {
	delegate T
}

func (l logObjectMarshaller[T]) MarshalZerologObject(e *zerolog.Event) {

}
