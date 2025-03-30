package strings

import (
	"fmt"
	"iter"
	stdslices "slices"
	stdstrings "strings"

	"github.com/jarrodhroberson/ossgo/seq"
	"github.com/jarrodhroberson/ossgo/slices"
)

func FirstNonEmpty(data ...string) string {
	idx, err := slices.FindFirst[string](data, func(t string) bool {
		return t != "" && t != NO_DATA
	})
	if err != nil {
		return NO_DATA
	}
	return data[idx]
}

func MapString[S fmt.Stringer](it iter.Seq[S]) iter.Seq[string] {
	return seq.Map[S, string](it, StringFunc[S])
}

func StringFunc[S fmt.Stringer](s S) string {
	return s.String()
}

func Join[S fmt.Stringer](it iter.Seq[S], sep string) string {
	ms := MapString[S](it)
	sls := stdslices.Collect[string](ms)
	return stdstrings.Join(sls, sep)
}