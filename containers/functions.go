package containers

import (
	"context"
)

func NewCache[T any](ctx context.Context) Cache[T] {
	return mapCache[T]{
		ctx,
		make(map[string]T),
	}
}

func AsCache[T any](ctx context.Context, m map[string]T) Cache[T] {
	return mapCache[T]{
		ctx:        ctx,
		backingMap: m,
	}
}

func WithLoading[T any](cache Cache[T], loadingFunction func(key string) (T, error)) Cache[T] {
	return loadingCache[T]{
		cache:       cache,
		loadingFunc: loadingFunction,
	}
}
