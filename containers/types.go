package containers

import (
	"iter"
	"sync"
)

type SyncSet[T comparable] interface {
	Set(v T) bool
	All() iter.Seq[T]
}

type syncSet[T comparable] struct {
	delegate sync.Map
}

func (s *syncSet[T]) Set(v T) bool {
	if _, ok := s.delegate.Load(v); ok {
		return false
	} else {
		s.delegate.Store(v, struct{}{})
		return true
	}
}

func (s *syncSet[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		s.delegate.Range(func(key, value any) bool {
			if !yield(key.(T)) {
				return false
			}
			return true
		})
	}
}

// NewSyncSet creates a new thread-safe set implementation.
// The set can store any comparable type T and prevents duplicate entries.
// It is safe for concurrent access by multiple goroutines.
func NewSyncSet[T comparable]() SyncSet[T] {
	return &syncSet[T]{}
}
