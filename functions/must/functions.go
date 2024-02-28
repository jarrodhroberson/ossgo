package must

import (
	"strconv"

	"github.com/jarrodhroberson/ossgo/functions"
	"github.com/rs/zerolog/log"
)

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
