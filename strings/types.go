package strings

import (
	"unique"
)

/*
NO_DATA is used as a marker value to show that no data was provided
and that the data is not just missing.
It is represented as unicode character NULL (\u0000)
*/
var NO_DATA = unique.Make("\u0000")

type Stringifier[T any] func(t T) string