package slices

import (
	"errors"
	"slices"

	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"

	errs "github.com/jarrodhroberson/ossgo/errors"
)

const none = -1

func Must(result int, err error) int {
	if err != nil {
		err = errors.Join(errs.MustNeverError.New("the wrapped function is expected to never fail; it failed with error:%s", err.Error()), err)
		log.Error().Stack().Err(err).Msg(err.Error())
		panic(err)
	}
	return result
}

func FindStructInSlice[T any](toSearch []T, find func(t T) bool) (int, error) {
	idx := slices.IndexFunc(toSearch, find)
	if idx == -1 {
		return idx, errs.StructNotFound.New("could not find struct in slice")
	}
	return idx, nil
}

func FindInSlice[T any](toSearch []T, find func(t T) bool) (int, error) {
	if len(toSearch) == 0 {
		return none, errorx.IllegalArgument.New("toSearch Argument can not be empty")
	}
	idx := slices.IndexFunc(toSearch, find)
	if idx == none {
		return idx, errs.NotFoundError.New("could not find instance of %T in slice", *new(T))
	}
	return idx, nil
}
