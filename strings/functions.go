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
	return !IsControl(IsOtherSurrogate(IsOtherPrivateUse(IsSepAny(func(r rune) bool { return false }))))(r)
}

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
func IsStringUnicodePrintable(s string) bool {
	for r := range seq.RuneSeq(s) {
		if !IsUnicodePrintable(r) {
			return false
		}
	}
	return true
}

// IsPrintable creates a RunePredicate that checks if a rune is not printable by combining multiple Unicode checks.
// It chains together several predicates to check for control characters, marks, line separators,
// paragraph separators, private use characters, and surrogate code points.
//
// Parameters:
//   - next: The next RunePredicate in the chain to be evaluated if none of the non-printable checks match
//
// Returns:
//   - RunePredicate: A function that returns true if the rune is not printable or if the next predicate returns true
func IsPrintable(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return IsLetter(IsNumber(IsPunct(IsSymbol(IsSpace(next)))))(r)
	}
}

func IsEmoji(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.So, r) || next(r)
	}
}

// IsOtherPrivateUse creates a RunePredicate that checks if a rune is a private use character (Co),
// or matches the next predicate.
func IsOtherPrivateUse(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.Co, r) || next(r)
	}
}

// IsOtherSurrogate creates a RunePredicate that checks if a rune is a surrogate character (Cs),
// or matches the next predicate.
func IsOtherSurrogate(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.Cs, r) || next(r)
	}
}

// IsSepAny creates a RunePredicate that checks if a rune is any type of separator (spaces, lines, paragraphs),
// by chaining together IsSepSpace, IsSepLine, and IsSepParagraph predicates.
//
// Parameters:
//   - next: The next RunePredicate in the chain to be evaluated if none of the separators checks match
//
// Returns:
//   - RunePredicate: A function that returns true if the rune is any type of separator or if the next predicate returns true
func IsSepAny(next RunePredicate) RunePredicate {
	return IsSepSpace(IsSepLine(IsSepParagraph(next)))
}

// IsSepSpace creates a RunePredicate that checks if a rune is a space separator (Zs),
// or matches the next predicate.
//
// Parameters:
//   - next: The next RunePredicate in the chain to be evaluated if the space separator check doesn't match
//
// Returns:
//   - RunePredicate: A function that returns true if the rune is a space separator or if the next predicate returns true
func IsSepSpace(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.Zs, r) || next(r)
	}
}

// IsSepLine creates a RunePredicate that checks if a rune is a line separator (Zl),
// or matches the next predicate.
func IsSepLine(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.Zl, r) || next(r)
	}
}

// IsSepParagraph creates a RunePredicate that checks if a rune is a paragraph separator (Zp),
// or matches the next predicate.
func IsSepParagraph(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.Zp, r) || next(r)
	}
}

// IsControl creates an IsRune that checks if a rune is a control character (Cc).
func IsControl(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.Cc, r) || next(r)
	}
}

// IsMark creates an IsRune that checks if a rune is a mark character (Mc, Me, Mn).
func IsMark(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.M, r) || next(r)
	}
}

// IsLetter creates an IsRune that checks if a rune is a letter character (Lu, Ll, Lt, Lm, Lo).
func IsLetter(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.L, r) || next(r)
	}
}

// IsNumber creates an IsRune that checks if a rune is a number character (Nd, Nl, No).
func IsNumber(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.N, r) || next(r)
	}
}

// IsPunct creates an IsRune that checks if a rune is a punctuation character (Pc, Pd, Ps, Pe, Pi, Pf, Po).
func IsPunct(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.P, r) || next(r)
	}
}

// IsSymbol creates an IsRune that checks if a rune is a symbol character (Sm, Sc, Sk, So).
func IsSymbol(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.Is(unicode.S, r) || next(r)
	}
}

// IsSpace creates an IsRune that checks if a rune is a space character (Zs, Zl, Zp).
func IsSpace(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.IsSpace(r) || next(r)
	}
}

// IsGraphic creates an IsRune that checks if a rune is a graphic character (includes letters, marks, numbers, punctuation, symbols, and spaces).
func IsGraphic(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.IsGraphic(r) || next(r)
	}
}

// IsPrint creates an IsRune that checks if a rune is a printable character (includes letters, marks, numbers, punctuation, symbols, spaces, and some control characters).
func IsPrint(next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.IsPrint(r) || next(r)
	}
}

// IsOneOf creates an IsRune that checks if a rune belongs to any of the given unicode ranges.
func IsOneOf(tables []*unicode.RangeTable, next RunePredicate) RunePredicate {
	return func(r rune) bool {
		return unicode.IsOneOf(tables, r) || next(r)
	}
}
