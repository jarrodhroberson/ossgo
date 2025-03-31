package structs

import (
	"fmt"
	"strings"
)

const ReadOnly = "readonly"
const Immutable = "immutable"
const Ignore = "ignore"

type Tag struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
}

func (t Tag) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("{\"%s\": [%s]}", t.Name, strings.Join(t.Values, ","))), nil
}

func (t Tag) String() string {
	return fmt.Sprintf("{\"%s\": [%s]}", t.Name, strings.Join(t.Values, ","))
}
