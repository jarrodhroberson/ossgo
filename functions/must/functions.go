package must

import (
	"encoding/json"
	"errors"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/jarrodhroberson/ossgo/seq"
	"github.com/joomcode/errorx"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	errs "github.com/jarrodhroberson/ossgo/errors"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

// MustWithErrFunc if err != nil use custom errFunc to handle the error and panic(), else return T
func MustWithErrFunc[T any](result T, err error, errFunc func(err error) error) T {
	if err != nil {
		panic(errorx.Panic(errorx.EnsureStackTrace(errs.MustNeverError.WrapWithNoMessage(errFunc(err)))))
	} else {
		return result
	}
}

func Must[T any](result T, err error) T {
	if err != nil {
		err = errs.MustNeverError.Wrap(err, "the wrapped function is expected to never fail, it failed %s", err.Error())
		panic(errorx.Panic(err))
	}
	return result
}

func ParseTime(format string, s string) (time.Time, error) {
	t, err := time.Parse(format, s)
	if err != nil {
		err = errs.MustNeverError.WrapWithNoMessage(errs.ParseError.Wrap(err, "could not parse %s as time %s", s, format))
		return time.Time{}, err
	}
	return t, nil
}

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		err = errs.MustNeverError.WrapWithNoMessage(errs.ParseError.Wrap(err, "could not parse %s as int", s))
		log.Error().Err(err).Str("arg", s).Msg(err.Error())

		panic(errorx.Panic(err))
	}
	return i
}

func UnMarshalJson(bytes []byte, o any) {
	err := json.Unmarshal(bytes, o)
	if err != nil {
		err = errors.Join(err, errs.UnMarshalError.New("could not unmarshal %s", string(bytes)))
		log.Error().Err(err).Msg(err.Error())

		panic(errorx.Panic(err))
	}
	return
}

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

func MarshallMap[T any](o T) map[string]interface{} {
	m := make(map[string]interface{})
	UnMarshalJson(MarshalJson(o), &m)
	return m
}

func UnmarshallMap[T any](m map[string]interface{}, o T) {
	UnMarshalJson(MarshalJson(m), o)
}

func MinInt(ints ...int) int {
	min := math.MaxInt
	return seq.First[int](slices.All(ints), func(i int) bool {
		if i < min {
			min == i
		}
	})
}
