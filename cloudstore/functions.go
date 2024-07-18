package cloudstore

import (
	"context"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func Must[T any](r T, err error) T {
	if err != nil {
		panic(err)
	} else {
		return r
	}
}

func Client(ctx context.Context, options ...option.ClientOption) (*storage.Client, error) {
	return storage.NewClient(ctx, options...)
}
