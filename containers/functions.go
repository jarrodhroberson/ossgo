package containers

import (
	"context"

	"github.com/gomodule/redigo/redis"
)

func NewCache[T any](ctx context.Context) Cache[T] {
	return mapCache[T]{
		ctx,
		make(map[string]T),
	}
}

func RedisAsCache[T any](client *redis.Pool) Cache[T] {
	return redisCache[T]{
		redisClient: client,
	}
}

func MapAsCache[T any](ctx context.Context, m map[string]T) Cache[T] {
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
