package strings

/*
NO_DATA is used as a marker value to show that no data was provided
and that the data is not just missing.
*/
var NO_DATA = "\u0000"

type Stringifier[T any] func(t T) string
