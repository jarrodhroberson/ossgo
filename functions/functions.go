package functions

import (
	"slices"
	"strings"

	"github.com/joomcode/errorx"
)

var none = -1
var functionNamespace = errorx.NewNamespace("Functions")
var notFoundError = errorx.NewType(functionNamespace, "Not Found", errorx.NotFound())

func FindStringInSlice(toSearch []string, target string) (int, error) {
	if len(toSearch) == 0 {
		return none, errorx.IllegalArgument.New("toSearch Argument can not be empty")
	}
	idx := slices.IndexFunc(toSearch, func(s string) bool {
		return s == target
	})
	if idx == none {
		return idx, notFoundError.New("could not find %s in %s", target, strings.Join(toSearch, ","))
	}
	return idx, nil
}

func FindInSlice[T any](toSearch []T, find func(t T) bool) (int, error) {
	if len(toSearch) == 0 {
		return none, errorx.IllegalArgument.New("toSearch Argument can not be empty")
	}
	idx := slices.IndexFunc(toSearch, find)
	if idx == none {
		return idx, notFoundError.New("could not find instance of %T in slice", *new(T))
	}
	return idx, nil
}
