package must

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/jarrodhroberson/ossgo/functions"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func Must[T any](result T, err error) T {
	if err != nil {
		log.Error().Stack().Err(err).Msg("this must not fail")
	}
	return result
}

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal().Err(err).Str("arg", s).Msg(err.Error())
		return -1
	}
	return i
}

func FindStructInSlice[T any](toSearch []T, find func(t T) bool) int {
	idx, err := functions.FindStructInSlice[T](toSearch, find)
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
		return -1
	}
	return idx
}

func FindStringInSlice(toSearch []string, target string) int {
	idx, err := functions.FindStringInSlice(toSearch, target)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("could not find %s in %s", target, strings.Join(toSearch, ",")))
		log.Fatal().Err(err).Msg(err.Error())
		return idx
	}
	return idx
}

func UnMarshalJson(bytes []byte, o any) {
	err := json.Unmarshal(bytes, o)
	if err != nil {
		log.Fatal().Err(err).Msgf("could not unmarshal %s", string(bytes))
		return
	}
	return
}

func MarshalJson(o any) []byte {
	bytes, err := json.Marshal(o)
	if err != nil {
		err = errors.Wrap(err, fmt.Sprintf("could not marshal %v", o))
		log.Fatal().Stack().Err(err).Msg(err.Error())
		return []byte{}
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
