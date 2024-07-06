package slices

import (
	"errors"
	"slices"

	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"

	errs "github.com/jarrodhroberson/ossgo/errors"
)

const NOT_FOUND = -1

func Must(result int, err error) int {
	if err != nil {
		err = errors.Join(errs.MustNeverError.New("the wrapped function is expected to never fail; it failed with error:%s", err.Error()), err)
		log.Error().Stack().Err(err).Msg(err.Error())
		panic(err)
	}
	return result
}

func Partition[T any](s []T, size int) [][]T {
	chunks := make([][]T, 0, len(s)/size+1)
	for i := 0; i < len(s); i += size {
		end := i + size
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}

func Map[F any, T any](in []F, f func(F) T) []T {
	m := make([]T, len(in))
	for i, el := range in {
		m[i] = f(el)
	}
	return m
}

func FindStructIn[T any](toSearch []T, find func(t T) bool) (int, error) {
	idx := slices.IndexFunc(toSearch, find)
	if idx == -1 {
		return idx, errs.NotFoundError.New("could not find struct in slice")
	}
	return idx, nil
}

func Filter[T any](toSearch []T, match func(t T) bool) ([]int, error) {
	var results []int
	for idx, v := range toSearch {
		if match(v) {
			results = append(results, idx)
		}
	}
	if len(results) == 0 {
		return nil, errs.NotFoundError.New("could not match any instance of %T in slice", *new(T))
	}
	return results, nil
}

func FindFirst[T any](toSearch []T, find func(t T) bool) (int, error) {
	if len(toSearch) == 0 {
		return NOT_FOUND, errorx.IllegalArgument.New("toSearch Argument can not be empty")
	}
	idx := slices.IndexFunc(toSearch, find)
	if idx == NOT_FOUND {
		return idx, errs.NotFoundError.New("could not find instance of %T in slice", *new(T))
	}
	return idx, nil
}

func FirstNonNilIn[T any](toSearch ...*T) (int, error) {
	for idx, v := range toSearch {
		if v != nil {
			return idx, nil
		}
	}
	return NOT_FOUND, errs.NotFoundError.New("could not find a non-nil value in the provided slice")
}
