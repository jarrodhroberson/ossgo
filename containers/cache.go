package containers

import (
	"context"
	"errors"

	"github.com/gomodule/redigo/redis"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"

	errs "github.com/jarrodhroberson/ossgo/errors"
	"github.com/jarrodhroberson/ossgo/functions/must"
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
