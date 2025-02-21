package containers

import (
	"iter"
)

type Keyer[T any] func(t *T) string

type Store[K string, V any] interface {
	All() (iter.Seq2[K, *V], error)
	Load(id string) (*V, error)
	Store(v *V) (*V, error)
	Remove(id string) error
}
