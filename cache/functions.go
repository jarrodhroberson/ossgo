// Package cache provides caching functionality between Valkey/Redis and Google Cloud Platform Firestore
package cache

import (
	"context"

	"github.com/jarrodhroberson/ossgo/containers"
	fs "github.com/jarrodhroberson/ossgo/firestore"
)

// NewCachedCollectionStore creates a new CachedCollectionStore with the specified options
func NewCachedCollectionStore[T any](
	backingStore fs.CollectionStore[T],
	redisClient RedisClient,
	collection string,
	keyer containers.Keyer[T],
	options CacheOptions,
) *CachedCollectionStore[T] {
	ctx, cancel := context.WithCancel(context.Background())

	store := &CachedCollectionStore[T]{
		backingStore:     backingStore,
		redisClient:      redisClient,
		options:          options,
		collection:       collection,
		keyer:            keyer,
		writeBehindQueue: make(map[string]*T),
		ctx:              ctx,
		cancel:           cancel,
	}

	// Start background worker for write-behind if that strategy is used
	if options.Strategy == WriteBehind {
		go store.processWriteBehindQueue()
	}

	return store
}
