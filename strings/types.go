package strings

/*
NO_DATA is used as a marker value to show that no data was provided
and that the data is not just missing.
It is represented as unicode character NULL (\u0000)
*/
var NO_DATA = "\u0000"

type Stringifier[T any] func(t T) string


// RunePredicate is a function that takes a rune and returns a boolean.
type RunePredicate func(r rune) bool
