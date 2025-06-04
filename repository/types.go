package repository

import (
	"context"

	errs "github.com/jarrodhroberson/ossgo/errors"
	fs "github.com/jarrodhroberson/ossgo/firestore"
	"github.com/jarrodhroberson/ossgo/functions/must"
	vk "github.com/jarrodhroberson/ossgo/valkey"
	"github.com/joomcode/errorx"
	"github.com/rs/zerolog/log"
	"github.com/valkey-io/valkey-go"
)

type Repository[T any] interface {
	Get(key string) (*T, error)
	Set(key string, value *T) error
	Delete(key string) error
}

type valKeyRepository[T any] struct {
	vkc     valkey.Client
	keyFunc func(key string) string
}

func (v *valKeyRepository[T]) Get(key string) (*T, error) {
	ctx := context.Background()
	vkey := v.keyFunc(key)
	vkr := v.vkc.Do(ctx, v.vkc.B().JsonGet().Key(vkey).Path("$").Build())
	err := vk.ValkeyResultErrors(vkr)
	if err != nil {
		return nil, err
	}

	var t T
	err = vkr.DecodeJSON(&t)
	if err != nil {
		return nil, errs.ParseError.WrapWithNoMessage(err)
	}
	return &t, err
}

func (v *valKeyRepository[T]) Set(key string, value *T) error {
	ctx := context.Background()
	vkey := v.keyFunc(key)
	vkr := v.vkc.Do(ctx, v.vkc.B().JsonSet().Key(vkey).Path("$").Value(string(must.MarshalJson(value))).Build())
	return vk.ValkeyResultErrors(vkr)
}

func (v *valKeyRepository[T]) Delete(key string) error {
	ctx := context.Background()
	vkey := v.keyFunc(key)
	vkr := v.vkc.Do(ctx, v.vkc.B().Del().Key(vkey).Build())
	return vk.ValkeyResultErrors(vkr)
}

type firestoreRepository[T any] struct {
	fsc fs.CollectionStore[T]
}

func (f *firestoreRepository[T]) Get(key string) (*T, error) {
	return f.fsc.Load(key)
}

func (f *firestoreRepository[T]) Set(key string, value *T) error {
	_, err := f.fsc.Store(value)
	return err
}

func (f *firestoreRepository[T]) Delete(key string) error {
	return f.fsc.Remove(key)
}

type wrapRepository[T any] struct {
	cache  Repository[T]
	source Repository[T]
}

func (w *wrapRepository[T]) Get(key string) (*T, error) {
	v, err := w.cache.Get(key)
	if err == nil {
		return v, nil
	}

	if errorx.HasTrait(err, errorx.NotFound()) {
		v, err = w.source.Get(key)
		if err != nil {
			return nil, err
		}
		err = w.cache.Set(key, v)
		if err != nil {
			log.Error().Err(err).Msgf("failed to store value in cache: %s", key)
		}
	} else {
		err = errs.NotReadError.New("failed to get value from cache: %s", key).WithUnderlyingErrors(err)
		err = errorx.Decorate(err, "error is other than NOT FOUND")
		log.Error().Err(err).Msg(err.Error())
	}
	return v, err
}

func (w *wrapRepository[T]) Set(key string, value *T) error {
	err := w.cache.Set(key, value)
	if err != nil {
		err = errs.NotWrittenError.New("failed to write value from cache: %s", key).WithUnderlyingErrors(err)
		err = w.cache.Delete(key)
		if err != nil {
			err = errs.NotWrittenError.New("failed to delete value from cache: %s", key).WithUnderlyingErrors(err)
			log.Warn().Err(err).Msg(err.Error())
		}
		log.Error().Err(err).Msgf("failed to store value in cache: %s", key)
	}
	return w.source.Set(key, value)
}

func (w *wrapRepository[T]) Delete(key string) error {
	err := w.cache.Delete(key)
	if err != nil {
		err = errs.NotWrittenError.New("failed to delete value from cache: %s", key).WithUnderlyingErrors(err)
		log.Warn().Err(err).Msg(err.Error())
	}
	return w.source.Delete(key)
}
