package cloudstore

import (
	"context"

	"cloud.google.com/go/storage"
)

func Must(client *storage.Client, err error) *storage.Client {
	if err != nil {
		panic(err)
	} else {
		return client
	}
}

func Client(ctx context.Context) (*storage.Client, error) {
	return storage.NewClient(ctx)
}
