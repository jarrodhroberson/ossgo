package functions

import (
	"github.com/joomcode/errorx"
)

func InsteadOfNil[T any](a *T, b *T) *T {
	if b == nil {
		panic(errorx.IllegalArgument.New("second argument to function \"b\" can not be \"nil\""))
	}
	if a == nil {
		return b
	}
	return a
}
