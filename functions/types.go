package functions

type Provider[T any] func() T

type ClosingProvider[t any] interface {
	Close() error
}