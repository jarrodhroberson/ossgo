package firestore

import (
	"context"

	fs "cloud.google.com/go/firestore"
)

const MAX_BULK_WRITE_SIZE = 20

type DatabaseName string

const DEFAULT DatabaseName = fs.DefaultDatabaseID

const DocumentCreated = "google.cloud.firestore.v1.created"
const DocumentUpdated = "google.cloud.firestore.v1.updated"
const DocumentDeleted = "google.cloud.firestore.v1.deleted"
const DocumentWritten = "google.cloud.firestore.v1.written"

// The op argument must be one of "==", "!=", "<", "<=", ">", ">=",
// "array-contains", "array-contains-any", "in" or "not-in"
type QueryOp string

func (q QueryOp) String() string {
	return string(q)
}

var QueryOps = struct {
	Equals              QueryOp
	NotEquals           QueryOp
	LessThan            QueryOp
	LessThanOrEqual     QueryOp
	GreaterThan         QueryOp
	GreaterThanOrEquals QueryOp
	ArrayContains       QueryOp
	ArrayContainsAny    QueryOp
	In                  QueryOp
	NotIn               QueryOp
}{
	Equals:              "==",
	NotEquals:           "!=",
	LessThan:            "<",
	LessThanOrEqual:     "<=",
	GreaterThan:         ">",
	GreaterThanOrEquals: ">=",
	ArrayContains:       "array-contains",
	ArrayContainsAny:    "array-contains-any",
	In:                  "in",
	NotIn:               "not-in",
}

type WherePredicate func(q fs.Query) fs.Query

type FirestoreRepository[T any] struct {
	client *fs.Client
}

func (f FirestoreRepository[T]) Load(ctx context.Context, key string) (*T, error) {
	docSnapshot, err := f.client.Collection("account").Doc(key).Get(ctx)
	if err != nil {
		return nil, err
	}
	var data T
	err = docSnapshot.DataTo(data)
	return &data, err
}

func (f FirestoreRepository[T]) Store(ctx context.Context, key string, value *T) (*T, error) {
	_, err := f.client.Collection("account").Doc(key).Set(ctx, value)
	if err != nil {
		return nil, err
	}
	return f.Load(ctx, key)

}

func (f FirestoreRepository[T]) Remove(ctx context.Context, key string) error {
	_, err := f.client.Collection("account").Doc(key).Delete(ctx)
	return err
}
