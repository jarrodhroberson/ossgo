package functions

import (
	"fmt"
	"slices"
	"strings"

	"github.com/joomcode/errorx"

	"github.com/jarrodhroberson/ossgo/functions/must"
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

var struct_not_found = errorx.NewType(errorx.NewNamespace("SERVER"), "STRUCT NOT FOUND", errorx.NotFound())

func UnmarshallMap[T any](m map[string]interface{}, o T) {
	must.UnMarshalJson(must.MarshalJson(m), o)
}

func InsteadOfNil[T any](a *T, b *T) *T {
	if b == nil {
		panic(fmt.Errorf("second argument to function \"b\" can not be \"nil\""))
	}
	if a == nil {
		return b
	}
	return a
}

func FirstNonNil[T any](structs ...*T) *T {
	if len(structs) < 1 {
		return nil
	}
	for _, s := range structs {
		if s != nil {
			return s
		}
	}
	return nil
}
