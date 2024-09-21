package containers

import (
	"context"

	"github.com/joomcode/errorx"

	errs "github.com/jarrodhroberson/ossgo/errors"
)

// Cache .Get returns a errors.NotFound() with a custom error messages if the key was
// not found in the cache.
type Cache[T any] interface {
	Set(key string, value T)
	Get(key string) (T, error)
	Remove(key string) (T, error)
}

type mapCache[T any] struct {
	ctx        context.Context
	backingMap map[string]T
}

func (m mapCache[T]) Set(key string, value T) {
	m.backingMap[key] = value
}

func (m mapCache[T]) Get(key string) (T, error) {
	value, exists := m.backingMap[key]
	if !exists {
		return value, errs.NotFoundError.New("value for key [%s] did not exist", key)
	} else {
		return value, nil
	}
}

func (m mapCache[T]) Remove(key string) (T, error) {
	value, exists := m.backingMap[key]
	if !exists {
		return value, errs.NotFoundError.New("value for key [%s] did not exist, nothing removed", key)
	} else {
		delete(m.backingMap, key)
		return value, nil
	}
}

type loadingCache[T any] struct {
	cache       Cache[T]
	loadingFunc func(key string) (T, error)
}

func (l loadingCache[T]) Set(key string, value T) {
	l.cache.Set(key, value)
}

func (l loadingCache[T]) Get(key string) (T, error) {
	value, err := l.cache.Get(key)
	if errorx.IsNotFound(err) {
		return l.load(key)
	} else {
		return value, err
	}
}

func (l loadingCache[T]) Remove(key string) (T, error) {
	return l.cache.Remove(key)
}

func (l loadingCache[T]) load(key string) (T, error) {
	value, err := l.loadingFunc(key)
	if err != nil {
		return value, err
	} else {
		l.cache.Set(key, value)
		return value, nil
	}
}
