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

type BulkStoreErrorHandling string

func (bs BulkStoreErrorHandling) ErrGroup(ctx context.Context) *errgroup.Group {
	if bs == FAIL_ON_FIRST_ERROR {
		if ctx == nil {
			ctx = context.Background()
		}
		eg, _ := errgroup.WithContext(ctx)
		return eg
	} else { // COLLECT_ERRORS
		return &errgroup.Group{}
	}
}

func (bs BulkStoreErrorHandling) String() string {
	return string(bs)
}

const (
	FAIL_ON_FIRST_ERROR BulkStoreErrorHandling = "fail_on_first_error"
	COLLECT_ERRORS      BulkStoreErrorHandling = "collect_errors"
)

// BulkStoreErrorHandling stores multiple items in batches using Firestore BulkWriter.
// It utilizes errgroup.Group to concurrently process batches of documents up to firestore.MAX_BULK_WRITE_SIZE.
// The errgroup ensures all goroutines complete and collects any errors that occur during batch processing.
// If any goroutine returns an error, BulkStoreErrorHandling will return that error after all goroutines complete.
//
// The errorHandling parameter determines how errors are handled:
// - FAIL_ON_FIRST_ERROR: Stop processing batches as soon as any error occurs
// - COLLECT_ERRORS: Continue processing remaining batches even if some fail report errors after iterator is complete
func (c collectionStore[T]) BulkStore(iter iter.Seq[*T], errorHandling BulkStoreErrorHandling) error {
	ctx := context.Background()
	client := c.clientProvider()
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg(err.Error())
		}
	}(client)

	bw := client.BulkWriter(ctx)
	defer func() {
		bw.Flush()
		bw.End()
	}()

	eg := errorHandling.ErrGroup(ctx)

	batches := seq.Chunk(iter, MAX_BULK_WRITE_SIZE)
	for batch := range batches {
		eg.Go(func() error {
			for item := range batch {
				docRef := client.Collection(c.collection).Doc(c.keyer(item))
				m := must.MarshallMap(item)
				m["last_updated_at"] = timestamp.Now()
				containers.RemoveKeys(m, "created_at")
				_, err := bw.Set(docRef, m)
				if err != nil {
					return BulkWriterError.Wrap(err, "error storing document \"%s\"", docRef)
				}
			}
			bw.Flush()
			return nil
		})
	}
	return eg.Wait()
}

func (c collectionStore[T]) Remove(id string) error {
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

func (c collectionStore[T]) BulkRemove(iter iter.Seq[string], errorHandling BulkStoreErrorHandling) error {
	ctx := context.Background()
	client := c.clientProvider()
	defer func(client *firestore.Client) {
		err := client.Close()
		if err != nil {
			log.Err(err).Msg(err.Error())
		}
	}(client)
	bw := client.BulkWriter(ctx)
	defer func() {
		bw.Flush()
		bw.End()
	}()

	eg := errorHandling.ErrGroup(ctx)
	batches := seq.Chunk(iter, MAX_BULK_WRITE_SIZE)
	for batch := range batches {
		eg.Go(func() error {
			for id := range batch {
				docRef := client.Collection(c.collection).Doc(id)
				_, err := bw.Delete(docRef)
				if err != nil {
					return BulkWriterError.Wrap(err, "error deleting document \"%s\"", docRef)
				}
			}
			bw.Flush()
			return nil
		})
	}
	return eg.Wait()
}

type CollectionStore[T any] interface {
	All() iter.Seq2[string, *T]
	Load(id string) (*T, error)
	Store(v *T) (*T, error)
	Remove(id string) error
}
