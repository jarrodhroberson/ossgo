package repository

import (
	fs "github.com/jarrodhroberson/ossgo/firestore"
	"github.com/valkey-io/valkey-go"
)

func NewFirestoreRepository[T any](fsc fs.CollectionStore[T]) Repository[T] {
	return &firestoreRepository[T]{
		fsc: fsc,
	}
}

func NewValKeyRepository[T any](client valkey.Client, keyFunc func(string) string) Repository[T] {
	return &valKeyRepository[T]{
		vkc: client,
		keyFunc: keyFunc,
	}
}

func NewWrapRepository[T any](cache Repository[T], source Repository[T]) Repository[T] {
	return &wrapRepository[T]{
		cache:  cache,
		source: source,
	}
}