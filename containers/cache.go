package containers

import (
	"context"
	"errors"
	"regexp"

	"github.com/gomodule/redigo/redis"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"

	errs "github.com/jarrodhroberson/ossgo/errors"
	"github.com/jarrodhroberson/ossgo/functions/must"
)

//type Cache interface {
//	Set(entry Entry) error
//	Get(path string, labels map[string]string) ([]byte, error)
//	Keys() []string
//	Find(path string, labels map[string]string) (Entry, error)
//}

//type MapCache struct {
//	cache map[string]Entry
//}
//
//// Keys returns a sorted []string slice of all the keys in the cache
//func (m MapCache) Keys() []string {
//	return slices.Sorted(maps.Keys(m.cache))
//}
//
//func (m MapCache) Set(entry Entry) error {
//	m.cache[entry.Path] = entry
//	return nil
//}
//
//func (m MapCache) Get(path string, labels map[string]string) ([]byte, error) {
//	return m.cache[path].Data, nil
//}

type KeyMatcher func(key string) bool

func RegExKeyMatcher(expr string) KeyMatcher {
	matcher := regexp.MustCompile(expr)
	return func(key string) bool {
		return matcher.MatchString(key)
	}
}

func AllKeysMatcher() KeyMatcher {
	return func(key string) bool {
		return true
	}
}

func AllValuesMatcher[T any]() ValueMatcher[T] {
	return func(value T) bool {
		return true
	}
}

type ValueMatcher[T any] func(value T) bool

// Cache .Get returns a errors.NotFound() with a custom error messages if the key was
// not found in the cache.
type Cache[T any] interface {
	Set(key string, value T)
	Get(key string) (T, error)
	FindKeys(searchFunc KeyMatcher) ([]string, error)
	FindValues(searchFunc ValueMatcher[T]) ([]T, error)
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
		value, err = l.loadingFunc(key)
		if err != nil {
			return value, err
		} else {
			l.cache.Set(key, value)
			return value, nil
		}
	} else {
		return value, err
	}
}

func (l loadingCache[T]) Remove(key string) (T, error) {
	return l.cache.Remove(key)
}

type redisCache[T any] struct {
	redisClient *redis.Pool
}

func (r redisCache[T]) Set(key string, value T) {
	conn := r.redisClient.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			log.Warn().Err(err).Msgf("Error closing redis connection")
		}
	}(conn)

	b := must.MarshalJson(value)

	_, err := conn.Do("SET", key, b)
	if err != nil {
		log.Error().Err(err).Msgf("Error setting value for key [%s]", key)
	}
}

func (r redisCache[T]) Get(key string) (T, error) {
	conn := r.redisClient.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			log.Warn().Err(err).Msgf("Error closing redis connection")
		}
	}(conn)

	b, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return *new(T), errors.Join(err, errs.NotFoundError.New("value for key [%s] did not exist", key))
	}
	value := new(T)
	must.UnMarshalJson(b, value)
	return *value, nil
}

func (r redisCache[T]) Remove(key string) (T, error) {
	conn := r.redisClient.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			log.Warn().Err(err).Msgf("Error closing redis connection")
		}
	}(conn)

	b, err := redis.Bytes(conn.Do("DEL", key))
	if err != nil {
		return *new(T), err
	} else {
		value := new(T)
		must.UnMarshalJson(b, value)
		return *value, nil
	}
}
