package strings

import (
	"fmt"
	"iter"
	stdslices "slices"
	stdstrings "strings"
	"unicode"

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

func IsUnicodePrintable(r rune) bool {
	if unicode.IsControl(r) {
		return false // Exclude control characters
	}
	if unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Me, r) || unicode.Is(unicode.Mc, r) {
		return false // Exclude nonspacing marks, enclosing marks, and combining marks.
	}
	if unicode.Is(unicode.Zs, r) {
		return true // Keep spaces
	}
	if unicode.Is(unicode.Zl, r) || unicode.Is(unicode.Zp, r) {
		return false // Exclude line and paragraph separators.
	}
	if unicode.Is(unicode.Cs,r) || unicode.Is(unicode.Co, r) {
		return false // Exclude surrogates and private use
	}

	// Include everything else as printable
	return true
}

//
// IsStringUnicodePrintable checks if all runes in a given string `s` are considered printable according to the Unicode standard.
//
// The function iterates through each rune in the string using a sequence generator (`seq.RuneSeq`)
// and determines its printability using `IsUnicodePrintable`.
//
// A string is considered printable if all of its constituent runes satisfy the criteria defined in `IsUnicodePrintable`.
//
// Parameters:
//   - s: The string to check for printable runes.
//
// Returns:
//   - bool: Returns true if all runes in the string are printable; otherwise, returns false.
//
func IsStringUnicodePrintable(s string) bool {
	for r := range seq.RuneSeq(s) {
		if !IsUnicodePrintable(r) {
			return false
		}
	}
	return true
}