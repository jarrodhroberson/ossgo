package structs

import (
	"github.com/jarrodhroberson/ossgo/functions/must"
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
