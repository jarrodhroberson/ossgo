package must

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/jarrodhroberson/ossgo/functions"
	"github.com/kofalt/go-memoize"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func Must[T any](result T, err error) T {
	if err != nil {
		err = mustNeverError.New("the wrapped function is expected to never fail; it failed with error:%s", err.Error())
		log.Error().Stack().Err(err).Msg(err.Error())
		panic(err)
	}
	return result
}

func ParseTime(format string, s string) time.Time {
	t, err := time.Parse(format, s)
	if err != nil {
		err = errors.Join(err, parseError.New("Could not parse %s as time %s", s, format))
		log.Error().Err(err).Msgf(err.Error())
		panic(err)
	}
	return t
}

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		err = errors.Join(err, parseError.New("Could not parse %s as int", s))
		log.Error().Err(err).Str("arg", s).Msg(err.Error())
		panic(err)
	}
	return i
}

func FindInSlice[T any](toSearch []T, find func(t T) bool) int {
	idx, err := functions.FindInSlice[T](toSearch, find)
	if err != nil {
		err = notFoundError.New("Could not find instance of %T in slice", *new(T))
		log.Error().Err(err).Msg(err.Error())
		panic(err)
	}
	return idx
}

func FindStringInSlice(toSearch []string, target string) int {
	idx, err := functions.FindStringInSlice(toSearch, target)
	if err != nil {
		err = errors.Join(err, notFoundError.New("could not find %s in %s", target, strings.Join(toSearch, ",")))
		log.Error().Err(err).Msg(err.Error())
		panic(err)
	}
	return idx
}

func UnMarshalJson(bytes []byte, o any) {
	err := json.Unmarshal(bytes, o)
	if err != nil {
		err = errors.Join(err, unMarshalError.New("could not unmarshal %s", string(bytes)))
		log.Error().Err(err).Msg(err.Error())
		panic(err)
	}
	return
}

func MarshalJson(o any) []byte {
	bytes, err := json.Marshal(o)
	if err != nil {
		err = errors.Join(err, marshalError.New("could not marshal %v", o))
		log.Error().Stack().Err(err).Msg(err.Error())
		panic(err)
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

func Identity[T any](t T) string {
	switch reflect.ValueOf(t).Kind() {
	case reflect.Pointer:
		return fmt.Sprintf("%T_%p", t, t)
	case reflect.Struct:
		return fmt.Sprintf("%T_%p", t, &t)
	default:

		panic("unhandled default case")
	}
}

func Call[T any](m *memoize.Memoizer, key string, f memoize.MemoizedFunction[T]) T {
	result, err, cached := memoize.Call(m, key, f)
	log.Debug().Msgf("%s:%s:%t", result, err, cached)
	if err != nil {
		panic(err)
	}
	return result
}
