package persistence

import (
	"context"
)

type Repository[T any] interface {
	Load(ctx context.Context, key string) (*T, error)
	Store(ctx context.Context, key string, value *T) error
	Remove(ctx context.Context, key string) error
}
