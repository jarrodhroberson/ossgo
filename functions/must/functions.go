package must

import (
	"encoding/json"
	"errors"
	"github.com/joomcode/errorx"
	"math"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	errs "github.com/jarrodhroberson/ossgo/errors"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

// MustErrFunc if err != nil use custom errFunc to handle the error and panic(), else return T
func MustErrFunc[T any](result T, err error, errFunc func(err error) error) T {
	if err != nil {
		panic(errors.Join(errFunc(err), err))
	} else {
		return result
	}
}

func Must[T any](result T, err error) T {
	if err != nil {
		err = errs.MustNeverError.New("the wrapped function is expected to never fail; it failed with error:%s", err.Error())
		log.Error().Stack().Err(err).Msg(err.Error())
		panic(errorx.Panic(err))
	}
	return result
}

func ParseTime(format string, s string) time.Time {
	t, err := time.Parse(format, s)
	if err != nil {
		err = errors.Join(err, errs.ParseError.New("Could not parse %s as time %s", s, format))
		log.Error().Err(err).Msgf(err.Error())

		panic(errorx.Panic(err))
	}
	return t
}

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		err = errors.Join(err, errs.ParseError.New("Could not parse %s as int", s))
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

func MinInt(first int, second int) int {
	return int(math.Min(float64(first), float64(second)))
}
