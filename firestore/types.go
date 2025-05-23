// Package firestore provides functionality for interacting with Google Cloud Firestore database
package firestore

import (
	"context"
	"iter"

	"cloud.google.com/go/firestore"
	fs "cloud.google.com/go/firestore"

	"github.com/jarrodhroberson/ossgo/containers"
	errs "github.com/jarrodhroberson/ossgo/errors"
	"github.com/jarrodhroberson/ossgo/functions"
	"github.com/jarrodhroberson/ossgo/functions/must"
	"github.com/jarrodhroberson/ossgo/timestamp"
	"github.com/rs/zerolog/log"
)

// MAX_BULK_WRITE_SIZE defines the maximum number of operations that can be performed in a single bulk write
// [maximum number of field transformations that can be performed on a single document in a Commit operation or in a transaction] : https://firebase.google.com/docs/firestore/quotas
const MAX_BULK_WRITE_SIZE = 500

// DatabaseName represents the name of a Firestore database
type DatabaseName string

// CollectionName represents the name of a Firestore collection
type CollectionName string

const DEFAULT DatabaseName = fs.DefaultDatabaseID

const DocumentCreated = "google.cloud.firestore.v1.created"
const DocumentUpdated = "google.cloud.firestore.v1.updated"
const DocumentDeleted = "google.cloud.firestore.v1.deleted"
const DocumentWritten = "google.cloud.firestore.v1.written"

// The op argument must be one of "==", "!=", "<", "<=", ">", ">=",
// "array-contains", "array-contains-any", "in" or "not-in"
// QueryOp represents a Firestore query operation
type QueryOp string

// String returns the string representation of the QueryOp
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

// WherePredicate is a function type that applies filtering conditions to a Firestore query
type WherePredicate func(q fs.Query) fs.Query

// Projection defines the fields to be returned in a query result
type Projection string

// String returns the string representation of the Projection
func (proj Projection) String() string {
	return string(proj)
}

const (
	// OnlyId returns only the document ID
	OnlyId Projection = "id_only"
	// OnlyIdLastUpdatedAt returns the document ID and last updated timestamp
	OnlyIdLastUpdatedAt Projection = "id_last_updated_at"
	// All returns all fields of the document
	All Projection = "all"
)

// CollectionStore provides CRUD operations for a specific Firestore collection
type CollectionStore[T any] struct {
	clientProvider functions.Provider[*firestore.Client]
	collection     string
	keyer          containers.Keyer[T]
}

func (c CollectionStore[T]) All() (iter.Seq2[string, *T], error) {
	ctx := context.Background()
	client := c.clientProvider()
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg(err.Error())
		}
	}(client)
	docIter := client.Collection(c.collection).Documents(ctx)
	dssSeq2 := DocSnapShotSeq2ToType[T](DocumentIteratorToSeq2(docIter))
	return ClosingWhenDoneSeq2(dssSeq2, client), nil
}

func (c CollectionStore[T]) Load(id string) (*T, error) {
	ctx := context.Background()
	client := c.clientProvider()
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg(err.Error())
		}
	}(client)

	docSnapshot, err := client.Collection(c.collection).Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}
	var t T
	err = docSnapshot.DataTo(&t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (c CollectionStore[T]) Store(v *T) (*T, error) {
	client := c.clientProvider()
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg(err.Error())
		}
	}(client)

	docRef := client.Collection(c.collection).Doc(c.keyer(v))
	m := must.MarshallMap(v)
	m["last_updated_at"] = timestamp.Now()
	containers.RemoveKeys(m, "created_at")
	ctx := context.Background()
	_, err := docRef.Set(ctx, m)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func (c CollectionStore[T]) Remove(id string) error {
	ctx := context.Background()
	client := c.clientProvider()
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg(err.Error())
		}
	}(client)

	_, err := client.Collection(c.collection).Doc(id).Delete(ctx)
	if err != nil {
		err = errs.NotDeletedError.Wrap(err, "failed to delete %s/%s", c.collection, id)
	}
	return err
}
