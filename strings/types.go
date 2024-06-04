package strings

import (
	"errors"
)

var NUL = `\u0000`

type Rope []string

func (r Rope) Join(sep string) string {
	panic(errors.New("not implemented yet"))
}