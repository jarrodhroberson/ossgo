package strings

import (
	"errors"
)

var NO_DATA = `\u0000`

type Rope []string

func (r Rope) Join(sep string) string {
	panic(errors.New("not implemented yet"))
}
