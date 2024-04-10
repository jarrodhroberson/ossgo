package functions

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/joomcode/errorx"
)

var struct_not_found = errorx.NewType(errorx.NewNamespace("SERVER"), "STRUCT NOT FOUND", errorx.NotFound())

func MinInt(first int, second int) int {
	return int(math.Min(float64(first), float64(second)))
}

func FindStringInSlice(toSearch []string, target string) (int, error) {
	idx := slices.IndexFunc(toSearch, func(s string) bool {
		return s == target
	})
	if idx == -1 {
		return idx, fmt.Errorf("could not find %s in %s", target, strings.Join(toSearch, ","))
	}
	return idx, nil
}

func FindStructInSlice[T any](toSearch []T, find func(t T) bool) (int, error) {
	idx := slices.IndexFunc(toSearch, find)
	if idx == -1 {
		return idx, struct_not_found.New("could not find struct in slice")
	}
	return idx, nil
}

func UnmarshallMap[T any](m map[string]interface{}, o T) {
	must.UnMarshalJson(must.MarshalJson(m), o)
}
