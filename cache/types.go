package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"slices"
	"sync"
	"time"

	"github.com/jarrodhroberson/ossgo/containers"
	fs "github.com/jarrodhroberson/ossgo/firestore"
	"github.com/rs/zerolog/log"
)

// CacheStrategy defines the caching strategy to use
type CacheStrategy string

const (
	// WriteThrough writes data to both cache and backing store synchronously
	WriteThrough CacheStrategy = "write-through"
	// WriteBehind writes data to cache immediately and asynchronously to backing store
	WriteBehind CacheStrategy = "write-behind"
	// ReadThrough reads data from cache, and if not found, fetches from backing store and caches it
	ReadThrough CacheStrategy = "read-through"
)

// CacheOptions defines configuration options for the cache
type CacheOptions struct {
	// TTL is the time-to-live for cache entries
	TTL time.Duration
	// Capacity is the maximum number of items in the cache
	Capacity uint64
	// Strategy is the caching strategy to use
	Strategy CacheStrategy
	// WriteBehindInterval is the interval at which write-behind operations are flushed to the backing store
	WriteBehindInterval time.Duration
}

// DefaultCacheOptions provides sensible defaults for cache options
var DefaultCacheOptions = CacheOptions{
	TTL:                 1 * time.Hour,
	Capacity:            10000,
	Strategy:            ReadThrough,
	WriteBehindInterval: 5 * time.Second,
}

// RedisClient is an interface for Redis/Valkey operations
type RedisClient interface {
	// Get retrieves a value from Redis
	Get(ctx context.Context, key string) (string, error)
	// Set stores a value in Redis
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Del deletes keys from Redis
	Del(ctx context.Context, keys ...string) error
	// Keys returns all keys matching the pattern
	Keys(ctx context.Context, pattern string) ([]string, error)
	// Close closes the Redis client
	Close() error
}

// CachedCollectionStore provides a caching layer between Redis/Valkey and Firestore
type CachedCollectionStore[T any] struct {
	// backingStore is the Firestore collection store
	backingStore fs.CollectionStore[T]
	// redisClient is the Redis/Valkey client
	redisClient RedisClient
	// options contains the cache configuration
	options CacheOptions
	// collection is the name of the collection
	collection string
	// keyer is a function that extracts the key from an item
	keyer containers.Keyer[T]
	// writeBehindQueue holds items that need to be written to the backing store
	writeBehindQueue map[string]*T
	// mutex protects the writeBehindQueue
	mutex sync.Mutex
	// ctx is the context for background operations
	ctx context.Context
	// cancel is the function to cancel background operations
	cancel context.CancelFunc
}

// processWriteBehindQueue periodically flushes the write-behind queue to the backing store
func (c *CachedCollectionStore[T]) processWriteBehindQueue() {
	ticker := time.NewTicker(c.options.WriteBehindInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.flushWriteBehindQueue()
		case <-c.ctx.Done():
			// Flush any remaining items before exiting
			c.flushWriteBehindQueue()
			return
		}
	}
}

// flushWriteBehindQueue writes all items in the write-behind queue to the backing store
func (c *CachedCollectionStore[T]) flushWriteBehindQueue() {
	c.mutex.Lock()
	// Create a copy of the queue and clear the original
	queue := make(map[string]*T, len(c.writeBehindQueue))
	for k, v := range c.writeBehindQueue {
		queue[k] = v
	}
	c.writeBehindQueue = make(map[string]*T)
	c.mutex.Unlock()

	// Process the queue
	for _, item := range queue {
		_, err := c.backingStore.Store(item)
		if err != nil {
			log.Error().Err(err).
				Str("collection", c.collection).
				Str("id", c.keyer(item)).
				Msg("Failed to write item to backing store in write-behind process")

			// Put the item back in the queue for retry
			c.mutex.Lock()
			c.writeBehindQueue[c.keyer(item)] = item
			c.mutex.Unlock()
		}
	}
}

// Close stops background processes and closes resources
func (c *CachedCollectionStore[T]) Close() error {
	c.cancel()
	return c.redisClient.Close()
}

// All returns all items from the collection
func (c *CachedCollectionStore[T]) All() iter.Seq2[string, *T] {
	// For All operations, we always go to the backing store
	// as it's impractical to cache all items
	return c.backingStore.All()
}

// Load retrieves an item by ID, using the configured caching strategy
func (c *CachedCollectionStore[T]) Load(id string) (*T, error) {
	// Check Redis cache
	if c.redisClient != nil {
		jsonData, err := c.redisClient.Get(c.ctx, c.cacheKey(id))
		if err == nil && jsonData != "" {
			var item T
			if err := json.Unmarshal([]byte(jsonData), &item); err == nil {
				return &item, nil
			}
		}
	}

	// If using ReadThrough or not found in cache, get from backing store
	if c.options.Strategy == ReadThrough || c.options.Strategy == WriteThrough {
		item, err := c.backingStore.Load(id)
		if err != nil {
			return nil, err
		}

		// Update Redis cache
		if c.redisClient != nil {
			jsonData, err := json.Marshal(item)
			if err == nil {
				c.redisClient.Set(c.ctx, c.cacheKey(id), string(jsonData), c.options.TTL)
			}
		}

		return item, nil
	}

	return nil, fmt.Errorf("item not found in cache and read-through not enabled")
}

