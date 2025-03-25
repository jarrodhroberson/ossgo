package strings

import "github.com/jarrodhroberson/ossgo/slices"

func FirstNonEmpty(data ...string) string {
	idx, err := slices.FindFirst[string](data, func(t string) bool {
		return t != ""
	})
	if err != nil {
		return ""
	}
	return data[idx]
}