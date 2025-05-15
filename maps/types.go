package maps

import (
	"sync"
	"time"
)

type Map[K comparable, V any] interface {
	Put(k string, v V) error
	Get(k string) (V, bool)
	Delete(k string) error
}

// item is a struct that holds the value and the last access time
type item[V any] struct {
	value      V
	lastAccess int64
}

// You can have a single map for an application or few maps for different purposes
type TTLMap[V any] struct {
	m map[string]*item[V]
	// For safe access to the map
	mu sync.Mutex
}

func New[V any](size int, maxTTL time.Duration, purgeInterval time.Duration) (m *TTLMap[V]) {
	// map is created with the given length
	m = &TTLMap[V]{m: make(map[string]*item[V], size)}

	// this goroutine will clean up the map from old items
	go func() {
		if len(m.m) > 0 {
			// You can adjust this ticker to be more or less frequent
			for now := range time.Tick(purgeInterval) {
				m.mu.Lock()
				for k, v := range m.m {
					if now.Unix()-v.lastAccess > int64(maxTTL) {
						delete(m.m, k)
					}
				}
				m.mu.Unlock()
			}
		}
	}()

	return
}

// Put adds a new item to the map or updates the existing one
func (m *TTLMap[T]) Put(k string, v T) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	it, ok := m.m[k]
	if !ok {
		it = &item[T]{
			value: v,
		}
	}
	it.value = v
	it.lastAccess = time.Now().Unix()
	m.m[k] = it
	return nil
}

// Get returns the value of the given key if it exists
func (m *TTLMap[V]) Get(k string) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if it, ok := m.m[k]; ok {
		it.lastAccess = time.Now().Unix()
		return it.value, true
	}

	return *new(V), false
}

// Delete removes the item from the map
func (m *TTLMap[T]) Delete(k string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.m[k]; ok {
		delete(m.m, k)
	}
	return nil
}
