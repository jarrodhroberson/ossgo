package containers

import (
	"context"
	"errors"
)

var Done = errors.New("no more items in iterator")

type Iterator[T []T] struct {
	iterable []T
	index    int
}

func (it Iterator[T]) Next(ctx context.Context) (*T, error) {
	if it.index >= len(it.iterable) {
		return nil, Done
	}
	r := it.iterable[it.index]
	it.index++
	return &r, nil
}
