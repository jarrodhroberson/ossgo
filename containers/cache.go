package containers

import (
	"time"

	"github.com/joomcode/errorx"
)

type Cache[T any] interface {
	Set(key string, value T)
	Get(key string) (T, error)
	Remove(key string) error
}

type LoadingCache[T any] interface {
	Load(key string) (T, error)
	LoadFunction(func(key string) (T, error))
}

type ExpiringCache[T any] interface {
	AbsoluteTimeToLive(ttl time.Duration)
	SinceLastAccess(ttl time.Duration)
}

type BoundedCache[T any] interface {
	AbsoluteMaxItemCount(maxItems int)
	MaxItemsWaiting(maxItems int, ttl time.Duration)
}

type MapCache[T any] struct {
	backingMap map[string]T
}

func (m MapCache[T]) Set(key string, value T) {
	//TODO: implement me
	panic(errorx.NotImplemented)
}

func (m MapCache[T]) Get(key string) (T, error) {
	//TODO: implement me
	panic(errorx.NotImplemented)
}
