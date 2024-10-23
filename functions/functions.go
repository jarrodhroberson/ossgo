package functions

import (
	"github.com/jarrodhroberson/ossgo/slices"
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

func FirstNonEmpty(data ...string) string {
	idx, err := slices.FindFirst[string](data, func(t string) bool {
		return t != ""
	})
	if err != nil {
		return ""
	}
	return data[idx]
}
