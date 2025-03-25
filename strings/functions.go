package strings

import "github.com/jarrodhroberson/ossgo/slices"

func FirstNonEmpty(data ...string) string {
	idx, err := slices.FindFirst[string](data, func(t string) bool {
		return t != "" && t != NO_DATA
	})
	if err != nil {
		return NO_DATA
	}
	return data[idx]
}
