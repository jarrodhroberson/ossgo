package strings

import (
	"fmt"
	"slices"
	"strings"
)

func FindStringInSlice(toSearch []string, target string) (int, error) {
	idx := slices.IndexFunc(toSearch, func(s string) bool {
		return s == target
	})
	if idx == -1 {
		return idx, fmt.Errorf("could not find %s in %s", target, strings.Join(toSearch, ","))
	}
	return idx, nil
}
