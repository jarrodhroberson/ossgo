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
type ttlMap[V any] struct {
	m map[string]*item[V]
	// For safe access to the map
	mu sync.Mutex
}

func (m *ttlMap[V]) Put(k string, v V) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	it, ok := m.m[k]
	if !ok {
		it = &item[V]{
			value: v,
		}
	}
	it.value = v
	it.lastAccess = time.Now().Unix()
	m.m[k] = it
	return nil
}

func (m *ttlMap[V]) Get(k string) (V, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if it, ok := m.m[k]; ok {
		it.lastAccess = time.Now().Unix()
		return it.value, true
	}

	return *new(V), false
}

func (m *ttlMap[V]) Delete(k string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.m[k]; ok {
		delete(m.m, k)
	}
	return nil
}

func New[V any](size int, maxTTL time.Duration, purgeInterval time.Duration) (m Map[string, V]) {
	// map is created with the given length
	ttlm := ttlMap[V]{
		m: make(map[string]*item[V], size),
	}

	// this goroutine will clean up the map from old items
	go func() {
		if len(ttlm.m) > 0 {
			// You can adjust this ticker to be more or less frequent
			for now := range time.Tick(purgeInterval) {
				ttlm.mu.Lock()
				for k, v := range ttlm.m {
					if now.Unix()-v.lastAccess > int64(maxTTL) {
						delete(ttlm.m, k)
					}
				}
				ttlm.mu.Unlock()
			}
		}
	}()

	return &ttlm
}
