package containers

type Keyer[T any] func(t *T) string

type Repository[K string, V any] interface {
	Get(id string) (*V, error)
	Store(v *V) (*V, error)
	Remove(id string) error
}
