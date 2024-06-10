package strings

import (
	"slices"
	"strings"

	"github.com/jarrodhroberson/ossgo/errors"
)

func FindInSlice(toSearch []string, target string) (int, error) {
	idx := slices.IndexFunc(toSearch, func(s string) bool {
		return s == target
	})
	if idx == -1 {
		return idx, errors.NotFoundError.New("could not find %s in %s", target, strings.Join(toSearch, ","))
	}
	return idx, nil
}