// Find queries items with the given predicate and projection
func (c *CachedCollectionStore[T]) Find(where fs.WherePredicate, selectPaths fs.Projection) iter.Seq[*T] {
	// For Find operations, we always go to the backing store
	// as we can't effectively cache query results
	return c.backingStore.Find(where, selectPaths)
}

// Store stores an item using the configured caching strategy
func (c *CachedCollectionStore[T]) Store(v *T) (*T, error) {
	id := c.keyer(v)

	// Update Redis cache
	if c.redisClient != nil {
		jsonData, err := json.Marshal(v)
		if err == nil {
			c.redisClient.Set(c.ctx, c.cacheKey(id), string(jsonData), c.options.TTL)
		}
	}

	// Handle different write strategies
	switch c.options.Strategy {
	case WriteThrough:
		// Write directly to backing store
		return c.backingStore.Store(v)
	case WriteBehind:
		// Queue for async write
		c.mutex.Lock()
		c.writeBehindQueue[id] = v
		c.mutex.Unlock()
		return v, nil
	default:
		// For ReadThrough, we still write to backing store immediately
		return c.backingStore.Store(v)
	}
}

// BulkStore stores multiple items using the configured caching strategy
func (c *CachedCollectionStore[T]) BulkStore(iter iter.Seq[*T], errorHandling fs.BulkStoreErrorHandling) error {
	// For WriteBehind, we need to handle each item individually
	if c.options.Strategy == WriteBehind {
		var lastErr error
		for item := range iter {
			_, err := c.Store(item)
			if err != nil && errorHandling == fs.FAIL_ON_FIRST_ERROR {
				return err
			} else if err != nil {
				lastErr = err
			}
		}
		return lastErr
	}

	// For other strategies, update cache and delegate to backing store
	cacheItems := make([]*T, 0)
	for item := range iter {
		id := c.keyer(item)

		if c.redisClient != nil {
			jsonData, err := json.Marshal(item)
			if err == nil {
				c.redisClient.Set(c.ctx, c.cacheKey(id), string(jsonData), c.options.TTL)
			}
		}

		cacheItems = append(cacheItems, item)
	}

	return c.backingStore.BulkStore(slices.Values(cacheItems), errorHandling)
}

// Remove removes an item by ID
func (c *CachedCollectionStore[T]) Remove(id string) error {
	// Remove from Redis cache
	if c.redisClient != nil {
		c.redisClient.Del(c.ctx, c.cacheKey(id))
	}

	// For WriteBehind, we need to remove from the queue
	if c.options.Strategy == WriteBehind {
		c.mutex.Lock()
		delete(c.writeBehindQueue, id)
		c.mutex.Unlock()
	}

	// Remove from backing store
	return c.backingStore.Remove(id)
}

// BulkRemove removes multiple items by ID
func (c *CachedCollectionStore[T]) BulkRemove(iter iter.Seq[string], errorHandling fs.BulkStoreErrorHandling) error {
	// Collect IDs for batch operations
	ids := make([]string, 0)
	for id := range iter {
		// Remove from write-behind queue if applicable
		if c.options.Strategy == WriteBehind {
			c.mutex.Lock()
			delete(c.writeBehindQueue, id)
			c.mutex.Unlock()
		}

		ids = append(ids, id)
	}

	// Remove from Redis cache
	if c.redisClient != nil && len(ids) > 0 {
		cacheKeys := make([]string, len(ids))
		for i, id := range ids {
			cacheKeys[i] = c.cacheKey(id)
		}
		c.redisClient.Del(c.ctx, cacheKeys...)
	}

	// Remove from backing store
	return c.backingStore.BulkRemove(slices.Values(ids), errorHandling)
}

// BulkLoad loads multiple items by ID
func (c *CachedCollectionStore[T]) BulkLoad(iter iter.Seq[string]) iter.Seq2[*T, error] {
	return func(yield func(*T, error) bool) {
		for id := range iter {
			item, err := c.Load(id)
			if !yield(item, err) {
				return
			}
		}
	}
}

// cacheKey generates a Redis key for an item
func (c *CachedCollectionStore[T]) cacheKey(id string) string {
	return fmt.Sprintf("%s:%s", c.collection, id)
}

// PendingWrite represents an item that is pending write to the backing store
type PendingWrite[T any] struct {
	// Key is the item's key
	Key string
	// Value is the item's value
	Value *T
	// Timestamp is when the item was queued for writing
	Timestamp time.Time
}
