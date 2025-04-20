// Package must provides utility functions for handling errors and type conversions with panic behavior
package must

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/joomcode/errorx"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	errs "github.com/jarrodhroberson/ossgo/errors"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

// MustWithErrFunc executes a custom error handling function and panics if an error occurs.
// It takes a result of type T, an error, and a custom error handling function.
// Returns the result if no error occurs, otherwise panics with the processed error.
func MustWithErrFunc[T any](result T, err error, errFunc func(err error) error) T {
	if err != nil {
		panic(errorx.Panic(errorx.EnsureStackTrace(errs.MustNeverError.WrapWithNoMessage(errFunc(err)))))
	} else {
		return result
	}
}

// Must is a generic function that panics if an error occurs, otherwise returns the result.
// It takes a result of type T and an error as parameters.
// Returns the result if no error occurs, otherwise panics with the wrapped error.
func Must[T any](result T, err error) T {
	if err != nil {
		err = errs.MustNeverError.Wrap(err, "the wrapped function is expected to never fail, it failed %s", err.Error())
		panic(errorx.Panic(err))
	}
	return result
}

// ParseTime attempts to parse a time string using the specified format.
// It returns the parsed time.Time and an error if parsing fails.
// The error is wrapped with additional context information.
func ParseTime(format string, s string) (time.Time, error) {
	t, err := time.Parse(format, s)
	if err != nil {
		err = errs.MustNeverError.WrapWithNoMessage(errs.ParseError.Wrap(err, "could not parse %s as time %s", s, format))
		return time.Time{}, err
	}
	return t, nil
}

// ParseInt converts a string to an integer.
// It panics if the conversion fails, logging the error before panicking.
// Returns the parsed integer value.
func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		err = errs.MustNeverError.WrapWithNoMessage(errs.ParseError.Wrap(err, "could not parse %s as int", s))
		log.Error().Err(err).Str("arg", s).Msg(err.Error())

		panic(errorx.Panic(err))
	}
	return i
}

// ParseIntOr attempts to parse a string into an integer, returning a default value if parsing fails.
//
// Parameters:
//   - s: The string to parse.
//   - or: The default integer value to return if parsing fails.
//
// Returns:
//   - The parsed integer value if successful, otherwise the default value `or`.
//
// Example:
//
//	parsedValue := ParseIntOr("42", 0) // returns 42
//	parsedValue := ParseIntOr("invalid", 10) // returns 10
func ParseIntOr(s string, or int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return or
	}
	return i
}

// ParseInt64 converts a string to an int64.
// It panics if the conversion fails, logging the error before panicking.
// Returns the parsed int64 value.
func ParseInt64(s string) int64 {
	i, err := strconv.Atoi(s)
	if err != nil {
		err = errs.MustNeverError.WrapWithNoMessage(errs.ParseError.Wrap(err, "could not parse %s as int", s))
		log.Error().Err(err).Str("arg", s).Msg(err.Error())

		panic(errorx.Panic(err))
	}
	return int64(i)
}

// UnMarshalJson unmarshals JSON bytes into the provided object.
// It panics if unmarshaling fails, logging the error before panicking.
// The object parameter must be a pointer to the target type.
func UnMarshalJson(bytes []byte, o any) {
	err := json.Unmarshal(bytes, o)
	if err != nil {
		err = errors.Join(err, errs.UnMarshalError.New("could not unmarshal %s", string(bytes)))
		log.Error().Err(err).Msg(err.Error())

		panic(errorx.Panic(err))
	}
	return
}

// MarshalJson marshals an object into JSON bytes.
// It panics if the object is nil or if marshaling fails, logging the error before panicking.
// Returns the marshaled JSON bytes.
func MarshalJson(o any) []byte {
	if o == nil {
		err := errors.Join(errs.MustNotBeNil.New("cannot marshal nil"), errs.MarshalError.New("could not marshal %v", o))
		log.Error().Stack().Err(err).Msg(err.Error())

		panic(errorx.Panic(err))
	}

	bytes, err := json.Marshal(o)
	if err != nil {
		err = errors.Join(err, errs.MarshalError.New("could not marshal %v", o))
		log.Error().Stack().Err(err).Msg(err.Error())

		panic(errorx.Panic(err))
	}
	return bytes
}

// MarshallMap converts an object of type T to a map[string]interface{} through JSON marshaling and unmarshaling.
// It uses MarshalJson and UnMarshalJson internally for the conversion.
// Returns the resulting map.
func MarshallMap[T any](o T) map[string]interface{} {
	m := make(map[string]interface{})
	UnMarshalJson(MarshalJson(o), &m)
	return m
}

// UnmarshallMap converts a map[string]interface{} to an object of type T through JSON marshaling and unmarshaling.
// It uses MarshalJson and UnMarshalJson internally for the conversion.
// The object parameter must be a pointer to the target type.
func UnmarshallMap[T any](m map[string]interface{}, o T) {
	UnMarshalJson(MarshalJson(m), o)
}

// AsString converts a byte slice to a string.
// It panics if the input byte slice is nil, logging the error before panicking.
// Returns the resulting string.
func AsString(b []byte) string {
	if b == nil {
		err := errs.MustNotBeNil.New("cannot convert nil bytes to string")
		log.Error().Stack().Err(err).Msg(err.Error())
		panic(errorx.Panic(err))
	}
	return string(b)
}
